package main

import (
	"net/url"
	"os"
	"runtime"
	"sync"

	"github.com/yosebyte/passport/pkg/log"
)

var version = "dev"

func readme() {
	log.Info(`Version: %v %v/%v

Usage: <core_scheme>://<link_address>/<target_address>#<auth_url>

Examples:
    # Run as a server
    passport server://10.0.0.1:10101/:10022#http://:80/secret

    # Run as a client
    passport client://10.0.0.1:10101/127.0.0.1:22

    # Run as a broker
    passport broker://:8080/10.0.0.1:8080#https://:443/secret

Arguments:
    <core_scheme>    Select from "server", "client" or "broker"
    <link_address>   Tunneling or forwarding address to connect
    <target_address> Service address to be exposed or forwarded
    <auth_url>       Optional authorizing options in URL format
`, version, runtime.GOOS, runtime.GOARCH)
}

func main() {
	if len(os.Args) < 2 {
		readme()
		os.Exit(1)
	}
	rawURL := os.Args[1]
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		log.Fatal("Error parsing raw URL: %v", err)
	}
	var whiteList sync.Map
	authSetups(parsedURL, &whiteList)
	coreSelect(parsedURL, rawURL, &whiteList)
}
