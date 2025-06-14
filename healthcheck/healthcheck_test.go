package healthcheck

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHealthcheck(t *testing.T) {
	router := gin.Default()
	router.GET("/healthcheck", Healthcheck)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, httptest.NewRequest("GET", "/healthcheck", nil))

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := map[string]interface{}{
		"status":  "ok",
		"message": "Firelink is running",
	}
	var actual map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &actual); err != nil {
		t.Fatalf("could not unmarshal response: %v", err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("handler returned unexpected body: got %v want %v", actual, expected)
	}
}
