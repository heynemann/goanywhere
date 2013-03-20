package main

import (
    "github.com/gorilla/mux"
    "net/http"
    "log"
    "labix.org/v2/mgo"
    web "goanywhere/web"
)

var (
    Session *mgo.Session
)

func main() {
    session, err := mgo.Dial("localhost:8888")
    if err != nil {
        panic(err)
    }
    defer session.Close()

    session.SetMode(mgo.Monotonic, true)

    router := mux.NewRouter()
    router.HandleFunc("/", web.Index(session))
    router.HandleFunc("/route", web.Router(session))

    http.Handle("/", router)

    err = http.ListenAndServe(":12345", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
