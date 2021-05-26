package rest

import (
	"encoding/json"
	distrybute "github.com/mmichaelb/distrybute/internal"
	"github.com/rs/zerolog/log"
	"net/http"
	"regexp"
	"unicode"
)

type UserCreateRequest struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	Password []byte `json:"password"`
}

type UserCreateResponse struct {
	Username           string `json:"username"`
	AuthorizationToken string `json:"authorization_token"`
}

var usernameRegex = regexp.MustCompile("^\\w{4,}$")

const passwordMinLength = 8

func (r *Router) handleUserCreate(w http.ResponseWriter, req *http.Request) {
	var parsedReq UserCreateRequest
	err := json.NewDecoder(req.Body).Decode(&parsedReq)
	if _, ok := err.(*json.UnmarshalTypeError); ok {
		r.sendAutomaticError(w, req, http.StatusBadRequest)
		return
	} else if err != nil {
		log.Err(err).Msg("could not unmarshal user request body")
		r.sendInternalServerError(w, req)
		return
	}
	if !validateUsername(parsedReq.Username) {
		r.sendError(w, req, http.StatusBadRequest, "username does not fulfill requirements")
		return
	}
	if !validatePassword(parsedReq.Password) {
		r.sendError(w, req, http.StatusBadRequest, "password does not fulfill requirements")
		return
	}
	user, err := r.userService.CreateNewUser(parsedReq.Username, distrybute.Role(parsedReq.Role), parsedReq.Password)
	if err != nil {
		log.Err(err).Str("createUsername", parsedReq.Username).Msg("could not create new user")
		r.sendInternalServerError(w, req)
		return
	}
	r.sendResponse(w, req, &UserCreateResponse{
		Username:           user.Username,
		AuthorizationToken: user.AuthorizationToken,
	})
}

func validateUsername(username string) bool {
	return username != "" && usernameRegex.MatchString(username)
}

func validatePassword(password []byte) bool {
	if len(password) < passwordMinLength {
		return false
	}
	check := 0
	for i := range password {
		toCheck := rune(password[i])
		if unicode.IsUpper(toCheck) && check&1 != 1 {
			check = check | 1
		}
		if unicode.IsLower(toCheck) && check&2 != 2 {
			check = check | 2
		}
		if unicode.IsDigit(toCheck) && check&4 != 4 {
			check = check | 4
		}
	}
	return check == 7
}
