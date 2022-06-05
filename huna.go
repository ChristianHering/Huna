package huna

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	goessentials "github.com/ChristianHering/GoEssentials"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

var dataDir string
var callback func(string) interface{}

var templates *template.Template

func init() {
	templates = template.Must(template.ParseGlob("../Huna/templates/*.html"))
}

func Huna(mux *mux.Router, d string, c func(string) interface{}) (t *template.Template) {
	dataDir = d
	callback = c

	authMiddleware := alice.New(authHandler)

	mux.Handle("/huna", authMiddleware.ThenFunc(indexHandler))
	mux.Handle("/huna/download", authMiddleware.ThenFunc(downloadHandler))
	mux.Handle("/huna/update", authMiddleware.ThenFunc(updateHandler))

	mux.PathPrefix("/huna/js/").Handler(http.StripPrefix("/huna/js/", http.FileServer(http.Dir("./../Huna/js"))))
	mux.PathPrefix("/huna/css/").Handler(http.StripPrefix("/huna/css/", http.FileServer(http.Dir("./../Huna/css"))))
	mux.PathPrefix("/huna/asm/").Handler(http.StripPrefix("/huna/asm/", http.FileServer(http.Dir("./../Huna/asm"))))

	return template.Must(template.ParseFiles("../Huna/templates/huna.html"))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}
}

func authHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isUserAuthenticated := callback("IsUserAuthenticated").(func(r *http.Request) (bool, string, error))

		authenticated, _, err := isUserAuthenticated(r)
		if err != nil {
			panic(err)
		}

		if !authenticated {
			http.Redirect(w, r, "/login", http.StatusSeeOther)

			return
		}

		next.ServeHTTP(w, r)
	})
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	isUserAuthenticated := callback("IsUserAuthenticated").(func(r *http.Request) (bool, string, error))

	authenticated, username, err := isUserAuthenticated(r)
	if err != nil {
		panic(err)
	}

	if !authenticated {
		panic("d")
	}

	err = os.MkdirAll(filepath.Join(dataDir, username, "huna"), 0700)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	if goessentials.FileNotExist(filepath.Join(dataDir, username, "huna", "database.kdbx")) {
		w.WriteHeader(http.StatusNotFound)

		return
	}

	b, err := ioutil.ReadFile(filepath.Join(dataDir, username, "huna", "database.kdbx"))
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(b)
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(256000000) //256MB max db size
	if err != nil {
		fmt.Println(err)
	}

	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, _, err := r.FormFile("dbFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer file.Close()

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	c, err := r.Cookie("username")
	if err != nil {
		fmt.Println(err)
	}

	username := c.Value

	kdbxFile, _ := os.Create(filepath.Join(dataDir, username, "huna", "database.kdbx"))

	// write this byte array to our temporary file
	_, err = kdbxFile.Write(fileBytes)
	if err != nil {
		fmt.Println(err)
	}
	// return that we have successfully uploaded our file!
	fmt.Fprintf(w, "Successfully Uploaded File\n")
}
