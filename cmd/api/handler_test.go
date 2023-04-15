package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
)

func TestLoginRegisterhandler(t *testing.T) {
	app := &application{}
	test := []struct {
		name   string
		status int
		err    string
	}{
		{name: "login", status: http.StatusOK},
	}
	for _, ts := range test {

		rec := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "", nil)

		if err != nil {
			t.Error(err)
		}

		app.loginregisterHandler(rec, req)

		res := rec.Result()

		if res.StatusCode != ts.status {
			t.Fatalf("got %v, expected OK", res.Status)
		}
		defer res.Body.Close()

		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Error(err)
		}
		t.Logf("worked fine, %v", string(b))
		if string(b) != ts.name {
			t.Fatalf("couldn't read response body, expected: %s, got %s", ts.name, string(b))
		}
	}
}

func TestActivatedHandler(t *testing.T) {
	app := &application{}
	test := []struct {
		name   string
		status int
		err    string
	}{
		{name: "activated", status: http.StatusOK},
	}
	for _, ts := range test {

		rec := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "", nil)

		if err != nil {
			t.Error(err)
		}

		app.activateUserHandler(rec, req)

		res := rec.Result()

		if res.StatusCode != ts.status {
			t.Fatalf("got %v, expected OK", res.Status)
		}
		defer res.Body.Close()

		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Error(err)
		}
		t.Logf("worked fine, %v", string(b))
		if string(b) != ts.name {
			t.Fatalf("couldn't read response body, expected: %s, got %s", ts.name, string(b))
		}
	}
}

func TestRouting(t *testing.T) {
	test := []struct {
		name   string
		status int
		route  string
		err    string
	}{
		{name: "login", status: http.StatusOK, route: "login-register"},
		{name: "home", status: http.StatusOK, route: ""},
		{name: "activated", status: http.StatusOK, route: "activated/:token"},
	}
	app := &application{}
	for _, ts := range test {
		s := httptest.NewServer(app.routes())
		defer s.Close()

		res, err := http.Get(fmt.Sprintf("%v/%s", s.URL, ts.route))
		if err != nil {
			t.Fatalf("couldn't get GET method, %v", err)
		}
		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Error(err)
		}
		if string(b) != ts.name {
			t.Fatalf("couldn't read response body, expected: %s, got %s", ts.name, string(b))
		}
		t.Logf("worked fine, %v", string(b))
	}
}

func TestAuthenticate(t *testing.T) {
	app := &application{}
	test := []struct {
		name   string
		status int
		err    string
	}{
		{name: "authenticate", status: http.StatusCreated},
	}

	for _, ts := range test {
		router := httprouter.New()
		input := struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}{
			Email:    "aruzhanism@mail.ru",
			Password: "password",
		}

		bodyRequest, _ := json.Marshal(input)

		router.HandlerFunc(http.MethodPost, "/myacc", app.createAuthenticationTokenHandler)
		rec := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/myacc", bytes.NewBuffer(bodyRequest))

		if err != nil {
			t.Fatal(err)
		}
		app.createAuthenticationTokenHandler(rec, req)

		res := rec.Result()

		if res.StatusCode != ts.status {
			t.Fatalf("got %v, expected Created", res.Status)
		}
		defer res.Body.Close()
	}
}

func TestRegisterHandler(t *testing.T) {
	app := &application{}
	test := []struct {
		name   string
		status int
		err    string
	}{
		{name: "welcome", status: http.StatusCreated},
	}

	for _, ts := range test {
		router := httprouter.New()
		input := struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}{
			Name:     "aruzhan",
			Email:    "aruzhanism@mail.ru",
			Password: "password",
		}

		bodyRequest, _ := json.Marshal(input)

		router.HandlerFunc(http.MethodPost, "/welcome", app.createAuthenticationTokenHandler)
		rec := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/welcome", bytes.NewBuffer(bodyRequest))

		if err != nil {
			t.Fatal(err)
		}
		app.welcomeHandler(rec, req)

		res := rec.Result()

		if res.StatusCode != ts.status {
			t.Fatalf("got %v, expected Created", res.Status)
		}
		defer res.Body.Close()

		t.Logf("data:%v", input)
	}
}

func TestHomehandler(t *testing.T) {
	app := &application{}
	test := []struct {
		name   string
		status int
		err    string
	}{
		{name: "home", status: http.StatusOK},
	}
	for _, ts := range test {

		rec := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/", nil)

		if err != nil {
			t.Error(err)
		}

		app.homeHandler(rec, req)

		res := rec.Result()

		if res.StatusCode != ts.status {
			t.Fatalf("got %v, expected OK", res.Status)
		}
		defer res.Body.Close()

		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Error(err)
		}
		t.Logf("worked fine, %v", string(b))
		if string(b) != ts.name {
			t.Fatalf("couldn't read response body, expected: %s, got %s", ts.name, string(b))
		}
	}
}
