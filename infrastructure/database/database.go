package database

import (
	"fmt"
	"os"
	"strings"
	"time"

	stdlog "log"

	"gin-boilerplate/infrastructure/config"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// jejakSQL menyusun logger GORM.
//
// Dibangun sendiri alih-alih memakai logger.Default karena dua bawaannya keliru
// untuk layanan yang lognya dikirim ke agregator, bukan dibaca di terminal:
//
//  1. IgnoreRecordNotFoundError bawaannya false, sehingga ErrRecordNotFound
//     dicatat sebagai ERROR lengkap dengan SQL dan nilai yang diikat. Padahal
//     pada banyak kueri "tidak ketemu" adalah hasil NORMAL yang sudah ditangani
//     pemanggilnya. Terlihat di produksi 21 Agu 2026: pemeriksaan berkas kembar
//     memakai First() dan menangani ErrRecordNotFound sebagai "tidak ada
//     kembaran", tetapi SETIAP unggahan yang bukan duplikat tetap menulis baris
//     error berisi id_sdm dan sidik berkasnya.
//  2. Prefiks bawaannya carriage-return + newline, yang di agregator
//     menjadi satu BARIS KOSONG tersendiri sebelum tiap jejak.
//     Dipakai prefiks kosong di sini.
//  3. Colorful bawaannya true, sehingga tiap baris membawa escape ANSI
//     ([31;1m dan kawan-kawan). Berguna di terminal, sampah di Loki.
func jejakSQL(logMode bool, logLevel string) logger.Interface {
	return logger.New(
		stdlog.New(os.Stdout, "", stdlog.LstdFlags),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  levelJejakSQL(logMode, logLevel),
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)
}

// levelJejakSQL menentukan seberapa banyak GORM mencatat kueri.
//
// Hanya DEBUG yang mendapat logger.Info, karena level itu mencatat SETIAP kueri
// LENGKAP DENGAN NILAI YANG DIIKAT — isi tabelnya, bukan sekadar bentuk
// kuerinya. Pada turunan boilerplate ini nilai tersebut adalah data pribadi:
// NIP, nama pegawai, nama berkas, alamat IP. Dua akibatnya sama-sama tidak
// diinginkan: aliran log dibanjiri kueri yang BERHASIL sampai kegagalan
// sungguhan tenggelam, dan data pribadi tersalin ke tempat dengan retensi serta
// hak akses yang berbeda dari basis data asalnya.
//
// Silent bukan jawabannya untuk selain DEBUG: ia ikut membungkam kueri lambat
// dan error basis data, justru dua hal yang paling dibutuhkan di produksi.
// Silent hanya dipakai bila operator memang mematikan logMode.
func levelJejakSQL(logMode bool, logLevel string) logger.LogLevel {
	switch {
	case !logMode:
		return logger.Silent
	case strings.EqualFold(strings.TrimSpace(logLevel), "DEBUG"):
		return logger.Info
	default:
		return logger.Warn
	}
}

// namaLevelJejakSQL menerjemahkan level GORM ke nama yang terbaca di log,
// beserta apa yang sebenarnya ikut tercatat pada level itu.
func namaLevelJejakSQL(level logger.LogLevel) string {
	switch level {
	case logger.Silent:
		return "silent (tidak ada)"
	case logger.Info:
		return "info (SEMUA kueri beserta nilainya)"
	case logger.Warn:
		return "warn (kueri lambat + error)"
	default:
		return "error (error saja)"
	}
}

type DatabaseClient struct {
	DbConn *gorm.DB
}

func NewDatabaseClient(config *config.Config) (*DatabaseClient, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		config.Postgres.Host,
		config.Postgres.User,
		config.Postgres.Password,
		config.Postgres.DBName,
		config.Postgres.Port,
		config.Postgres.SSLMode,
		config.Postgres.TimeZone,
	)

	level := levelJejakSQL(config.LogMode, config.LogLevel)
	gormLogger := jejakSQL(config.LogMode, config.LogLevel)

	// Dilaporkan saat boot supaya pertanyaan "kenapa log penuh kueri?" terjawab
	// dari log itu sendiri, bukan dari menebak isi config produksi yang tidak
	// dapat dibaca dari luar. Menyebut juga NILAI ASAL, karena penyebab tersering
	// adalah logLevel yang tidak disetel sama sekali sehingga jatuh ke bawaan.
	log.WithFields(log.Fields{
		"log_mode":  config.LogMode,
		"log_level": config.LogLevel,
		"jejak_sql": namaLevelJejakSQL(level),
	}).Info("Jejak SQL GORM disetel")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying sql.DB to configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Test the connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Info("Successfully connected to PostgreSQL database")

	return &DatabaseClient{
		DbConn: db,
	}, nil
}

func (dc *DatabaseClient) Close() error {
	sqlDB, err := dc.DbConn.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// AutoMigrate runs database migrations for the given models
func (dc *DatabaseClient) AutoMigrate(models ...interface{}) error {
	return dc.DbConn.AutoMigrate(models...)
}
