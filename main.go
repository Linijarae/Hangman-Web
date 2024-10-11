package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"regexp"
)

var vues int

func main() {

	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./css/"))))

	type Etudiants struct {
		Nom    string
		Prenom string
		Age    int
		Sexe   bool
	}

	type Promo struct {
		Nomdeclasse string
		Filiere     string
		Niveau      string
		Nbetudiants int
		Etudiants   []Etudiants
	}

	type Change struct {
		Vues   int
		Modulo bool
	}

	type InfoPerso struct {
		Nom    string
		Prenom string
		Birth  string
		Sexe   string
	}

	type StockForm struct {
		CheckNom     bool
		Nom          string
		CheckPrenom  bool
		Prenom       string
		CheckBirth   bool
		Birth        string
		CheckSexe    bool
		Sexe         string
		ErrorMessage string
	}

	var stockForm = StockForm{false, "", false, "", false, "", false, "", ""}
	//gestion des erreurs
	temp, err := template.ParseGlob("templates/*.html")
	if err != nil {
		fmt.Printf("ERREUR => %s", err.Error())
		os.Exit(02)
	}
	//Excecution du template promo
	http.HandleFunc("/promo", func(w http.ResponseWriter, r *http.Request) {
		dataPage := Promo{"B1 Cybersécurité", "Cybersécurité", "B1", 5, []Etudiants{{"Al", "Capone", 20, true}, {"Ali", "Baba", 18, true}, {"Jude", "Holy", 25, false}, {"Marie", "Joseph", 32, false}, {"Tani", "Turnher", 21, false}}}
		temp.ExecuteTemplate(w, "promo", dataPage)
	})

	//Excecution du template change
	http.HandleFunc("/change", func(w http.ResponseWriter, r *http.Request) {
		// Incrémentation de la variable vue ( nombres de vue de la page)
		vues++
		//Pour savoir si le nombre de vue est pair ou impair j'utilise un modulo.
		modulo := vues%2 == 0
		data := Change{
			Vues:   vues,
			Modulo: modulo,
		}
		temp.ExecuteTemplate(w, "change", data)
	})
	//Excecution du formulaire
	http.HandleFunc("/form", func(w http.ResponseWriter, r *http.Request) {
		temp.ExecuteTemplate(w, "form", stockForm)
	})
	//Traitement des informations entrées par l'utilisateur dans le formulaire
	http.HandleFunc("/form/treatment", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/erreur?code=400&message=La méthode utilisée est incorrecte", http.StatusSeeOther)
			return
		}
		//Vérification des informations pour éviter les erreurs
		CheckNom, _ := regexp.MatchString("^[a-zA-Z-]{1,64}$", r.FormValue("name"))
		CheckPrenom, _ := regexp.MatchString("^[a-zA-Z-]{1,64}$", r.FormValue("prenom"))
		Birth := r.FormValue("birth")
		CheckBirth := len(Birth) > 0
		Sexe := r.FormValue("sexe")
		CheckSexe := Sexe == "Homme" || Sexe == "Femme"

		// Initialiser stockForm avec les valeurs saisies
		stockForm = StockForm{
			CheckNom:    CheckNom,
			Nom:         r.FormValue("name"),
			CheckPrenom: CheckPrenom,
			Prenom:      r.FormValue("prenom"),
			CheckBirth:  CheckBirth,
			Birth:       Birth,
			CheckSexe:   CheckSexe,
			Sexe:        Sexe,
		}
		//Erreur dans le prénom et le nom
		if !CheckPrenom && !CheckNom {
			stockForm.ErrorMessage = "Nom et Prénom incorrect !"
			temp.ExecuteTemplate(w, "form", stockForm)
			return
		}
		//Erreur dans le nom
		if !CheckNom {
			stockForm.ErrorMessage = "Nom incorrect !"
			temp.ExecuteTemplate(w, "form", stockForm)
			return
		}
		//Erreur dans le prenom
		if !CheckPrenom {
			stockForm.ErrorMessage = "Prénom incorrect !"
			temp.ExecuteTemplate(w, "form", stockForm)
			return
		}
		//Si tout est bon dans le formulaire, redirige l'utilisateur vers display qui affiche le résultat
		http.Redirect(w, r, "/form/display", http.StatusSeeOther)
	})

	type Affichage struct {
		CheckNom    bool
		Nom         string
		CheckPrenom bool
		Prenom      string
		CheckBirth  bool
		Birth       string
		CheckSexe   bool
		Sexe        string
	}

	http.HandleFunc("/form/display", func(w http.ResponseWriter, r *http.Request) {
		dataform := Affichage{
			CheckNom:    stockForm.CheckNom,
			Nom:         stockForm.Nom,
			CheckPrenom: stockForm.CheckPrenom,
			Prenom:      stockForm.Prenom,
			CheckBirth:  stockForm.CheckBirth,
			Birth:       stockForm.Birth,
			CheckSexe:   stockForm.CheckSexe,
			Sexe:        stockForm.Sexe,
		}
		temp.ExecuteTemplate(w, "display", dataform)
	})
	//Démarrage du serveur
	http.ListenAndServe("localhost:8080", nil)
}
