package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

const johnToken = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJPbmxpbmUgSldUIEJ1aWxkZXIiLCJpYXQiOjE3MTMyNjM5MDAsImV4cCI6MTc0NDc5OTkwMCwiYXVkIjoid3d3LmV4YW1wbGUuY29tIiwic3ViIjoiam9obiJ9._Mdu2Gvz6QsApNpACSZfwIJTOP1ZoKJADmXAHGqHJMc"
const wadeToken = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJPbmxpbmUgSldUIEJ1aWxkZXIiLCJpYXQiOjE3MTMyNjM5MDAsImV4cCI6MTc0NDc5OTkwMCwiYXVkIjoid3d3LmV4YW1wbGUuY29tIiwic3ViIjoid2FkZSJ9.EYx6AXb8dOLIYfDLhhrVEb0VJxeDM_RVTZp82tZ934w"

var testRoutes = server()

func TestJohnCanViewProjectTest(t *testing.T) {
	req, err := http.NewRequest("GET", fmt.Sprintf("/projects/test?tenant=%s", tenant), nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+johnToken)

	rr := httptest.NewRecorder()
	testRoutes.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestJohnCannotViewProjects(t *testing.T) {
	req, err := http.NewRequest("GET", "/projects", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+johnToken)
	rr := httptest.NewRecorder()
	testRoutes.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusForbidden)
	}
}

func TestWadeCanViewAndCreateProjects(t *testing.T) {
	req, err := http.NewRequest("GET", "/projects", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+wadeToken)
	rr := httptest.NewRecorder()
	testRoutes.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	req, err = http.NewRequest("POST", "/projects", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+wadeToken)
	rr = httptest.NewRecorder()
	testRoutes.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	req, err = http.NewRequest("GET", "/projects/x", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+wadeToken)
	rr = httptest.NewRecorder()
	testRoutes.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
