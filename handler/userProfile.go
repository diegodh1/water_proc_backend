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
		respondJSON(w, http.StatusBadRequest, JSONResponse{Message: "Error al recibir la petición"})
		return
	}
	defer r.Body.Close()
	UserTemp := getUserProfileOrNull(db, userProfile.ProfileID, userProfile.AppUserID)
	if UserTemp != nil {
		if err := db.Model(&userProfile).Where("app_user_id = ? and profile_id = ?", userProfile.AppUserID, userProfile.ProfileID).Omit("UserProfileCreationDate").Save(userProfile).Error; err != nil {
			respondJSON(w, http.StatusInternalServerError, JSONResponse{Message: "Error No se pudo realizar el registro"})
			return
		}
	}
	if UserTemp == nil {
		if result := db.Create(&userProfile); result.Error != nil || result.RowsAffected == 0 {
			if result.Error != nil {
				respondJSON(w, http.StatusBadRequest, JSONResponse{Message: "Error No se pudo realizar el registro"})
				return
			}
			respondJSON(w, http.StatusInternalServerError, JSONResponse{Message: "Error No se pudo realizar el registro"})
			return
		}
	}
	respondJSON(w, http.StatusCreated, JSONResponse{Message: "Operación realizada con éxito!"})
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
	userProfileTemp := getUserProfileOrNull(db, userProfile.ProfileID, userProfile.AppUserID)
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
func getUserProfileOrNull(db *gorm.DB, profileID string, appUserID string) *models.UserProfile {
	userProfile := models.UserProfile{}
	if err := db.First(&userProfile, models.UserProfile{ProfileID: profileID, AppUserID: appUserID}).Error; err != nil {
		return nil
	}
	return &userProfile
}
