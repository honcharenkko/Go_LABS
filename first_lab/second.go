package main

import (
	"html/template"
	"net/http"
	"strconv"
)

type InputData struct {
	H, C, S, Q, O, W, A, V float64
	Results                *ResultData
}

type ResultData struct {
	CValc, HValc, OValc, SValc, AValc, Q_wm, VValc float64
}

func handler(w http.ResponseWriter, r *http.Request) {
	data := InputData{}

	if r.Method == http.MethodPost {
		r.ParseForm()
		data.H, _ = strconv.ParseFloat(r.FormValue("H"), 64)
		data.C, _ = strconv.ParseFloat(r.FormValue("C"), 64)
		data.S, _ = strconv.ParseFloat(r.FormValue("S"), 64)
		data.Q, _ = strconv.ParseFloat(r.FormValue("Q"), 64)
		data.O, _ = strconv.ParseFloat(r.FormValue("O"), 64)
		data.W, _ = strconv.ParseFloat(r.FormValue("W"), 64)
		data.A, _ = strconv.ParseFloat(r.FormValue("A"), 64)
		data.V, _ = strconv.ParseFloat(r.FormValue("V"), 64)

		if data.H >= 0 && data.C >= 0 && data.S >= 0 && data.Q >= 0 && data.O >= 0 && data.W >= 0 && data.A >= 0 && data.V >= 0 {
			data.Results = &ResultData{
				Q_wm:  (data.Q * (100 - data.W - data.A)) / 100,
				HValc: data.H * (100 - data.W - data.A) / 100,
				CValc: data.C * (100 - data.W - data.A) / 100,
				SValc: data.S * (100 - data.W - data.A) / 100,
				OValc: data.O * (100 - data.W - data.A) / 100,
				AValc: data.A * (100 - data.W) / 100,
				VValc: data.V * (100 - data.W) / 100,
			}
		}
	}

	tmpl, _ := template.New("index").Parse(htmlTemplate)
	tmpl.Execute(w, data)
}

const htmlTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>Fuel Composition Calculator</title>
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
        <label>H: <input type="text" name="H"></label><br>
        <label>C: <input type="text" name="C"></label><br>
        <label>S: <input type="text" name="S"></label><br>
        <label>Q: <input type="text" name="Q"></label><br>
        <label>O: <input type="text" name="O"></label><br>
        <label>W: <input type="text" name="W"></label><br>
        <label>A: <input type="text" name="A"></label><br>
        <label>V: <input type="text" name="V"></label><br>
        <input type="submit" value="Calculate">
        <button type="submit" name="clear" value="true">Clear Results</button>
    </form>
    
    {{if .Results}}
    <div class="results">
        <h2>Results:</h2>
        <p>Вуглець: {{.Results.CValc}}</p>
        <p>Водень: {{.Results.HValc}}</p>
        <p>Кисень: {{.Results.OValc}}</p>
        <p>Сірка: {{.Results.SValc}}</p>
        <p>Зола: {{.Results.AValc}}</p>
        <p>Нижча теплота згоряння: {{.Results.Q_wm}} МДж/кг</p>
        <p>Вміст ванадію: {{.Results.VValc}} мг/кг</p>
    </div>
    {{end}}
</body>
</html>`

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
