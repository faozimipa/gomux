package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
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

func (server *Server) Refresh(w http.ResponseWriter, r *http.Request) {
	var rt map[string]string
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rt); err != nil {
		responses.ERROR(w, http.StatusBadRequest, errors.New("Invalid request payload"))
		return
	}
	defer r.Body.Close()
	refreshToken, ok := rt["refresh_token"]
	if ok {
		token, err := auth.IsValidRefreshToken(refreshToken)

		if err != nil {
			responses.ERROR(w, http.StatusBadRequest, errors.New("Invalid request payload"))
			return
		}
		//Since token is valid, get the uuid:
		claims, ok := token.Claims.(jwt.MapClaims)
		if ok && token.Valid {
			refreshUuid, ok := claims["refresh_uuid"].(string)
			if !ok {
				responses.ERROR(w, http.StatusUnprocessableEntity, err)
				return
			}
			userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
			if err != nil {
				responses.ERROR(w, http.StatusUnprocessableEntity, errors.New("Error occurred"))
				return
			}
			//Delete the previous Refresh Token
			deleted, delErr := DeleteAuth(refreshUuid)
			if delErr != nil || deleted == 0 {
				responses.ERROR(w, http.StatusUnauthorized, errors.New("unauthorized"))
				return
			}
			//Create new pairs of refresh and access tokens
			ts, createErr := auth.CreateToken(uint32(userId))
			if createErr != nil {
				responses.ERROR(w, http.StatusForbidden, createErr)
				return
			}
			//save the tokens metadata to redis
			saveErr := auth.CreateAuth(uint32(userId), ts)
			if saveErr != nil {
				responses.ERROR(w, http.StatusForbidden, saveErr)
				return
			}
			tokens := map[string]string{
				"access_token":  ts.AccessToken,
				"refresh_token": ts.RefreshToken,
			}
			responses.JSON(w, http.StatusCreated, tokens)
		} else {
			responses.ERROR(w, http.StatusUnauthorized, errors.New("refresh expired"))
		}
	} else {
		responses.ERROR(w, http.StatusBadRequest, errors.New("Invalid token"))
		return
	}
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
