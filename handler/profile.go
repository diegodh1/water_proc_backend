package handler

import (
	"encoding/json"
	"net/http"
	models "water_proccesing/model"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

//CreateProfile creates a new profile at the database
func CreateProfile(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	profile := models.AppProfile{}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&profile); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	profileTemp := getProfileOrNull(db, profile.ProfileID, profile.CompanyID)
	if profileTemp != nil {
		respondError(w, http.StatusBadRequest, "Ya existe un perfil con este ID para esta empresa")
		return
	}
	defer r.Body.Close()
	if result := db.Create(&profile); result.Error != nil || result.RowsAffected == 0 {
		if result.Error != nil {
			respondError(w, http.StatusBadRequest, result.Error.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, "Error No se pudo realizar el registro")
		return
	}
	respondJSON(w, http.StatusCreated, profile)
}

//UpdateProfile this function change the internal states of a profile
func UpdateProfile(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	profileID := vars["ProfileID"]
	profile := models.AppProfile{}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&profile); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	profile.ProfileID = profileID
	defer r.Body.Close()
	profileTemp := getProfileOrNull(db, profile.ProfileID, profile.CompanyID)
	if profileTemp == nil {
		respondError(w, http.StatusBadRequest, "Este Perfil No existe")
		return
	}
	if err := db.Model(&profile).Select("ProfileDescription", "ProfileStatus").Updates(profile).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, profile)
}

// get a user whose primary key is equal to the param
func getProfileOrNull(db *gorm.DB, profileID string, companyID string) *models.AppProfile {
	profile := models.AppProfile{}
	if err := db.First(&profile, models.AppProfile{ProfileID: profileID, CompanyID: companyID}).Error; err != nil {
		return nil
	}
	return &profile
}
