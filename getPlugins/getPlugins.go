package main

import (
	"encoding/json"
	//"fmt"
	//"github.com/julienschmidt/httprouter"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	//"reflect"
	"strings"
)

// struct holding repo info
type Plugin struct {
	Name        string
	FullName    string
	Type        string
	Owner       string
	Description string
	URL         string
	Forks       int
	Stars       int
	Watchers    int
	Issues      int
}

func JSONString(key string, value string) {
	retVal := '"' + key + '": "' + value + '"'
	return retVal
}

func JSONString(key string, value int) {
	retVal := '"' + key + '": ' + string(value)
	return retVal
}

func repoToPlugin(repo github.Repository) {
	plugin := new(Plugin)

	name := strings.Split(repo.Name, "-")
	plugin.Name =  JSONString("name", name[len(name) - 1])
	plugin.FullName = JSONString("full_name", repo.Name)
	plugin.Type = JSONString("type", name[len(name) - 2])
	plugin.Owner = JSONString("owner", repo.Owner.Login)
	plugin.Description = JSONString("description", repo.Description)
	plugin.URL = JSONString("url", repo.HTMLURL)
	plugin.Forks = JSONInt("fork_count", repo.ForksCount)
	plugin.Stars = JSONInt("star_count", repo.StargazersCount)
	plugin.Watchers = JSONInt("watch_count", repo.WatchersCount)
	plugin.Issues = JSONInt("issues_count", repo.OpenIssuesCount)

	return plugin
}

func main() {
	// Get Oauth token from environment.
	token := os.Getenv("GHPATOKEN")

	// Authenticate with github.
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tokenClient := oauth2.NewClient(oauth2.NoContext, tokenSource)

	// Client for interacting with github API
	client := github.NewClient(tokenClient)

	// List of links
	var links []string
	var plugins []Plugin
	var repos []string

	// First we go out and get the plugin catalog.
	response, err := http.Get("https://raw.githubusercontent.com/intelsdi-x/snap/master/docs/PLUGIN_CATALOG.md")
	if err != nil {
		log.Fatal(err)
	}
	bCatalog, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	sCatalog := strings.Split(string(bCatalog), "\n")

	// This is where we grab the links to the plugins. Very easily broken.
	// We currently rely heavily on the formatting of the md file.
	// Assumes that any line that starts with | is a table line.
	// Assumes that the 5th item in that line is the link.
	// Assumes that the link is formatted [displayName](url)
	// FIXME: If the format of the md file ever changes, this entire script will break. A better option would be to create
	// a csv file that both the md file and this script pull from.
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

	// This is where we get the repo information with the API.
	for _, link := range links {
		link := strings.Split(link, "/")
		repo, _, err := client.Repositories.Get(link[len(link)-2], link[len(link)-1])
		if err != nil {
			log.Fatal(err)
		}
		plugins := append(plugins, repoToPlugin(repo)
	}

	//repoToString := "[\n" + strings.Join(repos, ",\n") + "\n]"

	//ioutil.WriteFile("plugins.json", []byte(repoToString), 0644)
}
