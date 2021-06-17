package rest

import (
	"encoding/json"
	distrybute "github.com/mmichaelb/distrybute/internal"
	"github.com/rs/zerolog/hlog"
	"net/http"
	"regexp"
	"unicode"
)

type UserLoginState string
type UserCreateState string

const (
	userNotFoundState        = UserLoginState("USER_NOT_FOUND")
	invalidPasswordState     = UserLoginState("INVALID_PASSWORD")
	loginSuccessfulState     = UserLoginState("LOGIN_SUCCESSFUL")
	userCreatedState         = UserCreateState("USER_CREATED")
	userAlreadyExistentState = UserCreateState("USER_ALREADY_EXISTENT")
	usernameInvalidState     = UserCreateState("USERNAME_INVALID")
	passwordInvalidState     = UserCreateState("PASSWORD_INVALID")
)

type UserRequest struct {
	Username string `json:"username"`
	Password []byte `json:"password"`
}

type UserLoginResponse struct {
	State UserLoginState `json:"state"`
}

type UserCreateResponse struct {
	Username           string          `json:"username,omitempty"`
	AuthorizationToken string          `json:"authorization_token,omitempty"`
	State              UserCreateState `json:"state"`
}

var usernameRegex = regexp.MustCompile("^\\w{4,}$")

const passwordMinLength = 8

func (r *router) handleUserLogin(w *responseWriter, req *http.Request) {
	var parsedReq UserRequest
	err := json.NewDecoder(req.Body).Decode(&parsedReq)
	if _, ok := err.(*json.UnmarshalTypeError); ok {
		w.WriteAutomaticErrorResponse(http.StatusBadRequest, nil, req)
		return
	} else if err != nil {
		hlog.FromRequest(req).Err(err).Msg("could not unmarshal user login request body due to an unknown error")
		w.WriteAutomaticErrorResponse(http.StatusInternalServerError, nil, req)
		return
	}
	ok, err := r.userService.CheckPassword(parsedReq.Username, parsedReq.Password)
	if err == distrybute.ErrUserAlreadyExists {
		hlog.FromRequest(req).Info().Str("username", parsedReq.Username).Msg("failed login with non-existent username")
		w.WriteResponse(http.StatusNotFound, "", &UserLoginResponse{userNotFoundState}, req)
		return
	} else if err != nil {
		hlog.FromRequest(req).Err(err).Str("username", parsedReq.Username).Msg("could not check password")
		w.WriteAutomaticErrorResponse(http.StatusInternalServerError, nil, req)
		return
	} else if !ok {
		hlog.FromRequest(req).Info().Str("username", parsedReq.Username).Msg("failed login attempt")
		w.WriteResponse(http.StatusUnauthorized, "", &UserLoginResponse{invalidPasswordState}, req)
		return
	}
	// TODO add jwt logic
	hlog.FromRequest(req).Info().Str("username", parsedReq.Username).Msg("successful login attempt")
	w.WriteSuccessfulResponse(&UserLoginResponse{loginSuccessfulState}, req)
}

func (r *router) handleUserCreate(w *responseWriter, req *http.Request) {
	var parsedReq UserRequest
	err := json.NewDecoder(req.Body).Decode(&parsedReq)
	if _, ok := err.(*json.UnmarshalTypeError); ok {
		w.WriteAutomaticErrorResponse(http.StatusBadRequest, nil, req)
		return
	} else if err != nil {
		hlog.FromRequest(req).Err(err).Msg("could not unmarshal user request body due to an unknown error")
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
		hlog.FromRequest(req).Err(err).Str("createUsername", parsedReq.Username).Msg("could not create new user")
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
