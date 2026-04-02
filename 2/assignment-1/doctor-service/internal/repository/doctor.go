package repository

import (
	"database/sql"

	"github.com/pythonsogood/ap-assignment1/doctor/internal/model"
)

type DoctorRepository interface {
	FindByID(ID string) (*model.Doctor, error)
	All() ([]*model.Doctor, error)
	Insert(doctor *model.Doctor) error
}

type SQLiteDoctorRepository struct {
	db *sql.DB
}

func NewSQLiteDoctorRepository(db *sql.DB) DoctorRepository {
	return &SQLiteDoctorRepository{
		db: db,
	}
}

func (r *SQLiteDoctorRepository) FindByID(ID string) (*model.Doctor, error) {
	statement, err := r.db.Prepare("SELECT id, full_name, specialization, email FROM doctors WHERE id = ?")

	if err != nil {
		return nil, err
	}

	defer statement.Close()

	var id, full_name, specialization, email string

	if err := statement.QueryRow(ID).Scan(&id, &full_name, &specialization, &email); err != nil {
		return nil, err
	}

	return &model.Doctor{
		ID:             id,
		FullName:       full_name,
		Specialization: specialization,
		Email:          email,
	}, nil
}

func (r *SQLiteDoctorRepository) All() ([]*model.Doctor, error) {
	rows, err := r.db.Query("SELECT id, full_name, specialization, email FROM doctors")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	doctors := make([]*model.Doctor, 0)

	for rows.Next() {
		var id, full_name, specialization, email string

		if err := rows.Scan(&id, &full_name, &specialization, &email); err == nil {
			doctor := &model.Doctor{
				ID:             id,
				FullName:       full_name,
				Specialization: specialization,
				Email:          email,
			}

			doctors = append(doctors, doctor)
		}
	}

	return doctors, nil
}

func (r *SQLiteDoctorRepository) Insert(doctor *model.Doctor) error {
	statement, err := r.db.Prepare("INSERT INTO doctors(id, full_name, specialization, email) VALUES(?, ?, ?, ?)")

	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(doctor.ID, doctor.FullName, doctor.Specialization, doctor.Email)

	if err != nil {
		return err
	}

	return nil
}
