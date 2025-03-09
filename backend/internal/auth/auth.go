package auth

import (
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth/gothic"
	"os"

	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
)

// NewAuth initializes the authentication provider using Google OAuth 2.0.
//
// It sets up a Google OAuth 2.0 provider with client ID, client secret, and callback URL obtained from environment
// variables. The callback URL is constructed by combining the BACKEND_URL environment variable with "/auth/google/callback".
// The provider requests "email" and "profile" scopes.
//
// Parameters:
//   - none
//
// Returns:
//   - none
//
// Error types:
//   - error: if environment variables GOOGLE_CLIENT_ID, GOOGLE_CLIENT_SECRET, or BACKEND_URL are not set, or if there's
//     an issue creating the Google provider. The specific error will depend on the underlying cause.
func NewAuth() {
	store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))
	store.MaxAge(86400 * 30) // 30 days
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = os.Getenv("ENVIRONMENT") == "production"

	gothic.Store = store

	googleCallback := fmt.Sprintf("%s/auth/google/callback", os.Getenv("BACKEND_URL"))

	goth.UseProviders(
		google.New(os.Getenv("GOOGLE_CLIENT_ID"), os.Getenv("GOOGLE_CLIENT_SECRET"), googleCallback, "email", "profile"),
	)
}
