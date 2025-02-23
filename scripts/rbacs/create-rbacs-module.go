package rbac_script

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type StringSet map[string]struct{}

func NewStringSet() StringSet {
	return make(StringSet)
}

func NewStringSetFromList(items []string) StringSet {
	set := make(StringSet)
	for _, item := range items {
		set[item] = struct{}{}
	}
	return set
}

func (s StringSet) Add(value string) {
	s[value] = struct{}{}
}

func (s StringSet) Remove(value string) {
	delete(s, value)
}

func (s StringSet) Contains(value string) bool {
	_, exists := s[value]
	return exists
}

func (s StringSet) List() []string {
	keys := make([]string, 0, len(s))
	for key := range s {
		keys = append(keys, key)
	}
	return keys
}

func (s StringSet) Equals(other StringSet) bool {
	if len(s) != len(other) {
		return false
	}

	for key := range s {
		if _, exists := other[key]; !exists {
			return false
		}
	}

	return true
}

func (s StringSet) SortedList() []string {
	keys := s.List()
	sort.Strings(keys)
	return keys
}

func convertToPlural(kind string) string {
	if kind == "" {
		return kind
	}
	if strings.HasSuffix(kind, "s") {
		return kind
	}
	return kind + "s"
}

func ExtractRBACRules(yamlContent []byte) (map[schema.GroupVersionResource]StringSet, error) {
	// We create a map of GroupVersionResource to verbs where the verbs is a map just to simulate a set really
	resources := make(map[schema.GroupVersionResource]StringSet)

	decoder := yaml.NewDecoder(bytes.NewReader(yamlContent))
	for {
		var raw map[string]interface{}
		if err := decoder.Decode(&raw); err != nil { // End of file
			break
		}

		// Convert into an Unstructured Kubernetes object
		obj := &unstructured.Unstructured{Object: raw}
		apiVersion := obj.GetAPIVersion()
		kind := strings.ToLower(obj.GetKind())
		// namespace := obj.GetNamespace()

		if kind == "" && apiVersion == "" {
			continue
		}

		// Parse group and version
		gVer, err := schema.ParseGroupVersion(apiVersion)
		if err != nil {
			return nil, err
		}
		resourceName := convertToPlural(kind)
		gv := schema.GroupVersionResource{
			Group:   gVer.Group,
			Version: gVer.Version,
		}

		defaultVerbs := NewStringSetFromList([]string{"get", "list", "watch", "create", "update", "patch", "delete"})
		// if namespace == "" {
		// 	defaultVerbs.Add("deletecollection")
		// }

		// Special case: If the object is a Role or ClusterRole, extract its rules and
		// add them to the resources map
		if kind == "Role" || kind == "ClusterRole" {
			rules, found, _ := unstructured.NestedSlice(obj.Object, "rules")
			if !found {
				continue
			}
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

							for _, v := range verbsList {
								if verb, ok := v.(string); ok {
									resources[gvr].Add(verb)
								}
							}
						}
					}
				}
			}

		} else {
			gvr := schema.GroupVersionResource{
				Group:    gv.Group,
				Version:  gv.Version,
				Resource: resourceName,
			}
			resources[gvr] = defaultVerbs
		}
	}

	return resources, nil
}

func GenerateRBACMarkers(rules map[schema.GroupVersionResource]StringSet) []string {
	s := []string{}
	for gvr, verbs := range rules {
		group := gvr.Group
		if group == "" {
			group = "core"
		}
		s = append(s, fmt.Sprintf("//+kubebuilder:rbac:groups=%s,resources=%s,verbs=%s\n",
			group, gvr.Resource, strings.Join(verbs.SortedList(), ",")))
	}
	return s
}
