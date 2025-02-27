package auth

import (
	"os"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

func NewAuth() {
	store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))
	store.MaxAge(86400 * 30) // 30 days
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = os.Getenv("ENVIRONMENT") == "production"

	gothic.Store = store

	goth.UseProviders(
		google.New(os.Getenv("GOOGLE_CLIENT_ID"), os.Getenv("GOOGLE_CLIENT_SECRET"), "http://localhost:8081/auth/google/callback", "email", "profile"),
	)
}
