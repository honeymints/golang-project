package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHomeHandler(t *testing.T) {

	/* req, err := http.NewRequest("GET", "http://localhost:3000/login-register/", nil)
	if err != nil {
		t.Fatalf("Couldn't create request %v", err)
	}
	rec := httptest.NewRecorder()

	loginregisterHandler(rec, req)

	t.Logf("Request have been created, %v", req.Header)

	res := rec.Result()

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status OK, got %v", res.Status)
	}
	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Couldn't read response body %v", err)
	} */
	server := httptest.NewServer(http.HandlerFunc(loginregisterHandler))
	responce, err := http.Get(server.URL)
	if err != nil {
		t.Error(err)
	}
	if responce.StatusCode != http.StatusOK {
		t.Errorf("got %v, expected OK", responce.Status)
	}
}
