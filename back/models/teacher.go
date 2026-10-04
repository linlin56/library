package models

type Teacher struct {
	User
	DepartmentId int         `json:"department_id"`
	Department   *Department `json:"department,omitempty"`
}

func NewTeacher(lastName string, firstName string, login string, password string, department *Department) *Teacher {
	return &Teacher{
		User: User{
			LastName:  lastName,
			FirstName: firstName,
			Login:     login,
			Password:  password,
		},
		Department: department,
	}
}
