package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/eatmore01/light/internal/config"
	"github.com/eatmore01/light/internal/shared/client/keycloak"
	"github.com/eatmore01/light/internal/shared/constants"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var (
	CookieName = "auth_token"
)

const jwtExpiryHours = 24

type AuthApi struct {
	AppConfig *config.Config
}

func NewAuthApi(c *config.Config) *AuthApi {
	return &AuthApi{
		AppConfig: c,
	}
}

// Custom claims structure for our JWT
type CustomClaims struct {
	Username     string `json:"username"`
	IDToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
	jwt.RegisteredClaims
}

func (ap *AuthApi) LoginPage(c *gin.Context) {
	// Check if user already has a valid auth token
	token, err := c.Cookie(CookieName)
	if err == nil && token != "" {
		if ap.validateToken(token) {
			c.Redirect(http.StatusFound, constants.Routes["home"])
			return
		}
	}

	c.HTML(http.StatusOK, "login.html", gin.H{
		"loginPath": constants.Routes["login"],
	})
}

func (a *AuthApi) LoginHandler(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	if username == "" || password == "" {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"error": "Username and password are required",
		})
		return
	}

	kkc := keycloak.NewKeycloakCLient(a.AppConfig)

	// Login to Keycloak and get tokens
	ctx := context.Background()
	token, err := kkc.KeycloakClient.Login(
		ctx,
		a.AppConfig.ClientID,
		a.AppConfig.ClientSecret,
		a.AppConfig.KeycloakRealm,
		username,
		password,
	)

	if err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"error": "Authentication failed: " + err.Error(),
		})
		return
	}

	// Create a JWT containing user info and tokens
	jwtToken, err := a.createJWT(username, token.IDToken, token.RefreshToken)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "login.html", gin.H{
			"error": "Failed to create session: " + err.Error(),
		})
		return
	}

	// Set cookie with JWT
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     CookieName,
		Value:    jwtToken,
		Expires:  time.Now().Add(jwtExpiryHours * time.Hour),
		HttpOnly: true,
		Secure:   a.AppConfig.CookieSecure,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})

	c.Redirect(http.StatusFound, "/home")
}

// LogoutHandler clears the auth cookie
func (a *AuthApi) LogoutHandler(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     CookieName,
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour), // expire immediately
		HttpOnly: true,
		Secure:   a.AppConfig.CookieSecure,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})

	c.Redirect(http.StatusFound, constants.Routes["login"])
}

func (a *AuthApi) HomeHandler(c *gin.Context) {
	tokenString, err := c.Cookie(CookieName)
	if err != nil {
		c.Redirect(http.StatusFound, constants.Routes["login"])
		return
	}

	// Parse and validate token
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.AppConfig.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		c.Redirect(http.StatusFound, constants.Routes["login"])
		return
	}

	// Get claims from token
	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		c.Redirect(http.StatusFound, constants.Routes["login"])
		return
	}

	c.HTML(http.StatusOK, "home.html", gin.H{
		"Username": claims.Username,
	})
}

func (a *AuthApi) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie(CookieName)
		if err != nil {
			c.Redirect(http.StatusFound, constants.Routes["login"])
			c.Abort()
			return
		}

		if !a.validateToken(tokenString) {
			c.Redirect(http.StatusFound, constants.Routes["login"])
			c.Abort()
			return
		}

		c.Next()
	}
}

// helpers  Functions
// create token with additional info
func (a *AuthApi) createJWT(username, idToken, refreshToken string) (string, error) {
	claims := CustomClaims{
		Username:     username,
		IDToken:      idToken,
		RefreshToken: refreshToken,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtExpiryHours * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "auth-ui",
			Subject:   username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(a.AppConfig.JWTSecret))
}

// validate token
func (a *AuthApi) validateToken(tokenString string) bool {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.AppConfig.JWTSecret), nil
	})

	if err != nil {
		return false
	}

	_, ok := token.Claims.(*CustomClaims)
	return ok && token.Valid
}

func AddAuthHandlers(r *gin.Engine, aa *AuthApi) {
	authGroup := r.Group("/auth")

	authGroup.GET("/login", aa.LoginPage)
	authGroup.POST("/login", aa.LoginHandler)

	authGroup.POST("/logout", aa.LogoutHandler)
}
