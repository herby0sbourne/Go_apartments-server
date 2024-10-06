package routes

import (
	"appartments-server/database"
	"appartments-server/models"
	"appartments-server/utils"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/kataras/iris/v12"
	"golang.org/x/crypto/bcrypt"
)

// Create/Register a new User
func Register(ctx iris.Context) {
	var userInput RegisterUser

	err := ctx.ReadJSON(&userInput)

	if err != nil {
		utils.HandleValidationErrors(err, ctx)
		// fmt.Println(err.Error())
		return
	}

	var newUser models.User
	usersExists, userError := checkUserExists(&newUser, userInput.Email)

	if userError != nil {
		utils.CreateInternalServerError(ctx)
		// fmt.Println(err.Error())
		return
	}

	if usersExists {
		utils.UserRegisterAlready(ctx)
		return
	}

	hashedPasword, hashErr := hashPasword(userInput.Password)

	if hashErr != nil {
		utils.CreateInternalServerError(ctx)
		return
	}

	newUser = models.User{
		FirstName:   userInput.FirstName,
		LastName:    userInput.LastName,
		Email:       strings.ToLower(userInput.Email),
		Password:    hashedPasword,
		SocialLogin: false,
	}

	database.DB.Create(&newUser)

	ctx.JSON(iris.Map{
		"ID":        newUser.ID,
		"firstName": newUser.FirstName,
		"lastName":  newUser.LastName,
		"email":     newUser.Email,
	})
}

func Login(ctx iris.Context) {
	var userInput LoginUserStruct

	err := ctx.ReadJSON(&userInput)
	if err != nil {
		utils.HandleValidationErrors(err, ctx)
		return
	}

	var existingUser models.User
	statusTitle := "Credentials Error"

	errorMsg := "Invalid email or password"

	userExists, userExistsErr := checkUserExists(&existingUser, userInput.Email)

	if userExistsErr != nil {
		utils.CreateInternalServerError(ctx)
		return
	}

	if userExists == false {
		utils.CreateError(iris.StatusUnauthorized, statusTitle, errorMsg, ctx)
		return
	}

	if existingUser.SocialLogin {
		utils.CreateError(iris.StatusUnauthorized, statusTitle, "Social Login Account", ctx)
		return
	}

	passwordErr := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(userInput.Password))
	if passwordErr != nil {
		utils.CreateError(iris.StatusUnauthorized, statusTitle, errorMsg, ctx)
		return
	}

	ctx.JSON(iris.Map{
		"ID":        existingUser.ID,
		"firstName": existingUser.FirstName,
		"lastName":  existingUser.LastName,
		"email":     existingUser.Email,
	})
}

func FacebookBookLoginOrSignUp(ctx iris.Context) {
	var userInput FacebookOrGoogleOAuth

	err := ctx.ReadJSON(&userInput)
	if err != nil {
		utils.HandleValidationErrors(err, ctx)
		return
	}

	endpoint := fmt.Sprintf("https://graph.facebook.com/me?fields=id,name,email&access_token=%s", userInput.AccessToken)
	client := &http.Client{}

	req, _ := http.NewRequest("GET", endpoint, nil)
	res, facebookErr := client.Do(req)

	if facebookErr != nil {
		utils.CreateInternalServerError(ctx)
		return
	}

	defer res.Body.Close()
	body, bodyErr := io.ReadAll(res.Body)

	if bodyErr != nil {
		log.Panic(bodyErr)
		utils.CreateInternalServerError(ctx)
		return
	}

	var facebookBody FacebookUserRes
	json.Unmarshal(body, &facebookBody)

	if facebookBody.Email != "" {
		var user models.User

		userExists, userExistsErr := checkUserExists(&user, facebookBody.Email)

		if userExistsErr != nil {
			utils.CreateInternalServerError(ctx)
			return

		}

		if userExists == false {
			firstLastName := strings.SplitN(facebookBody.Name, " ", 2)

			user = models.User{
				FirstName:      firstLastName[0],
				LastName:       firstLastName[1],
				Email:          facebookBody.Email,
				SocialLogin:    true,
				SocialProvider: "Facebook",
			}

			database.DB.Create(&user)

			ctx.JSON(iris.Map{
				"ID":        user.ID,
				"firstName": user.FirstName,
				"lastName":  user.LastName,
				"email":     user.Email,
			})

			return
		}

		if user.SocialLogin == true && user.SocialProvider == "Facebook" {
			ctx.JSON(iris.Map{
				"ID":        user.ID,
				"firstName": user.FirstName,
				"lastName":  user.LastName,
				"email":     user.Email,
			})

			return
		}

		utils.UserRegisterAlready(ctx)
		return
	}
}

func GoogleLoginOrSignUp(ctx iris.Context) {
	var userInput FacebookOrGoogleOAuth

	err := ctx.ReadJSON(&userInput)
	if err != nil {
		utils.HandleValidationErrors(err, ctx)
		return
	}

	endpoint := fmt.Sprintf("https://www.googleapis.com/oauth2/v3/userinfo?access_token=%s", userInput.AccessToken)
	client := &http.Client{}

	req, _ := http.NewRequest("GET", endpoint, nil)
	res, googleErr := client.Do(req)

	if googleErr != nil {
		utils.CreateInternalServerError(ctx)
		return
	}

	defer res.Body.Close()
	body, bodyErr := io.ReadAll(res.Body)

	if bodyErr != nil {
		log.Panic(bodyErr)
		utils.CreateInternalServerError(ctx)
		return
	}

	var googleBody GoogleUserRes
	json.Unmarshal(body, &googleBody)

	if googleBody.Email != "" {
		var user models.User

		userExists, userExistsErr := checkUserExists(&user, googleBody.Email)

		if userExistsErr != nil {
			utils.CreateInternalServerError(ctx)
			return

		}

		if userExists == false {
			firstLastName := strings.SplitN(googleBody.Name, " ", 2)

			user = models.User{
				FirstName:      firstLastName[0],
				LastName:       firstLastName[1],
				Email:          googleBody.Email,
				SocialLogin:    true,
				SocialProvider: "Google",
			}

			database.DB.Create(&user)

			ctx.JSON(iris.Map{
				"ID":        user.ID,
				"firstName": user.FirstName,
				"lastName":  user.LastName,
				"email":     user.Email,
			})

			return
		}

		if user.SocialLogin == true && user.SocialProvider == "Google" {
			ctx.JSON(iris.Map{
				"ID":        user.ID,
				"firstName": user.FirstName,
				"lastName":  user.LastName,
				"email":     user.Email,
			})

			return
		}

		utils.UserRegisterAlready(ctx)
		return
	}

}

func checkUserExists(user *models.User, email string) (exists bool, err error) {
	doesUserExist := database.DB.Where("email = ?", strings.ToLower(email)).Limit(1).Find(&user)

	if doesUserExist.Error != nil {
		return false, doesUserExist.Error
	}

	userExist := doesUserExist.RowsAffected > 0

	if userExist {
		return true, nil
	}

	return false, nil

}

func hashPasword(password string) (hashedPassword string, err error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

type RegisterUser struct {
	FirstName string `json:"firstName" validate:"required,max=250"`
	LastName  string `json:"LastName" validate:"required,max=250"`
	Email     string `json:"email" validate:"required,max=250,email"`
	Password  string `json:"password" validate:"required,min=8,max=250"`
}

type LoginUserStruct struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type FacebookOrGoogleOAuth struct {
	AccessToken string `json:"accessToken" validate:"required"`
}

type FacebookUserRes struct {
	ID    string `json:"id"`
	Name  string `json:"name" `
	Email string `json:"email" `
}

type GoogleUserRes struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
}
