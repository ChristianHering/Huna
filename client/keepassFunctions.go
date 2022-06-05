package main

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/tobischo/gokeepasslib"
)

//updateDatabase posts the given locked(!) database to the server
//
//If updateDatabase is passed a keepass database that isn't locked,
//passwords and other fields could be sent to the server IN PLAIN TEXT!
func updateDatabase(db *gokeepasslib.Database) (response *http.Response, err error) {
	b := bytes.Buffer{}

	err = gokeepasslib.NewEncoder(&b).Encode(db)
	if err != nil {
		return nil, err
	}

	body := bytes.Buffer{}

	writer := multipart.NewWriter(&body)

	part, err := writer.CreateFormFile("dbFile", "database.kdbx")
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(part, &b)
	if err != nil {
		return nil, err
	}

	err = writer.Close() //Don't move this
	if err != nil {
		return nil, err
	}

	response, err = http.Post(constructURL("/huna/update"), writer.FormDataContentType(), &body)
	if err != nil {
		return nil, err
	}

	return response, err
}

//parseAndUnlockDatabase parses a database from the given io.Reader and unlocks it with the given password
func parseAndUnlockDatabase(b []byte, password string) (db *gokeepasslib.Database, err error) {
	db = gokeepasslib.NewDatabase()
	db.Credentials = gokeepasslib.NewPasswordCredentials(password)

	err = gokeepasslib.NewDecoder(bytes.NewReader(b)).Decode(db)
	if err != nil {
		return nil, err
	}

	err = db.UnlockProtectedEntries()
	if err != nil {
		return nil, err
	}

	return db, nil
}

//newDatabase returns a new, locked database with the given password
func newDatabase(password string) (db *gokeepasslib.Database, err error) {
	db = gokeepasslib.NewDatabase()
	db.Credentials = gokeepasslib.NewPasswordCredentials(password)

	b := bytes.Buffer{}

	err = gokeepasslib.NewEncoder(&b).Encode(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}
