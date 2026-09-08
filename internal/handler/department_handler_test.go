package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"ticket-api/internal/testutil"
)

func TestDepartmentHandler_Endpoints(t *testing.T) {
	app := testutil.SetupTestApp(t, nil)
	defer app.DB.Close()

	t.Run("GET /api/v1/departments/GetAllActiveDepartments/ -> 200", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/departments/GetAllActiveDepartments/", nil)
		w := httptest.NewRecorder()

		app.Engine.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200 OK, got %d", w.Code)
		}
	})
}
