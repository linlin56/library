package controllers

import (
	"library/db"
	"library/repositories"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CopySummary struct {
	Id      int    `json:"id"`
	Barcode string `json:"barcode"`
}

type CopyDetail struct {
	Id         int    `json:"id"`
	Barcode    string `json:"barcode"`
	IsBorrowed bool   `json:"is_borrowed"`
	DueDate    string `json:"loan_due_date,omitempty"`
}

// GetAvailableCopies: GET /api/items/:id/copies
// Returns the available (not borrowed) copies of an item.
// @Summary      List available copies
// @Description  Returns the copies of an item that are not on loan.
// @Tags         Copies
// @Produce      json
// @Param        id   path      int  true  "Item ID"
// @Success      200  {array}   CopySummary
// @Success      204  {string}  string  "No Content"
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /items/{id}/copies [get]
func GetAvailableCopies(c *gin.Context) {
	idStr := c.Param("id")
	itemId, err := strconv.Atoi(idStr)
	if err != nil || itemId <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": tr(c, "common.invalid_id")})
		return
	}

	repo := repositories.NewCopyRepository(db.GlobalDBO)
	copies, err := repo.FindAvailableByItemId(itemId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": tr(c, "copies.load_failed")})
		return
	}

	result := make([]CopySummary, 0, len(copies))
	for _, cp := range copies {
		result = append(result, CopySummary{Id: cp.Id, Barcode: cp.Barcode})
	}
	if len(result) == 0 {
		c.Status(http.StatusNoContent)
		return
	}
	c.JSON(http.StatusOK, result)
}

// GetAllCopies: GET /api/items/:id/copies/all
// Returns every copy (available and on loan) of an item.
// @Summary      List all copies
// @Description  Returns every copy of an item.
// @Tags         Copies
// @Produce      json
// @Param        id   path      int  true  "Item ID"
// @Success      200  {array}   CopyDetail
// @Success      204  {string}  string  "No Content"
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /items/{id}/copies/all [get]
func GetAllCopies(c *gin.Context) {
	idStr := c.Param("id")
	itemId, err := strconv.Atoi(idStr)
	if err != nil || itemId <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": tr(c, "common.invalid_id")})
		return
	}

	repo := repositories.NewCopyRepository(db.GlobalDBO)
	copies, err := repo.FindAllByItemId(itemId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": tr(c, "copies.load_failed")})
		return
	}

	result := make([]CopyDetail, 0, len(copies))
	for _, cp := range copies {
		d := CopyDetail{
			Id:         cp.Id,
			Barcode:    cp.Barcode,
			IsBorrowed: cp.IsBorrowed,
		}
		if cp.IsBorrowed && !cp.LoanDueDate.IsZero() {
			d.DueDate = cp.LoanDueDate.Format("2006-01-02")
		}
		result = append(result, d)
	}
	if len(result) == 0 {
		c.Status(http.StatusNoContent)
		return
	}
	c.JSON(http.StatusOK, result)
}

// CreateCopyForItem: POST /api/items/:id/copies
// Creates a new copy of an item.
// @Summary      Create a copy
// @Description  Creates a copy of a given item.
// @Tags         Copies
// @Accept       json
// @Produce      json
// @Param        id       path      int                true  "Item ID"
// @Param        payload  body      CreateCopyRequest  true  "Copy data"
// @Success      201      {object}  IDResponse
// @Failure      400      {object}  ErrorResponse
// @Router       /items/{id}/copies [post]
func CreateCopyForItem(c *gin.Context) {
	idStr := c.Param("id")
	itemId, err := strconv.Atoi(idStr)
	if err != nil || itemId <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": tr(c, "common.invalid_id")})
		return
	}

	var payload CreateCopyRequest
	if err := c.ShouldBindJSON(&payload); err != nil || payload.Barcode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": tr(c, "copies.barcode_required")})
		return
	}
	if payload.LoanPeriodDays <= 0 {
		payload.LoanPeriodDays = 15
	}

	repo := repositories.NewCopyRepository(db.GlobalDBO)
	newID, err := repo.CreateForItem(itemId, payload.Barcode, payload.LoanPeriodDays)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": trErr(c, err, "copies.create_failed")})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": newID})
}
