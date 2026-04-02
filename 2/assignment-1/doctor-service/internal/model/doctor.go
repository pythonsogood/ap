package model

const DoctorTableName = "doctors"
const DoctorTableCreateQuery = `
	CREATE TABLE IF NOT EXISTS ` + DoctorTableName + ` (
		id TEXT UNIQUE NOT NULL,
		full_name TEXT NOT NULL,
		specialization TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL
	)
`

type Doctor struct {
	ID             string `json:"id"`
	FullName       string `json:"full_name"`
	Specialization string `json:"specialization"`
	Email          string `json:"email"`
}

func (*Doctor) TableName() string {
	return DoctorTableName
}

func (*Doctor) TableCreateQuery() string {
	return DoctorTableCreateQuery
}

func NewDoctor(id string, full_name string, specialization string, email string) *Doctor {
	return &Doctor{
		ID:             id,
		FullName:       full_name,
		Specialization: specialization,
		Email:          email,
	}
}
