package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
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

// Creates the oauthConf variable.
// To replicate: Create developer oauth application at github.com/setting/applications
// ClientID and ClientSecret are supplied by github after creating the developer oauth application token.
// This gives us access to their username:email and repo information.
var (
	oauthConf = &oauth2.Config{
		ClientID:     "2502f1186a0b687847a8",
		ClientSecret: "3269896c65db7131a0bd473b1361777234aa4b86",
		Scopes:       []string{"user:email", "repo"},
		Endpoint:     githuboauth.Endpoint,
	}
	oauthStateString = "asdfghjkl"
)

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

<p>Log in with <a href="/login">Github</a> for access to external plugins.</p>

<p>/plugins</p>
<p>/plugins/collector</p>
<p>/plugins/processor</p>
<p>/plugins/publisher</p>
<p>/plugin/:name</p>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(htmlOut))
}

// Logs external users into Github in order to grant us access to their repo information.
func loginToGitHub(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	url := oauthConf.AuthCodeURL(oauthStateString, oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)

}

// Returns us from github after successful OAuth request.
// TODO: Perhaps this is where we could add any of their snap plugins to our list?
func callBack(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	state := r.FormValue("state")
	// TODO: if the State string is incorrect we should probably handle that.
	if state != oauthStateString {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	token, err := oauthConf.Exchange(oauth2.NoContext, code)
	// TODO: We might want to do some error handling?
	if err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	oauthClient := oauthConf.Client(oauth2.NoContext, token)
	client := github.NewClient(oauthClient)
	_, _, err = client.Users.Get("")
	// TODO: And again - Error Handling.
	if err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
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
	router.GET("/login", loginToGitHub)
	router.GET("/snap_github_cb", callBack)

	var port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(":"+port, router))
}
