package handler

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	models "water_proccesing/model"

	"gorm.io/gorm"
)

//ItemReport struct for post request
type ItemReport struct {
	Label string  `json:"label"`
	Value float64 `json:"value"`
}

//ItemList struct for post request
type ItemList struct {
	Values []ItemReport `json:"valores"`
	Titulo string       `json:"titulo"`
	Tipo   string       `json:"tipo"`
	Limite int          `json:"limite"`
}

//ImageResponse struct for post request
type ImageResponse struct {
	File string `json:"File"`
}

//GetRegistersDisponibilityDay get all the registers from the view by day
func GetRegistersDisponibilityDay(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	registerFilter := models.RegisterFilter{}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&registerFilter); err != nil {
		respondJSON(w, http.StatusBadRequest, JSONResponse{Message: "Error interno del servidor"})
		return
	}
	registers := []models.PretratamientoDisp{}

	if err := db.Group("for_date, on_hour, tag_index, total").Order("for_date,on_hour asc").Having("tag_index = ? AND for_date BETWEEN ? AND ?", registerFilter.TagIndex, registerFilter.StartDate, registerFilter.FinalDate).Find(&registers).Error; err != nil {
		respondJSON(w, http.StatusNotFound, JSONResponse{Message: "No se pueden obtener los registros"})
		return
	}
	mes := convertToSpanish(registerFilter.StartDate.Month().String())
	t := mes + " " + strconv.Itoa(registerFilter.StartDate.Day()) + " de " + strconv.Itoa(registerFilter.StartDate.Year())
	base64 := createChartsByDay(registers, t, "Hora", 120)
	respondJSON(w, http.StatusOK, JSONResponse{Payload: base64, Message: "Ok, registros encontrados"})
}

//GetRegistersDisponibilityWeek get all the registers from the view by week
func GetRegistersDisponibilityWeek(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	registerFilter := models.RegisterFilter{}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&registerFilter); err != nil {
		respondJSON(w, http.StatusBadRequest, JSONResponse{Message: "Error interno del servidor"})
		return
	}
	registers := []models.PretratamientoDispWeek{}

	if err := db.Group("for_date, week_day, tag_index, total").Order("for_date asc").Having("tag_index = ? AND for_date BETWEEN ? AND ?", registerFilter.TagIndex, registerFilter.StartDate, registerFilter.FinalDate).Find(&registers).Error; err != nil {
		respondJSON(w, http.StatusNotFound, JSONResponse{Message: "No se pueden obtener los registros"})
		return
	}
	mes := convertToSpanish(registerFilter.StartDate.Month().String())
	t := mes + " " + strconv.Itoa(registerFilter.StartDate.Day()) + " - " + strconv.Itoa(registerFilter.FinalDate.Day()) + " de " + strconv.Itoa(registerFilter.StartDate.Year())
	base64 := createChartsByWeek(registers, t, "Día", 2880)
	respondJSON(w, http.StatusOK, JSONResponse{Payload: base64, Message: "Ok, registros encontrados"})
}

//GetRegistersDisponibilityMonth get all the registers from the view by week
func GetRegistersDisponibilityMonth(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	registerFilter := models.RegisterFilter{}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&registerFilter); err != nil {
		respondJSON(w, http.StatusBadRequest, JSONResponse{Message: "Error interno del servidor"})
		return
	}
	registers := []models.PretratamientoDispMonth{}

	if err := db.Group("on_year, on_week, tag_index, total").Order("on_year, on_week asc").Having("tag_index = ? AND on_week >= ? AND on_week <= ? and on_year = ?", registerFilter.TagIndex, registerFilter.StartWeek, registerFilter.FinalWeek, registerFilter.Year).Find(&registers).Error; err != nil {
		respondJSON(w, http.StatusNotFound, JSONResponse{Message: "No se pueden obtener los registros"})
		return
	}
	t := "Semanas " + strconv.Itoa(registerFilter.StartWeek) + " - " + strconv.Itoa(registerFilter.FinalWeek) + " del año " + strconv.Itoa(registerFilter.Year)
	base64 := createChartsByMonth(registers, t, "Semana", 20160)
	respondJSON(w, http.StatusOK, JSONResponse{Payload: base64, Message: "Ok, registros encontrados"})
}

func createChartsByDay(lista []models.PretratamientoDisp, titulo string, tipo string, limite int) string {
	jsonBody := ItemList{Titulo: titulo, Tipo: tipo, Limite: limite}
	for i := 0; i < len(lista); i++ {
		s := strconv.Itoa(lista[i].OnHour)
		jsonBody.Values = append(jsonBody.Values, ItemReport{Value: lista[i].Total, Label: s})
	}
	jsonReq, err := json.Marshal(jsonBody)
	resp, err := http.Post("http://0.0.0.0:5000/crear_imagen", "application/json; charset=utf-8", bytes.NewBuffer(jsonReq))
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	// Convert response body to a struct
	var imageResponse ImageResponse
	json.Unmarshal(bodyBytes, &imageResponse)
	return imageResponse.File
}

func createChartsByWeek(lista []models.PretratamientoDispWeek, titulo string, tipo string, limite int) string {
	jsonBody := ItemList{Titulo: titulo, Tipo: tipo, Limite: limite}
	for i := 0; i < len(lista); i++ {
		jsonBody.Values = append(jsonBody.Values, ItemReport{Value: lista[i].Total, Label: lista[i].WeekDay})
	}
	jsonReq, err := json.Marshal(jsonBody)
	resp, err := http.Post("http://0.0.0.0:5000/crear_imagen", "application/json; charset=utf-8", bytes.NewBuffer(jsonReq))
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	// Convert response body to a struct
	var imageResponse ImageResponse
	json.Unmarshal(bodyBytes, &imageResponse)
	return imageResponse.File
}

func createChartsByMonth(lista []models.PretratamientoDispMonth, titulo string, tipo string, limite int) string {
	jsonBody := ItemList{Titulo: titulo, Tipo: tipo, Limite: limite}
	for i := 0; i < len(lista); i++ {
		s := strconv.Itoa(lista[i].OnWeek)
		jsonBody.Values = append(jsonBody.Values, ItemReport{Value: lista[i].Total, Label: s})
	}
	jsonReq, err := json.Marshal(jsonBody)
	resp, err := http.Post("http://0.0.0.0:5000/crear_imagen", "application/json; charset=utf-8", bytes.NewBuffer(jsonReq))
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	// Convert response body to a struct
	var imageResponse ImageResponse
	json.Unmarshal(bodyBytes, &imageResponse)
	return imageResponse.File
}

func convertToSpanish(month string) string {
	switch month {
	case "January":
		return "Enero"
	case "February":
		return "Febrero"
	case "March":
		return "Marzo"
	case "April":
		return "Abril"
	case "May":
		return "Mayo"
	case "June":
		return "Junio"
	case "July":
		return "Julio"
	case "August":
		return "Agosto"
	case "September":
		return "Septiembre"
	case "October":
		return "Octubre"
	case "November":
		return "Noviembre"
	case "December":
		return "Diciembre"
	default:
		return "Sin Mes"
	}
}
