package handler

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	models "water_proccesing/model"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

//GetCompanies get all the companies available at hte database
func GetCompanies(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	companies := []models.Company{}
	db.Where("company_status = ?", 1).Find(&companies)
	respondJSON(w, http.StatusOK, companies)
}

//UserLogin this function allows the user log in if the userID an password given are correct
func UserLogin(db *gorm.DB, w http.ResponseWriter, r *http.Request, seed string) {
	user := models.AppUser{}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&user); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	userTemp := getUserOrNull(db, user.AppUserID, w, r)
	if userTemp == nil {
		respondJSON(w, http.StatusUnauthorized, JSONResponse{Message: "el usuario no está registrado ó está inactivo"})
		return
	}
	errPass := bcrypt.CompareHashAndPassword([]byte(userTemp.AppUserPassword), []byte(user.AppUserPassword))
	if errPass != nil {
		respondJSON(w, http.StatusUnauthorized, JSONResponse{Message: "usuario y/o contraseña incorrecta"})
		return
	}
	v, err := createToken(user.AppUserID, seed)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
	}
	defer r.Body.Close()
	respondJSON(w, http.StatusOK, JSONResponse{Payload: v, Message: "Ingreso Realizado!"})
}

//This function creates a new token when a user log in to te app
func createToken(userid string, seed string) (string, error) {
	var err error
	//Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userid
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv(seed)))
	if err != nil {
		return "", err
	}
	return token, nil
}
