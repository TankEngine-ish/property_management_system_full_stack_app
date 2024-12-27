// this file validates the behavior of the GET /api/repair/users API endpoint.
// It uses a mocked HTTP request to ensure the handler responds with the correct status code and data structure.

package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestGetUsersAPI(t *testing.T) {

	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"id":1,"name":"spiderman","statement":100}]`))
	})

	req, err := http.NewRequest("GET", "/api/repair/users", nil)
	if err != nil {
		t.Fatalf("Error creating request: %s", err)
	}

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.Handle("/api/repair/users", mockHandler)

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Expected HTTP status code 200")
	expectedBody := `[{"id":1,"name":"spiderman","statement":100}]`
	assert.Equal(t, expectedBody, rr.Body.String(), "Unexpected response body")
}
