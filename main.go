package main

import (
	"fmt"
	"html/template"
	"net/http"
	"regexp"
)

type Utilisateur struct {
	Nom           string
	Prenom        string
	DateNaissance string
	Sexe          string
	Erreur        string
}

type Student struct {
	Prenom string
	Nom    string
	Age    int
	Genre  string
}

type Class struct {
	ClasseName  string
	Filiere     string
	Niv         string
	NbrStudents int
	Students    []Student
}

var viewCount int
var utilisateur Utilisateur

func main() {
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))

	templates := template.Must(template.ParseGlob("templates/*.html"))

	http.HandleFunc("/promo", func(w http.ResponseWriter, r *http.Request) {

		students := []Student{
			{"Bastien", "Jouffre", 16, "M"},
			{"Lilian", "Lepiver", 21, "M"},
			{"Nans", "Moll", 20, "M"},
			{"Diane", "Lefevre", 21, "F"},
		}

		class := Class{
			ClasseName:  "B1 Informatique",
			Filiere:     "Informatique",
			Niv:         "Bachelor 1",
			NbrStudents: len(students),
			Students:    students,
		}

		err := templates.ExecuteTemplate(w, "Challenge1", class)
		if err != nil {
			http.Error(w, "Erreur lors de l'exécution du template", http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/change", func(w http.ResponseWriter, r *http.Request) {

		viewCount++

		message := "Le nombre de vues est "
		if viewCount%2 == 0 {
			message += "pair"
		} else {
			message += "impair"
		}
		message += fmt.Sprintf(" : %d", viewCount)

		err := templates.ExecuteTemplate(w, "Challenge2", message)
		if err != nil {
			http.Error(w, "Erreur lors de l'exécution du template", http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/user/form", func(w http.ResponseWriter, r *http.Request) {

		err := templates.ExecuteTemplate(w, "Challenge3", nil)
		if err != nil {
			http.Error(w, "Erreur lors de l'exécution du template", http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/user/treatment", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			nom := r.FormValue("nom")
			prenom := r.FormValue("prenom")
			dateNaissance := r.FormValue("dateNaissance")
			sexe := r.FormValue("sexe")

			valide, messageErreur := validerFormulaire(nom, prenom, sexe)
			if !valide {
				utilisateur.Erreur = messageErreur
				http.Redirect(w, r, "/user/error", http.StatusSeeOther)
				return
			}

			utilisateur = Utilisateur{
				Nom:           nom,
				Prenom:        prenom,
				DateNaissance: dateNaissance,
				Sexe:          sexe,
			}

			http.Redirect(w, r, "/user/display", http.StatusSeeOther)
		}
	})

	http.HandleFunc("/user/display", func(w http.ResponseWriter, r *http.Request) {
		err := templates.ExecuteTemplate(w, "display", utilisateur)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/user/error", func(w http.ResponseWriter, r *http.Request) {
		err := templates.ExecuteTemplate(w, "error", utilisateur.Erreur)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		fmt.Println("Erreur lors du démarrage du serveur :", err)
	}
}

func validerFormulaire(nom, prenom, sexe string) (bool, string) {

	nomPrenomRegex := regexp.MustCompile("^[a-zA-ZÀ-ÿ\\s'-]{1,32}$")

	if !nomPrenomRegex.MatchString(nom) {
		return false, "Nom invalide : uniquement des lettres, entre 1 et 32 caractères"
	}
	if !nomPrenomRegex.MatchString(prenom) {
		return false, "Prénom invalide : uniquement des lettres, entre 1 et 32 caractères"
	}

	if sexe != "masculin" && sexe != "féminin" && sexe != "autre" {
		return false, "Sexe invalide. Choisissez parmi 'masculin', 'féminin' ou 'autre'."
	}

	return true, ""
}
