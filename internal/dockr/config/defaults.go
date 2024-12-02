package config

import "math/rand"

func DefaultConfigDB() ContainerConfig {
	return defaultConfigs()["DB"]
}

func DefaultConfigCache() ContainerConfig {
	return defaultConfigs()["Cache"]
}

func DefaultConfigServerMain() ContainerConfig {
	return defaultConfigs()["Server_main"]
}

func DefaultConfigServerAdd() ContainerConfig {
	return defaultConfigs()["Server_add"]
}

func DefaultConfigLB() ContainerConfig {
	return defaultConfigs()["LB"]
}

// DefaultConfigs returns default configurations for all container types.
func defaultConfigs() map[string]ContainerConfig {
	ind := randString(5)

	return map[string]ContainerConfig{
		"DB": {
			IsDefault:     true,
			ContainerType: "DB",
			Image:         "postgres:latest",
			ContainerName: "db_container_" + ind,
			Network:       "db-net",
			Ports:         []string{"5432:5432"},
			Volumes:       []string{"db_data_" + ind + ":/var/lib/postgresql/data"},
			EnvVars: map[string]string{
				"POSTGRES_USER":     "admin",
				"POSTGRES_PASSWORD": "password",
				"POSTGRES_DB":       "mydb",
			},
			RestartPolicy: "always",
		},
		"Cache": {
			IsDefault:     true,
			ContainerType: "Cache",
			Image:         "redis:latest",
			ContainerName: "cache_container_" + ind,
			Network:       "cache-net",
			Ports:         []string{"6379:6379"},
			Volumes:       []string{"cache_data_" + ind + ":/data"},
			EnvVars: map[string]string{
				"REDIS_PASSWORD": "secret",
			},
			RestartPolicy: "always",
		},
		"Server_main": {
			IsDefault:     true,
			ContainerType: "Server_main",
			Image:         "nginx:latest",
			ContainerName: "server_main_" + ind,
			Network:       "web-net",
			Ports:         []string{"80:80", "443:443"},
			Volumes:       []string{"nginx_config_" + ind + ":/etc/nginx", "nginx_logs_" + ind + ":/var/log/nginx"},
			EnvVars:       map[string]string{},
			RestartPolicy: "always",
		},
		"Server_add": {
			IsDefault:     true,
			ContainerType: "Server_add",
			Image:         "httpd:latest",
			ContainerName: "server_add_" + ind,
			Network:       "web-net",
			Ports:         []string{"8080:80"},
			Volumes:       []string{"httpd_data_" + ind + ":/usr/local/apache2/htdocs"},
			EnvVars:       map[string]string{},
			RestartPolicy: "always",
		},
		"LB": {
			IsDefault:     true,
			ContainerType: "LB",
			Image:         "haproxy:latest",
			ContainerName: "load_balancer_" + ind,
			Network:       "lb-net",
			Ports:         []string{"5000:5000", "5001:5001"},
			Volumes:       []string{"haproxy_config_" + ind + ":/usr/local/etc/haproxy", "haproxy_logs_" + ind + ":/var/log/haproxy"},
			EnvVars:       map[string]string{},
			RestartPolicy: "always",
		},
		"Other": {
			IsDefault:     true,
			ContainerType: "Other",
			Image:         "alpine:latest",
			ContainerName: "other_container_" + ind,
			Network:       "default-net",
			Ports:         []string{},
			Volumes:       []string{"other_data_" + ind + ":/data"},
			EnvVars:       map[string]string{},
			RestartPolicy: "on-failure",
		},
	}
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
