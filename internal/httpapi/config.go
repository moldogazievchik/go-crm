package httpapi

import (
	"os"
	"strings"
)

type Config struct {
	DataPath    string
	CorsOrigins []string
}

func LoadConfig() Config {
	dataPath := strings.TrimSpace(os.Getenv("DATA_PATH"))
	if dataPath == "" {
		dataPath = "./data.json"
	}

	originsEnv := strings.TrimSpace(os.Getenv("CORS_ORIGINS"))
	origins := []string{"http://localhost:5173", "http://localhost:3000"}
	if originsEnv != "" {
		parts := strings.Split(originsEnv, ",")
		origins = origins[:0]
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p != "" {
				origins = append(origins, p)
			}
		}
		if len(origins) == 0 {
			origins = []string{"http://localhost:5173", "http://localhost:3000"}
		}
	}

	return Config{
		DataPath:    dataPath,
		CorsOrigins: origins,
	}
}
