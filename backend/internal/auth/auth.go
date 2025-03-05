package auth

import (
	"fmt"
	"os"

	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
)

func NewAuth() {
	googleCallback := fmt.Sprintf("%s/auth/google/callback", os.Getenv("BACKEND_URL"))

	goth.UseProviders(
		google.New(os.Getenv("GOOGLE_CLIENT_ID"), os.Getenv("GOOGLE_CLIENT_SECRET"), googleCallback, "email", "profile"),
	)
}
