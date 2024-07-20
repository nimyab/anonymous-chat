package handlers

import (
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
)

type UserDto struct {
	Name string `json:"name"`
}

type User struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

var authUsers = make(map[string]*User)

func GetUserById(id string) *User {
	return authUsers[id]
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}

	var userDto UserDto
	if err := json.NewDecoder(r.Body).Decode(&userDto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if userDto.Name == "" {
		http.Error(w, "User name is required", http.StatusBadRequest)
		return
	}

	newUser := User{Name: userDto.Name, Id: uuid.New().String()}
	authUsers[newUser.Id] = &newUser

	if err := json.NewEncoder(w).Encode(newUser); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
