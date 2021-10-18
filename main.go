package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

var tpl *template.Template

func init() {
	//  Get files ready for work
	tpl = template.Must(template.ParseGlob("templates/*.html"))
}

type Artists []struct {
	Id           int                 `json:"id"`
	Image        string              `json:"image"`
	Name         string              `json:"name"`
	Members      []string            `json:"members"`
	CreationDate int                 `json:"creationDate"`
	FirstAlbum   string              `json:"firstAlbum"`
	Relations    string              `json:"relations"`
	Connection   map[string][]string `json:"connection"`
}

type Relation struct {
	Id             int                 `json: "id`
	DatesLocations map[string][]string `json:"datesLocations`
}

//  Think , its only vor declaring a type
type (
	ApiRootURL string
)

var apis = ApiRootURL("https://groupietrackers.herokuapp.com/api").getFromApi()

func main() {
	fmt.Println("Starting server on port :8000")
	fs := http.FileServer(http.Dir("templates"))
	http.Handle("/templates", fs)
	http.Handle("/css/", fs)
	http.HandleFunc("/", loadMainPage)
	fmt.Println("Started\n visit http://localhost:8000/ ")

	// http://localhost:8000/
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func loadMainPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		fmt.Fprint(w, "This site is unavailable", http.StatusNotFound)
		return
	}
	body := getUrlBody(apis["artists"])
	var artists Artists
	err := json.Unmarshal(body, &artists)
	checkError(err, "Error with json parsing.")

	// GETS RELATION/ID and assign rel(structure of Relation) to current artist.Connections by getting from json DatesLocations
	for i := range artists {
		var rel Relation
		relation := getUrlBody(artists[i].Relations)
		json.Unmarshal(relation, &rel)
		artists[i].Connection = rel.DatesLocations
	}
	tpl.ExecuteTemplate(w, "index.html", artists)
}
func getUrlBody(URL string) []byte {
	response, err := http.Get(URL)
	checkError(err, "Error with GET response. ")
	body, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	return body
}

func (apiURL ApiRootURL) getFromApi() map[string]string {
	body := getUrlBody(string(apiURL))
	var apis map[string]string
	err := json.Unmarshal(body, &apis)
	checkError(err, "Error with json parsing (APIS).")
	return apis
}

func checkError(err error, msg string) {
	if err != nil {
		log.Fatal(msg, err)
	}
}
