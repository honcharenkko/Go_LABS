package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// HTML-шаблон для форми вводу та відображення результатів
var tmpl = `
<!DOCTYPE html>
<html>
<head>
    <title>Розрахунок втрат електроенергії</title>
</head>
<body>
    <h1>Розрахунок втрат електроенергії</h1>
    <form method="POST" action="/">
        <label for="Pwt">Навантаження (МВт):</label>
        <input type="text" id="Pwt" name="Pwt" required value="{{.Pwt}}"><br>

        <label for="Kp">Коефіцієнт використання:</label>
        <input type="text" id="Kp" name="Kp" required value="{{.Kp}}"><br>

        <label for="T">Час роботи в році (год):</label>
        <input type="text" id="T" name="T" required value="{{.T}}"><br>

        <label for="Wvt">Параметр W (Вт):</label>
        <input type="text" id="Wvt" name="Wvt" required value="{{.Wvt}}"><br>

        <input type="submit" value="Розрахувати">
    </form>

    {{if .Calculated}}
    <h2>Результати</h2>
    <p><strong>Втрати електроенергії автотрансформатора:</strong> {{.Mav}} кВт⋅год</p>
    <p><strong>Математичне очікування втрат електропередачі:</strong> {{.MWvt}} кВт⋅год</p>
    <p><strong>Загальні втрати електроенергії:</strong> {{.Mtotal}} кВт⋅год</p>
    {{end}}
</body>
</html>
`

// Функція розрахунку втрат електроенергії
func calculateLosses(Pwt, Kp float64, T int, Wvt float64) (float64, float64, float64) {
	Mav := Kp * Pwt * float64(T)                     // Втрати автотрансформатора
	MWvt := Kp * float64(T) * 4e-3 * 5.12e2 * Wvt    // Математичне очікування втрат
	Mtotal := (23.6 * Mav) + (17.6 * MWvt) - 2682000 // Загальні втрати
	return Mav, MWvt, Mtotal
}

// Обробник HTTP-запитів
func handler(w http.ResponseWriter, r *http.Request) {
	// Початкові значення
	data := struct {
		Pwt, Kp, Wvt      float64
		T                 int
		Mav, MWvt, Mtotal float64
		Calculated        bool
	}{
		Pwt: 23.6,
		Kp:  0.7,
		T:   51200,
		Wvt: 6451,
	}

	if r.Method == http.MethodPost {
		// Зчитуємо введені користувачем дані
		Pwt, _ := strconv.ParseFloat(r.FormValue("Pwt"), 64)
		Kp, _ := strconv.ParseFloat(r.FormValue("Kp"), 64)
		T, _ := strconv.Atoi(r.FormValue("T"))
		Wvt, _ := strconv.ParseFloat(r.FormValue("Wvt"), 64)

		// Обчислення
		data.Pwt, data.Kp, data.T, data.Wvt = Pwt, Kp, T, Wvt
		data.Mav, data.MWvt, data.Mtotal = calculateLosses(Pwt, Kp, T, Wvt)
		data.Calculated = true
	}

	// Відображення сторінки
	t, err := template.New("form").Parse(tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, data)
}

func main() {
	http.HandleFunc("/", handler)
	log.Println("Сервер запущено на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
