package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/faozimipa/gomux/api/auth"
	"github.com/faozimipa/gomux/api/models"
	"github.com/faozimipa/gomux/api/responses"
	"github.com/faozimipa/gomux/api/utils/formaterror"
	"golang.org/x/crypto/bcrypt"
)

func (server *Server) Login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	user.Prepare()
	err = user.Validate("login")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	ts, err := server.SignIn(user.Email, user.Password)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusUnprocessableEntity, formattedError)
		return
	}

	saveErr := auth.CreateAuth(user.ID, ts)

	if saveErr != nil {
		formattedError2 := formaterror.FormatError(saveErr.Error())
		responses.ERROR(w, http.StatusUnprocessableEntity, formattedError2)
	}

	tokens := map[string]string{
		"access_token":  ts.AccessToken,
		"refresh_token": ts.RefreshToken,
	}

	responses.JSON(w, http.StatusOK, tokens)
}

func (server *Server) SignIn(email, password string) (*models.TokenDetails, error) {

	var err error

	user := models.User{}

	err = server.DB.Debug().Model(models.User{}).Where("email = ?", email).Take(&user).Error
	if err != nil {
		return nil, err
	}
	err = models.VerifyPassword(user.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return nil, err
	}
	return auth.CreateToken(user.ID)
}
