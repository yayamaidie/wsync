package webserver

import (
	"crypto/tls"
	"encoding/pem"
	"fmt"
	"log"
	"net/http"

	"wsync/conf"
	"wsync/global"

	"github.com/gin-gonic/gin"
)

func InitWebServer() {
	go func() {
		engine := gin.Default()
		gin.SetMode(gin.ReleaseMode)
		initRouter(engine)

		// err := engine.Run(conf.GlobalConfig.Puller.ListenAddr)
		// err := engine.RunTLS(conf.GlobalConfig.Puller.ListenAddr, string(global.CertBytes), string(global.KeyBytes))
		// if err != nil {
		// 	global.ReleaseSources()
		// 	log.Fatalln(err.Error())
		// }

		if conf.GlobalConfig.Puller.HTTPS {
			if err := runTlsWebServer(engine); err != nil {
				global.ReleaseSources()
				log.Fatalln(err.Error())
			}
		} else {
			log.Printf("Listening and serving HTTP on %s\n", conf.GlobalConfig.Puller.ListenAddr)
			if err := http.ListenAndServe(conf.GlobalConfig.Puller.ListenAddr, engine); err != nil {
				global.ReleaseSources()
				log.Fatalln(err.Error())
			}
		}
	}()
}

func initRouter(router *gin.Engine) {
	api := router.Group("")
	{
		api.POST("/watch", watchHandle)
	}
}

func runTlsWebServer(_engine *gin.Engine) error {
	certblock, _ := pem.Decode(global.CertBytes)
	if certblock == nil || certblock.Type != "CERTIFICATE" {
		return fmt.Errorf("failed to decode certificate PEM")
	}
	keyblock, _ := pem.Decode(global.KeyBytes)
	if keyblock == nil || keyblock.Type != "PRIVATE KEY" {
		return fmt.Errorf("failed to decode key PEM")
	}

	// x509cert, err := x509.ParseCertificate(certblock.Bytes)
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }
	// x509key, err := x509.ParsePKCS8PrivateKey(keyblock.Bytes)
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }

	tlscert, err := tls.X509KeyPair(global.CertBytes, pem.EncodeToMemory(keyblock))
	if err != nil {
		return fmt.Errorf("failed to create x509 key pair: %v", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{tlscert},
	}

	server := &http.Server{
		Addr:      conf.GlobalConfig.Puller.ListenAddr,
		Handler:   _engine.Handler(),
		TLSConfig: tlsConfig,
	}

	log.Printf("Listening and serving HTTPS on %s\n", conf.GlobalConfig.Puller.ListenAddr)
	err = server.ListenAndServeTLS("", "")
	if err != nil {
		return fmt.Errorf("failed to start server: %v", err)
	}
	return nil
}
