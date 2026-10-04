package controllers

import (
	"library/db"
	"library/repositories"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// ListLoans: GET /api/loans?user_id=X
// Returns a user's active loans.
// @Summary      List a user's loans
// @Description  Returns the active loans of a given user.
// @Tags         Loans
// @Produce      json
// @Param        user_id  query     int  true  "User ID"
// @Success      200      {array}   repositories.LoanSummary
// @Success      204      {string}  string  "No Content"
// @Failure      400      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Router       /loans [get]
func ListLoans(c *gin.Context) {
	idStr := c.Query("user_id")
	userId, err := strconv.Atoi(idStr)
	if err != nil || userId <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": tr(c, "loans.invalid_user_id")})
		return
	}

	repo := repositories.NewLoanRepository(db.GlobalDBO)
	loans, err := repo.GetLoans(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": tr(c, "loans.load_failed")})
		return
	}
	if len(loans) == 0 {
		c.Status(http.StatusNoContent)
		return
	}
	c.JSON(http.StatusOK, loans)
}

// ListOverdueLoans: GET /api/loans/overdue
// Librarians only. Returns every loan whose due date has passed.
// @Summary      List overdue loans
// @Description  Returns every overdue loan.
// @Tags         Loans
// @Produce      json
// @Success      200  {array}   repositories.OverdueLoan
// @Success      204  {string}  string  "No Content"
// @Failure      500  {object}  ErrorResponse
// @Router       /loans/overdue [get]
func ListOverdueLoans(c *gin.Context) {
	repo := repositories.NewLoanRepository(db.GlobalDBO)
	loans, err := repo.GetOverdueLoans()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": tr(c, "loans.overdue_load_failed")})
		return
	}
	if len(loans) == 0 {
		c.Status(http.StatusNoContent)
		return
	}
	c.JSON(http.StatusOK, loans)
}

// SearchUsers: GET /api/users/search?last_name=&first_name=&postal_code=&phone_number=
// Librarians only. At least one non-empty field is required.
// @Summary      Search users
// @Description  Searches users by last name, first name, postal code or phone number.
// @Tags         Users
// @Produce      json
// @Param        last_name     query     string  false  "Last name"
// @Param        first_name    query     string  false  "First name"
// @Param        postal_code   query     string  false  "Postal code"
// @Param        phone_number  query     string  false  "Phone number"
// @Success      200           {array}   repositories.UserSummary
// @Success      204           {string}  string  "No Content"
// @Failure      400           {object}  ErrorResponse
// @Failure      500           {object}  ErrorResponse
// @Router       /users/search [get]
func SearchUsers(c *gin.Context) {
	lastName := strings.TrimSpace(c.Query("last_name"))
	firstName := strings.TrimSpace(c.Query("first_name"))
	postalCode := strings.TrimSpace(c.Query("postal_code"))
	phoneNumber := strings.TrimSpace(c.Query("phone_number"))

	if lastName == "" && firstName == "" && postalCode == "" && phoneNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": tr(c, "users.search_criteria_required")})
		return
	}

	repo := repositories.NewUserRepository(db.GlobalDBO)
	users, err := repo.SearchUsers(lastName, firstName, postalCode, phoneNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": tr(c, "items.search_failed")})
		return
	}
	if len(users) == 0 {
		c.Status(http.StatusNoContent)
		return
	}
	c.JSON(http.StatusOK, users)
}
