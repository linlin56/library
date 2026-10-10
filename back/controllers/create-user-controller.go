package controllers

import (
	"library/builder"
	"library/db"
	"library/models"
	"library/repositories"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var phoneRegexp = regexp.MustCompile(`^0[1-9][0-9]{8}$`)

type CreateUserRequest struct {
	LastName     string `json:"last_name" binding:"required"`
	FirstName    string `json:"first_name" binding:"required"`
	PhoneNumber  string `json:"phone_number" binding:"required"`
	BirthDate    string `json:"birth_date" binding:"required"`
	Email        string `json:"email" binding:"required"`
	Status       string `json:"status" binding:"required,oneof=student teacher individual"`
	StudyYear    string `json:"study_year"`
	DepartmentId int    `json:"department_id"`
	StreetNumber string `json:"street_number" binding:"required"`
	StreetName   string `json:"street_name" binding:"required"`
	PostalCode   string `json:"postal_code" binding:"required"`
	City         string `json:"city" binding:"required"`
}

type CreateUserResponse struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Message  string `json:"message"`
}

// GetDepartments godoc
// @Summary      List departments
// @Description  Returns the list of school departments.
// @Tags         Users
// @Produce      json
// @Success      200  {array}   models.Department
// @Failure      500  {object}  ErrorResponse
// @Router       /departments [get]
func GetDepartments(c *gin.Context) {
	repo := repositories.NewDepartmentRepository(db.GlobalDBO)
	departments, err := repo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": tr(c, "users.departments_load_failed")})
		return
	}
	c.JSON(http.StatusOK, departments)
}

// CreateUser godoc
// @Summary      Create a user
// @Description  Creates a student, a teacher or an individual.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        payload  body      CreateUserRequest  true  "User data"
// @Success      201      {object}  CreateUserResponse
// @Failure      400      {object}  ErrorResponse
// @Failure      409      {object}  LoginConflictResponse
// @Failure      500      {object}  ErrorResponse
// @Router       /users [post]
func CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": tr(c, "users.missing_fields")})
		return
	}

	req.LastName = strings.TrimSpace(req.LastName)
	req.FirstName = strings.TrimSpace(req.FirstName)
	req.Email = strings.TrimSpace(req.Email)
	req.PhoneNumber = strings.TrimSpace(req.PhoneNumber)
	req.StreetName = capitalize(strings.TrimSpace(req.StreetName))
	req.City = capitalize(strings.TrimSpace(req.City))

	if !phoneRegexp.MatchString(req.PhoneNumber) {
		c.JSON(http.StatusBadRequest, gin.H{"error": tr(c, "users.invalid_phone")})
		return
	}

	birthDate, err := time.Parse("2006-01-02", req.BirthDate)
	// This should not happen since the format is handled by the front end
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": tr(c, "users.invalid_date")})
		return
	}

	if req.Status == "student" && strings.TrimSpace(req.StudyYear) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": tr(c, "users.year_required")})
		return
	}
	if req.Status == "teacher" && req.DepartmentId == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": tr(c, "users.department_required")})
		return
	}

	login := generateLogin(req.LastName, req.FirstName)
	// To keep things simple, the password is always "password".
	// In a real application, a random password would of course be generated and hashed
	const password = "password"

	userRepo := repositories.NewUserRepository(db.GlobalDBO)
	exists, err := userRepo.LoginExists(login)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": tr(c, "users.account_check_failed")})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{
			"message":        tr(c, "users.login_exists", login),
			"existing_login": login,
		})
		return
	}

	var addressId int
	if err := db.GlobalDBO.ExecReturning(`
		INSERT INTO addresses (city, postal_code, street_name, street_number)
		VALUES ($1, $2, $3, $4) RETURNING id`,
		req.City, req.PostalCode, req.StreetName, req.StreetNumber,
	).Scan(&addressId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": tr(c, "users.address_failed")})
		return
	}

	// Users are assembled with the builders, then persisted by the matching repository
	const initialBalance = 20.0
	switch req.Status {
	case "student":
		s := builder.NewStudentBuilder().
			WithLastName(req.LastName).
			WithFirstName(req.FirstName).
			WithPhoneNumber(req.PhoneNumber).
			WithDepositBalance(initialBalance).
			WithEmail(req.Email).
			WithLogin(login).
			WithPassword(password).
			WithBirthDate(birthDate).
			WithAddressId(addressId).
			WithStudyYear(req.StudyYear).
			Build()
		studentRepo := repositories.NewStudentRepository(db.GlobalDBO)
		if _, err := studentRepo.Create(s); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": tr(c, "users.student_create_failed")})
			return
		}

	case "teacher":
		t := builder.NewTeacherBuilder().
			WithLastName(req.LastName).
			WithFirstName(req.FirstName).
			WithPhoneNumber(req.PhoneNumber).
			WithDepositBalance(initialBalance).
			WithEmail(req.Email).
			WithLogin(login).
			WithPassword(password).
			WithBirthDate(birthDate).
			WithAddressId(addressId).
			WithDepartment(&models.Department{Id: req.DepartmentId}).
			Build()
		teacherRepo := repositories.NewTeacherRepository(db.GlobalDBO)
		if _, err := teacherRepo.Create(t); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": tr(c, "users.teacher_create_failed")})
			return
		}

	case "individual":
		u := builder.NewUserBuilder().
			WithLastName(req.LastName).
			WithFirstName(req.FirstName).
			WithPhoneNumber(req.PhoneNumber).
			WithDepositBalance(initialBalance).
			WithEmail(req.Email).
			WithLogin(login).
			WithPassword(password).
			WithBirthDate(birthDate).
			WithAddressId(addressId).
			Build()
		if _, err := userRepo.Create(u); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": tr(c, "users.individual_create_failed")})
			return
		}
	}

	c.JSON(http.StatusCreated, CreateUserResponse{
		Login:    login,
		Password: password,
		Message:  tr(c, "users.created"),
	})
}

// --- helpers ---
// To keep things simple, accents are not stored
func normalizeString(s string) string {
	r := strings.NewReplacer(
		"à", "a", "â", "a", "ä", "a",
		"é", "e", "è", "e", "ê", "e", "ë", "e",
		"î", "i", "ï", "i",
		"ô", "o", "ö", "o",
		"ù", "u", "û", "u", "ü", "u",
		"ç", "c",
		"À", "A", "Â", "A", "Ä", "A",
		"É", "E", "È", "E", "Ê", "E", "Ë", "E",
		"Î", "I", "Ï", "I",
		"Ô", "O", "Ö", "O",
		"Ù", "U", "Û", "U", "Ü", "U",
		"Ç", "C",
	)
	return r.Replace(s)
}

func generateLogin(lastName, firstName string) string {
	normLastName := strings.ToLower(strings.ReplaceAll(normalizeString(lastName), " ", ""))
	normFirstName := strings.ToLower(strings.ReplaceAll(normalizeString(firstName), " ", ""))
	if len(normFirstName) == 0 {
		return normLastName
	}
	return string(normFirstName[0]) + normLastName
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(string(s[0])) + strings.ToLower(s[1:])
}
