package handlers

import (
	"auth-service/models"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var secretKey = []byte("spider-Man")

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func RefreshTokensHandler(w http.ResponseWriter, r *http.Request) {

	var request struct {
		RefreshToken string `json:"refresh_token"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil || request.RefreshToken == "" {
		http.Error(w, "Refresh token is required", http.StatusBadRequest)
		return
	}

	ip := r.RemoteAddr

	refreshToken, err := models.GetRefreshTokenByToken(request.RefreshToken)
	if err != nil {
		http.Error(w, "Invalid or expired refresh token", http.StatusUnauthorized)
		return
	}

	if refreshToken.IPAdress != ip {
		http.Error(w, "IP address mismatch", http.StatusUnauthorized)
		return
	}

	if refreshToken.ExpiresAt.Before(time.Now()) {
		http.Error(w, "Refresh token expired", http.StatusUnauthorized)
		return
	}

	newAccessToken := createAccessToken(refreshToken.UserID, ip)
	newRefreshToken := generateRandomToken()
	hashedNewRefreshToken, _ := bcrypt.GenerateFromPassword([]byte(newRefreshToken), bcrypt.DefaultCost)

	refreshToken.RefreshTokenHash = string(hashedNewRefreshToken)
	refreshToken.ExpiresAt = time.Now().Add(24 * time.Hour)
	err = refreshToken.Save()
	if err != nil {
		http.Error(w, "Failed to update refresh token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(TokenPair{AccessToken: newAccessToken, RefreshToken: newRefreshToken})
}

func IssueTokenHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "Id required", http.StatusBadGateway)
		return
	}

	ip := r.RemoteAddr

	access_token := createAccessToken(userID, ip)

	refreshed_tpken := generateRandomToken()

	hashed_token, _ := bcrypt.GenerateFromPassword([]byte(refreshed_tpken), bcrypt.DefaultCost)
	refresh := models.RefreshToken{
		UserID:           userID,
		RefreshTokenHash: string(hashed_token),
		IPAdress:         ip,
		ExpiresAt:        time.Now().Add(24 * time.Hour),
	}

	err := refresh.Save()
	if err != nil {
		http.Error(w, "Failed to save refresh token", http.StatusInternalServerError)
		log.Fatalf("mistake: %v", err)
		return
	}

	json.NewEncoder(w).Encode(TokenPair{AccessToken: access_token, RefreshToken: refreshed_tpken})
}

func createAccessToken(userID, ip string) string {
	expiration := time.Now().Add(15 * time.Minute)

	claims := jwt.MapClaims{
		"user_id": userID,
		"ip":      ip,
		"exp":     expiration.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, _ := token.SignedString(secretKey)
	return tokenString
}

func generateRandomToken() string {
	bits := make([]byte, 32)
	_, err := rand.Read(bits)
	if err != nil {
		// В случае ошибки возвращаем пустую строку
		fmt.Println("Error generating random bytes:", err)
		return ""
	}
	return base64.URLEncoding.EncodeToString(bits)
}
