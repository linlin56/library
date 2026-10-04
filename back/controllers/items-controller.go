package controllers

import (
	"library/db"
	"library/repositories"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetItems godoc
// @Summary      Search items
// @Description  Advanced search for items (books/journals), or initial list without filters.
// @Tags         Items
// @Produce      json
// @Param        title         query     string   false  "Title"
// @Param        author        query     string   false  "Author"
// @Param        isbn          query     string   false  "ISBN"
// @Param        barcode       query     string   false  "Barcode"
// @Param        issue_number  query     string   false  "Journal issue number"
// @Param        available     query     boolean  false  "Available only"
// @Success      200           {array}   repositories.ItemResult
// @Success      204           {string}  string  "No Content"
// @Failure      400           {object}  ErrorResponse
// @Failure      500           {object}  ErrorResponse
// @Router       /items [get]
func GetItems(c *gin.Context) {
	title := c.Query("title")
	author := c.Query("author")
	isbn := c.Query("isbn")
	barcode := c.Query("barcode")
	issueNumber := c.Query("issue_number")
	availableStr := c.Query("available")
	available, _ := strconv.ParseBool(availableStr)

	repo := repositories.NewItemRepository(db.GlobalDBO)

	atLeastOneField := title != "" || author != "" || isbn != "" || barcode != "" || issueNumber != ""

	// No parameter: return every item (initial load)
	if !atLeastOneField && !available {
		items, err := repo.FindAll()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": tr(c, "items.search_failed")})
			return
		}
		if len(items) == 0 {
			c.Status(http.StatusNoContent)
			return
		}
		c.JSON(http.StatusOK, items)
		return
	}

	// "available" alone is not enough
	if !atLeastOneField {
		c.JSON(http.StatusBadRequest, gin.H{"error": tr(c, "items.search_field_required")})
		return
	}

	filters := repositories.SearchFilters{
		Title:       title,
		Author:      author,
		Isbn:        isbn,
		Barcode:     barcode,
		IssueNumber: issueNumber,
		Available:   available,
	}

	results, err := repo.Search(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": tr(c, "items.search_failed")})
		return
	}
	if len(results) == 0 {
		c.Status(http.StatusNoContent)
		return
	}
	c.JSON(http.StatusOK, results)
}
