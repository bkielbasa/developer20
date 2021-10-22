package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %s", err)
	}

	conf := &oauth2.Config{
		ClientID:     envOrPanic("CLIENT_ID"),
		ClientSecret: envOrPanic("CLIENT_SECRET"),
		RedirectURL:  "http://localhost:8080/redirect",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://api.raindrop.io/v1/oauth/authorize",
			TokenURL: "https://raindrop.io/oauth/access_token",
		},
	}

	h := oauthHandler{conf: conf}

	http.HandleFunc("/auth", h.auth)
	http.HandleFunc("/redirect", h.handleCode)

	log.Print("visit: http://localhost:8080/auth")
	err = http.ListenAndServe(":8080", nil)
	log.Fatal(err)
}

func envOrPanic(name string) string {
	v := os.Getenv(name)
	if v == "" {
		panic("the env " + name + " shouldn't be empty")
	}

	return v
}

type oauthHandler struct {
	conf *oauth2.Config
}

func (h oauthHandler) auth(w http.ResponseWriter, r *http.Request) {
	url := h.conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	body := fmt.Sprintf(`<p>To authenticate please <a href="%v">click here</a></p>`, url)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, body)
}

func (h oauthHandler) handleCode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	code := r.URL.Query()["code"][0]
	log.Printf("code fetched: %s", code)
	tok, err := h.conf.Exchange(ctx, code)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("exchanged for the token: %s", tok)

	b, err := json.Marshal(tok)

	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("token.json", b, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
