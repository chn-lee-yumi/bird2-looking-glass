package main

import (
	"net/http"
	"os/exec"
	"strings"
)

func tracerouteHandler(httpW http.ResponseWriter, httpR *http.Request) {
	query := string(httpR.URL.Query().Get("q"))
	query = strings.TrimSpace(query)
	if !isIP(query) {
		invalidHandler(httpW, httpR)
	} else {
		cmd := exec.Command("bash", "-c", "traceroute -w1 " + query + " 2>&1")
		result, err := cmd.Output()
		if err != nil {
			httpW.WriteHeader(http.StatusInternalServerError)
		}
		httpW.Write(result)
	}
}
