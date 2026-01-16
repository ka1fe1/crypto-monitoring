package main

import (
	"os"
	"strings"

	"github.com/ka1fe1/crypto-monitoring/config"
	"github.com/ka1fe1/crypto-monitoring/pkg/logger"
	"gopkg.in/yaml.v3"
)

func main() {
	configPath := config.GetConfigPath()
	tempPath := config.GetConfigTempPath()

	// Read config.yaml
	data, err := os.ReadFile(configPath)
	if err != nil {
		logger.Error("Failed to read config file: %v", err)
		return
	}

	var node yaml.Node
	if err := yaml.Unmarshal(data, &node); err != nil {
		logger.Error("Failed to unmarshal config: %v", err)
		return
	}

	maskSensitiveData(&node)

	// Write to config.yaml.temp
	tempData, err := yaml.Marshal(&node)
	if err != nil {
		logger.Error("Failed to marshal temp config: %v", err)
		return
	}

	if err := os.WriteFile(tempPath, tempData, 0644); err != nil {
		logger.Error("Failed to write config.yaml.temp: %v", err)
		return
	}

	logger.Info("Successfully generated config.yaml.temp")
}

func maskSensitiveData(node *yaml.Node) {
	if node.Kind == yaml.DocumentNode {
		for _, child := range node.Content {
			maskSensitiveData(child)
		}
		return
	}

	if node.Kind == yaml.MappingNode {
		for i := 0; i < len(node.Content); i += 2 {
			keyNode := node.Content[i]
			valueNode := node.Content[i+1]

			if isSensitive(keyNode.Value) {
				valueNode.Value = "YOUR_" + strings.ToUpper(keyNode.Value) + "_HERE"
			} else {
				maskSensitiveData(valueNode)
			}
		}
	}
}

func isSensitive(key string) bool {
	sensitiveKeys := []string{"api_key", "access_token", "secret"}
	for _, k := range sensitiveKeys {
		if k == key {
			return true
		}
	}
	return false
}
