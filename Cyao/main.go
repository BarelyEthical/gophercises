package main

import (
	"fmt"
	"os"
	"log"
	"encoding/json"
	"bytes"
	"io"
	"strings"
	"html/template"
	"net/http"
)

type StoryArc struct {
	Title string `json:title`
	Story []string `json:story`
	Options []struct {
		Text string `json:text`
		Arc string `json:arc`
	} `json:options`
}

type Result map[string]StoryArc
var webPage = template.Must(template.ParseFiles("cyoa.html"))

func (result Result) httpHandler(path string, w http.ResponseWriter) {
	t, err := webPage.Clone()
	if err != nil {
		return
	}
	story := result[path[1:]]
	t.Execute(w, story)
}

func (result Result) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path;
	if path == "/" {
		path =  "/intro"
	}
	result.httpHandler(path, w)
}

func ReadJsonFile(filename string) []byte {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal("Error opening the file %s", filename)
	}
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(file)
	if err != nil {
		log.Fatal("Error parsing the file $s", filename)
	}
	return buf.Bytes()
}

func StartYourAdventure(result Result, adv string) {
	fmt.Fprintf(os.Stdout, "%s\n", result[adv].Title)
	for _, value := range result[adv].Story {
		fmt.Fprintf(os.Stdout, "%s", value)
	}
	fmt.Fprintf(os.Stdout, "\n")
	flag := strings.Compare(result[adv].Title, "Home Sweet Home")
	if flag == 0 {
		return
	}
	nextOptions := result[adv].Options;
	for index, nextArc := range nextOptions {
		fmt.Fprintf(os.Stdout, "%d %s\n", index+1, nextArc.Text)
	}
	fmt.Fprintf(os.Stdout, "Enter your choice for the next adventure\n\n")
	var input int
	_, err := fmt.Fscanf(os.Stdin, "%d\n", &input)
	if err != nil {
		if err == io.EOF {
			return 
		}
		log.Fatal("Error reading the input from the user", err)
	}
	StartYourAdventure(result, result[adv].Options[input-1].Arc)
}

func main() {
	jsonFile := "gopher.json";
	jsonData := ReadJsonFile(jsonFile);
	var result Result
	json.Unmarshal(jsonData, &result)
	//StartYourAdventure(result, "intro")
	r := http.DefaultServeMux
	r.Handle("/", result)
	http.ListenAndServe(":8080", nil)
}
