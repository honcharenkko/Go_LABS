package main

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"text/template"
)

var CosPhi float64

func init() {
	CosPhi = math.Sqrt(3)
}

type Results struct {
	Izm    float64
	IzmMax float64
	Sek    float64
	Smin   float64
	Valid  bool
}

func calculate(Sm, Unom, Kz, Ft, Jek, Ct float64) Results {
	Izm := Sm / (CosPhi * Unom) * 1000.0 // Convert to Amps
	IzmMax := 2 * Izm
	Sek := Izm / Jek
	Smin := (Kz * 1000.0 * math.Sqrt(Ft)) / Ct
	return Results{
		Izm:    Izm,
		IzmMax: IzmMax,
		Sek:    Sek,
		Smin:   Smin,
		Valid:  Sek >= Smin,
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		Sm, _ := strconv.ParseFloat(r.FormValue("Sm"), 64)
		Unom, _ := strconv.ParseFloat(r.FormValue("Unom"), 64)
		Kz, _ := strconv.ParseFloat(r.FormValue("Kz"), 64)
		Ft, _ := strconv.ParseFloat(r.FormValue("Ft"), 64)
		Jek, _ := strconv.ParseFloat(r.FormValue("Jek"), 64)
		Ct, _ := strconv.ParseFloat(r.FormValue("Ct"), 64)

		results := calculate(Sm, Unom, Kz, Ft, Jek, Ct)

		tmpl := template.Must(template.New("result").Parse(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Power Calculation</title>
			<style>
				body { font-family: Arial, sans-serif; margin: 20px; text-align: center; }
				.container { max-width: 500px; margin: auto; padding: 20px; border: 1px solid #ccc; border-radius: 10px; box-shadow: 2px 2px 10px rgba(0,0,0,0.1); }
				h2 { color: #333; }
				.result { font-size: 18px; }
				.valid { color: green; }
				.invalid { color: red; }
			</style>
		</head>
		<body>
			<div class="container">
				<h2>Power Calculation Results</h2>
				<p class="result">Operational Current: <b>{{printf "%.2f" .Izm}} A</b></p>
				<p class="result">Max Starting Current: <b>{{printf "%.2f" .IzmMax}} A</b></p>
				<p class="result">Economic Cable Section: <b>{{printf "%.2f" .Sek}} mm²</b></p>
				<p class="result">Min Cable Section (Thermal Stability): <b>{{printf "%.2f" .Smin}} mm²</b></p>
				<p class="result {{if .Valid}}valid{{else}}invalid{{end}}">
					{{if .Valid}}Selected cable section meets requirements.{{else}}Cable section needs to be increased!{{end}}
				</p>
				<a href="/">Back</a>
			</div>
		</body>
		</html>
		`))
		tmpl.Execute(w, results)
		return
	}

	tmpl := template.Must(template.New("form").Parse(`
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Power Calculation</title>
		<style>
			body { font-family: Arial, sans-serif; margin: 20px; text-align: center; }
			.container { max-width: 500px; margin: auto; padding: 20px; border: 1px solid #ccc; border-radius: 10px; box-shadow: 2px 2px 10px rgba(0,0,0,0.1); }
		</style>
	</head>
	<body>
		<div class="container">
			<h2>Enter Data for Calculation</h2>
			<form method="post">
				<label>Sm (kVA): <input type="text" name="Sm" required></label><br>
				<label>Unom (kV): <input type="text" name="Unom" required></label><br>
				<label>Kz (kA): <input type="text" name="Kz" required></label><br>
				<label>Ft (sec): <input type="text" name="Ft" required></label><br>
				<label>Jek (A/mm²): <input type="text" name="Jek" required></label><br>
				<label>Ct (A*sqrt(sec)/mm²): <input type="text" name="Ct" required></label><br>
				<button type="submit">Calculate</button>
			</form>
		</div>
	</body>
	</html>
	`))
	tmpl.Execute(w, nil)
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
