package sqlite

import (
	"database/sql"

	"github.com/crackedngineer/go-interview/internal/config"
	"github.com/crackedngineer/go-interview/internal/types"
	_ "github.com/mattn/go-sqlite3"
)

type SqliteStorage struct {
	Db *sql.DB
}

func NewDb(config *config.Config) (*SqliteStorage, error) {
	db, err := sql.Open("sqlite3", config.StoragePath)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT NOT NULL,
		age INTEGER
	)`)

	if err != nil {
		return nil, err
	}

	return &SqliteStorage{Db: db}, nil
}

func (s *SqliteStorage) CreateStudent(name string, email string, age int) (int64, error) {
	stmt, err := s.Db.Prepare("INSERT INTO students (name, email, age) VALUES (?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	result, err := stmt.Exec(name, email, age)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *SqliteStorage) GetStudent(id int64) (types.Student, error) {
	var student types.Student
	stmt, err := s.Db.Prepare("SELECT id, name, email, age FROM students WHERE id = ?")
	if err != nil {
		return types.Student{}, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(id).Scan(&student.ID, &student.Name, &student.Email, &student.Age)
	if err != nil {
		return types.Student{}, err
	}
	return student, nil
}

func (s *SqliteStorage) GetAllStudents() ([]types.Student, error) {
	rows, err := s.Db.Query("SELECT id, name, email, age FROM students")
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var students []types.Student

	for rows.Next() {
		var student types.Student
		err := rows.Scan(&student.ID, &student.Name, &student.Email, &student.Age)

		if err != nil {
			return nil, err
		}
		students = append(students, student)
	}
	return students, nil
}
