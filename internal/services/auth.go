package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"btc-trading-bot/internal/database"
	"btc-trading-bot/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	db        *database.Database
	jwtSecret string
}

type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func NewAuthService(db *database.Database, jwtSecret string) *AuthService {
	return &AuthService{
		db:        db,
		jwtSecret: jwtSecret,
	}
}

func (s *AuthService) Register(userData *models.UserRegister) (*models.User, error) {
	var existingUser models.User
	err := s.db.Get(&existingUser, "SELECT id FROM users WHERE username = $1 OR email = $2", userData.Username, userData.Email)
	if err == nil {
		return nil, fmt.Errorf("user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}

	user := &models.User{
		Username:  userData.Username,
		Password:  string(hashedPassword),
		Email:     userData.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	query := `
		INSERT INTO users (username, password, email, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, username, email, created_at, updated_at
	`

	err = s.db.QueryRowx(query, user.Username, user.Password, user.Email, user.CreatedAt, user.UpdatedAt).StructScan(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	return user, nil
}

func (s *AuthService) Login(loginData *models.UserLogin) (string, error) {
	var user models.User
	err := s.db.Get(&user, "SELECT id, username, password FROM users WHERE username = $1", loginData.Username)
	if err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password))
	if err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	token, err := s.generateJWT(user.ID, user.Username)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}

	return token, nil
}

func (s *AuthService) generateJWT(userID int, username string) (string, error) {
	claims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func (s *AuthService) GetUserByID(userID int) (*models.User, error) {
	var user models.User
	err := s.db.Get(&user, "SELECT id, username, email, created_at, updated_at FROM users WHERE id = $1", userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}

	return &user, nil
}

func (s *AuthService) GenerateAPIKey() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
