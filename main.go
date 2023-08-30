package main

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type ViewSchema struct {
	Title *string `validate:"required,alphanum"`
	FolderID *int `validate:"omitempty,gte=0"`
	OwnerID *int `validate:"omitempty,gte=0"`
	Description *string `validate:"omitempty"`
}

type PersonSchema struct {
	FirstName *string `validate:"required,alpha"`
	LastName *string `validate:"required,alpha"`
}

type Validator[T any] struct {
	Schema T
}

func (v *Validator[T]) validate(data T) (bool) {
	valid := validator.New()
	err := valid.Struct(data)
	if err != nil {
		fmt.Println("Error:", err.Error())
		for _, err := range err.(validator.ValidationErrors) {
			fmt.Println("Field:", err.Field())
			fmt.Println("Tag:", err.Tag())
			fmt.Println("Param:", err.Param())
			fmt.Println("Value:", err.Value())
		}
		return false
	}
	return true
}


func main() {
	ownerId := 1
	folderId := 1
	title := "Dummy title"
	desc := "Description"


	viewData := ViewSchema{
		Title: &title,
		FolderID: &folderId,
		OwnerID: &ownerId,
		Description: &desc,
	}

	viewValidator := new(Validator[ViewSchema])
	viewValid := viewValidator.validate(viewData)

	if viewValid {
		// If we got here, data is valid
		fmt.Println("View data is valid")
	}

	firstName := "Angela"
	//lastName := ""
	personData := PersonSchema{
		FirstName: &firstName,
		LastName: nil,
	}

	personValidator := new(Validator[PersonSchema])
	personValid := personValidator.validate(personData)

	if personValid {
		// If we got here, data is valid
		fmt.Println("Person data is valid")
	}
}
