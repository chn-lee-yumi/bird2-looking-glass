package main

import (
	"fmt"
	"github.com/gorilla/handlers"
	"html"
	"net/http"
	"os"
	"strings"
)

func webHandlerWhois(w http.ResponseWriter, r *http.Request) {
	var target string = r.URL.Path[len("/whois/"):]

	renderTemplate(
		w, r,
		"<h2>whois "+html.EscapeString(target)+"</h2>"+smartFormatter(whois(target)),
	)
}

func webBackendCommunicator(endpoint string, command string) func(w http.ResponseWriter, r *http.Request) {
	backendCommandPrimitive, commandPresent := (map[string]string{
		"summary":         "show protocols",
		"detail":          "show protocols all %s",
		"route":           "show route for %s",
		"route_all":       "show route for %s all",
		"route_where":     "show route where net ~ [ %s ]",
		"route_where_all": "show route where net ~ [ %s ] all",
		"route_generic":   "show route %s",
		"generic":         "show %s",
		"traceroute":      "%s",
		"ping":            "%s",
	})[command]

	if !commandPresent {
		panic("invalid command: " + command)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		split := strings.SplitN(r.URL.Path[1:], "/", 3)
		var urlCommands string
		if len(split) >= 3 {
			urlCommands = split[2]
		}

		var backendCommand string
		if strings.Contains(backendCommandPrimitive, "%") {
			backendCommand = fmt.Sprintf(backendCommandPrimitive, urlCommands)
		} else {
			backendCommand = backendCommandPrimitive
		}
		backendCommand = strings.TrimSpace(backendCommand)

		var servers []string = strings.Split(split[1], "+")
		var responses []string = batchRequest(servers, endpoint, backendCommand)
		var result string
		for i, response := range responses {
			result += "<h2>" + html.EscapeString(servers[i]) + ": " + html.EscapeString(backendCommand) + "</h2>"
			if endpoint == "bird" && backendCommand == "show protocols" && len(response) > 4 && strings.ToLower(response[0:4]) == "name" {
				result += summaryTable(response, servers[i])
			} else {
				result += smartFormatter(response)
			}
		}

		renderTemplate(w, r, result)
	}
}

func webHandlerBGPMap(endpoint string, command string) func(w http.ResponseWriter, r *http.Request) {
	backendCommandPrimitive, commandPresent := (map[string]string{
		"route_bgpmap":       "show route for %s all",
		"route_where_bgpmap": "show route where net ~ [ %s ] all",
	})[command]

	if !commandPresent {
		panic("invalid command: " + command)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		split := strings.Split(r.URL.Path, "/")
		urlCommands := strings.Join(split[3:], "/")
		var backendCommand string
		if strings.Contains(backendCommandPrimitive, "%") {
			backendCommand = fmt.Sprintf(backendCommandPrimitive, urlCommands)
		} else {
			backendCommand = backendCommandPrimitive
		}
		var servers []string = strings.Split(split[2], "+")
		var responses []string = batchRequest(servers, endpoint, backendCommand)
		renderTemplate(
			w, r,
			`<script>
			var viz = new Viz();
			viz.renderSVGElement(`+"`"+birdRouteToGraphviz(servers, responses, urlCommands)+"`"+`)
			.then(element => {
				document.body.appendChild(element);
			})
			.catch(error => {
				document.body.innerHTML = "<pre>"+error+"</pre>"
			});
			</script>`,
		)
	}
}

func webHandlerNavbarFormRedirect(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	if query.Get("action") == "whois" {
		http.Redirect(w, r, "/"+query.Get("action")+"/"+query.Get("target"), 302)
	} else if query.Get("action") == "summary" {
		http.Redirect(w, r, "/"+query.Get("action")+"/"+query.Get("server"), 302)
	} else {
		http.Redirect(w, r, "/"+query.Get("action")+"/"+query.Get("server")+"/"+query.Get("target"), 302)
	}
}

func webServerStart() {
	// Start HTTP server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/summary/"+strings.Join(setting.servers, "+"), 302)
	})
	http.HandleFunc("/summary/", webBackendCommunicator("bird", "summary"))
	http.HandleFunc("/detail/", webBackendCommunicator("bird", "detail"))
	http.HandleFunc("/route/", webBackendCommunicator("bird", "route"))
	http.HandleFunc("/route_all/", webBackendCommunicator("bird", "route_all"))
	http.HandleFunc("/route_bgpmap/", webHandlerBGPMap("bird", "route_bgpmap"))
	http.HandleFunc("/route_where/", webBackendCommunicator("bird", "route_where"))
	http.HandleFunc("/route_where_all/", webBackendCommunicator("bird", "route_where_all"))
	http.HandleFunc("/route_where_bgpmap/", webHandlerBGPMap("bird", "route_where_bgpmap"))
	// http.HandleFunc("/route_generic/", webBackendCommunicator("bird", "route_generic"))
	// http.HandleFunc("/generic/", webBackendCommunicator("bird", "generic"))
	http.HandleFunc("/traceroute/", webBackendCommunicator("traceroute", "traceroute"))
	http.HandleFunc("/ping/", webBackendCommunicator("ping", "ping"))
	http.HandleFunc("/whois/", webHandlerWhois)
	http.HandleFunc("/redir", webHandlerNavbarFormRedirect)
	// http.HandleFunc("/telegram/", webHandlerTelegramBot)
	http.ListenAndServe(setting.listen, handlers.LoggingHandler(os.Stdout, http.DefaultServeMux))
}
