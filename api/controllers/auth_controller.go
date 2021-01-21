package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/faozimipa/gomux/api/auth"
	"github.com/faozimipa/gomux/api/models"
	"github.com/faozimipa/gomux/api/responses"
	"github.com/faozimipa/gomux/api/utils/formaterror"
	"github.com/faozimipa/gomux/api/utils/redismanager"
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

	getUser, _ := user.FindUserByEmail(server.DB, user.Email)
	saveErr := auth.CreateAuth(getUser.ID, ts)

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

func DeleteAuth(givenUuid string) (int64, error) {
	c := redismanager.InitRedisClient()
	defer c.Close()
	deleted, err := c.Del(ctx, givenUuid).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}

func (server *Server) Logout(w http.ResponseWriter, r *http.Request) {

	au, err := auth.ExtractTokenMetadata(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	deleted, delErr := DeleteAuth(au.AccessUuid)
	if delErr != nil || deleted == 0 {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	responses.JSON(w, http.StatusOK, "Successfully logged out")
}
