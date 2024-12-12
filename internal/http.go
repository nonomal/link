package internal

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"sync"

	"github.com/yosebyte/passport/pkg/log"
)

func HandleHTTP(parsedURL *url.URL, whiteList *sync.Map, tlsConfig *tls.Config) error {
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
			log.Error("Unable to serve HTTP: %v", err)
			return err
		}
	} else {
		authServer := &http.Server{
			Addr:      parsedURL.Host,
			TLSConfig: tlsConfig,
			ErrorLog:  log.NewLogger(),
		}
		if err := authServer.ListenAndServeTLS("", ""); err != nil {
			log.Error("Unable to serve HTTPS: %v", err)
			return err
		}
	}
	return nil
}
