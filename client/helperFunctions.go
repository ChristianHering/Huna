package main

import (
	"net/url"
	"syscall/js"
)

//constructURL creates a new url off our
//current host and adds the given path to it
func constructURL(path string) string {
	host := js.Global().Get("window").Get("location").Get("origin").String()

	u, err := url.Parse(host)
	if err != nil {
		panic(err)
	}

	u.Path = path

	return u.String()
}
