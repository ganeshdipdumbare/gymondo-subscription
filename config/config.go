package config

import "github.com/ganeshdipdumbare/goenv"

type EnvVar struct {
	MongoUri           string `json:"mongo_uri"`
	MongoDb            string `json:"mongo_db"`
	Port               string `json:"port"`
	MigrationFilesPath string `json:"migration_files_path"`
}

var (
	envVars = &EnvVar{
		Port:               "8080",
		MongoDb:            "gymondodb",
		MigrationFilesPath: "file://migration",
	}
)

func init() {
	goenv.SyncEnvVar(&envVars)
}

func Get() *EnvVar {
	return envVars
}
