package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
)

type InputData struct {
	Name      string  `json:"name"`
	Eta       float64 `json:"eta"`
	CosPhi    float64 `json:"cos_phi"`
	Voltage   float64 `json:"voltage"`
	Count     int     `json:"count"`
	Power     float64 `json:"power"`
	UtilCoeff float64 `json:"util_coeff"`
	TgPhi     float64 `json:"tg_phi"`
	Kv        float64 `json:"kv"`
}

type ResultData struct {
	Ip float64 `json:"Ip (Calculation Current in A)"`
	Kv float64 `json:"Kv (Group Usage Coefficient)"`
	Ne float64 `json:"Ne (Effective Number of EP)"`
	Kr float64 `json:"Kr (Calculated Power Coefficient)"`
	Pp float64 `json:"Pp (Calculated Active Load in kW)"`
	Qp float64 `json:"Qp (Calculated Reactive Load in kVAr)"`
	Sp float64 `json:"Sp (Total Power in kVA)"`
	Ig float64 `json:"Ig (Calculated Group Current in A)"`
}

func calculateLoad(data InputData) (ResultData, error) {
	if data.Voltage <= 0 || data.CosPhi <= 0 || data.Eta <= 0 || data.Power <= 0 || data.Count <= 0 {
		return ResultData{}, fmt.Errorf("invalid input values, all numeric values must be positive")
	}

	pTotal := float64(data.Count) * data.Power
	ip := (1000 * pTotal) / (math.Sqrt(3) * data.Voltage * data.CosPhi * data.Eta)
	kvGroup := data.Kv
	ne := math.Pow(float64(data.Count)*data.Power, 2) / (float64(data.Count) * math.Pow(data.Power, 2))
	kr := 1.25
	pp := kr * kvGroup * pTotal
	qp := kr * kvGroup * pTotal * data.TgPhi
	sp := math.Sqrt(pp*pp + qp*qp)
	ig := (1000 * pp) / (math.Sqrt(3) * data.Voltage)

	return ResultData{
		Ip: ip,
		Kv: kvGroup,
		Ne: ne,
		Kr: kr,
		Pp: pp,
		Qp: qp,
		Sp: sp,
		Ig: ig,
	}, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is supported", http.StatusMethodNotAllowed)
		return
	}

	var input InputData
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	result, err := calculateLoad(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("index").Parse(`
	<html>
	<head><title>Electrical Load Calculator</title></head>
	<body>
		<h2>Enter Input Data</h2>
		<form id="calcForm">
			<label>Count: <input type="number" name="count" required></label><br>
			<label>Power (kW): <input type="number" step="any" name="power" required></label><br>
			<label>Voltage (V): <input type="number" step="any" name="voltage" required></label><br>
			<label>CosPhi: <input type="number" step="any" name="cos_phi" required></label><br>
			<label>Eta: <input type="number" step="any" name="eta" required></label><br>
			<label>Utilization Coeff.: <input type="number" step="any" name="util_coeff" required></label><br>
			<label>TgPhi: <input type="number" step="any" name="tg_phi" required></label><br>
			<label>Kv: <input type="number" step="any" name="kv" required></label><br>
			<button type="submit">Calculate</button>
		</form>
		<h3>Result:</h3>
		<table border="1" id="resultTable" style="display:none;">
			<tr><th>Parameter</th><th>Value</th></tr>
		</table>
		<script>
		document.getElementById("calcForm").addEventListener("submit", function(event) {
			event.preventDefault();
			const formData = new FormData(event.target);
			const data = {};
			formData.forEach((value, key) => { data[key] = parseFloat(value) || 0; });
			fetch("/calculate", {
				method: "POST",
				headers: { "Content-Type": "application/json" },
				body: JSON.stringify(data)
			})
			.then(response => response.json())
			.then(result => {
				const table = document.getElementById("resultTable");
				table.innerHTML = "<tr><th>Parameter</th><th>Value</th></tr>";
				Object.keys(result).forEach(key => {
					table.innerHTML += "<tr><td>" + key + "</td><td>" + parseFloat(result[key]).toFixed(2) + "</td></tr>";
				});
				table.style.display = "block";
			})
			.catch(error => console.error("Error:", error));
		});
		</script>
	</body>
	</html>`)

	if err != nil {
		http.Error(w, "Template parsing error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Template execution error", http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/calculate", handler)
	log.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
