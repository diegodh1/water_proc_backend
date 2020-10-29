package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"
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

//GetRegistersDisponibility get all the registers from the view by day
func GetRegistersDisponibility(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	registerFilter := models.RegisterFilter{}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&registerFilter); err != nil {
		respondJSON(w, http.StatusBadRequest, JSONResponse{Message: "Error interno del servidor"})
		return
	}
	switch {
	case registerFilter.TypeReport == 0:
		fmt.Println("Dia")
		startDate := registerFilter.SelectedDate
		finalDate := startDate.Add(time.Hour*23 + time.Minute*59 + time.Second*59)
		base64, err := GetRegistersDisponibilityDay(db, startDate, finalDate, registerFilter.TagIndex, 120)
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, JSONResponse{Message: "No se pueden obtener los registros"})
			return
		}
		respondJSON(w, http.StatusOK, JSONResponse{Payload: base64, Message: "ok! registros obtenidos"})
		return
	case registerFilter.TypeReport == 1:
		fmt.Println("Semana")
		dayOfWeek := -1 * int(registerFilter.SelectedDate.Weekday())
		startDate := registerFilter.SelectedDate.AddDate(0, 0, dayOfWeek)
		lastDate := registerFilter.SelectedDate.AddDate(0, 0, 6+dayOfWeek)
		base64, err := GetRegistersDisponibilityWeek(db, startDate, lastDate, registerFilter.TagIndex, 2880)
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, JSONResponse{Message: "No se pueden obtener los registros"})
			return
		}
		respondJSON(w, http.StatusOK, JSONResponse{Payload: base64, Message: "ok! registros obtenidos"})
		return
	case registerFilter.TypeReport == 2:
		fmt.Println("Semanas")
		currentLocation := registerFilter.SelectedDate.Location()
		firstOfMonth := time.Date(registerFilter.SelectedDate.Year(), registerFilter.SelectedDate.Month(), 1, 0, 0, 0, 0, currentLocation)
		lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
		limite := 2 * 60 * 24 * 7
		base64, err := GetRegistersDisponibilityMonth(db, firstOfMonth, lastOfMonth, registerFilter.TagIndex, limite)
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, JSONResponse{Message: "No se pueden obtener los registros"})
			return
		}
		respondJSON(w, http.StatusOK, JSONResponse{Payload: base64, Message: "ok! registros obtenidos"})
		return
	default:
		respondJSON(w, http.StatusInternalServerError, JSONResponse{Message: "No se pueden obtener los registros"})
	}
}

//GetRegistersDisponibilityDay get all the registers from the view by day
func GetRegistersDisponibilityDay(db *gorm.DB, startDate time.Time, finalDate time.Time, tag int, limite int) (string, error) {

	registers := []models.PretratamientoDisp{}

	if err := db.Group("for_date, on_hour, tag_index, total").Order("for_date,on_hour asc").Having("tag_index = ? AND for_date BETWEEN ? AND ?", tag, startDate, finalDate).Find(&registers).Error; err != nil {
		return "", nil
	}
	t := "Disponibilidad del día " + strconv.Itoa(startDate.Year()) + "/" + strconv.Itoa(int(startDate.Month())) + "/" + strconv.Itoa(startDate.Day())
	base64 := createChartsByDay(registers, t, "Hora", limite)
	return base64, nil
}

//GetRegistersDisponibilityWeek get all the registers from the view by week
func GetRegistersDisponibilityWeek(db *gorm.DB, startDate time.Time, finalDate time.Time, tag int, limite int) (string, error) {

	registers := []models.PretratamientoDispWk{}
	fmt.Println("entro")
	if err := db.Where("(for_date BETWEEN ? AND ?) AND tag = ?", startDate, finalDate, tag).Find(&registers).Error; err != nil {
		return "", nil
	}
	sort.Slice(registers, func(i, j int) bool {
		return registers[i].ForDate.Before(registers[j].ForDate)
	})
	t := "Disponibilidad entre los días " + strconv.Itoa(startDate.Year()) + "/" + strconv.Itoa(int(startDate.Month())) + "/" + strconv.Itoa(startDate.Day()) + " - " + strconv.Itoa(finalDate.Year()) + "/" + strconv.Itoa(int(finalDate.Month())) + "/" + strconv.Itoa(finalDate.Day())
	base64 := createChartsByWeek(registers, t, "Día", limite)
	return base64, nil
}

//GetRegistersDisponibilityMonth get all the registers from the view by week
func GetRegistersDisponibilityMonth(db *gorm.DB, startDate time.Time, finalDate time.Time, tag int, limite int) (string, error) {

	registers := []models.PretratamientoDispMonth{}

	if err := db.Raw(`SELECT on_year,on_week,tag_index,COUNT(*) as total FROM (SELECT        CAST(DateAndTime AS date) AS for_date,DATEPART(YEAR, DateAndTime) AS on_year, DATEPART(WEEK, DateAndTime) AS on_week, TagIndex AS tag_index
	FROM            dbo.Pretratamiento) AS Temp
	WHERE Temp.for_date BETWEEN ? AND ?
	GROUP BY on_year,on_week,tag_index HAVING tag_index = ?
	ORDER BY on_year,on_week
	`, startDate, finalDate, tag).Scan(&registers).Error; err != nil {
		return "", err
	}
	t := "Disponibilidad de las semanas entre los días " + strconv.Itoa(startDate.Year()) + "/" + strconv.Itoa(int(startDate.Month())) + "/" + strconv.Itoa(startDate.Day()) + " - " + strconv.Itoa(finalDate.Year()) + "/" + strconv.Itoa(int(finalDate.Month())) + "/" + strconv.Itoa(finalDate.Day())
	base64 := createChartsByMonth(registers, t, "Semana", limite)
	return base64, nil
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

func createChartsByWeek(lista []models.PretratamientoDispWk, titulo string, tipo string, limite int) string {
	jsonBody := ItemList{Titulo: titulo, Tipo: tipo, Limite: limite}
	for i := 0; i < len(lista); i++ {
		s := strconv.Itoa(lista[i].ForDate.Year()) + "/" + strconv.Itoa(int(lista[i].ForDate.Month())) + "/" + strconv.Itoa(lista[i].ForDate.Day())
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
