package main

import (
	"bytes"
	"strconv"
	"syscall/js"

	"github.com/tobischo/gokeepasslib"
)

func newDatabaseHandler(data interface{}) error {
	buf := new(bytes.Buffer)

	err := embeddedTemplates.ExecuteTemplate(buf, "new.html", data)
	if err != nil {
		return err
	}

	displayHTML(buf.String(), jsBody)

	return nil
}

func unlockHandler(data interface{}) error {
	buf := new(bytes.Buffer)

	err := embeddedTemplates.ExecuteTemplate(buf, "unlock.html", data)
	if err != nil {
		return err
	}

	displayHTML(buf.String(), jsBody)

	return nil
}

func databaseHandler(data interface{}, root bool) error {
	buf := new(bytes.Buffer)

	if root {
		err := embeddedTemplates.ExecuteTemplate(buf, "databaseRoot.html", data)
		if err != nil {
			return err
		}
	} else {
		err := embeddedTemplates.ExecuteTemplate(buf, "database.html", data)
		if err != nil {
			return err
		}
	}

	displayHTML(buf.String(), jsBody)

	if root {
		//DeletedObjects click handler
		for i := 0; i < len(data.(*gokeepasslib.RootData).DeletedObjects); i++ {
			index := strconv.Itoa(i)

			jsDocument.Call("getElementById", "deleted-"+index).Call("addEventListener", "click", js.FuncOf(jsDeletedObjectsCallback))
		}

		//Groups click handler
		for i := 0; i < len(data.(*gokeepasslib.RootData).Groups); i++ {
			index := strconv.Itoa(i)

			jsDocument.Call("getElementById", "group-"+index).Call("addEventListener", "click", js.FuncOf(jsGroupCallback))
		}
	} else {
		//Entries click handler
		for i := 0; i < len(data.(*gokeepasslib.Group).Entries); i++ {
			for n := 0; n < len(data.(*gokeepasslib.Group).Entries[i].Values); n++ {
				if data.(*gokeepasslib.Group).Entries[i].Values[n].Key != "Title" {
					continue
				}

				index := strconv.Itoa(i)

				jsDocument.Call("getElementById", "entry-"+index).Call("addEventListener", "click", js.FuncOf(jsEntryCallback))
			}
		}

		//Groups click handler
		for i := 0; i < len(data.(*gokeepasslib.Group).Groups); i++ {
			index := strconv.Itoa(i)

			jsDocument.Call("getElementById", "group-"+index).Call("addEventListener", "click", js.FuncOf(jsGroupCallback))
		}
	}

	return nil
}

func jsDeletedObjectsCallback(this js.Value, args []js.Value) interface{} {
	event <- args[0].Get("target").Get("id").String()

	return nil
}

func jsEntryCallback(this js.Value, args []js.Value) interface{} {
	event <- args[0].Get("target").Get("id").String()

	return nil
}
func jsGroupCallback(this js.Value, args []js.Value) interface{} {
	event <- args[0].Get("target").Get("id").String()

	return nil
}
