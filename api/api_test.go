package api

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestGet(t *testing.T) {
	// Test cache miss
	func() {
		new()

		param := make(url.Values)
		param["key"] = []string{"name"}
		req, err := http.NewRequest("GET", "/get?"+param.Encode(), nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(Get)

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("Handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
	}()

	// Test cache hit
	func() {
		KEY := "name"
		VALUE := "Andrew"
		new()

		s.Set(KEY, VALUE, 0)

		param := make(url.Values)
		param["key"] = []string{KEY}
		req, err := http.NewRequest("GET", "/get?"+param.Encode(), nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(Get)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		expected := VALUE
		if rr.Body.String() != expected {
			t.Errorf("Handler returned unexpected body: got %v want %v",
				rr.Body.String(), expected)
		}
	}()
}

func TestSet(t *testing.T) {
	KEY := "name"
	VALUE := "Alice"
	new()

	param := make(url.Values)
	param["key"] = []string{KEY}
	param["value"] = []string{VALUE}
	param["expire"] = []string{"0"}
	req, err := http.NewRequest("GET", "/set?"+param.Encode(), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Set)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := ""
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	if curValue, ok := s.Get(KEY); ok == false && curValue == VALUE {
		t.Errorf("Value wasn't set in cache")
	}
}