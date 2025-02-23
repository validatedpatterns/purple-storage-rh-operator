package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func convertToPlural(kind string) string {
	if strings.HasSuffix(kind, "s") {
		return kind
	}
	return kind + "s"
}

func extractRBACRules(filePath string) (map[schema.GroupVersionResource][]string, error) {
	file, err := os.Open(filePath)
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
	resources := make(map[schema.GroupVersionResource][]string)

	decoder := yaml.NewDecoder(bytes.NewReader(yamlContent.Bytes()))
	for {
		var raw map[string]interface{}
		if err := decoder.Decode(&raw); err != nil { // End of file
			break
		}

		// Convert into an Unstructured Kubernetes object
		obj := &unstructured.Unstructured{Object: raw}
		apiVersion := obj.GetAPIVersion()
		kind := obj.GetKind()
		namespace := obj.GetNamespace()

		// Parse group and version
		gv, err := schema.ParseGroupVersion(apiVersion)
		if err != nil {
			return nil, err
		}
		resourceName := convertToPlural(kind)
		verbs := []string{"get", "list", "watch", "create", "update", "patch", "delete"}
		if namespace == "" {
			verbs = append(verbs, "deletecollection")
		}

		// Special case: If the object is a Role or ClusterRole, extract its rules
		if kind == "Role" || kind == "ClusterRole" {
			rules, found, _ := unstructured.NestedSlice(obj.Object, "rules")
			if found {
				for _, rule := range rules {
					if ruleMap, ok := rule.(map[string]interface{}); ok {
						apiGroups, _ := ruleMap["apiGroups"].([]interface{})
						resourcesList, _ := ruleMap["resources"].([]interface{})
						verbsList, _ := ruleMap["verbs"].([]interface{})

						for _, res := range resourcesList {
							resStr, _ := res.(string)
							for _, group := range apiGroups {
								groupStr, _ := group.(string)
								gvr := schema.GroupVersionResource{
									Group:    groupStr,
									Version:  gv.Version,
									Resource: resStr,
								}

								// Convert interface{} to []string for verbs
								var extractedVerbs []string
								for _, v := range verbsList {
									if verb, ok := v.(string); ok {
										extractedVerbs = append(extractedVerbs, verb)
									}
								}

								// Store rules
								if existing, exists := resources[gvr]; exists {
									resources[gvr] = append(existing, extractedVerbs...)
								} else {
									resources[gvr] = extractedVerbs
								}
							}
						}
					}
				}
			}

			// Ensure operator can write Roles & ClusterRoles
			gvr := schema.GroupVersionResource{
				Group:    gv.Group,
				Version:  gv.Version,
				Resource: resourceName,
			}
			resources[gvr] = []string{"get", "list", "watch", "create", "update", "patch", "delete"}
		} else {
			// Handle normal resources
			gvr := schema.GroupVersionResource{
				Group:    gv.Group,
				Version:  gv.Version,
				Resource: resourceName,
			}
			resources[gvr] = verbs
		}
	}

	return resources, nil
}

func generateRBACMarkers(rules map[schema.GroupVersionResource][]string) {
	for gvr, verbs := range rules {
		group := gvr.Group
		if group == "" {
			group = "core"
		}
		fmt.Printf("// +kubebuilder:rbac:groups=%s,resources=%s,verbs=%s\n",
			group, gvr.Resource, strings.Join(verbs, ","))
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run create-rbacs.go <path-to-yaml>")
	}
	filePath := os.Args[1]

	rules, err := extractRBACRules(filePath)
	if err != nil {
		log.Fatalf("Failed to extract RBAC rules: %v", err)
	}
	generateRBACMarkers(rules)
}
