package httplib

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
)

func konteks(t *testing.T) (*gin.Context, *test.Hook) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("GET", "/api/v1/contoh", nil)
	kait := test.NewLocal(logrus.StandardLogger())
	return c, kait
}

// Kegagalan 5xx yang dijawab langsung lewat SetErrorResponse tidak pernah
// melewati ErrorHandler, sehingga sebelumnya hilang tanpa jejak.
func TestCatatKegagalanServer_5xxTanpaErrorHandlerDicatat(t *testing.T) {
	c, kait := konteks(t)
	catatKegagalanServer(c, 500, "gagal menghubungi basis data")
	if n := len(kait.AllEntries()); n != 1 {
		t.Fatalf("%d baris log, mau 1 — kegagalan server harus terlihat", n)
	}
	if lv := kait.LastEntry().Level; lv != logrus.ErrorLevel {
		t.Errorf("level=%v, mau error", lv)
	}
}

// Bila error sudah ditaruh di c.Errors, ErrorHandler yang mencatatnya.
// Mencatat lagi di sini hanya melipatgandakan baris log.
func TestCatatKegagalanServer_TidakGandaSaatLewatErrorHandler(t *testing.T) {
	c, kait := konteks(t)
	_ = c.Error(gin.Error{Err: errContoh{}, Type: gin.ErrorTypePrivate})
	catatKegagalanServer(c, 500, "gagal menghubungi basis data")
	if n := len(kait.AllEntries()); n != 0 {
		t.Errorf("%d baris log, mau 0 — ErrorHandler sudah mencatatnya", n)
	}
}

// 4xx adalah kesalahan pemanggil. Mencatatnya membuat UUID salah ketik ikut
// menjadi baris ERROR.
func TestCatatKegagalanServer_4xxTidakDicatat(t *testing.T) {
	c, kait := konteks(t)
	catatKegagalanServer(c, 400, "id tidak valid")
	catatKegagalanServer(c, 404, "tidak ditemukan")
	if n := len(kait.AllEntries()); n != 0 {
		t.Errorf("%d baris log, mau 0", n)
	}
}

type errContoh struct{}

func (errContoh) Error() string { return "contoh" }

// Uji di atas memanggil catatKegagalanServer LANGSUNG, sehingga tidak menjaga
// apakah SetErrorResponse benar-benar memanggilnya. Dibuktikan lewat kontrol
// negatif: sambungan itu diputus, dan seluruh uji di atas tetap lolos.
//
// Dua uji berikut memakai pintu masuk yang sesungguhnya, jadi memutus
// sambungannya membuat keduanya gagal.
func TestSetErrorResponse_TersambungKeCatatan(t *testing.T) {
	c, kait := konteks(t)
	SetErrorResponse(c, 0, 500, "gagal menghubungi basis data")
	if n := len(kait.AllEntries()); n != 1 {
		t.Fatalf("%d baris log, mau 1 — SetErrorResponse harus mencatat 5xx", n)
	}
}

func TestSetErrorResponseWithError_TersambungKeCatatan(t *testing.T) {
	c, kait := konteks(t)
	SetErrorResponseWithError(c, 0, 503, "layanan tidak tersedia", errContoh{})
	if n := len(kait.AllEntries()); n != 1 {
		t.Fatalf("%d baris log, mau 1 — SetErrorResponseWithError harus mencatat 5xx", n)
	}
}
