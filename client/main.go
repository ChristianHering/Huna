package main

import (
	"bytes"
	"embed"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"syscall/js"
	"time"

	"github.com/tobischo/gokeepasslib"
)

//go:embed templates/*.html
var vFS embed.FS

var embeddedTemplates *template.Template

var event chan string

func init() {
	embeddedTemplates = template.Must(template.ParseFS(vFS, "templates/*.html"))
}

func main() {
	event = make(chan string)

	getCookies()

	c, err := getCookie("username")
	if err != nil {
		log.Fatalln("Error getting username cookie:", err)
	}

	username := c.Value

download:

	response, err := http.Get(constructURL("/huna/download"))
	if err != nil {
		log.Println("Error downloading database:", err)

		time.Sleep(time.Second * 5)

		goto download
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNotFound {
		err = newDatabaseHandler(struct{ Name string }{Name: username})
		if err != nil {
			log.Fatalln("Error executing new database template:", err)
		}

		jsDocument.Call("getElementById", "submit").Call("addEventListener", "click", js.FuncOf(jsSubmissionCallback))

		<-event

		password := jsDocument.Call("getElementById", "password").Get("value").String()

		db, err := newDatabase(password)
		if err != nil {
			log.Fatalln("Error creating new database:", err)
		}

		response, err := updateDatabase(db)
		if err != nil {
			log.Panicln("Error updating database:", err)

			time.Sleep(time.Second * 5)

			goto download
		}

		if response.StatusCode != http.StatusOK {
			log.Println("uploading our new database file failed:", response.Status)
		}

		goto download
	}

	err = unlockHandler(struct{ Name string }{Name: username})
	if err != nil {
		log.Fatalln("Error executing unlock template:", err)
	}

	jsDocument.Call("getElementById", "submit").Call("addEventListener", "click", js.FuncOf(jsSubmissionCallback))

	buf := bytes.Buffer{}

	_, err = (&buf).ReadFrom(response.Body)
	if err != nil {
		log.Fatalln("Error reading in bytes from response body:", err)
	}

	var db *gokeepasslib.Database
	var dbBytes = buf.Bytes()

	for {
		<-event

		password := jsDocument.Call("getElementById", "password").Get("value").String()

		db, err = parseAndUnlockDatabase(dbBytes, password)
		if err != nil {
			log.Println("parseAndUnlockDatabase failed:", err)
			log.Println("You likely entered the wrong password.")
		} else {
			break
		}
	}

	for {
		display(db.Content.Root)
	}

	<-make(chan bool)
}

func display(data interface{}) {
	for {
		switch data.(type) {
		case *gokeepasslib.RootData:
			err := databaseHandler(data, true)
			if err != nil {
				log.Fatalln("Error executing root database template:", err)
			}
		case *gokeepasslib.Group:
			err := databaseHandler(data, false)
			if err != nil {
				log.Fatalln("Error executing database template:", err)
			}
		}

		callbackEvent := <-event

		if strings.Contains(callbackEvent, "back") {
			break
		} else if strings.Contains(callbackEvent, "group-") {
			i, err := strconv.Atoi(callbackEvent[6:])
			if err != nil {
				log.Fatalln("Error converting string to int:", err)
			}

			switch data.(type) {
			case *gokeepasslib.RootData:
				display(&data.(*gokeepasslib.RootData).Groups[i])
			case *gokeepasslib.Group:
				display(&data.(*gokeepasslib.Group).Groups[i])
			}
		} else if strings.Contains(callbackEvent, "entry-") {
			break // TODO: Implement entry display
		} else if strings.Contains(callbackEvent, "deleted-") {
			break // TODO: Implement deleted entry display
		}
	}
}

//main window where password groups are
//add button at bottom

//Add pop-up window for entering information

//right click menu for copying stuff

//A window with your password listed oldest -> newest

//copy to clipboard navigator.clipboard.writeText('Copy this text')

//TODO: remove GOOS=js GOARCH=wasm go build -o ./../asm/bin.wasm
