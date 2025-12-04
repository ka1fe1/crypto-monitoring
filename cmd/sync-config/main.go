package main

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"gopkg.in/yaml.v3"
)

func main() {
	_, filename, _, _ := runtime.Caller(0)
	rootDir := filepath.Dir(filepath.Dir(filepath.Dir(filename)))
	configPath := filepath.Join(rootDir, "config", "config.yaml")
	tempPath := filepath.Join(rootDir, "config", "config.yaml.temp")

	// Read config.yaml
	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	var node yaml.Node
	if err := yaml.Unmarshal(data, &node); err != nil {
		log.Fatalf("Failed to unmarshal config: %v", err)
	}

	maskSensitiveData(&node)

	// Write to config.yaml.temp
	tempData, err := yaml.Marshal(&node)
	if err != nil {
		log.Fatalf("Failed to marshal temp config: %v", err)
	}

	if err := os.WriteFile(tempPath, tempData, 0644); err != nil {
		log.Fatalf("Failed to write temp config: %v", err)
	}

	log.Println("Successfully generated config.yaml.temp")
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
