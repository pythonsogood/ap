package repository

import (
	"database/sql"

	"github.com/pythonsogood/ap-assignment1/doctor/model"
)

type DoctorRepository interface {
	FindByID(ID string) (*model.Doctor, error)
	All() ([]*model.Doctor, error)
	Insert(doctor *model.Doctor) error
}

type SQLiteDoctorRepository struct {
	db *sql.DB
}

func NewSQLiteDoctorRepository(db *sql.DB) *SQLiteDoctorRepository {
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

	doctor := &model.Doctor{}

	if err := statement.QueryRow(ID).Scan(&doctor); err != nil {
		return nil, err
	}

	return doctor, nil
}

func (r *SQLiteDoctorRepository) All() ([]*model.Doctor, error) {
	rows, err := r.db.Query("SELECT ID, full_name, specialization, email FROM doctors")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var doctors []*model.Doctor

	for rows.Next() {
		doctor := &model.Doctor{}

		err := rows.Scan(&doctor)

		if err == nil {
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
