package models

type Student struct {
	User
	StudyYear string `json:"study_year"`
}

func NewStudent(lastName string, firstName string, login string, password string, studyYear string) *Student {
	return &Student{
		User: User{
			LastName:  lastName,
			FirstName: firstName,
			Login:     login,
			Password:  password,
		},
		StudyYear: studyYear,
	}
}
