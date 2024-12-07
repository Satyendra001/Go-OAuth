package auth

import (
	"log"
	"os"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
)

func NewAuth() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("error Loading the envs ==> ", err.Error())
	}

	githubKey := os.Getenv("GITHUB_KEY")
	githubSecret := os.Getenv("GITHUB_SECRET")

	googleKey := os.Getenv("GOOGLE_KEY")
	googleSecret := os.Getenv("GOOGLE_SECRET")

	store := sessions.NewCookieStore([]byte("randomData"))
	store.Options.HttpOnly = true

	gothic.Store = store

	goth.UseProviders(
		github.New(githubKey, githubSecret, "http://dmt.localhost:3000/auth/github/callback"),
		google.New(googleKey, googleSecret, "http://dmt.localhost:3000/auth/google/callback"),
	)
}
