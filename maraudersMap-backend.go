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

var region_name = map[string]string {
  "11A29583-9A74-4EDC-91B3-0A06A45DC539": "Syltherin common room",
}
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
  fmt.Println("printig out decode stuff: ", location.Name)
  location.RegionName = region_name[location.Uuid]
  if _, err := os.Stat("marauderMap/" + location.Name + ".json"); os.IsNotExist(err) {
    updateLocation(location)
  } else {
    current := Location{}
    db.Read("MarauderMap", location.Name, current)
    if (current.Accuracy < location.Accuracy) {
      updateLocation(location)
    }
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
  router.HandleFunc("/updateLocation/", UpdateLocation)
  router.HandleFunc("/maraudersMap/", MarauderMap)

  log.Fatal(http.ListenAndServe(":1234", router))
}
