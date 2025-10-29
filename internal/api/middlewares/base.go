package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/victor-lima-142/oak-bank/internal/api/security"
)

// UserClaims representa as informações do usuário extraídas do token
type UserClaims struct {
	UserID   interface{} `json:"user_id"`
	Username string      `json:"username"`
	Email    string      `json:"email"`
	Role     string      `json:"role"`
	Claims   jwt.MapClaims
}

// Context keys para acessar dados do usuário
const (
	UserClaimsKey = "user_claims"
	UserIDKey     = "user_id"
	UsernameKey   = "username"
	EmailKey      = "email"
	RoleKey       = "role"
	UserKey       = "user"
)

// AuthMiddleware é o middleware de autenticação básico
// Valida o token e adiciona os claims ao contexto
func AuthMiddleware(jwtService security.JwtService) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := jwtService.ValidateTokenFromHeader(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token inválido ou ausente",
			})
			c.Abort()
			return
		}

		// Extrai informações do usuário dos claims
		user := extractUserFromClaims(claims)

		// Adiciona ao contexto
		c.Set(UserClaimsKey, claims)
		c.Set(UserKey, user)
		c.Set(UserIDKey, user.UserID)
		c.Set(UsernameKey, user.Username)
		c.Set(EmailKey, user.Email)
		c.Set(RoleKey, user.Role)

		c.Next()
	}
}

// AuthMiddlewareOptional é um middleware que tenta autenticar, mas não bloqueia se falhar
// Útil para rotas que funcionam tanto autenticadas quanto não autenticadas
func AuthMiddlewareOptional(jwtService security.JwtService) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := jwtService.ValidateTokenFromHeader(c)
		if err == nil {
			// Token válido, adiciona ao contexto
			user := extractUserFromClaims(claims)
			c.Set(UserClaimsKey, claims)
			c.Set(UserKey, user)
			c.Set(UserIDKey, user.UserID)
			c.Set(UsernameKey, user.Username)
			c.Set(EmailKey, user.Email)
			c.Set(RoleKey, user.Role)
		}
		// Continua mesmo se o token for inválido ou não existir
		c.Next()
	}
}

// RequireRole é um middleware que verifica se o usuário tem uma role específica
// Deve ser usado após AuthMiddleware
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get(RoleKey)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Usuário não autenticado",
			})
			c.Abort()
			return
		}

		role, ok := userRole.(string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Role inválida",
			})
			c.Abort()
			return
		}

		// Verifica se a role do usuário está na lista de roles permitidas
		for _, allowedRole := range roles {
			if strings.EqualFold(role, allowedRole) {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"error": "Acesso negado: permissões insuficientes",
		})
		c.Abort()
	}
}

// RequireRoleStrict é como RequireRole mas case-sensitive
func RequireRoleStrict(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get(RoleKey)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Usuário não autenticado",
			})
			c.Abort()
			return
		}

		role, ok := userRole.(string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Role inválida",
			})
			c.Abort()
			return
		}

		for _, allowedRole := range roles {
			if role == allowedRole {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"error": "Acesso negado: permissões insuficientes",
		})
		c.Abort()
	}
}

// extractUserFromClaims extrai informações do usuário dos claims JWT
func extractUserFromClaims(claims jwt.MapClaims) *UserClaims {
	user := &UserClaims{
		Claims: claims,
	}

	// Extrai user_id (pode ser int, float64 ou string)
	if userID, ok := claims["user_id"]; ok {
		user.UserID = userID
	}

	// Extrai username
	if username, ok := claims["username"].(string); ok {
		user.Username = username
	}

	// Extrai email
	if email, ok := claims["email"].(string); ok {
		user.Email = email
	}

	// Extrai role
	if role, ok := claims["role"].(string); ok {
		user.Role = role
	}

	return user
}

// GetUser retorna o usuário do contexto
func GetUser(c *gin.Context) (*UserClaims, bool) {
	if user, exists := c.Get(UserKey); exists {
		if u, ok := user.(*UserClaims); ok {
			return u, true
		}
	}
	return nil, false
}

// GetUserID retorna o ID do usuário do contexto
func GetUserID(c *gin.Context) (interface{}, bool) {
	return c.Get(UserIDKey)
}

// GetUserIDAsInt retorna o ID do usuário como int
func GetUserIDAsInt(c *gin.Context) (int, bool) {
	userID, exists := c.Get(UserIDKey)
	if !exists {
		return 0, false
	}

	// Tenta converter de diferentes tipos
	switch v := userID.(type) {
	case int:
		return v, true
	case int64:
		return int(v), true
	case float64:
		return int(v), true
	case string:
		// Tenta converter string para int
		var id int
		if _, err := fmt.Sscanf(v, "%d", &id); err == nil {
			return id, true
		}
	}

	return 0, false
}

// GetUserIDAsString retorna o ID do usuário como string
func GetUserIDAsString(c *gin.Context) (string, bool) {
	userID, exists := c.Get(UserIDKey)
	if !exists {
		return "", false
	}

	return fmt.Sprintf("%v", userID), true
}

// GetUsername retorna o username do contexto
func GetUsername(c *gin.Context) (string, bool) {
	if username, exists := c.Get(UsernameKey); exists {
		if u, ok := username.(string); ok {
			return u, true
		}
	}
	return "", false
}

// GetEmail retorna o email do contexto
func GetEmail(c *gin.Context) (string, bool) {
	if email, exists := c.Get(EmailKey); exists {
		if e, ok := email.(string); ok {
			return e, true
		}
	}
	return "", false
}

// GetRole retorna a role do contexto
func GetRole(c *gin.Context) (string, bool) {
	if role, exists := c.Get(RoleKey); exists {
		if r, ok := role.(string); ok {
			return r, true
		}
	}
	return "", false
}

// GetClaims retorna todos os claims do contexto
func GetClaims(c *gin.Context) (jwt.MapClaims, bool) {
	if claims, exists := c.Get(UserClaimsKey); exists {
		if cl, ok := claims.(jwt.MapClaims); ok {
			return cl, true
		}
	}
	return nil, false
}

// MustGetUser retorna o usuário ou entra em panic (use apenas quando tem certeza que existe)
func MustGetUser(c *gin.Context) *UserClaims {
	user, exists := GetUser(c)
	if !exists {
		panic("user not found in context")
	}
	return user
}

// MustGetUserID retorna o user ID ou entra em panic
func MustGetUserID(c *gin.Context) interface{} {
	userID, exists := GetUserID(c)
	if !exists {
		panic("user_id not found in context")
	}
	return userID
}

// IsAuthenticated verifica se há um usuário autenticado no contexto
func IsAuthenticated(c *gin.Context) bool {
	_, exists := c.Get(UserKey)
	return exists
}

// HasRole verifica se o usuário tem uma role específica
func HasRole(c *gin.Context, role string) bool {
	userRole, exists := GetRole(c)
	if !exists {
		return false
	}
	return strings.EqualFold(userRole, role)
}

// HasAnyRole verifica se o usuário tem qualquer uma das roles especificadas
func HasAnyRole(c *gin.Context, roles ...string) bool {
	userRole, exists := GetRole(c)
	if !exists {
		return false
	}

	for _, role := range roles {
		if strings.EqualFold(userRole, role) {
			return true
		}
	}
	return false
}
