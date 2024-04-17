package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

// This sample shows how standard http can be used with golang-jwt
// to create a middleware validating the incoming requests for the users.
//
// This sample requires that you have these roles in your tenant:
// - projects.viewer: projects.read
// - projects.owner: projects.read | projects.write
//
// It also requires these bindings to be present:
// IDENTITY |      ROLES      | RESOURCE
// ---------------------------------------------
//   john   | projects.viewer | /projects/test
//   wade   | projects.owner  | /projects
//
// To understand how these work, please refer to the documentation at:
// https://docs.accessbox.dev/basics/designing_authorization

var tenant = os.Getenv("ACCESSBOX_TENANT")
var apiKey = os.Getenv("ACCESSBOX_API_KEY")

func main() {
	http.ListenAndServe(":8080", server())
}

func server() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /projects", protect("projects.read", getProjects))
	mux.HandleFunc("GET /projects/{id}", protect("projects.read", getProject))
	mux.HandleFunc("POST /projects", protect("projects.write", createProject))

	return mux
}

func protect(permission string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if len(r.Header.Get("Authorization")) < 7 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tok, err := jwt.Parse(r.Header.Get("Authorization")[7:], func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		sub, err := tok.Claims.GetSubject()
		if err != nil {
			http.Error(w, "Invalid token: subject", http.StatusUnauthorized)
			return
		}

		reqBody := []byte(fmt.Sprintf(`{"identity": "%s","resource": "%s","permission": "%s"}`, sub, r.URL.Path, permission))

		req, err := http.NewRequest("POST", fmt.Sprintf("https://api.accessbox.dev/v1/authorize?tenant=%s", tenant), bytes.NewBuffer(reqBody))
		if err != nil {
			return
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+apiKey)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			slog.Error("Error checking permission", slog.Any("error", err))
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		defer resp.Body.Close()

		var result map[string]bool
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			slog.Error("Error decoding response", slog.Any("error", err))
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if !result["allow"] {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

func getProjects(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func getProject(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func createProject(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
