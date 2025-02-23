package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"

	rbac_script "github.com/darkdoc/purple-storage-rh-operator/scripts/rbacs"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run create-rbacs.go <path-to-yaml>")
	}
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	var yamlContent bytes.Buffer
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		yamlContent.WriteString(scanner.Text() + "\n")
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading file: %v", err)
	}
	rules, err := rbac_script.ExtractRBACRules(yamlContent.Bytes())
	if err != nil {
		log.Fatalf("Failed to extract RBAC rules: %v", err)
	}
	markers := rbac_script.GenerateRBACMarkers(rules)
	for _, marker := range markers {
		fmt.Println(marker)
	}
}
