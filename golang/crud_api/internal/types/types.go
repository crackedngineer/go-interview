package types

type Student struct {
	ID    int64
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
	Age   int    `json:"age" validate:"required,min=0"`
}
