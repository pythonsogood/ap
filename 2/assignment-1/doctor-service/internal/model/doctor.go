package model

type Doctor struct {
	ID             string `json:"id"`
	FullName       string `json:"full_name"`
	Specialization string `json:"specialization"`
	Email          string `json:"email"`
}

func NewDoctor(id string, full_name string, specialization string, email string) *Doctor {
	return &Doctor{
		ID:             id,
		FullName:       full_name,
		Specialization: specialization,
		Email:          email,
	}
}
