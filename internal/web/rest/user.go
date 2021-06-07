package rest

import (
	"encoding/json"
	distrybute "github.com/mmichaelb/distrybute/internal"
	"github.com/rs/zerolog"
	"net/http"
	"regexp"
	"unicode"
)

type UserCreateState string

const (
	userCreatedState         = UserCreateState("USER_CREATED")
	userAlreadyExistentState = UserCreateState("USER_ALREADY_EXISTENT")
	usernameInvalidState     = UserCreateState("USERNAME_INVALID")
	passwordInvalidState     = UserCreateState("PASSWORD_INVALID")
)

type UserCreateRequest struct {
	Username string `json:"username"`
	Password []byte `json:"password"`
}

type UserCreateResponse struct {
	Username           string          `json:"username,omitempty"`
	AuthorizationToken string          `json:"authorization_token,omitempty"`
	State              UserCreateState `json:"state"`
}

var usernameRegex = regexp.MustCompile("^\\w{4,}$")

const passwordMinLength = 8

func (r *router) handleUserCreate(w *responseWriter, req *http.Request) {
	var parsedReq UserCreateRequest
	err := json.NewDecoder(req.Body).Decode(&parsedReq)
	if _, ok := err.(*json.UnmarshalTypeError); ok {
		w.WriteAutomaticErrorResponse(http.StatusBadRequest, nil, req)
		return
	} else if err != nil {
		r.log(zerolog.ErrorLevel, req).Err(err).Msg("could not unmarshal user request body due to an unknown error")
		w.WriteAutomaticErrorResponse(http.StatusInternalServerError, nil, req)
		return
	}
	if !validateUsername(parsedReq.Username) {
		w.WriteAutomaticErrorResponse(http.StatusBadRequest, &UserCreateResponse{State: usernameInvalidState}, req)
		return
	}
	if !validatePassword(parsedReq.Password) {
		w.WriteAutomaticErrorResponse(http.StatusBadRequest, &UserCreateResponse{State: passwordInvalidState}, req)
		return
	}
	user, err := r.userService.CreateNewUser(parsedReq.Username, parsedReq.Password)
	if err == distrybute.ErrUserAlreadyExists {
		w.WriteAutomaticErrorResponse(http.StatusBadGateway, &UserCreateResponse{State: userAlreadyExistentState}, req)
		return
	} else if err != nil {
		r.log(zerolog.ErrorLevel, req).Err(err).Str("createUsername", parsedReq.Username).Msg("could not create new user")
		w.WriteAutomaticErrorResponse(http.StatusInternalServerError, nil, req)
		return
	}
	w.WriteResponse(http.StatusOK, "", &UserCreateResponse{
		Username:           user.Username,
		AuthorizationToken: user.AuthorizationToken,
		State:              userCreatedState,
	}, req)
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
