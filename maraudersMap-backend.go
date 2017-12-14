package main

import (
    "log"
    "fmt"
    "net/http"
    "github.com/nanobox-io/golang-scribble"
    "github.com/gorilla/mux"
    "encoding/json"
    "os"
)

type Location struct {
  Name string `json:"Name"`
  Uuid string `json:"Uuid"`
  Accuracy float32 `json:"Accuracy"`
  RegionName string `json:"RegionName"`
}

var db *scribble.Driver

func welcomeHandler(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, "Welcome to the HomePage!")
    fmt.Println("Endpoint Hit: homePage")
}


func Welcome(w http.ResponseWriter, r *http.Request){
  fmt.Fprintln(w, "Welcome to wizard world!")
}

func NewWizard(w http.ResponseWriter, r *http.Request){
  vars := mux.Vars(r)
  name := vars["name"]
  fmt.Fprintln(w, "input name is: ", name)
}

func UpdateLocation(w http.ResponseWriter, r *http.Request){
  location := Location{}

  err := json.NewDecoder(r.Body).Decode(&location)
  defer r.Body.Close()
  if err != nil {
    panic(err)
  }
  if _, err := os.Stat("marauderMap/" + location.Name + ".json"); os.IsNotExist(err) {
    updateLocation(location)
  } else {
    current := Location{}
    db.Read("MarauderMap", location.Name, &current)
    fmt.Println("Wizard Exist");
    fmt.Println("Current reading is: ", location.RegionName)
    fmt.Println("Current accuracy is: ", location.Accuracy)
    fmt.Println("Accuracy from db is: ", current.Accuracy)
    if (location.Accuracy > 0) {
      if (current.RegionName == location.RegionName){
        fmt.Println("update location: same region but need to keep track of current proximity")
        updateLocation(location)
      }else {
        if (current.Accuracy > location.Accuracy) {
          fmt.Println("update location: different region and we are closer to new region now")
          updateLocation(location)
        }
      }
    }
    fmt.Println("========================");
  }
  fmt.Fprintln(w,location.Uuid)
}

func MarauderMap(w http.ResponseWriter, r *http.Request) {
  aLotOfWizards := []Location{}
  wizards, _ := db.ReadAll("MarauderMap")

  for _, wizard := range wizards {
    l := Location{}
    json.Unmarshal([]byte(wizard), &l)
   aLotOfWizards = append(aLotOfWizards, l)
  }

  js, err := json.Marshal(aLotOfWizards)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  w.Header().Set("Content-Type", "application/json")
  w.Write(js);
}

func updateLocation(location Location) {
  db.Write("marauderMap", location.Name, location)
}

func initialiseDb() (driver *scribble.Driver){
  dir := "./"
  db, err := scribble.New(dir, nil);
  if err != nil {
    fmt.Println("Error", err)
    }
  return db
}

func main() {
  db = initialiseDb()

  router := mux.NewRouter().StrictSlash(true)
  router.HandleFunc("/", Welcome)
  router.HandleFunc("/newWizard/{name}/", NewWizard)
  router.HandleFunc("/updateLocation/", UpdateLocation).Methods("POST")
  router.HandleFunc("/maraudersMap/", MarauderMap)

  log.Fatal(http.ListenAndServe(":1234", router))
}
