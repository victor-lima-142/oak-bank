package security

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var (
	ErrInvalidToken      = errors.New("token inválido")
	ErrExpiredToken      = errors.New("token expirado")
	ErrMissingToken      = errors.New("token não encontrado no header")
	ErrInvalidAuthHeader = errors.New("formato de autorização inválido")
)

type JwtService interface {
	// CreateToken cria um novo token JWT com o payload fornecido
	CreateToken(payload map[string]interface{}, expiration time.Duration) (string, error)

	// CreateRefreshToken cria um token de refresh com duração maior
	CreateRefreshToken(payload map[string]interface{}) (string, error)

	// ValidateToken valida um token e retorna os claims se válido
	ValidateToken(tokenString string) (jwt.MapClaims, error)

	// ValidateTokenFromHeader extrai e valida o token do header Authorization da request
	ValidateTokenFromHeader(c *gin.Context) (jwt.MapClaims, error)

	// RefreshToken revalida um refresh token e gera um novo access token
	RefreshToken(refreshToken string) (string, error)

	// ExtractTokenFromHeader extrai o token do header Authorization
	ExtractTokenFromHeader(c *gin.Context) (string, error)
}

type jwtService struct {
	secretKey         []byte
	refreshSecretKey  []byte
	defaultExpiration time.Duration
	refreshExpiration time.Duration
	issuer            string
}

// JwtConfig contém as configurações para criar um novo JwtService
type JwtConfig struct {
	SecretKey         string
	RefreshSecretKey  string
	DefaultExpiration time.Duration
	RefreshExpiration time.Duration
	Issuer            string
}

func NewJwtService(config *JwtConfig) JwtService {
	_ = godotenv.Load()

	if config == nil {
		config = &JwtConfig{}
	}

	if config.SecretKey == "" {
		config.SecretKey = getEnv("JWT_SECRET_KEY", "")
	}
	if config.RefreshSecretKey == "" {
		config.RefreshSecretKey = getEnv("JWT_REFRESH_SECRET_KEY", "")
	}
	if config.DefaultExpiration == 0 {
		config.DefaultExpiration = getEnvDuration("JWT_DEFAULT_EXPIRATION", 15*time.Minute)
	}
	if config.RefreshExpiration == 0 {
		config.RefreshExpiration = getEnvDuration("JWT_REFRESH_EXPIRATION", 7*24*time.Hour)
	}
	if config.Issuer == "" {
		config.Issuer = getEnv("JWT_ISSUER", "my-app")
	}

	if config.SecretKey == "" {
		panic("JWT_SECRET_KEY é obrigatório. Configure no .env ou passe na JwtConfig")
	}

	if config.RefreshSecretKey == "" {
		config.RefreshSecretKey = config.SecretKey + "-refresh"
	}

	return &jwtService{
		secretKey:         []byte(config.SecretKey),
		refreshSecretKey:  []byte(config.RefreshSecretKey),
		defaultExpiration: config.DefaultExpiration,
		refreshExpiration: config.RefreshExpiration,
		issuer:            config.Issuer,
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	if duration, err := time.ParseDuration(value); err == nil {
		return duration
	}

	if strings.HasSuffix(value, "d") {
		days := strings.TrimSuffix(value, "d")
		if d, err := strconv.Atoi(days); err == nil {
			return time.Duration(d) * 24 * time.Hour
		}
	}

	if seconds, err := strconv.Atoi(value); err == nil {
		return time.Duration(seconds) * time.Second
	}

	return defaultValue
}

func (j *jwtService) CreateToken(payload map[string]interface{}, expiration time.Duration) (string, error) {
	if expiration == 0 {
		expiration = j.defaultExpiration
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"iss": j.issuer,
		"iat": now.Unix(),
		"exp": now.Add(expiration).Unix(),
	}

	for key, value := range payload {
		claims[key] = value
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}

func (j *jwtService) CreateRefreshToken(payload map[string]interface{}) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"iss":  j.issuer,
		"iat":  now.Unix(),
		"exp":  now.Add(j.refreshExpiration).Unix(),
		"type": "refresh",
	}

	for key, value := range payload {
		claims[key] = value
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.refreshSecretKey)
}

func (j *jwtService) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de assinatura inesperado: %v", token.Header["alg"])
		}
		return j.secretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func (j *jwtService) ValidateTokenFromHeader(c *gin.Context) (jwt.MapClaims, error) {
	tokenString, err := j.ExtractTokenFromHeader(c)
	if err != nil {
		return nil, err
	}

	return j.ValidateToken(tokenString)
}

func (j *jwtService) ExtractTokenFromHeader(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", ErrMissingToken
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", ErrInvalidAuthHeader
	}

	return parts[1], nil
}

func (j *jwtService) RefreshToken(refreshToken string) (string, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de assinatura inesperado: %v", token.Header["alg"])
		}
		return j.refreshSecretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return "", ErrExpiredToken
		}
		return "", fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	if !token.Valid {
		return "", ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", ErrInvalidToken
	}

	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "refresh" {
		return "", errors.New("não é um refresh token válido")
	}

	payload := make(map[string]interface{})
	excludeKeys := map[string]bool{
		"iss":  true,
		"iat":  true,
		"exp":  true,
		"type": true,
	}

	for key, value := range claims {
		if !excludeKeys[key] {
			payload[key] = value
		}
	}

	return j.CreateToken(payload, j.defaultExpiration)
}
