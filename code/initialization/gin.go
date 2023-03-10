package initialization

import (
	"crypto/tls"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

func loadCertificate(config Config) (cert tls.Certificate, err error) {
	cert, err = tls.LoadX509KeyPair(config.CertFile, config.KeyFile)
	if err != nil {
		return cert, fmt.Errorf("failed to load certificate: %v", err)
	}
	// check certificate expiry
	certExpiry := cert.Leaf.NotAfter
	if certExpiry.Before(time.Now()) {
		return cert, fmt.Errorf("certificate expired on %v", certExpiry)
	}
	return cert, nil
}
func startHTTPServer(config Config, r *gin.Engine) (err error) {
	log.Printf("http server started: http://localhost:%d/webhook/event\n", config.HttpPort)
	err = r.Run(fmt.Sprintf(":%d", config.HttpPort))
	if err != nil {
		return fmt.Errorf("failed to start http server: %v", err)
	}
	return nil
}
func startHTTPSServer(config Config, r *gin.Engine) (err error) {
	cert, err := loadCertificate(config)
	if err != nil {
		return fmt.Errorf("failed to load certificate: %v", err)
	}
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.HttpsPort),
		Handler: r,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
		},
	}
	fmt.Printf("https server started: https://localhost:%d/webhook/event\n", config.HttpsPort)
	err = server.ListenAndServeTLS("", "")
	if err != nil {
		return fmt.Errorf("failed to start https server: %v", err)
	}
	return nil
}
func StartServer(config Config, r *gin.Engine) (err error) {
	if config.UseHttps {
		err = startHTTPSServer(config, r)
	} else {
		err = startHTTPServer(config, r)
	}
	return err
}
