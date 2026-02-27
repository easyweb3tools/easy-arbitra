package auth

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"easy-arbitra/backend/config"
	"easy-arbitra/backend/internal/model"
	"easy-arbitra/backend/pkg/response"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm"
)

type Handler struct {
	db          *gorm.DB
	cfg         config.AuthConfig
	oauthConfig *oauth2.Config
}

func NewHandler(db *gorm.DB, cfg config.AuthConfig) *Handler {
	var oauthCfg *oauth2.Config
	if cfg.GoogleClientID != "" {
		oauthCfg = &oauth2.Config{
			ClientID:     cfg.GoogleClientID,
			ClientSecret: cfg.GoogleSecret,
			RedirectURL:  cfg.GoogleRedirectURL,
			Scopes:       []string{"openid", "email", "profile"},
			Endpoint:     google.Endpoint,
		}
	}
	return &Handler{db: db, cfg: cfg, oauthConfig: oauthCfg}
}

func (h *Handler) Register(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid body")
		return
	}
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.Name = strings.TrimSpace(req.Name)
	if req.Email == "" || !strings.Contains(req.Email, "@") {
		response.BadRequest(c, "valid email is required")
		return
	}
	if len(req.Password) < 8 {
		response.BadRequest(c, "password must be at least 8 characters")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		response.Internal(c, "failed to hash password")
		return
	}

	user := model.User{
		Email:        req.Email,
		PasswordHash: string(hash),
		Name:         req.Name,
		Provider:     "email",
	}
	if err := h.db.Create(&user).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			response.Conflict(c, "email already registered")
			return
		}
		response.Internal(c, "failed to create user")
		return
	}

	token, err := GenerateToken(user.ID, user.Email, h.cfg.JWTSecret, h.cfg.JWTExpiry)
	if err != nil {
		response.Internal(c, "failed to generate token")
		return
	}
	setAuthCookie(c, token, h.cfg.JWTExpiry)
	response.Created(c, user)
}

func (h *Handler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid body")
		return
	}
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	var user model.User
	if err := h.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		response.Unauthorized(c, "invalid email or password")
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		response.Unauthorized(c, "invalid email or password")
		return
	}

	token, err := GenerateToken(user.ID, user.Email, h.cfg.JWTSecret, h.cfg.JWTExpiry)
	if err != nil {
		response.Internal(c, "failed to generate token")
		return
	}
	setAuthCookie(c, token, h.cfg.JWTExpiry)
	response.OK(c, user)
}

func (h *Handler) GoogleLogin(c *gin.Context) {
	if h.oauthConfig == nil {
		response.BadRequest(c, "Google OAuth not configured")
		return
	}
	state := randomState()
	c.SetCookie("oauth_state", state, 600, "/", "", false, true)
	url := h.oauthConfig.AuthCodeURL(state)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *Handler) GoogleCallback(c *gin.Context) {
	if h.oauthConfig == nil {
		response.BadRequest(c, "Google OAuth not configured")
		return
	}
	storedState, err := c.Cookie("oauth_state")
	if err != nil || storedState == "" || storedState != c.Query("state") {
		response.BadRequest(c, "invalid oauth state")
		return
	}
	// Clear the state cookie
	c.SetCookie("oauth_state", "", -1, "/", "", false, true)

	code := c.Query("code")
	if code == "" {
		response.BadRequest(c, "missing code")
		return
	}

	token, err := h.oauthConfig.Exchange(c.Request.Context(), code)
	if err != nil {
		response.Internal(c, "failed to exchange code")
		return
	}

	client := h.oauthConfig.Client(c.Request.Context(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		response.Internal(c, "failed to get user info")
		return
	}
	defer resp.Body.Close()

	var info struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		response.Internal(c, "failed to decode user info")
		return
	}

	// Find or create user
	var user model.User
	result := h.db.Where("email = ?", info.Email).First(&user)
	if result.Error != nil {
		// Create new user
		user = model.User{
			Email:      info.Email,
			Name:       info.Name,
			AvatarURL:  info.Picture,
			Provider:   "google",
			ProviderID: info.ID,
		}
		if err := h.db.Create(&user).Error; err != nil {
			response.Internal(c, "failed to create user")
			return
		}
	} else {
		// Update avatar and provider info if needed
		h.db.Model(&user).Updates(map[string]any{
			"avatar_url":  info.Picture,
			"provider":    "google",
			"provider_id": info.ID,
		})
	}

	jwtToken, err := GenerateToken(user.ID, user.Email, h.cfg.JWTSecret, h.cfg.JWTExpiry)
	if err != nil {
		response.Internal(c, "failed to generate token")
		return
	}
	setAuthCookie(c, jwtToken, h.cfg.JWTExpiry)

	redirectURL := h.cfg.FrontendURL + "/auth/callback"
	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

func (h *Handler) Me(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "not authenticated")
		return
	}
	var user model.User
	if err := h.db.First(&user, userID).Error; err != nil {
		response.NotFound(c, "user not found")
		return
	}
	response.OK(c, user)
}

func (h *Handler) Logout(c *gin.Context) {
	c.SetCookie("auth_token", "", -1, "/", "", false, true)
	response.OK(c, gin.H{"logged_out": true})
}

func setAuthCookie(c *gin.Context, token string, expiry time.Duration) {
	c.SetCookie("auth_token", token, int(expiry.Seconds()), "/", "", false, true)
}

func randomState() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func (h *Handler) RegisterRoutes(group *gin.RouterGroup) {
	group.POST("/register", h.Register)
	group.POST("/login", h.Login)
	group.GET("/google", h.GoogleLogin)
	group.GET("/google/callback", h.GoogleCallback)
	group.GET("/me", AuthRequired(h.cfg.JWTSecret), h.Me)
	group.POST("/logout", h.Logout)
}
