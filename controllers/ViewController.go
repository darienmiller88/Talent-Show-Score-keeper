package controllers

import (
	"fmt"
	"net/http"
	"os"
	"text/template"

	"github.com/go-chi/chi/v5"
)

// Hardcoded users
var allowedUsers = map[string]bool{
	"jennamandelricci": true,
	"asisatmuldoon":    true,
	"midrenelamy":      true,
	"adonisbrown":      true,
	"lindalaul":        true,
}

var actualNames = map[string]string{
	"jennamandelricci": "Jenna Mandel-ricci",
	"asisatmuldoon":    "Asisat Muldoon",
	"midrenelamy":      "Midrene Lamy",
	"adonisbrown":      "Adonis Brown",
	"lindalaul":        "Linda Laul",
}

// judge -> talent -> score
var scores = map[string]map[string]float64{}

var talentNames = []string{
	"Kiefer Inson",
	"CAYENNE NO_LUCK aka Justin Jacob",
	"Grand Concourse TOP",
	"Money",
	"Rachel Fonseca & Sophie Thurschwell",
	"Woody Tanor",
	
}

type TalentCard struct {
	ID    int
	Name  string
	Score float64
}

type TalentTotal struct {
	Index int
	Name  string
	Score float64
}

type ViewController struct{
	templates        map[string]*template.Template
	// logInTemplate    *template.Template
	Router           *chi.Mux
}

func (v *ViewController) Init() {
	v.templates = make(map[string]*template.Template)
	v.templates["home"] = template.Must(template.ParseFiles("views/home.html"))
	// v.logInTemplate = template.Must(template.ParseFiles("views/login.html"))	
	v.Router = chi.NewRouter()
	v.Router.Get("/", v.Home)
	v.Router.Get("/log-in", v.LogIn)
	v.Router.Get("/sign-out", v.SignOut)
	v.Router.Post("/log-in", v.HandleLogIn)
	// v.Router.Get("/talent-show", v.TalentShow)
	// v.Router.Post("/talent-show/score", v.UpdateTalentScore)
	// v.Router.Get("/point-totals", v.PointTotals)
	// v.Router.Get("/login", v.LogIn)
}

func getUser(r *http.Request) string {
	cookie, err := r.Cookie("auth")
	if err != nil {
		return ""
	}
	return cookie.Value
}

func (v *ViewController) Home(res http.ResponseWriter, req *http.Request) {
	if err := v.templates["home"].Execute(res, nil); err != nil {
		http.Error(res, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}	

func (v *ViewController) TalentShow(res http.ResponseWriter, req *http.Request) {
	user := getUser(req)
	actualName := actualNames[user]

	if scores[user] == nil {
		scores[user] = make(map[string]float64)
	}

	cards := make([]TalentCard, len(talentNames))
	for i, name := range talentNames {
		cards[i] = TalentCard{
			ID:    i,
			Name:  name,
			Score: getScore(user, name),
		}
	}

	data := map[string]any{
		"Cards": cards,
		"user":  actualName,
	}

	if err := v.templates["TalentShow"].Execute(res, data); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func (v *ViewController) UpdateTalentScore(w http.ResponseWriter, r *http.Request) {
	user := getUser(r)
	name := r.FormValue("name")
	action := r.FormValue("action")

	if user == "" || name == "" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if scores[user] == nil {
		scores[user] = make(map[string]float64)
	}

	current := getScore(user, name)
	switch action {
	case "plus":
		current += 0.5
	case "minus":
		current -= 0.5
		if current < 0 {
			current = 0
		}
	}
	scores[user][name] = current

	id := 0
	for i, n := range talentNames {
		if n == name {
			id = i
			break
		}
	}

	card := TalentCard{ID: id, Name: name, Score: current}

	tmpl := template.Must(template.ParseFiles("templates/partials/talentcard.html"))
	tmpl.ExecuteTemplate(w, "talentcard", card)
}

func (v *ViewController) PointTotals(res http.ResponseWriter, req *http.Request) {
	totals := make([]TalentTotal, len(talentNames))
	for i, name := range talentNames {
		var sum float64
		for judge := range allowedUsers {
			sum += getScore(judge, name)
		}
		totals[i] = TalentTotal{Index: i + 1, Name: name, Score: sum}
	}

	sort.Slice(totals, func(i, j int) bool {
		return totals[i].Score > totals[j].Score
	})

	data := map[string]any{
		"Totals": totals,
	}

	if err := v.templates["PointTotals"].Execute(res, data); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func (v *ViewController) LogIn(res http.ResponseWriter, req *http.Request) {
	if err := v.logInTemplate.Execute(res, nil); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func (v *ViewController) HandleLogIn(res http.ResponseWriter, req *http.Request) {
	username := req.FormValue("username")
	password := req.FormValue("password")

	if !allowedUsers[username] {
		fmt.Println("Invalid username")
		http.Redirect(res, req, "/log-in", http.StatusSeeOther)
		return
	}

	fmt.Println("username:", username, "password:", password)
	correctPassword := os.Getenv("PASSWORD")

	fmt.Println("value:", correctPassword == password)
	if password != correctPassword {
		fmt.Println("Invalid password")
		http.Redirect(res, req, "/log-in", http.StatusSeeOther)
		return
	}

	http.SetCookie(res, &http.Cookie{
		Name:     "auth",
		Value:    username,
		Path:     "/",
		HttpOnly: true,
	})

	http.Redirect(res, req, "/talent-show", http.StatusSeeOther)
}

func (v *ViewController) SignOut(res http.ResponseWriter, req *http.Request) {
	http.SetCookie(res, &http.Cookie{
		Name:     "auth",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	http.Redirect(res, req, "/log-in", http.StatusSeeOther)
}