package controllers

import (
	"library/db"
	"library/repositories"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetDeposit: GET /api/users/:id/deposit
// Returns a user's deposit_balance and total_deposit.
// @Summary      Get a user's deposit
// @Description  Returns the available deposit balance and the total deposit.
// @Tags         Users
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  repositories.DepositInfo
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Router       /users/{id}/deposit [get]
func GetDeposit(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": tr(c, "common.invalid_id")})
		return
	}

	repo := repositories.NewUserRepository(db.GlobalDBO)
	info, err := repo.GetDeposit(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": trErr(c, err, "users.not_found")})
		return
	}

	c.JSON(http.StatusOK, info)
}

// UpdateTotalDeposit: PUT /api/users/:id/deposit
// Updates a user's total_deposit (librarians only).
// Recomputes deposit_balance while keeping the borrowed amount.
// Returns 400 if the new value is lower than the borrowed amount.
// @Summary      Update the total deposit
// @Description  Updates a user's total deposit and recomputes their balance.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id       path      int                   true  "User ID"
// @Param        payload  body      UpdateDepositRequest  true  "New deposit"
// @Success      200      {object}  repositories.DepositInfo
// @Failure      400      {object}  ErrorResponse
// @Router       /users/{id}/deposit [put]
func UpdateTotalDeposit(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": tr(c, "common.invalid_id")})
		return
	}

	var payload UpdateDepositRequest
	if err := c.ShouldBindJSON(&payload); err != nil || payload.TotalDeposit < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": tr(c, "deposit.invalid_value")})
		return
	}

	repo := repositories.NewUserRepository(db.GlobalDBO)
	if err := repo.UpdateTotalDeposit(id, payload.TotalDeposit); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": trErr(c, err, "deposit.update_failed")})
		return
	}

	// Return the new deposit information
	info, err := repo.GetDeposit(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": tr(c, "deposit.updated")})
		return
	}
	c.JSON(http.StatusOK, info)
}
