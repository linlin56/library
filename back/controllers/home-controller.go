package controllers

import (
	"library/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

// The home screen displays a few statistics about the library
// It shows that the database has been seeded
type HomeStats struct {
	BookCount          int `json:"book_count"`
	JournalCount       int `json:"journal_count"`
	AvailableCopyCount int `json:"available_copy_count"`
	UserCount          int `json:"user_count"`
}

// Home godoc
// @Summary      Home statistics
// @Description  Returns global statistics about the library.
// @Tags         Home
// @Produce      json
// @Success      200  {object}  HomeStats
// @Router       /home [get]
func Home(c *gin.Context) {
	var stats HomeStats
	_ = db.GlobalDBO.QueryRow(`SELECT COUNT(*) FROM books`).Scan(&stats.BookCount)
	_ = db.GlobalDBO.QueryRow(`SELECT COUNT(*) FROM journals`).Scan(&stats.JournalCount)
	_ = db.GlobalDBO.QueryRow(`SELECT COUNT(*) FROM copies WHERE is_borrowed = FALSE`).Scan(&stats.AvailableCopyCount)
	_ = db.GlobalDBO.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&stats.UserCount)

	c.JSON(http.StatusOK, stats)
}
