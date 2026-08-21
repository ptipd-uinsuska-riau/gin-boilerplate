package database

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Komentar pada jejakSQL menjanjikan empat hal, dan sebelumnya tidak satu pun
// diuji: hanya DEBUG yang mencatat isi kueri, ErrRecordNotFound tidak dicatat
// sebagai kegagalan, tanpa escape ANSI, dan tanpa baris kosong di depan.
// Ditambah satu yang tersirat — kueri lambat TETAP tercatat di luar DEBUG,
// yang justru alasan Silent bukan jawabannya untuk produksi.
//
// Uji levelJejakSQL hanya menjaga nilai kembalian, bukan satu pun dari kelimanya.
// Janji anti-bocor data pribadi — inti dari perbaikan ini — selama ini hanya
// tertulis di komentar, sehingga penyederhanaan jejakSQL suatu hari tidak akan
// meledakkan apa pun.
//
// Logger diambil DARI konfigurasiGorm, bukan dibangun ulang: membangun ulang
// berarti kembali menguji potongan implementasi, dan menghapus logger dari
// konfigurasiGorm tidak akan tertangkap.
func loggerUji(t *testing.T, logLevel string) (logger.Interface, *bytes.Buffer) {
	t.Helper()
	var buf bytes.Buffer
	var cfg *gorm.Config = konfigurasiGorm(jejakSQLKe(&buf, true, logLevel))
	if cfg.Logger == nil {
		t.Fatal("konfigurasiGorm tidak memasang logger — GORM akan memakai bawaannya, dan seluruh janji di atas batal")
	}
	return cfg.Logger, &buf
}

const sqlRahasia = `SELECT * FROM pegawai WHERE nip = '198501012010011001' AND nama = 'BUDI SANTOSO'`

func telusur(l logger.Interface, sejak time.Duration, err error) {
	l.Trace(context.Background(), time.Now().Add(-sejak),
		func() (string, int64) { return sqlRahasia, 1 }, err)
}

// JANJI 1 — dan yang paling penting: di luar DEBUG, isi kueri TIDAK boleh keluar.
func TestJejakSQL_IsiKueriHanyaBocorPadaDEBUG(t *testing.T) {
	l, buf := loggerUji(t, "INFO")
	telusur(l, time.Millisecond, nil)
	if strings.Contains(buf.String(), "198501012010011001") {
		t.Errorf("NIP ikut tercatat pada level INFO — janji anti-bocor batal:\n%s", buf.String())
	}

	l2, buf2 := loggerUji(t, "DEBUG")
	telusur(l2, time.Millisecond, nil)
	if !strings.Contains(buf2.String(), "198501012010011001") {
		t.Error("DEBUG seharusnya mencatat kueri lengkap; kalau tidak, penelusuran lokal jadi mustahil")
	}
}

// JANJI 2 — ErrRecordNotFound bukan kegagalan. Sengaja diuji di luar DEBUG:
// pada level Info SETIAP kueri ditrace sebagai kueri biasa terlepas dari
// galatnya, sehingga perbedaannya tidak terlihat di sana.
func TestJejakSQL_RecordNotFoundBukanKegagalan(t *testing.T) {
	l, buf := loggerUji(t, "INFO")
	telusur(l, time.Millisecond, gorm.ErrRecordNotFound)
	if buf.Len() != 0 {
		t.Errorf("record not found ikut tercatat, padahal itu hasil normal:\n%s", buf.String())
	}
}

// JANJI TERSIRAT — kueri lambat TETAP tercatat di luar DEBUG. Inilah alasan
// Silent bukan jawabannya untuk produksi.
func TestJejakSQL_KueriLambatTetapTercatatDiLuarDEBUG(t *testing.T) {
	l, buf := loggerUji(t, "INFO")
	telusur(l, 300*time.Millisecond, nil)
	if buf.Len() == 0 {
		t.Error("kueri lambat tidak tercatat — kemampuan mendeteksi kelambatan di produksi hilang")
	}
}

// JANJI 3 & 4 — keluarannya untuk agregator, bukan terminal.
func TestJejakSQL_TanpaANSIDanTanpaBarisKosong(t *testing.T) {
	l, buf := loggerUji(t, "DEBUG")
	telusur(l, time.Millisecond, nil)
	keluaran := buf.String()
	if strings.Contains(keluaran, "\x1b[") {
		t.Errorf("keluaran memuat escape ANSI — sampah di agregator log:\n%q", keluaran)
	}
	if strings.HasPrefix(keluaran, "\n") || strings.HasPrefix(keluaran, "\r\n") {
		t.Errorf("keluaran diawali baris kosong — menjadi entri kosong tersendiri di agregator:\n%q", keluaran)
	}
}
