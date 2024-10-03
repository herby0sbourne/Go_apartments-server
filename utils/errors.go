package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
)

func CreateError(statusCode int, title string, detail string, ctx iris.Context) {
	ctx.StopWithProblem(statusCode, iris.NewProblem().Title(title).Detail(detail))
}

func CreateInternalServerError(ctx iris.Context) {
	CreateError(iris.StatusInternalServerError, "Internal Server Error", "Internal Server Error", ctx)
}

func UserRegisterAlready(ctx iris.Context) {

	CreateError(iris.StatusConflict, "Conflict", "Email Already Exisits", ctx)

}

func HandleValidationErrors(err error, ctx iris.Context) {
	if errs, ok := err.(validator.ValidationErrors); ok {
		validationErrors := wrapValidationErrors(errs)

		fmt.Println("validationErrors", validationErrors)
		ctx.StopWithProblem(
			iris.StatusBadRequest,
			iris.NewProblem().
				Title("Validation Errors").
				Detail("one or more fields are missing, Validation Error").
				Key("errors", validationErrors))

		return
	}

	CreateInternalServerError(ctx)
}

func wrapValidationErrors(errs validator.ValidationErrors) []validationErrorStruct {
	validationErrors := make([]validationErrorStruct, 0, len(errs))

	for _, validationErr := range errs {
		validationErrors = append(validationErrors, validationErrorStruct{
			ActualTag: validationErr.ActualTag(),
			NameSoace: validationErr.Namespace(),
			Kind:      validationErr.Kind().String(),
			Type:      validationErr.Type().String(),
			Value:     fmt.Sprintf("%v", validationErr.Value()),
			Param:     validationErr.Param(),
		})
	}

	return validationErrors
}

type validationErrorStruct struct {
	ActualTag string `json:"tag"`
	NameSoace string `json:"namespace"`
	Kind      string `json:"kind"`
	Type      string `json:"type"`
	Value     string `json:"value"`
	Param     string `json:"params"`
}
