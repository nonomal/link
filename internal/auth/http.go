package auth

import (
	"net"
	"net/http"
	"net/url"
	"sync"

	"github.com/yosebyte/passport/pkg/log"
	"github.com/yosebyte/passport/pkg/tls"
)

func HandleHTTP(parsedURL *url.URL, whiteList *sync.Map) error {
	http.HandleFunc(parsedURL.Path, func(w http.ResponseWriter, r *http.Request) {
		clientIP, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			log.Error("Invalid client IP address: [%v]", clientIP)
			return
		}
		if _, err := w.Write([]byte(clientIP + "\n")); err != nil {
			log.Error("Unable to write client IP address: [%v]", clientIP)
			return
		}
		whiteList.Store(clientIP, struct{}{})
		log.Info("Authorized IP address added: [%v]", clientIP)
	})
	if parsedURL.Scheme == "http" {
		if err := http.ListenAndServe(parsedURL.Host, nil); err != nil {
			log.Error("Error serving HTTP: %v", err)
			return err
		}
	} else {
		tlsConfig, err := tls.NewTLSconfig(parsedURL.Hostname())
		if err != nil {
			log.Error("Error generating TLS config: %v", err)
			return err
		}
		authServer := &http.Server{
			Addr:      parsedURL.Host,
			TLSConfig: tlsConfig,
			ErrorLog:  log.NewLogger(),
		}
		if err := authServer.ListenAndServeTLS("", ""); err != nil {
			log.Error("Error serving HTTPS: %v", err)
			return err
		}
	}
	return nil
}
