package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

type Equipment struct {
	Name  string
	Omega float64 // Інтенсивність відмов (рік⁻¹)
	Tfail float64 // Час відновлення (год)
}

type InputData struct {
	Pv    float64
	Kp    float64
	T     float64
	Equip Equipment
}

type Results struct {
	Qo       float64
	Tavg     float64
	Ka       float64
	M_energy float64
}

var equipmentData = map[string]Equipment{
	"ПЛ-110 кВ": {"ПЛ-110 кВ", 0.007, 10},
	"ПЛ-35 кВ":  {"ПЛ-35 кВ", 0.02, 8},
	"Т-110 кВ":  {"Т-110 кВ", 0.015, 100},
	"Т-35 кВ":   {"Т-35 кВ", 0.02, 80},
}

func calculateReliability(data InputData) Results {
	Qo := data.Equip.Omega
	Tavg := data.Equip.Tfail
	Ka := (Qo * Tavg) / 8760
	M_energy := data.Kp * data.Pv * data.T

	return Results{Qo, Tavg, Ka, M_energy}
}

const htmlTemplate = `
<!DOCTYPE html>
<html lang="uk">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Розрахунок Надійності</title>
    <script src="https://unpkg.com/htmx.org@1.9.4"></script>
</head>
<body>
    <h2>Введіть початкові дані</h2>
    <form hx-post="/calculate" hx-target="#results" hx-swap="innerHTML">
        <label>Навантаження (МВт): <input type="text" name="Pv"></label><br>
        <label>Коефіцієнт використання: <input type="text" name="Kp"></label><br>
        <label>Час роботи (год): <input type="text" name="T"></label><br>
        <label>Тип обладнання:
            <select name="equipment">
                <option value="ПЛ-110 кВ">ПЛ-110 кВ</option>
                <option value="ПЛ-35 кВ">ПЛ-35 кВ</option>
                <option value="Т-110 кВ">Т-110 кВ</option>
                <option value="Т-35 кВ">Т-35 кВ</option>
            </select>
        </label><br>
        <button type="submit">Розрахувати</button>
    </form>
    <div id="results"></div>
</body>
</html>`

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("index").Parse(htmlTemplate))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, nil)
}

func calculateHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	var input InputData
	input.Pv, _ = strconv.ParseFloat(r.FormValue("Pv"), 64)
	input.Kp, _ = strconv.ParseFloat(r.FormValue("Kp"), 64)
	input.T, _ = strconv.ParseFloat(r.FormValue("T"), 64)
	equipType := r.FormValue("equipment")
	input.Equip = equipmentData[equipType]

	results := calculateReliability(input)
	tmpl := template.Must(template.New("results").Parse(`
	<h3>Результати:</h3>
	<p>Частота відмов: {{ .Qo }}</p>
	<p>Середня тривалість відмови: {{ .Tavg }}</p>
	<p>Коефіцієнт простою: {{ .Ka }}</p>
	<p>Втрати енергії: {{ .M_energy }} МВт·год</p>
	`))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, results)
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/calculate", calculateHandler)
	fmt.Println("Сервер запущено на порту 8080...")
	http.ListenAndServe(":8080", nil)
}
