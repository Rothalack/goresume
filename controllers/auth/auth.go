package auth

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"goresume/config"
	"goresume/config/entities"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type loginForm struct {
	UserEmail string `json:"user_email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
}

type registerForm struct {
	UserEmail string `json:"user_email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	UserName  string `json:"user_name" binding:"required"`
}

func Login(c *gin.Context) {
	var form loginForm
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user entities.User
	err := config.DB.QueryRow(
		"SELECT id, email, password_hash FROM users WHERE email = ?",
		form.UserEmail,
	).Scan(&user.ID, &user.Email, &user.PasswordHash)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(form.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	sessionToken := generateSessionToken()
	sessionKey := config.RedisKeyPrefix + "session:" + sessionToken

	err = config.RedisClient.Set(sessionKey, user.ID, 24*time.Hour).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Session error"})
		return
	}

	secure := c.Request.TLS != nil
	c.SetCookie("session", sessionToken, 86400, "/", "", secure, true)

	c.JSON(http.StatusOK, gin.H{"message": "Logged in successfully"})
}

func Register(c *gin.Context) {
	var form registerForm
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var exists bool
	err := config.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)", form.UserEmail).Scan(&exists)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Password hashing error"})
		return
	}

	result, err := config.DB.Exec(
		"INSERT INTO users (email, password_hash, user_name) VALUES (?, ?, ?)",
		form.UserEmail, hashedPassword, form.UserName,
	)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	userID, _ := result.LastInsertId()

	sessionToken := generateSessionToken()
	sessionKey := config.RedisKeyPrefix + "session:" + sessionToken

	err = config.RedisClient.Set(sessionKey, userID, 24*time.Hour).Err()
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Session error"})
		return
	}

	secure := c.Request.TLS != nil
	c.SetCookie("session", sessionToken, 86400, "/", "", secure, true)

	c.JSON(http.StatusCreated, gin.H{"message": "Registered successfully"})
}

func Logout(c *gin.Context) {
	sessionToken, _ := c.Cookie("session")
	if sessionToken != "" {
		config.RedisClient.Del(config.RedisKeyPrefix + "session:" + sessionToken)
	}
	c.SetCookie("session", "", -1, "/", "", c.Request.TLS != nil, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func generateSessionToken() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
