package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Plugin struct {
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	Type        string `json:"type"`
	Owner       string `json:"owner"`
	Description string `json: "desription"`
	URL         string `json:"url"`
	Forks       int    `json:"fork_count"`
	Stars       int    `json:"star_count"`
	Watchers    int    `json:"watch_count"`
	Issues      int    `json:"issues_count"`
}

func Filter(vs []Plugin, f func(Plugin) bool) []Plugin {
	vsf := make([]Plugin, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, `Snap Plugin API Server:

/plugins
/plugins/collector
/plugins/processor
/plugins/publisher
/plugin/:name`)
}

func ListPlugin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

func ListPlugins(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	file, e := ioutil.ReadFile("./plugins.json")
	if e != nil {
		fmt.Fprintf(w, "File error: %v\n", e)
	}

	plugins := make([]Plugin, 0)
	err := json.Unmarshal(file, &plugins)
	if err != nil {
		fmt.Fprint(w, err)
	} else {
		plugin_type := strings.ToLower(ps.ByName("type"))
		if plugin_type != "" {
			plugins = Filter(plugins, func(v Plugin) bool {
				return strings.Contains(v.Type, plugin_type)
			})
		}

		output, _ := json.MarshalIndent(plugins, "", "    ")
		fmt.Fprint(w, string(output))
	}
}

func main() {
	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/plugins", ListPlugins)
	router.GET("/plugins/:type", ListPlugins)
	router.GET("/plugin/:name", ListPlugin)

	log.Fatal(http.ListenAndServe(":8080", router))
}
