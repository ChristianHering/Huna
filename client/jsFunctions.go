package main

import (
	"log"
	"syscall/js"
)

var (
	jsDocument js.Value
	jsBody     js.Value
)

func init() {
	jsDocument = js.Global().Get("document")
	jsBody = jsDocument.Get("body")
}

func jsSubmissionCallback(this js.Value, args []js.Value) interface{} {
	event <- ""

	return nil
}

func displayHTML(compiledTemplate string, htmlElement js.Value) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("JS Sucks:", err)
		}
	}()

	htmlElement.Set("innerHTML", compiledTemplate)

	return
}
