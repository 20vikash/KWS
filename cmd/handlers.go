package main

import "net/http"

func (app *Application) HelloWorld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world"))
}
