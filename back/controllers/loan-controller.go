package controllers

import (
	"library/db"
	"library/repositories"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type LoanRequest struct {
	UserId  int    `json:"user_id" binding:"required"`
	Barcode string `json:"barcode" binding:"required"`
}

// CheckLoan: GET /api/loans/check?user_id=X&barcode=Y
// @Summary      Check a loan
// @Description  Checks whether a copy can be borrowed and returns a preview.
// @Tags         Loans
// @Produce      json
// @Param        user_id  query     int     true  "User ID"
// @Param        barcode  query     string  true  "Copy barcode"
// @Success      200      {object}  repositories.LoanPreview
// @Failure      400      {object}  ErrorResponse
// @Failure      422      {object}  ErrorResponse
// @Router       /loans/check [get]
func CheckLoan(c *gin.Context) {
	userIdStr := c.Query("user_id")
	barcode := c.Query("barcode")

	userId, err := strconv.Atoi(userIdStr)
	if err != nil || userId <= 0 || barcode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": tr(c, "common.missing_params")})
		return
	}

	repo := repositories.NewLoanRepository(db.GlobalDBO)
	preview, err := repo.Check(userId, barcode)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": trErr(c, err, "loans.verify_failed")})
		return
	}
	c.JSON(http.StatusOK, preview)
}

// Borrow: POST /api/loans
// @Summary      Record a loan
// @Description  Records the loan of a copy for a user.
// @Tags         Loans
// @Accept       json
// @Produce      json
// @Param        payload  body      LoanRequest  true  "Loan request"
// @Success      201      {object}  MessageResponse
// @Failure      400      {object}  ErrorResponse
// @Failure      422      {object}  ErrorResponse
// @Router       /loans [post]
func Borrow(c *gin.Context) {
	var req LoanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": tr(c, "common.missing_params")})
		return
	}

	repo := repositories.NewLoanRepository(db.GlobalDBO)
	if err := repo.Borrow(req.UserId, req.Barcode); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": trErr(c, err, "loans.create_failed")})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": tr(c, "loans.created")})
}

// ReturnLoan: DELETE /api/loans/:id
// @Summary      Return a copy
// @Description  Marks the copy as returned and refunds the deposit.
// @Tags         Loans
// @Produce      json
// @Param        id   path      int  true  "Copy ID"
// @Success      200  {object}  MessageResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      422  {object}  ErrorResponse
// @Router       /loans/{id} [delete]
func ReturnLoan(c *gin.Context) {
	copyId, err := strconv.Atoi(c.Param("id"))
	if err != nil || copyId <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": tr(c, "common.invalid_id")})
		return
	}

	repo := repositories.NewLoanRepository(db.GlobalDBO)
	if err := repo.ReturnCopy(copyId); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": trErr(c, err, "loans.return_failed")})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": tr(c, "loans.returned")})
}
