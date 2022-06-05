package main

import (
	"errors"
	"strings"
	"syscall/js"
)

type Cookie struct {
	Key   string
	Value string
}

var ErrorFindCookieFailed = errors.New("Unable to locate cookie")

var cookies []Cookie

func getCookie(cookieName string) (cookie Cookie, err error) {
	for i := 0; i < len(cookies); i++ {
		if cookies[i].Key == cookieName {
			return cookies[i], nil
		}
	}

	return Cookie{}, ErrorFindCookieFailed
}

func getCookies() {
	cookieString := js.Global().Get("document").Get("cookie").String()

	cookies = parseCookies(cookieString)
}

//https://datatracker.ietf.org/doc/html/rfc6265
func parseCookies(cookieString string) (parsedCookies []Cookie) {
	cookies := strings.Split(cookieString, ";")

	for i := 0; i < len(cookies); i++ {
		keyValue := strings.Split(strings.TrimSpace(cookies[i]), "=")

		c := Cookie{
			Key:   keyValue[0],
			Value: keyValue[1],
		}

		parsedCookies = append(parsedCookies, c)
	}

	return parsedCookies
}
