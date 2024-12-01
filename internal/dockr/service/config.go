package service

// SrvType represent a type of service for docker container,
// it can be:
// DB || Cache || ServerMain || ServerAdd || LB
type SrvType string

const (
	DB         SrvType = "service_db"
	Cache      SrvType = "service_cache"
	ServerMain SrvType = "service_server_main"
	ServerAdd  SrvType = "service_server_add"
	LB         SrvType = "service_lb"
)

// SrvStatus represent a status of service,
// it can be running / stopped
type SrvStatus string

func SrvStatusRun() SrvStatus  { return "running" }
func SrvStatusStop() SrvStatus { return "stopped" }

// SrvConfig is a representation of docker container config it UNI for all
type SrvConfig struct {
	IsDefault     bool              // default or manual config
	Type          SrvType           // type of service
	Image         string            // docker Dim
	Status        SrvStatus         // status of container
	ContainerName string            // docker container name
	Network       string            // docker network
	Ports         []string          // docker ports
	Volumes       []string          // docker volume
	EnvVars       map[string]string // docker env variables
	RestartPolicy string            // docker restart policy
}
