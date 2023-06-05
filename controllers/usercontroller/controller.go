package usercontroller

import (
	"bankai/models"
	"bankai/services/userService"
	"bankai/utils"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo"
)

type UserController struct {
	UserService userService.UserService
}

type loginReq struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

func (u *UserController) Signup(c echo.Context) error {
	signupReq := &loginReq{}
	c.Bind(&signupReq)
	// todo check is user name is already registered
	newUser := models.User{
		Username: signupReq.UserName,
		Password: signupReq.Password,
		Admin:    true,
	}
	err := u.UserService.CreateUser(&newUser)
	if err != nil {
		return echo.ErrInternalServerError
	}
	return c.JSON(http.StatusCreated, nil)
}

func (u *UserController) Login(c echo.Context) error {
	loginReq := &loginReq{}
	c.Bind(&loginReq)
	user, err := u.UserService.GetUser(loginReq.UserName)
	if err != nil {
		return echo.ErrInternalServerError
	}
	fmt.Println(utils.HashPassword(loginReq.Password))
	if utils.ValidatePassword(user.Password, loginReq.Password) {
		token, err := utils.GenerateTokenPair(user)
		if err != nil {
			return echo.ErrInternalServerError
		}
		return c.JSON(http.StatusOK, token)
	}
	return echo.ErrUnauthorized
}

func (u *UserController) GetTime(c echo.Context) error {
	return c.JSON(http.StatusOK, time.Now())
}

// This is the api to refresh tokens
// Most of the code is taken from the jwt-go package's sample codes
// https://godoc.org/github.com/dgrijalva/jwt-go#example-Parse--Hmac
func (u *UserController) GetToken(c echo.Context) error {
	type tokenReqBody struct {
		RefreshToken string `json:"refresh_token"`
	}
	tokenReq := tokenReqBody{}
	c.Bind(&tokenReq)

	// Parse takes the token string and a function for looking up the key.
	// The latter is especially useful if you use multiple keys for your application.
	// The standard is to use 'kid' in the head of the token to identify
	// which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(tokenReq.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte("secret"), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims)
		// Get the user record from database or
		// run through your business logic to verify if the user can log in
		if int(claims["sub"].(float64)) == 1 {
			newTokenPair, err := utils.GenerateTokenPair(nil)
			if err != nil {
				return err
			}

			return c.JSON(http.StatusOK, newTokenPair)
		}

		return echo.ErrUnauthorized
	}

	return err
}
