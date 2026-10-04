package controllers

import (
	"library/db"
	"library/models"
	"library/repositories"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// GetAuthors: GET /api/authors
// @Summary      List authors
// @Description  Returns the list of available authors.
// @Tags         Items
// @Produce      json
// @Success      200  {array}   models.Author
// @Failure      500  {object}  ErrorResponse
// @Router       /authors [get]
func GetAuthors(c *gin.Context) {
	repo := repositories.NewAuthorRepository(db.GlobalDBO)
	authors, err := repo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": tr(c, "authors.load_failed")})
		return
	}
	c.JSON(http.StatusOK, authors)
}

// GetLocations: GET /api/locations
// @Summary      List locations
// @Description  Returns the list of available locations.
// @Tags         Items
// @Produce      json
// @Success      200  {array}   repositories.LocationSummary
// @Failure      500  {object}  ErrorResponse
// @Router       /locations [get]
func GetLocations(c *gin.Context) {
	repo := repositories.NewLocationRepository(db.GlobalDBO)
	locations, err := repo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": tr(c, "locations.load_failed")})
		return
	}
	c.JSON(http.StatusOK, locations)
}

type CreateBookRequest struct {
	Title           string  `json:"title"`
	Deposit         float64 `json:"deposit"`
	Isbn            string  `json:"isbn"`
	AuthorId        int     `json:"author_id"`
	AuthorLastName  string  `json:"author_last_name"`
	AuthorFirstName string  `json:"author_first_name"`
	LocationId      int     `json:"location_id"`
}

// CreateBook: POST /api/books
// @Summary      Create a book
// @Description  Creates a new book with an existing author or a new one.
// @Tags         Items
// @Accept       json
// @Produce      json
// @Param        payload  body      CreateBookRequest  true  "Book data"
// @Success      201      {object}  CreateItemResponse
// @Failure      400      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Router       /books [post]
func CreateBook(c *gin.Context) {
	var req CreateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": tr(c, "common.invalid_body")})
		return
	}

	req.Title = strings.TrimSpace(req.Title)
	req.Isbn = strings.TrimSpace(req.Isbn)
	req.AuthorLastName = strings.TrimSpace(req.AuthorLastName)
	req.AuthorFirstName = strings.TrimSpace(req.AuthorFirstName)

	if req.Title == "" || req.Isbn == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": tr(c, "books.title_isbn_required")})
		return
	}
	if req.LocationId <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": tr(c, "items.location_required")})
		return
	}
	if req.Deposit < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": tr(c, "items.deposit_negative")})
		return
	}

	authorId := req.AuthorId
	if authorId == 0 {
		if req.AuthorLastName == "" || req.AuthorFirstName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": tr(c, "books.author_required")})
			return
		}
		authorRepo := repositories.NewAuthorRepository(db.GlobalDBO)
		id, err := authorRepo.Create(&models.Author{LastName: req.AuthorLastName, FirstName: req.AuthorFirstName})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": tr(c, "books.author_create_failed")})
			return
		}
		authorId = id
	}

	itemRepo := repositories.NewItemRepository(db.GlobalDBO)
	id, err := itemRepo.CreateBook(req.Title, req.Deposit, req.Isbn, authorId, req.LocationId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": tr(c, "books.create_failed")})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id, "message": tr(c, "books.created")})
}

type CreateJournalRequest struct {
	Title           string  `json:"title"`
	Deposit         float64 `json:"deposit"`
	IssueNumber     int     `json:"issue_number"`
	PublicationDate string  `json:"publication_date"`
	LocationId      int     `json:"location_id"`
}

// CreateJournal: POST /api/journals
// @Summary      Create a journal
// @Description  Creates a new journal issue.
// @Tags         Items
// @Accept       json
// @Produce      json
// @Param        payload  body      CreateJournalRequest  true  "Journal data"
// @Success      201      {object}  CreateItemResponse
// @Failure      400      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Router       /journals [post]
func CreateJournal(c *gin.Context) {
	var req CreateJournalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": tr(c, "common.invalid_body")})
		return
	}

	req.Title = strings.TrimSpace(req.Title)
	req.PublicationDate = strings.TrimSpace(req.PublicationDate)

	if req.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": tr(c, "items.title_required")})
		return
	}
	if req.IssueNumber <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": tr(c, "journals.number_positive")})
		return
	}
	if req.Deposit < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": tr(c, "items.deposit_negative")})
		return
	}
	if req.PublicationDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": tr(c, "journals.date_required")})
		return
	}
	if req.LocationId <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": tr(c, "items.location_required")})
		return
	}

	itemRepo := repositories.NewItemRepository(db.GlobalDBO)
	id, err := itemRepo.CreateJournal(req.Title, req.Deposit, req.IssueNumber, req.PublicationDate, req.LocationId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": tr(c, "journals.create_failed")})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id, "message": tr(c, "journals.created")})
}
