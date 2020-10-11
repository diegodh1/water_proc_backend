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
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := db.Where("app_user_status = ?", 1).First(&user, userTemp).Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, user)
}

//CreateUser creates a new AppUser
func CreateUser(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	user := models.AppUser{}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&user); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	userTemp := getUserOrNull(db, user.AppUserID, w, r)
	if userTemp != nil {
		respondError(w, http.StatusBadRequest, "Ya existe un usuario con este ID")
		return
	}
	defer r.Body.Close()
	//hashing the password
	pass := user.AppUserPassword
	hashPass, err := bcrypt.GenerateFromPassword([]byte(pass), 10)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	s := bytes.NewBuffer(hashPass).String()
	user.AppUserPassword = s
	//end hashing

	if result := db.Create(&user); result.Error != nil || result.RowsAffected == 0 {
		if result.Error != nil {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, "Error No se pudo realizar el registro")
		return
	}
	respondJSON(w, http.StatusCreated, user)
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
