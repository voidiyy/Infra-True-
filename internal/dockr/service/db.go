package service

import (
	"log"
	"sync"
)

var _ Service = &SrvDB{}

// List of top DB images with latest version
const (
	imagePostgresLatest = "postgres:latest"
	imageMySQLLatest    = "mysql:latest"
	imageMongoDBLatest  = "mongo:latest"
	imageCassandra      = "cassandra:latest"
)

// defaultDBImages is a list of docker DB images
var defaultDBImages = []string{
	imagePostgresLatest,
	imageMySQLLatest,
	imageMongoDBLatest,
	imageCassandra,
}

// SrvDB represent an entity for docker container
// this is implementation of Service interface
type SrvDB struct {
	logic *SrvLogic
	mu    *sync.RWMutex
}

func DefaultSrvLogicDB() *SrvDB {
	return &SrvDB{
		logic: &SrvLogic{
			IsDefault:     true,
			Type:          DB,
			Image:         "postgres:latest",
			ContainerName: "db-container",
			Status:        SrvStatusRun(),
			Network:       "default-network",
			Ports:         []string{"5432:5432"},
			Volumes:       []string{"/var/lib/postgresql/data:/var/lib/postgresql/data"},
			EnvVars:       map[string]string{"POSTGRES_PASSWORD": "defaultpassword"},
			RestartPolicy: "always",
		},
		mu: &sync.RWMutex{},
	}
}

// NewSrvLogicDB creates a new instance of SrvDB with given config.
// It copies all fields from given SrvConfig to SrvLogic.
// If SrvConfig is nil, it will panic.
// It is used to create a new instance of SrvDB with given config.
func NewSrvLogicDB(conf *SrvConfig) *SrvDB {
	return &SrvDB{
		logic: &SrvLogic{
			IsDefault:     false,
			Type:          conf.Type,
			Image:         conf.Image,
			ContainerName: conf.ContainerName,
			Status:        SrvStatusRun(),
			Network:       conf.Network,
			Ports:         conf.Ports,
			Volumes:       conf.Volumes,
			EnvVars:       conf.EnvVars,
			RestartPolicy: conf.RestartPolicy,
		},
		mu: &sync.RWMutex{},
	}
}

//-------------------------------------------------------------------------

func (db *SrvDB) SetStatus(status SrvStatus) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.logic.Status = status
}

func (db *SrvDB) GetLogic() *SrvLogic {
	db.mu.Lock()
	defer db.mu.Unlock()
	return db.logic
}

func (db *SrvDB) GetType() SrvType {
	db.mu.Lock()
	defer db.mu.Unlock()
	return db.logic.Type
}

func (db *SrvDB) ListImages() []string {
	return defaultDBImages
}

func (db *SrvDB) GetStatus() SrvStatus {
	db.mu.Lock()
	defer db.mu.Unlock()
	return db.logic.Status
}

func (db *SrvDB) Info() {
	db.mu.Lock()
	defer db.mu.Unlock()

	log.Println("----- DB Info -----")
	log.Printf("Service Type: %s", db.logic.Type)
	log.Printf("Current Image: %s", db.logic.Image)
	log.Printf("Container Name: %s", db.logic.ContainerName)
	log.Printf("Status: %s", db.logic.Status)
	log.Printf("Network: %s", db.logic.Network)
	log.Printf("Ports: %v", db.logic.Ports)
	log.Printf("Volumes: %v", db.logic.Volumes)
	log.Printf("Environment Variables: %v", db.logic.EnvVars)
	log.Printf("Restart Policy: %s", db.logic.RestartPolicy)
	log.Println("-------------------")
}

func (db *SrvDB) GetImage() string {
	db.mu.Lock()
	defer db.mu.Unlock()
	return db.logic.Image
}
