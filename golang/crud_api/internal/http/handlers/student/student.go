package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/crackedngineer/go-interview/internal/storage"
	"github.com/crackedngineer/go-interview/internal/types"
	"github.com/crackedngineer/go-interview/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var student types.Student
		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			response.WriteJSON(w, http.StatusBadRequest, response.GenericError(fmt.Errorf("empty body")))
			return
		}
		if err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GenericError(err))
			return
		}

		// Validation
		if err := validator.New().Struct(student); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJSON(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		id, err := storage.CreateStudent(student.Name, student.Email, student.Age)
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, response.GenericError(err))
			return
		}

		// Process the student data (e.g., save to database)
		response.WriteJSON(w, http.StatusCreated, map[string]int64{"id": id})
	}
}

func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			response.WriteJSON(w, http.StatusBadRequest, response.GenericError(errors.New("missing id parameter")))
			return
		}

		_id, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GenericError(fmt.Errorf("invalid id parameter: %v", err)))
			return
		}

		student, err := storage.GetStudent(_id)
		if err != nil {
			response.WriteJSON(w, http.StatusNotFound, response.GenericError(err))
			return
		}

		response.WriteJSON(w, http.StatusOK, student)
	}
}

func GetAll(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		students, err := storage.GetAllStudents()
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, response.GenericError(err))
			return
		}
		response.WriteJSON(w, http.StatusOK, students)
	}
}
