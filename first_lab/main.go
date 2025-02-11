package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

type FuelData struct {
	HP, CP, SP, NP, OP, WP, AP                                float64
	KRS, KRG, HC, CC, SC, NC, OC, AC, HG, CG, SG, NG, OG, QrH float64
	ShowResults                                               bool
}

func calculateFuelData(fd *FuelData) {
	fd.KRS = 100 / (100 - fd.WP)
	fd.KRG = 100 / (100 - fd.WP - fd.AP)

	fd.HC = fd.HP * fd.KRS
	fd.CC = fd.CP * fd.KRS
	fd.SC = fd.SP * fd.KRS
	fd.NC = fd.NP * fd.KRS
	fd.OC = fd.OP * fd.KRS
	fd.AC = fd.AP * fd.KRS

	fd.HG = fd.HP * fd.KRG
	fd.CG = fd.CP * fd.KRG
	fd.SG = fd.SP * fd.KRG
	fd.NG = fd.NP * fd.KRG
	fd.OG = fd.OP * fd.KRG

	fd.QrH = 339*fd.CP + 1030*fd.HP - 108.8*(fd.OP-fd.SP) - 25*fd.WP
	fd.ShowResults = true
}

func fuelHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}

	if r.FormValue("clear") != "" {
		tmpl.Execute(w, nil)
		return
	}

	fd := FuelData{}
	fd.HP, _ = strconv.ParseFloat(r.FormValue("HP"), 64)
	fd.CP, _ = strconv.ParseFloat(r.FormValue("CP"), 64)
	fd.SP, _ = strconv.ParseFloat(r.FormValue("SP"), 64)
	fd.NP, _ = strconv.ParseFloat(r.FormValue("NP"), 64)
	fd.OP, _ = strconv.ParseFloat(r.FormValue("OP"), 64)
	fd.WP, _ = strconv.ParseFloat(r.FormValue("WP"), 64)
	fd.AP, _ = strconv.ParseFloat(r.FormValue("AP"), 64)

	calculateFuelData(&fd)
	tmpl.Execute(w, fd)
}

var tmpl = template.Must(template.New("fuel").Parse(`
<!DOCTYPE html>
<html>
<head>
	<title>Fuel Calculator</title>
	<style>
		body { font-family: Arial, sans-serif; background-color: #f4f4f4; text-align: center; padding: 20px; }
		form { background: white; padding: 20px; max-width: 300px; margin: auto; border-radius: 5px; box-shadow: 0px 0px 10px rgba(0,0,0,0.1); }
		input { width: 100%; padding: 8px; margin: 5px 0; border: 1px solid #ccc; border-radius: 4px; }
		input[type="submit"] { background-color: #28a745; color: white; border: none; padding: 10px; cursor: pointer; }
		input[type="submit"]:hover { background-color: #218838; }
		button { background-color: #dc3545; color: white; border: none; padding: 10px; cursor: pointer; width: 100%; margin-top: 10px; }
		button:hover { background-color: #c82333; }
		.results { background: white; padding: 20px; max-width: 400px; margin: 20px auto; border-radius: 5px; box-shadow: 0px 0px 10px rgba(0,0,0,0.1); }
	</style>
</head>
<body>
	<h1>Fuel Composition Calculator</h1>
	<form method="post">
		<label>HP: <input type="text" name="HP"></label><br>
		<label>CP: <input type="text" name="CP"></label><br>
		<label>SP: <input type="text" name="SP"></label><br>
		<label>NP: <input type="text" name="NP"></label><br>
		<label>OP: <input type="text" name="OP"></label><br>
		<label>WP: <input type="text" name="WP"></label><br>
		<label>AP: <input type="text" name="AP"></label><br>
		<input type="submit" value="Calculate">
		<button type="submit" name="clear" value="true">Clear Results</button>
	</form>
	{{if .ShowResults}}
		<div class="results">
			<h2>Results:</h2>
			<p>KRS: {{.KRS}}</p>
			<p>KRG: {{.KRG}}</p>
			<p>HC: {{.HC}}, CC: {{.CC}}, SC: {{.SC}}, NC: {{.NC}}, OC: {{.OC}}, AC: {{.AC}}</p>
			<p>HG: {{.HG}}, CG: {{.CG}}, SG: {{.SG}}, NG: {{.NG}}, OG: {{.OG}}</p>
			<p>Lower Heating Value (QrH): {{.QrH}} MJ/kg</p>
		</div>
	{{end}}
</body>
</html>
`))

func main() {
	http.HandleFunc("/", fuelHandler)
	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}
