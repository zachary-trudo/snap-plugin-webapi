package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

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
	// This is our main page. We may want to update it to look a bit spiffier.
	htmlOut := `Snap Plugin API Server:

<p>/plugins</p>
<p>/plugins/collector</p>
<p>/plugins/processor</p>
<p>/plugins/publisher</p>
<p>/plugin/:name</p>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(htmlOut))
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

	var port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(":"+port, router))
}
