package main

import (
	"flag"
	"regexp"
	"strings"
)

type settingType struct {
	servers         []string
	domain          string
	proxyPort       int
	whoisServer     string
	listen          string
	dnsInterface    string
	netSpecificMode string
}

var setting settingType

func isDigit(str string) bool {
	matched, _ := regexp.MatchString("^\\d+$", str)
	return matched
}

func main() {
	var settingDefault = settingType{
		[]string{""},
		"",
		8000,
		"172.20.129.8",
		":5000",
		"172.23.0.53",
		"",
	}

	serversPtr := flag.String("servers", strings.Join(settingDefault.servers, ","), "server name prefixes, separated by comma")
	domainPtr := flag.String("domain", settingDefault.domain, "server name domain suffixes")
	proxyPortPtr := flag.Int("proxy-port", settingDefault.proxyPort, "port bird-lgproxy is running on")
	whoisPtr := flag.String("whois", settingDefault.whoisServer, "whois server for queries")
	listenPtr := flag.String("listen", settingDefault.listen, "address bird-lg is listening on")
	dnsInterfacePtr := flag.String("dns-interface", settingDefault.dnsInterface, "dns zone to query ASN information")
	netSpecificModePtr := flag.String("net-specific-mode", settingDefault.netSpecificMode, "network specific operation mode, [(none)|dn42]")
	flag.Parse()

	if *serversPtr == "" {
		panic("no server set")
	} else if *domainPtr == "" {
		panic("no base domain set")
	}

	setting = settingType{
		strings.Split(*serversPtr, ","),
		*domainPtr,
		*proxyPortPtr,
		*whoisPtr,
		*listenPtr,
		*dnsInterfacePtr,
		strings.ToLower(*netSpecificModePtr),
	}

	webServerStart()
}
