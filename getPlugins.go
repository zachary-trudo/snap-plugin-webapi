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

// struct holding user info
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
	opt := &github.RepositoryListByOrgOptions{Type: "public"}
	repos, _, _ := client.Repositories.ListByOrg("intelsdi-x", opt)
	for repo := range repos {
		fmt.Println(repo)
	}
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
		if strings.HasPrefix(line, "##") {
			fmt.Println(line)
		}
	}
}
