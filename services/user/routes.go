package user

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/phildehovre/go-complete-api/config"
	"github.com/phildehovre/go-complete-api/services/auth"
	"github.com/phildehovre/go-complete-api/types"
	"github.com/phildehovre/go-complete-api/utils"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", h.handleLogin).Methods("POST")
	router.HandleFunc("/register", h.handleRegister).Methods("POST")
	router.HandleFunc("/users/:id", h.handleGetUserByID).Methods("GET")
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var payload types.LoginUserPayload
	// parse payload
	if err := utils.ParseJSON(r, &payload); err != nil {
		fmt.Println(payload)
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate payload
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload of %v", errors))
		return
	}
	// check if user exists
	user, _ := h.store.GetUserByEmail(payload.Email)
	if user == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with email %s does not exist", payload.Email))
		return
	}

	if !auth.ComparePasswords(user.Password, []byte(payload.Password)) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("not found, invalid email or password"))
		return
	}

	secret := []byte(config.Envs.JWTSecret)
	token, err := auth.CreateJWT(secret, user.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("Couldnot authenticate at this time"))
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	// get json payload
	var payload types.RegisterUserPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		fmt.Println(payload)
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload of %v", errors))
		return
	}
	// check if user exists
	_, err := h.store.GetUserByEmail(payload.Email)
	if err == nil {
		// check for non-existent err means there already is a user with the email address
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with email %s already exists", payload.Email))
	}

	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
	}

	// if user does not exist, create new user
	err = h.store.CreateUser(types.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  hashedPassword,
	})

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)

}

func (h *Handler) handleGetUserByID(w http.ResponseWriter, r *http.Request) {
	id := strings.Split(r.URL.Path, "/")

	fmt.Println(id)
	// user, _ := h.store.GetUserById(id)
	// if user == nil {
	// 	utils.WriteError(w, http.StatusNotFound, fmt.Errorf("user with id %d not found", id))
	// }
	utils.WriteJSON(w, http.StatusFound, id)
}
