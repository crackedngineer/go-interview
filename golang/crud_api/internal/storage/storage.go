package storage

import "github.com/crackedngineer/go-interview/internal/types"

type Storage interface {
	// CreateStudent creates a new student record in the storage.
	CreateStudent(name string, email string, age int) (int64, error)
	// GetStudent retrieves a student record by its ID.
	GetStudent(id int64) (types.Student, error)
	// GetAllStudents retrieves all student records from the storage.
	GetAllStudents() ([]types.Student, error)
	// // UpdateStudent updates an existing student record by its ID.
	// UpdateStudent(id int64, student types.Student) error
	// // DeleteStudent deletes a student record by its ID.
	// DeleteStudent(id int64) error
}
