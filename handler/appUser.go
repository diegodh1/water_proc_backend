package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	models "water_proccesing/model"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

//SarchUser get a AppUser whose ID is equal to param  given
func SarchUser(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	user := models.AppUser{}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	userTemp := models.AppUser{}
	if err := decoder.Decode(&userTemp); err != nil {
		respondJSON(w, http.StatusBadRequest, JSONResponse{Message: "Error interno del servidor"})
		return
	}
	if err := db.First(&user, userTemp).Error; err != nil {
		respondJSON(w, http.StatusNotFound, JSONResponse{Message: "Usuario no registrado"})
		return
	}
	user.AppUserPassword = ""

	anonymousStruct := struct {
		User     models.AppUser
		Perfiles []models.UserProfile
	}{user, []models.UserProfile{}}
	respondJSON(w, http.StatusOK, JSONResponse{Payload: anonymousStruct, Message: "Usuario encontrado"})
}

//CreateUser creates a new AppUser
func CreateUser(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	user := models.AppUser{}
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&user); err != nil {
		respondJSON(w, http.StatusBadRequest, JSONResponse{models.AppUser{}, "Erro interno del servidor"})
		return
	}
	userTemp := getUserOrNull(db, user.AppUserID, w, r)
	if userTemp != nil {
		respondJSON(w, http.StatusBadRequest, JSONResponse{models.AppUser{}, "Ya existe un usuario con este ID"})
		return
	}
	//hashing the password
	pass := user.AppUserPassword
	hashPass, err := bcrypt.GenerateFromPassword([]byte(pass), 10)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, JSONResponse{models.AppUser{}, "Error Interno del servidor"})
		return
	}
	s := bytes.NewBuffer(hashPass).String()
	user.AppUserPassword = s
	//end hashing

	if result := db.Create(&user); result.Error != nil || result.RowsAffected == 0 {
		if result.Error != nil {
			respondJSON(w, http.StatusBadRequest, JSONResponse{models.AppUser{}, err.Error()})
			return
		}
		respondJSON(w, http.StatusInternalServerError, JSONResponse{models.AppUser{}, "Error No se pudo realizar el registro"})
		return
	}
	respondJSON(w, http.StatusCreated, JSONResponse{user, "Registro realizado"})
}

//UpdateUser this function change the internal states of a AppUser
func UpdateUser(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	UserID := vars["UserID"]
	user := getUserOrNull(db, UserID, w, r)
	if user == nil {
		respondError(w, http.StatusBadRequest, "Usuario no registrado")
		return
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&user); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	//hashing the password
	pass := user.AppUserPassword
	hashPass, err := bcrypt.GenerateFromPassword([]byte(pass), 10)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	s := bytes.NewBuffer(hashPass).String()
	user.AppUserID = UserID
	user.AppUserPassword = s
	//end hashing
	defer r.Body.Close()

	if err := db.Model(&user).Omit("AppUserID", "CompanyID", "AppUserCreationDate").Updates(user).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, user)
}

// get a user whose AppUserID is equal to the params given
func getUserOrNull(db *gorm.DB, appUserID string, w http.ResponseWriter, r *http.Request) *models.AppUser {
	user := models.AppUser{}
	if err := db.First(&user, models.AppUser{AppUserID: appUserID}).Error; err != nil {
		return nil
	}
	return &user
}
