package config

import (
	"crypto/rand"

	"github.com/docker/docker/api/types/container"
)

var PostgresConfig = ContainerConfig{
		LoadLevel:        2,
		IsDefault:        true,
		ContainerService: "DB",
		Image:            "postgres:latest",
		Hostname:         "postgres-db",
		EnvVars: map[string]string{
			"POSTGRES_USER":     "admin",
			"POSTGRES_PASSWORD": "securepassword",
			"POSTGRES_DB":       "exampledb",
		},
		WorkingDir:    "/var/lib/postgresql/data",
		Cmd:           []string{"postgres"},
		Volumes:       []string{"/host/data/postgres:/var/lib/postgresql/data"},
		NetworkMode:   "bridge",
		Ports:         []string{"5432:5432"},
		RestartPolicy: "always",
		NetworkID:     "db-network",
		HealthCheck: HealthCheckConfig{
			Interval:    "30s",
			Timeout:     "10s",
			StartPeriod: "10s",
			Retries:     3,
			Test:        []string{"CMD", "pg_isready", "-U", "admin"},
		},
	}

	var MongoConfig = ContainerConfig{
			LoadLevel:        2,
			IsDefault:        true,
			ContainerService: "DB",
			Image:            "mongo:latest",
			Hostname:         "mongo-db",
			EnvVars: map[string]string{
				"MONGO_INITDB_ROOT_USERNAME": "admin",
				"MONGO_INITDB_ROOT_PASSWORD": "securepassword",
			},
			WorkingDir:    "/data/db",
			Cmd:           []string{"mongod"},
			Volumes:       []string{"/host/data/mongo:/data/db"},
			NetworkMode:   "bridge",
			Ports:         []string{"27017:27017"},
			RestartPolicy: "always",
			NetworkID:     "db-network",
			HealthCheck: HealthCheckConfig{
				Interval:    "30s",
				Timeout:     "10s",
				StartPeriod: "10s",
				Retries:     3,
				Test:        []string{"CMD", "mongo", "--eval", "db.adminCommand('ping')"},
			},
		}

	
	// Default Redis configuration
var	RedisConfig = ContainerConfig{
		LoadLevel:        1,
		IsDefault:        true,
		ContainerService: "Cache",
		Image:            "redis:latest",
		Hostname:         "redis-cache",
		EnvVars:          map[string]string{},
		WorkingDir:       "",
		Cmd:              []string{"redis-server"},
		Volumes:          []string{"/host/data/redis:/data"},
		NetworkMode:      "bridge",
		Ports:            []string{"6379:6379"},
		RestartPolicy:    "always",
		NetworkID:        "cache-network",
		HealthCheck: HealthCheckConfig{
			Interval:    "30s",
			Timeout:     "10s",
			StartPeriod: "10s",
			Retries:     3,
			Test:        []string{"CMD", "redis-cli", "ping"},
		},
	}

		// Default Nginx load balancer configuration
		var NginxConfig = ContainerConfig{
			LoadLevel:        1,
			IsDefault:        true,
			ContainerService: "LB",
			Image:            "nginx:latest",
			Hostname:         "nginx-lb",
			EnvVars:          map[string]string{},
			WorkingDir:       "/etc/nginx",
			Cmd:              []string{"nginx", "-g", "daemon off;"},
			Volumes:          []string{"/host/config/nginx:/etc/nginx"},
			NetworkMode:      "bridge",
			Ports:            []string{"80:80", "443:443"},
			RestartPolicy:    "always",
			NetworkID:        "lb-network",
			HealthCheck: HealthCheckConfig{
				Interval:    "30s",
				Timeout:     "10s",
				StartPeriod: "10s",
				Retries:     3,
				Test:        []string{"CMD", "curl", "-f", "http://localhost"},
			},
		}

		// Default HAProxy load balancer configuration
		var HaproxyConfig = ContainerConfig{
			LoadLevel:        1,
			IsDefault:        true,
			ContainerService: "LB",
			Image:            "haproxy:latest",
			Hostname:         "haproxy-lb",
			EnvVars:          map[string]string{},
			WorkingDir:       "/usr/local/etc/haproxy",
			Cmd:              []string{"haproxy", "-f", "/usr/local/etc/haproxy/haproxy.cfg"},
			Volumes:          []string{"/host/config/haproxy:/usr/local/etc/haproxy"},
			NetworkMode:      "bridge",
			Ports:            []string{"8080:8080", "8443:8443"},
			RestartPolicy:    "always",
			NetworkID:        "lb-network",
			HealthCheck: HealthCheckConfig{
				Interval:    "30s",
				Timeout:     "10s",
				StartPeriod: "10s",
				Retries:     3,
				Test:        []string{"CMD", "curl", "-f", "http://localhost:8080"},
			},
		}


		
		var VoipConfig1 = ContainerConfig{
				LoadLevel:        2,
				IsDefault:        true,
				ContainerService: "Server_main",
				Image:            "mumble:latest", // Наприклад, Mumble сервер
				Hostname:         "voip-primary",
				EnvVars: map[string]string{
					"MUMBLE_PORT":    "64738",
					"MUMBLE_MAXUSERS": "100",
				},
				WorkingDir:    "/etc/mumble",
				Cmd:           []string{"murmurd", "-ini", "/etc/mumble/mumble.ini"},
				Volumes:       []string{"/host/config/mumble:/etc/mumble"},
				NetworkMode:   "bridge",
				Ports:         []string{"64738:64738", "64738:64738/udp"},
				RestartPolicy: "always",
				NetworkID:     "voice-network",
				HealthCheck: HealthCheckConfig{
					Interval:    "30s",
					Timeout:     "10s",
					StartPeriod: "10s",
					Retries:     3,
					Test:        []string{"CMD", "nc", "-z", "localhost", "64738"},
				},
			}

			var VoipConfig2 = ContainerConfig{
				LoadLevel:        1,
				IsDefault:        true,
				ContainerService: "Server_main",
				Image:            "teamspeak:latest",
				Hostname:         "voip-secondary",
				EnvVars: map[string]string{
					"TS3SERVER_LICENSE": "accept",
				},
				WorkingDir:    "/var/ts3server",
				Cmd:           []string{"ts3server"},
				Volumes:       []string{"/host/config/ts3server:/var/ts3server"},
				NetworkMode:   "bridge",
				Ports:         []string{"9987:9987/udp", "30033:30033", "10011:10011"},
				RestartPolicy: "on-failure",
				NetworkID:     "voice-network",
				HealthCheck: HealthCheckConfig{
					Interval:    "30s",
					Timeout:     "10s",
					StartPeriod: "10s",
					Retries:     3,
					Test:        []string{"CMD", "nc", "-z", "localhost", "9987"},
				},
			}

			// API Gateway configuration
			var ApiGatewayConfig = ContainerConfig{
				LoadLevel:        2,
				IsDefault:        true,
				ContainerService: "Server_add",
				Image:            "kong:latest",
				Hostname:         "api-gateway",
				EnvVars: map[string]string{
					"KONG_DATABASE":  "off",
					"KONG_PROXY_LISTEN": "0.0.0.0:8000",
				},
				WorkingDir:    "",
				Cmd:           []string{"kong", "start"},
				Volumes:       []string{"/host/config/kong:/etc/kong"},
				NetworkMode:   "bridge",
				Ports:         []string{"8000:8000", "8443:8443"},
				RestartPolicy: "always",
				NetworkID:     "api-network",
				HealthCheck: HealthCheckConfig{
					Interval:    "30s",
					Timeout:     "10s",
					StartPeriod: "10s",
					Retries:     3,
					Test:        []string{"CMD", "curl", "-f", "http://localhost:8000"},
				},
			}

			// Monitoring Service configuration
		var	MonitoringConfig = ContainerConfig{
				LoadLevel:        1,
				IsDefault:        true,
				ContainerService: "Other",
				Image:            "prom/prometheus:latest",
				Hostname:         "monitoring-server",
				EnvVars:          map[string]string{},
				WorkingDir:       "/etc/prometheus",
				Cmd:              []string{"prometheus", "--config.file=/etc/prometheus/prometheus.yml"},
				Volumes:          []string{"/host/config/prometheus:/etc/prometheus"},
				NetworkMode:      "bridge",
				Ports:            []string{"9090:9090"},
				RestartPolicy:    "always",
				NetworkID:        "monitoring-network",
				HealthCheck: HealthCheckConfig{
					Interval:    "30s",
					Timeout:     "10s",
					StartPeriod: "10s",
					Retries:     3,
					Test:        []string{"CMD", "curl", "-f", "http://localhost:9090"},
				},
			}



// DEFAULT RESOURCES #############################################################################

var HighLoadConfig = container.Resources{
    CPUShares:          1024,    // Higher CPU share for more processing power
    Memory:             4 * 1024 * 1024 * 1024,  // 4GB of memory
    NanoCPUs:           2000000000,  // 2 CPUs in nano units
    CPUPeriod:          100000,  // High CPU period
    CPUQuota:           200000,  // High CPU quota
    MemoryReservation:  2 * 1024 * 1024 * 1024,  // 2GB reserved memory
    MemorySwap:         4 * 1024 * 1024 * 1024,  // Total of 4GB memory + swap
    KernelMemory:       2 * 1024 * 1024 * 1024,  // 2GB kernel memory
    BlkioWeight:        1000,    // High block IO weight
    CPURealtimePeriod:  100000,  // Real-time CPU period
    CPURealtimeRuntime: 100000,  // Real-time CPU runtime
    OomKillDisable:     new(bool), // Disable OOM Killer
    PidsLimit:          new(int64), // Unlimited PIDs
}


var MediumLoadConfig = container.Resources{
    CPUShares:          512,     // Moderate CPU share
    Memory:             2 * 1024 * 1024 * 1024,  // 2GB memory
    NanoCPUs:           1000000000,  // 1 CPU in nano units
    CPUPeriod:          100000,  // Moderate CPU period
    CPUQuota:           100000,  // Moderate CPU quota
    MemoryReservation:  1 * 1024 * 1024 * 1024,  // 1GB reserved memory
    MemorySwap:         2 * 1024 * 1024 * 1024,  // Total of 2GB memory + swap
    KernelMemory:       1 * 1024 * 1024 * 1024,  // 1GB kernel memory
    BlkioWeight:        500,     // Moderate block IO weight
    CPURealtimePeriod:  100000,  // Real-time CPU period
    CPURealtimeRuntime: 100000,  // Real-time CPU runtime
    OomKillDisable:     nil,     // Default OOM Killer behavior
    PidsLimit:          new(int64), // Unlimited PIDs
}

var LowLoadConfig = container.Resources{
    CPUShares:          256,     // Lower CPU share
    Memory:             1 * 1024 * 1024 * 1024,  // 1GB memory
    NanoCPUs:           500000000,  // 0.5 CPU in nano units
    CPUPeriod:          100000,  // Lower CPU period
    CPUQuota:           50000,   // Lower CPU quota
    MemoryReservation:  512 * 1024 * 1024,  // 512MB reserved memory
    MemorySwap:         1 * 1024 * 1024 * 1024,  // Total of 1GB memory + swap
    KernelMemory:       512 * 1024 * 1024,  // 512MB kernel memory
    BlkioWeight:        300,     // Lower block IO weight
    CPURealtimePeriod:  100000,  // Real-time CPU period
    CPURealtimeRuntime: 50000,   // Real-time CPU runtime
    OomKillDisable:     nil,     // Default OOM Killer behavior
    PidsLimit:          new(int64), // Unlimited PIDs
}

func randString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	_, err := rand.Read(b)
	
	if err != nil {
		panic(err)
	}
	for i, v := range b {
		b[i] = letters[v%byte(len(letters))]
	}
	return string(b)
}
