package main

import (
    "encoding/json"
    "github.com/gorilla/mux"
    "github.com/bestmethod/logger"
    "net/http"
    "syscall"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
)

// The person Type (more like an object)
type Person struct {
    ID        string   `json:"id,omitempty"`
    Firstname string   `json:"firstname,omitempty"`
    Lastname  string   `json:"lastname,omitempty"`
    Address   *Address `json:"address,omitempty"`
}
type Address struct {
    City  string `json:"city,omitempty"`
    State string `json:"state,omitempty"`
}

var people []Person

// Display all from the people var
func GetPeople(w http.ResponseWriter, r *http.Request) {
    log.Debug("GetPeople request, threadId: %d", syscall.Gettid())
    json.NewEncoder(w).Encode(people)
}

// Display a single data
func GetPerson(w http.ResponseWriter, r *http.Request) {
    log.Debug("GetPerson request, threadId: %d", syscall.Gettid())
    params := mux.Vars(r)
    for _, item := range people {
        if item.ID == params["id"] {
            json.NewEncoder(w).Encode(item)
            return
        }
    }
    json.NewEncoder(w).Encode(&Person{})
}

// create a new item
func CreatePerson(w http.ResponseWriter, r *http.Request) {
    log.Debug("CreatePerson request, threadId: %d", syscall.Gettid())
    params := mux.Vars(r)
    var person Person
    _ = json.NewDecoder(r.Body).Decode(&person)
    person.ID = params["id"]
    people = append(people, person)
    json.NewEncoder(w).Encode(people)
}

// Delete an item
func DeletePerson(w http.ResponseWriter, r *http.Request) {
    log.Debug("DeletePerson request, threadId: %d", syscall.Gettid())
    params := mux.Vars(r)
    for index, item := range people {
        if item.ID == params["id"] {
            people = append(people[:index], people[index+1:]...)
            break
        }
        json.NewEncoder(w).Encode(people)
    }
}

var log *Logger.Logger

func init() {
  log = new(Logger.Logger)
  log.Init("Main", "TutorialRest",
    Logger.LEVEL_DEBUG | Logger.LEVEL_INFO |
    Logger.LEVEL_WARN, Logger.LEVEL_ERROR |
    Logger.LEVEL_CRITICAL, Logger.LEVEL_NONE)
}

// main function to boot up everything
func main() {
  argsWithoutProg := os.Args[1:]

  log.Debug("Command line arguments: %s", argsWithoutProg := os.Args[1:])

  host := argsWithoutProg[0]
  port := argsWithoutProg[1]
  user := argsWithoutProg[2]
  password := argsWithoutProg[3]

  log.Debug("Host: %s", host)
  log.Debug("Port: %s", port)
  log.Debug("User: %s", user)
  log.Debug("Password: %s", password)

  log.Debug("Connecting to MongoDB ...")
  session, mongoSessionError := mgo.Dial(host + ":" + port)

  if mongoSessionError != nil {
    log.Fatal(1, "Failed to connect to MongoDB: %s", mongoSessionError)
  }

  log.Debug("Getting database ...")
  database, mongoDatabaseError = session.DB("persondb")

  router := mux.NewRouter()
  router.HandleFunc("/people", GetPeople).Methods("GET")
  router.HandleFunc("/people/{id}", GetPerson).Methods("GET")
  router.HandleFunc("/people/{id}", CreatePerson).Methods("POST")
  router.HandleFunc("/people/{id}", DeletePerson).Methods("DELETE")

  log.Debug("Listening on port 8000, threadId: %d", syscall.Gettid())
  err := http.ListenAndServe(":8000", router)
  log.Fatalf(1, "Failed to listen for HTTP connections %s", err)
}
