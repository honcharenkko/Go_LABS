package main

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"text/template"
)

type Results struct {
	XSum   float64
	Ik0    float64
	XcPU   float64
	XtPU   float64
	XSumPU float64
	Ik0PU  float64
}

func calculate(Unom, Sk, Xc, Xt, Sb float64) Results {
	XSum := Xc + Xt
	Ik0 := (Unom * 1000) / (math.Sqrt(3) * XSum)
	XcPU := Xc * (Sb / Sk)
	XtPU := Xt * (Sb / Sk)
	XSumPU := XcPU + XtPU
	Ib := Sb / (math.Sqrt(3) * Unom)
	Ik0PU := Ik0 / Ib
	return Results{
		XSum:   XSum,
		Ik0:    Ik0,
		XcPU:   XcPU,
		XtPU:   XtPU,
		XSumPU: XSumPU,
		Ik0PU:  Ik0PU,
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		Unom, _ := strconv.ParseFloat(r.FormValue("Unom"), 64)
		Sk, _ := strconv.ParseFloat(r.FormValue("Sk"), 64)
		Xc, _ := strconv.ParseFloat(r.FormValue("Xc"), 64)
		Xt, _ := strconv.ParseFloat(r.FormValue("Xt"), 64)
		Sb, _ := strconv.ParseFloat(r.FormValue("Sb"), 64)

		results := calculate(Unom, Sk, Xc, Xt, Sb)

		tmpl := template.Must(template.New("result").Parse(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Short Circuit Calculation</title>
			<style>
				body { font-family: Arial, sans-serif; margin: 20px; text-align: center; }
				.container { max-width: 500px; margin: auto; padding: 20px; border: 1px solid #ccc; border-radius: 10px; box-shadow: 2px 2px 10px rgba(0,0,0,0.1); }
			</style>
		</head>
		<body>
			<div class="container">
				<h2>Short Circuit Calculation Results</h2>
				<p>Σ Impedance at K1: <b>{{printf "%.2f" .XSum}} Ω</b></p>
				<p>Initial Short-Circuit Current: <b>{{printf "%.2f" .Ik0}} kA</b></p>
				<p>Xc in PU: <b>{{printf "%.2f" .XcPU}}</b></p>
				<p>Xt in PU: <b>{{printf "%.2f" .XtPU}}</b></p>
				<p>Σ Impedance in PU: <b>{{printf "%.2f" .XSumPU}}</b></p>
				<p>Initial Short-Circuit Current in PU: <b>{{printf "%.2f" .Ik0PU}}</b></p>
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
		<title>Short Circuit Calculation</title>
		<style>
			body { font-family: Arial, sans-serif; margin: 20px; text-align: center; }
			.container { max-width: 500px; margin: auto; padding: 20px; border: 1px solid #ccc; border-radius: 10px; box-shadow: 2px 2px 10px rgba(0,0,0,0.1); }
		</style>
	</head>
	<body>
		<div class="container">
			<h2>Enter Data for Calculation</h2>
			<form method="post">
				<label>Unom (kV): <input type="text" name="Unom" required></label><br>
				<label>Sk (MVA): <input type="text" name="Sk" required></label><br>
				<label>Xc (Ohm): <input type="text" name="Xc" required></label><br>
				<label>Xt (Ohm): <input type="text" name="Xt" required></label><br>
				<label>Sb (MVA): <input type="text" name="Sb" required></label><br>
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
