package httpapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestCustomers_CreateAndGet(t *testing.T) {
	// отдельный data.json для теста
	tmp, err := os.CreateTemp("", "go-crm-data-*.json")
	if err != nil {
		t.Fatalf("temp file: %v", err)
	}
	path := tmp.Name()
	_ = tmp.Close()
	defer os.Remove(path)

	// ВАЖНО: Routes() сейчас жестко использует "./data.json".
	// Для тестов нам лучше сделать RoutesWithDataPath(path).
	// Пока сделаем быстрый фокус: запускаем тест в tmp-dir и создадим data.json там.
	dir := t.TempDir()
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	h := Routes()

	// 1) POST /customers
	body := []byte(`{"name":"Aktan","email":"aktan@example.com"}`)
	req := httptest.NewRequest(http.MethodPost, "/customers", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected %d got %d body=%s", http.StatusCreated, rr.Code, rr.Body.String())
	}

	var created map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &created); err != nil {
		t.Fatalf("unmarshal: %v body=%s", err, rr.Body.String())
	}
	id, _ := created["id"].(string)
	if id == "" {
		t.Fatalf("expected id in response, got: %v", created)
	}

	// 2) GET /customers/{id}
	req2 := httptest.NewRequest(http.MethodGet, "/customers/"+id, nil)
	rr2 := httptest.NewRecorder()
	h.ServeHTTP(rr2, req2)

	if rr2.Code != http.StatusOK {
		t.Fatalf("expected %d got %d body=%s", http.StatusOK, rr2.Code, rr2.Body.String())
	}
}
