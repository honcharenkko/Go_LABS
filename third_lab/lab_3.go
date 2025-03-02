package main

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
)

func main() {
	http.HandleFunc("/", homeHandler)
	http.ListenAndServe(":8080", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		fmt.Fprintf(w, `
		<html>
		<head>
		<title>Розрахунок прибутку</title>
		<style>
			body { font-family: Arial, sans-serif; text-align: center; background-color: #f4f4f4; padding: 50px; }
			.container { background: white; padding: 20px; border-radius: 10px; box-shadow: 0px 0px 10px rgba(0, 0, 0, 0.1); display: inline-block; }
			input, button { margin: 10px; padding: 10px; font-size: 16px; }
			button { background-color: #ff9800; color: white; border: none; cursor: pointer; }
			button:hover { background-color: #e68900; }
		</style>
		</head>
		<body>
		<div class="container">
		<h2>Розрахунок прибутку від сонячних електростанцій</h2>
		<form action="/calculate" method="post">
			<label>Середньодобова потужність (Pc) у МВт:</label>
			<input type="number" step="any" name="pc" required><br>
			<label>Похибка прогнозу (%):</label>
			<input type="number" step="any" name="delta" required><br>
			<button type="submit">Розрахувати</button>
		</form>
		</div>
		</body>
		</html>
		`)
	} else if r.Method == http.MethodPost {
		r.ParseForm()
		pc, _ := strconv.ParseFloat(r.FormValue("pc"), 64)
		delta, _ := strconv.ParseFloat(r.FormValue("delta"), 64)
		result := calculateProfitOrPenalty(pc, delta)
		fmt.Fprintf(w, "<pre>%s</pre>", result)
	}
}

func calculateProfitOrPenalty(pc, delta float64) string {
	B := 7.0
	sigma := 1.0
	newSigma := 0.25
	lowerBound := pc - 0.25
	upperBound := pc + 0.25

	energyPercentage := (normalDistributionCDF(upperBound, pc, sigma) - normalDistributionCDF(lowerBound, pc, sigma)) * 100
	profit1 := (pc * 24) * B * (energyPercentage / 100)
	penalty1 := (pc * 24) * B * ((100 - energyPercentage) / 100)

	newEnergyPercentage := (normalDistributionCDF(upperBound, pc, newSigma) - normalDistributionCDF(lowerBound, pc, newSigma)) * 100
	profit2 := (pc * 24) * B * (newEnergyPercentage / 100)
	penalty2 := (pc * 24) * B * ((100 - newEnergyPercentage) / 100)

	res := profit2 - penalty2

	return fmt.Sprintf(`
Середньодобова потужність: %.2f МВт
Похибка прогнозу: %.2f %%

Відсоток енергії: %.2f %%
Прибуток: %.2f тис. грн
Штраф: %.2f тис. грн

Відсоток енергії (новий σ): %.2f %%
Прибуток: %.2f тис. грн
Штраф: %.2f тис. грн

Можна отримати %.2f тис. грн прибутку!
`, pc, delta, energyPercentage, profit1, penalty1, newEnergyPercentage, profit2, penalty2, res)
}

func normalDistributionCDF(x, mean, stdDev float64) float64 {
	return 0.5 * (1 + erf((x-mean)/(stdDev*math.Sqrt2)))
}

func erf(x float64) float64 {
	sign := 1.0
	if x < 0 {
		sign = -1.0
	}
	x = math.Abs(x)
	t := 1.0 / (1.0 + 0.3275911*x)
	y := 1.0 - ((((1.061405429*t+-1.453152027)*t+1.421413741)*t+-0.284496736)*t+0.254829592)*t*math.Exp(-x*x)

	return sign * y
}
