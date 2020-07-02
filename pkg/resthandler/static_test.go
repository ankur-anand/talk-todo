// +build unit_tests all_tests

package resthandler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ankur-anand/prod-todo/pkg/logger"
)

func TestStaticHandler_Home(t *testing.T) {
	t.Parallel()
	l := logger.NewTesting(nil)
	defer l.Sync()
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	sh := newStaticHandler(l)
	handler := http.HandlerFunc(sh.home)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"message": "hello world from resthandler svc"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestStaticHandler_HealthAliveness(t *testing.T) {
	t.Parallel()
	l := logger.NewTesting(nil)
	defer l.Sync()
	req := httptest.NewRequest("GET", "/health/live", nil)
	rr := httptest.NewRecorder()
	sh := newStaticHandler(l)
	handler := http.HandlerFunc(sh.healthLive)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"alive": true}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
