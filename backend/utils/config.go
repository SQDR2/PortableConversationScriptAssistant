package utils

import (
	"encoding/json"
	"os"
)

type TencentCloudConfig struct {
	SecretId  string `json:"secret_id"`
	SecretKey string `json:"secret_key"`
	Region    string `json:"region"`
	ProjectId int64  `json:"project_id"`
}

type AppConfig struct {
	TencentCloud TencentCloudConfig `json:"tencent_cloud"`
}

func LoadConfig(path string) (*AppConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := &AppConfig{}
	err = decoder.Decode(config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
