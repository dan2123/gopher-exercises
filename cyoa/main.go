package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type StoryArc struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []Option `json:"options"`
}

type Option struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

func main() {
	// Read all story arcs into memory
	storyArcs, err := readDataFromFile("gopher.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Initialize template
	tmpl := template.Must(template.ParseFiles("layout.html"))

	// Call story arc handler
	h, err := storyArcHandler(storyArcs, tmpl)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Initialize server and pass in handler
	fmt.Println("Starting web server on port: 8080")
	http.ListenAndServe(":8080", h)
}

func storyArcHandler(storyArcs map[string]StoryArc, tmpl *template.Template) (http.HandlerFunc, error) {
	return func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path == "/" {
			tmpl.Execute(writer, storyArcs["intro"])
			return
		}

		requestStoryArc := strings.Trim(request.URL.Path, "/")
		if arc, ok := storyArcs[requestStoryArc]; ok {
			tmpl.Execute(writer, arc)
		}
	}, nil
}

func readDataFromFile(fileName string) (map[string]StoryArc, error) {
	out := make(map[string]StoryArc)
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}
