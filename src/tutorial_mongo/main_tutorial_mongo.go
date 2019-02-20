package main

import (
    "encoding/json"
    "github.com/gorilla/mux"
    "github.com/bestmethod/logger"
    "net/http"
    "syscall"
    "os"
    "crypto/x509"
    "crypto/tls"
    "io/ioutil"
    "net"
    mgo "gopkg.in/mgo.v2"
    bson "gopkg.in/mgo.v2/bson"
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

    queryResult := database.C("person").Find(bson.M{})
    count, _ := queryResult.Count()
    log.Debug("Queried database and found %d entries", count)
    iterator := queryResult.Iter()
    defer func() {
      log.Debug("Closing iterator for results")
      closeErr := iterator.Close()
      if closeErr != nil { log.Debug("Got error, %s", closeErr) }
    }() // Yeah, this is confusing syntax, declared and then called anonymous function

    var result Person
    encoder := json.NewEncoder(w)
    for iterator.Next(&result) {
      encoder.Encode(result)
    }
}

// Display a single data
func GetPerson(w http.ResponseWriter, r *http.Request) {
    log.Debug("GetPerson request, threadId: %d", syscall.Gettid())
    params := mux.Vars(r)
    idToSearchFor := params["id"]

    var result Person
    queryError := database.C("person").Find(bson.M{ "id": idToSearchFor }).One(&result)

    if queryError != nil {
      log.Error("Error finding row for '%s', failure: '%s'", idToSearchFor, queryError)
    } else {
      log.Debug("Found an entry for ID: %s", idToSearchFor)
      json.NewEncoder(w).Encode(result)
    }
}

// create a new item
func CreatePerson(w http.ResponseWriter, r *http.Request) {
    log.Debug("CreatePerson request, threadId: %d", syscall.Gettid())
    params := mux.Vars(r)
    var person Person
    _ = json.NewDecoder(r.Body).Decode(&person)
    person.ID = params["id"]

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
var database *mgo.Database

func init() {
  log = new(Logger.Logger)
  log.Init("Main", "TutorialRest",
    Logger.LEVEL_DEBUG | Logger.LEVEL_INFO |
    Logger.LEVEL_WARN, Logger.LEVEL_ERROR |
    Logger.LEVEL_CRITICAL, Logger.LEVEL_NONE)
}

func connectToMongo(url string, pemFile string) (*mgo.Session, error) {

  certsFromPemFile := x509.NewCertPool()

  log.Debug("Getting PEM file ...")
  if pemFileBytes, err := ioutil.ReadFile(pemFile); err == nil {
    certsFromPemFile.AppendCertsFromPEM(pemFileBytes)
  } else {
    log.Fatalf(1, "Failed to read certs from '%s', get errors: %s", pemFile, err)
  }

  tlsConfig := &tls.Config{}

  // TODO: I think this needs to filter out the client cert
  tlsConfig.RootCAs = certsFromPemFile

  dialInfo, err := mgo.ParseURL(url)
  dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
    log.Debug("DialServer function, connecting to %s", addr.String())
    conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
    log.Debug("Connection error: %s", err)
    return conn, err
  }

  log.Debug("Connecting to mongo server using SSL")
  session, err := mgo.DialWithInfo(dialInfo)
  log.Debug("Connecting complete")
  return session, err
}

// main function to boot up everything
func main() {
  argsWithoutProg := os.Args[1:]

  log.Debug("Command line arguments: %s", argsWithoutProg)

  url := argsWithoutProg[0]
  pemFile := argsWithoutProg[1]

  log.Debug("Url: %s", url)
  log.Debug("PemFile: %s", pemFile)

  log.Debug("Connecting to MongoDB ...")
  session, mongoSessionError := connectToMongo(url, pemFile)

  if mongoSessionError != nil {
    log.Fatalf(1, "Failed to connect to MongoDB: %s", mongoSessionError)
  }

  log.Debug("Getting database ...")
  database = session.DB("persondb")

  router := mux.NewRouter()
  router.HandleFunc("/people", GetPeople).Methods("GET")
  router.HandleFunc("/people/{id}", GetPerson).Methods("GET")
  router.HandleFunc("/people/{id}", CreatePerson).Methods("POST")
  router.HandleFunc("/people/{id}", DeletePerson).Methods("DELETE")

  log.Debug("Listening on port 8000, threadId: %d", syscall.Gettid())
  err := http.ListenAndServe(":8000", router)
  log.Fatalf(1, "Failed to listen for HTTP connections %s", err)
}
