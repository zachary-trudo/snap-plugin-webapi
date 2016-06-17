package main

import (
	//"encoding/json"
	"fmt"
	//"github.com/julienschmidt/httprouter"
	//"io/ioutil"
	//"log"
	//"net/http"
	//"os"
	//"strings"
	"github.com/google/go-github/github"
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
	for _, repo := range repos {
		fmt.Println(repo.String())
	}
}
