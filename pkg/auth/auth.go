package auth

import (
	"github.com/sean-miningah/sil-backend-assessment/pkg/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type AuthConfig struct {
	GoogleOAuth *oauth2.Config
	JWTSecret   []byte
}

func NewAuthConfig(cfg *config.Config) *AuthConfig {
	return &AuthConfig{
		GoogleOAuth: &oauth2.Config{
			ClientID:     cfg.GoogleClientID,
			ClientSecret: cfg.GoogleClientSecret,
			RedirectURL:  cfg.GoogleRedirectURL,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		},
		JWTSecret: []byte(cfg.JWTSecret),
	}
}
