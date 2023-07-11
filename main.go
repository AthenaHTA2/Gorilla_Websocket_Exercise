// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"Gorilla_Websocket_Exercise/gorilla"
	"flag"
	"log"
	"net/http"
	_ "github.com/mattn/go-sqlite3" 
	"Gorilla_Websocket_Exercise/database"
)
//The _ before the github.com/mattn/go-sqlite3 import is necessary 
//to ensure that the SQLite driver is registered with the database/sql package.



var addr = flag.String("addr", ":8080", "http service address")

func ServeHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "./views/home.html")
}

func main() {
	database.ConnectDB()

	flag.Parse()
	hub := gorilla.NewHub()
	go hub.Run()
	//Create the static folder and then use http.StripPrefix
	// so that the URL path will be independent of the actual folder structure
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", ServeHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		gorilla.ServeWs(hub, w, r)
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
