package database

import (
	"testing"

	"gorm.io/gorm/logger"
)

// Jejak SQL GORM pada level Info memuat NILAI YANG DIIKAT — isi tabel, bukan
// sekadar bentuk kuerinya. Pada turunan boilerplate ini nilai tersebut adalah
// data pribadi. Uji ini menjaga agar hanya DEBUG yang mendapatnya, supaya
// proyek baru tidak lagi mewarisi produksi yang membocorkannya diam-diam.
func TestLevelJejakSQL(t *testing.T) {
	kasus := []struct {
		nama     string
		logMode  bool
		logLevel string
		mau      logger.LogLevel
		alasan   string
	}{
		{"logMode mati", false, "DEBUG", logger.Silent, "operator sengaja membungkam"},
		{"DEBUG dapat semuanya", true, "DEBUG", logger.Info, "pengembangan butuh kueri lengkap"},
		{"debug huruf kecil", true, "debug", logger.Info, "level tidak boleh peka huruf besar-kecil"},
		{"DEBUG berspasi", true, " DEBUG ", logger.Info, "spasi di YAML tidak boleh mengubah arti"},
		{"INFO hanya lambat & error", true, "INFO", logger.Warn, "produksi tidak boleh mencatat data pribadi"},
		{"WARN hanya lambat & error", true, "WARN", logger.Warn, ""},
		{"ERROR hanya lambat & error", true, "ERROR", logger.Warn, ""},
		{"kosong dianggap bukan DEBUG", true, "", logger.Warn, "bawaan harus yang aman, bukan yang bocor"},
	}

	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			if dapat := levelJejakSQL(k.logMode, k.logLevel); dapat != k.mau {
				t.Errorf("levelJejakSQL(%t, %q) = %v, mau %v — %s", k.logMode, k.logLevel, dapat, k.mau, k.alasan)
			}
		})
	}
}

func TestNamaLevelJejakSQL_MenyebutApaYangTercatat(t *testing.T) {
	// Baris boot ini yang akan dibaca operator saat log penuh kueri, jadi
	// namanya harus menjawab "apa yang ikut tercatat", bukan sekadar angka level.
	if n := namaLevelJejakSQL(logger.Info); n != "info (SEMUA kueri beserta nilainya)" {
		t.Errorf("nama level Info = %q — harus memperingatkan bahwa nilai ikut tercatat", n)
	}
	if n := namaLevelJejakSQL(logger.Warn); n != "warn (kueri lambat + error)" {
		t.Errorf("nama level Warn = %q", n)
	}
	if n := namaLevelJejakSQL(logger.Silent); n != "silent (tidak ada)" {
		t.Errorf("nama level Silent = %q", n)
	}
}
