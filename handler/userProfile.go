package handler

import (
	"encoding/json"
	"net/http"
	models "water_proccesing/model"

	"gorm.io/gorm"
)

//CreateUserProfile assigns a profile to a user
func CreateUserProfile(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	userProfile := models.UserProfile{}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&userProfile); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	userProfileTemp := getUserProfileOrNull(db, userProfile.ProfileID, userProfile.CompanyID, userProfile.AppUserID)
	if userProfileTemp != nil {
		respondError(w, http.StatusBadRequest, "El usuario ya tiene asignado este perfil")
		return
	}
	defer r.Body.Close()
	if result := db.Create(&userProfile); result.Error != nil || result.RowsAffected == 0 {
		if result.Error != nil {
			respondError(w, http.StatusBadRequest, result.Error.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, "Error No se pudo realizar el registro")
		return
	}
	respondJSON(w, http.StatusCreated, userProfile)
}

//UpdateUserProfile this function change the internal states of a UserProfile
func UpdateUserProfile(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	userProfile := models.UserProfile{}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&userProfile); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()
	userProfileTemp := getUserProfileOrNull(db, userProfile.ProfileID, userProfile.CompanyID, userProfile.AppUserID)
	if userProfileTemp == nil {
		respondError(w, http.StatusBadRequest, "Este usuario no tiene asignado este perfil")
		return
	}
	if err := db.Model(&userProfile).Select("UserProfileStatus").Updates(userProfile).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, userProfile)
}

// get a UserProfile whose primary key is equal to the param
func getUserProfileOrNull(db *gorm.DB, profileID string, companyID string, appUserID string) *models.UserProfile {
	userProfile := models.UserProfile{}
	if err := db.First(&userProfile, models.UserProfile{ProfileID: profileID, CompanyID: companyID, AppUserID: appUserID}).Error; err != nil {
		return nil
	}
	return &userProfile
}
