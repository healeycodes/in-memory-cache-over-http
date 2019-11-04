package api

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestGetMiss(t *testing.T) {
	// Test cache miss
	new(100)

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
}

func TestGetHit(t *testing.T) {
	// Test cache hit
	KEY := "name"
	VALUE := "Andrew"

	new(100)

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
}

func TestSet(t *testing.T) {
	KEY := "name"
	VALUE := "Alice"

	new(100)

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

	if curValue, ok := s.Get(KEY); ok == false && curValue == VALUE {
		t.Errorf("Value wasn't set in cache")
	}
}

func TestDelete(t *testing.T) {
	KEY := "name"
	VALUE := "Alice"

	new(100)
	s.Set(KEY, VALUE, 0)

	param := make(url.Values)
	param["key"] = []string{KEY}

	req, err := http.NewRequest("GET", "/delete?"+param.Encode(), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Delete)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	if _, ok := s.Get(KEY); ok == true {
		t.Errorf("Value wasn't deleted")
	}
}

func TestCheckAndSetFail(t *testing.T) {
	KEY := "name"
	ORIGVALUE := "Mary"
	VALUE := "Alice"
	COMPARE := "NotAlice"

	new(100)
	s.Set(KEY, ORIGVALUE, 0)

	param := make(url.Values)
	param["key"] = []string{KEY}
	param["value"] = []string{VALUE}
	param["expire"] = []string{"0"}
	param["compare"] = []string{COMPARE}
	req, err := http.NewRequest("GET", "/checkandset?"+param.Encode(), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CheckAndSet)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	if curValue, _ := s.Get(KEY); curValue != ORIGVALUE {
		t.Errorf("Value was set even though the compare should have failed")
	}
}

func TestCheckAndSetOk(t *testing.T) {
	KEY := "name"
	ORIGVALUE := "Mary"
	VALUE := "Alice"
	COMPARE := ORIGVALUE

	new(100)
	s.Set(KEY, ORIGVALUE, 0)

	param := make(url.Values)
	param["key"] = []string{KEY}
	param["value"] = []string{VALUE}
	param["expire"] = []string{"0"}
	param["compare"] = []string{COMPARE}
	req, err := http.NewRequest("GET", "/checkandset?"+param.Encode(), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CheckAndSet)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusNoContent)
	}

	if curValue, _ := s.Get(KEY); curValue == ORIGVALUE {
		t.Errorf("Value wasn't set even though the compare should have passed")
	}
}

func TestIncrement(t *testing.T) {
	KEY := "hits"
	VALUE := "1"
	ADDVALUE := "1"
	RESULT := "2"

	new(100)
	s.Set(KEY, VALUE, 0)

	param := make(url.Values)
	param["key"] = []string{KEY}
	param["value"] = []string{ADDVALUE}
	req, err := http.NewRequest("GET", "/increment?"+param.Encode(), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Increment)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusNoContent)
	}

	if curValue, _ := s.Get(KEY); curValue != RESULT {
		t.Errorf("Value wasn't incremented")
	}
}

func TestDecrement(t *testing.T) {
	KEY := "hits"
	VALUE := "1"
	MINUSVALUE := "1"
	RESULT := "0"

	new(100)
	s.Set(KEY, VALUE, 0)

	param := make(url.Values)
	param["key"] = []string{KEY}
	param["value"] = []string{MINUSVALUE}
	req, err := http.NewRequest("GET", "/decrement?"+param.Encode(), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Decrement)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusNoContent)
	}

	if curValue, _ := s.Get(KEY); curValue != RESULT {
		t.Errorf("Value wasn't decremented")
	}
}

func TestAppend(t *testing.T) {
	KEY := "name"
	VALUE := "And"
	APPENDVALUE := "y"
	RESULT := "Andy"

	new(100)
	s.Set(KEY, VALUE, 0)

	param := make(url.Values)
	param["key"] = []string{KEY}
	param["value"] = []string{APPENDVALUE}
	req, err := http.NewRequest("GET", "/append?"+param.Encode(), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Append)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusNoContent)
	}

	if curValue, _ := s.Get(KEY); curValue != RESULT {
		t.Errorf("Value wasn't appended")
	}
}

func TestPrepend(t *testing.T) {
	KEY := "name"
	VALUE := "ndy"
	PREPENDVALUE := "A"
	RESULT := "Andy"

	new(100)
	s.Set(KEY, VALUE, 0)

	param := make(url.Values)
	param["key"] = []string{KEY}
	param["value"] = []string{PREPENDVALUE}
	req, err := http.NewRequest("GET", "/prepend?"+param.Encode(), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Prepend)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusNoContent)
	}

	if curValue, _ := s.Get(KEY); curValue != RESULT {
		t.Errorf("Value wasn't prepended")
	}
}

func TestFlush(t *testing.T) {
	KEY := "name"
	VALUE := "Alice"

	new(100)
	s.Set(KEY, VALUE, 0)

	req, err := http.NewRequest("GET", "/flush", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Flush)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusNoContent)
	}

	if _, exist := s.Get(KEY); exist != false {
		t.Errorf("Cache wasn't flushed")
	}
}

func TestStats(t *testing.T) {
	KEY := "name"
	VALUE := "Mary"

	new(100)
	s.Set(KEY, VALUE, 0)

	req, err := http.NewRequest("GET", "/stats", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Stats)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	//TODO (healeycodes)
	// Could parse the stats and check but for now just sanity check the path
}
