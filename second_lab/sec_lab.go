package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

// Структура для вхідних даних
type InputData struct {
	FuelType string
	FuelMass float64
}

// Структура для результатів
type Result struct {
	EmissionFactor float64
	TotalEmission  float64
}

// Дані про вугілля та мазут
var fuelData = map[string]struct {
	A    float64 // Зольність (%)
	W    float64 // Вологість (%)
	Q    float64 // Нижча теплота згоряння (МДж/кг)
	Type string  // Тип палива (вугілля чи мазут)
}{
	// Вугілля
	"Антрацитовий штиб АШ":        {A: 5.0, W: 3.0, Q: 33.24, Type: "coal"},
	"Пісне вугілля ТР":            {A: 12.0, W: 6.0, Q: 34.29, Type: "coal"},
	"Донецьке газове ГР":          {A: 25.20, W: 10.0, Q: 31.98, Type: "coal"},
	"Донецьке довгополуменеве ДР": {A: 35.0, W: 15.0, Q: 30.56, Type: "coal"},
	"Львівсько-волинське (ЛВ) ГР": {A: 18.0, W: 10.0, Q: 31.69, Type: "coal"},
	"Олександрійське буре БІР":    {A: 45.0, W: 25.0, Q: 26.96, Type: "coal"},
	// Мазут
	"Високосірчастий 40":  {A: 0.15, W: 2.00, Q: 40.40, Type: "mazut"},
	"Високосірчастий 100": {A: 0.15, W: 2.00, Q: 40.03, Type: "mazut"},
	"Високосірчастий 200": {A: 0.30, W: 1.00, Q: 39.77, Type: "mazut"},
	"Малосірчастий 40":    {A: 0.15, W: 2.00, Q: 41.24, Type: "mazut"},
	"Малосірчастий 100":   {A: 0.15, W: 2.00, Q: 40.82, Type: "mazut"},
}

// Дані про газ
var gasData = map[string]struct {
	Q  float64 // Нижча теплота згоряння (МДж/нм³)
	Ro float64 // Щільність (кг/нм³)
}{
	"Уренгой—Ужгород":    {Q: 33.08, Ro: 0.723},
	"Середня Азія—Центр": {Q: 34.21, Ro: 0.764},
}

// Функція для розрахунку викидів
func calculateEmission(fuelType string, fuelMass float64) Result {
	fuel, exists := fuelData[fuelType]
	if !exists {
		return Result{}
	}

	var k, E float64
	if fuel.Type == "coal" {
		qr := fuel.Q * (1 - ((fuel.W + fuel.A) / 100))
		first := (1000000 / qr) * 0.8
		second := (fuel.A / (100 - 1.5)) * (1 - 0.985)
		k = first * second
		E = 0.000001 * k * qr * fuelMass
	} else if fuel.Type == "mazut" {
		first := (1000000 / fuel.Q) * 1
		second := (fuel.A / 100) * (1 - 0.985)
		k = first * second
		E = 0.000001 * k * fuel.Q * fuelMass
	}

	return Result{
		EmissionFactor: k,
		TotalEmission:  E,
	}
}

// Функція для розрахунку викидів газу
func calculateGasEmission(gasType string, gasVolume float64) Result {
	gas, exists := gasData[gasType]
	if !exists {
		return Result{}
	}

	// Формула для розрахунку викидів
	first := (1000000 / gas.Q) * 0.8
	E := 0.000001 * first * gas.Q * gasVolume

	return Result{
		EmissionFactor: first,
		TotalEmission:  E,
	}
}

// Обробник форми
func formHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}

	fuelType := r.FormValue("fuelType")
	fuelMass, _ := strconv.ParseFloat(r.FormValue("fuelMass"), 64)

	var result Result
	if _, exists := fuelData[fuelType]; exists {
		result = calculateEmission(fuelType, fuelMass)
	} else if _, exists := gasData[fuelType]; exists {
		result = calculateGasEmission(fuelType, fuelMass)
	} else {
		// Якщо тип палива невідомий, повертаємо порожній результат
		result = Result{}
	}

	tmpl.Execute(w, struct {
		Input  InputData
		Result Result
	}{Input: InputData{FuelType: fuelType, FuelMass: fuelMass}, Result: result})
}

// Шаблон HTML
var tmpl = template.Must(template.New("form").Parse(`
<!DOCTYPE html>
<html>
<head>
	<title>Калькулятор викидів</title>
	<style>
		body { font-family: Arial, sans-serif; background-color: #f4f4f4; text-align: center; padding: 20px; }
		.container { background: white; padding: 20px; border-radius: 8px; box-shadow: 0px 0px 10px rgba(0, 0, 0, 0.1); width: 50%; margin: auto; }
		input, select { padding: 10px; margin: 10px; width: 80%; border-radius: 5px; border: 1px solid #ccc; }
		input[type="submit"] { background-color: #28a745; color: white; border: none; cursor: pointer; }
		input[type="submit"]:hover { background-color: #218838; }
	</style>
</head>
<body>
	<div class="container">
		<h1>Калькулятор викидів</h1>
		<form method="POST">
			<label>Тип палива: </label>
			<select name="fuelType">
				<optgroup label="Вугілля">
					<option value="Антрацитовий штиб АШ">Антрацитовий штиб АШ</option>
					<option value="Пісне вугілля ТР">Пісне вугілля ТР</option>
					<option value="Донецьке газове ГР">Донецьке газове ГР</option>
					<option value="Донецьке довгополуменеве ДР">Донецьке довгополуменеве ДР</option>
					<option value="Львівсько-волинське (ЛВ) ГР">Львівсько-волинське (ЛВ) ГР</option>
					<option value="Олександрійське буре БІР">Олександрійське буре БІР</option>
				</optgroup>
				<optgroup label="Мазут">
					<option value="Високосірчастий 40">Високосірчастий 40</option>
					<option value="Високосірчастий 100">Високосірчастий 100</option>
					<option value="Високосірчастий 200">Високосірчастий 200</option>
					<option value="Малосірчастий 40">Малосірчастий 40</option>
					<option value="Малосірчастий 100">Малосірчастий 100</option>
				</optgroup>
				<optgroup label="Газ">
					<option value="Уренгой—Ужгород">Уренгой—Ужгород</option>
					<option value="Середня Азія—Центр">Середня Азія—Центр</option>
				</optgroup>
			</select>
			<input type="text" name="fuelMass" placeholder="Кількість (т)"><br>
			<input type="submit" value="Розрахувати">
		</form>
		{{if .Result.TotalEmission}}
		<h2>Результати:</h2>
		<p>Викиди від {{.Input.FuelType}}: {{.Result.TotalEmission}} т</p>
		<p>Коефіцієнт викиду: {{.Result.EmissionFactor}}</p>
		{{end}}
	</div>
</body>
</html>
`))

func main() {
	http.HandleFunc("/", formHandler)
	fmt.Println("Сервер запущено на :8080")
	http.ListenAndServe(":8080", nil)
}
