package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/sean-miningah/sil-backend-assessment/internal/core/domain"
	"github.com/sean-miningah/sil-backend-assessment/internal/core/ports"
	"github.com/sean-miningah/sil-backend-assessment/pkg/auth"
	"github.com/sean-miningah/sil-backend-assessment/pkg/utils"
	"go.opentelemetry.io/otel"
)

type AuthHandler struct {
	authConfig   *auth.AuthConfig
	customerRepo ports.CustomerService
}

func NewAuthHandler(authConfig *auth.AuthConfig, customerRepo ports.CustomerService) *AuthHandler {
	return &AuthHandler{
		authConfig:   authConfig,
		customerRepo: customerRepo,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	url := h.authConfig.GoogleOAuth.AuthCodeURL("state")
	c.JSON(http.StatusOK, gin.H{"redirect_url": url})
}

func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	ctx, span := otel.Tracer("").Start(c.Request.Context(), "GoogleCallback")
	defer span.End()
	code := c.Query("code")
	token, err := h.authConfig.GoogleOAuth.Exchange(c, code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to exchange token"})
		return
	}

	client := h.authConfig.GoogleOAuth.Client(c, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get user info"})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response body"})
		return
	}

	fmt.Println("Raw response body:", string(body))

	// print out the returned object

	var userInfo utils.GoogleApiResponse
	if err := json.Unmarshal(body, &userInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to decode user info"})
		return
	}

	customer := &domain.Customer{
		Name:          userInfo.Name,
		Email:         userInfo.Email,
		VerifiedEmail: userInfo.VerifiedEmail,
		Picture:       userInfo.Picture,
	}

	// Create or update user in database
	user, err := h.customerRepo.UpsertCustomer(ctx, customer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upsert user"})
		return
	}

	// Generate JWT
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := jwtToken.SignedString(h.authConfig.JWTSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.SetCookie("auth_token", tokenString, 3600*24, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.SetCookie("auth_token", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
