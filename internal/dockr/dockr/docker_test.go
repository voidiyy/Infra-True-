package dockr_test

import (
	"log"
	"os"
	"testing"
)


func TestDockr(t *testing.T) {
	t.Run("NewDockr", func(t *testing.T) {
		log.Printf("DOCKER_HOST: %s", os.Getenv("DOCKER_HOST"))
		log.Printf("DOCKER_TLS_VERIFY: %s", os.Getenv("DOCKER_TLS_VERIFY"))
		log.Printf("DOCKER_CERT_PATH: %s", os.Getenv("DOCKER_CERT_PATH"))
	})
}