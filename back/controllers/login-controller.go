package controllers

import (
	"library/db"
	"library/repositories"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Id        int    `json:"id"`
	LastName  string `json:"last_name"`
	FirstName string `json:"first_name"`
	Role      string `json:"role"`
	Message   string `json:"message"`
}

// Login godoc
// @Summary      User login
// @Description  Authenticates a user (student/teacher/individual) or a librarian.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        credentials  body      LoginRequest  true  "Login credentials"
// @Success      200          {object}  LoginResponse
// @Failure      400          {object}  ErrorResponse
// @Failure      401          {object}  ErrorResponse
// @Router       /login [post]
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": tr(c, "auth.missing_fields")})
		return
	}

	req.Login = strings.TrimSpace(req.Login)
	req.Password = strings.TrimSpace(req.Password)
	if req.Login == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": tr(c, "auth.credentials_required")})
		return
	}

	userRepo := repositories.NewUserRepository(db.GlobalDBO)
	user, err := userRepo.FindByLogin(req.Login)
	if err == nil {
		if user.Password != req.Password {
			c.JSON(http.StatusUnauthorized, gin.H{"error": tr(c, "auth.invalid_credentials")})
			return
		}
		c.JSON(http.StatusOK, LoginResponse{
			Id:        user.Id,
			LastName:  user.LastName,
			FirstName: user.FirstName,
			Role:      "user",
			Message:   tr(c, "auth.success"),
		})
		return
	}

	librarianRepo := repositories.NewLibrarianRepository(db.GlobalDBO)
	librarian, err2 := librarianRepo.FindByLogin(req.Login)
	if err2 != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": tr(c, "auth.invalid_credentials")})
		return
	}
	if librarian.Password != req.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": tr(c, "auth.invalid_credentials")})
		return
	}
	c.JSON(http.StatusOK, LoginResponse{
		Id:        librarian.Id,
		LastName:  librarian.LastName,
		FirstName: librarian.FirstName,
		Role:      "librarian",
		Message:   tr(c, "auth.success"),
	})
}
