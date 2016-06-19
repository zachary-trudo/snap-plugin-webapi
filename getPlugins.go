package main

import (
	//"encoding/json"
	"fmt"
	//"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"log"
	"net/http"
	//"os"
	"github.com/google/go-github/github"
	"strings"
)

// struct holding repo info
type Plugin struct {
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	Type        string `json:"type"`
	Owner       string `json:"owner"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Forks       int    `json:"fork_count"`
	Stars       int    `json:"star_count"`
	Watchers    int    `json:"watch_count"`
	Issues      int    `json:"issues_count"`
}

func main() {
	client := github.NewClient(nil)
	var links []string

	response, err := http.Get("https://raw.githubusercontent.com/intelsdi-x/snap/master/docs/PLUGIN_CATALOG.md")
	if err != nil {
		log.Fatal(err)
	}
	bCatalog, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("%s", catalog)
	sCatalog := strings.Split(string(bCatalog), "\n")
	for _, line := range sCatalog {
		if strings.HasPrefix(line, "|") {
			splitLine := strings.Split(line, "|")
			if len(splitLine) > 3 {
				link := strings.Split(splitLine[4], "(")
				if len(link) > 1 {
					links = append(links, strings.Split(link[1], ")")[0])
				}
			}
		}
	}
	fmt.Println(links)
	for _, link := range links {
		link := strings.Split(link, "/")
		repo, _, err := client.Repositories.Get(link[len(link)-2], link[len(link)-1])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(repo.String())
	}
}
