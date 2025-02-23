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

var defaultVerbs = []string{"get", "list", "watch", "create", "update", "patch", "delete"}

func ConvertToPlural(kind string) string {
	if kind == "" {
		return kind
	}
	if strings.HasSuffix(kind, "s") {
		return kind
	}
	return kind + "s"
}

// AddStringUnique adds a string to a slice but panics if it already exists
func AddStringUnique(slice []string, value string) []string {
	for _, v := range slice {
		if v == value {
			return slice
		}
	}
	return append(slice, value)
}

type Permission struct {
	ApiVersion string
	Kind       string
	Group      string
	Version    string
	Resource   string
	Name       string
	Namespace  string
	Rules      []interface{}
	Verbs      map[string]bool
}

func NewPermission(raw map[string]interface{}) *Permission {
	obj := &unstructured.Unstructured{Object: raw}
	apiVersion := obj.GetAPIVersion()
	kind := strings.ToLower(obj.GetKind())
	name := obj.GetName()
	namespace := obj.GetNamespace()
	if kind == "" && apiVersion == "" {
		return nil
	}
	gVer, err := schema.ParseGroupVersion(apiVersion)
	if err != nil {
		panic("Could not parse group version")
	}
	resourceName := ConvertToPlural(kind)
	rules, _, err := unstructured.NestedSlice(raw, "rules")
	if err != nil {
		panic("Could not parse rules")
	}

	return &Permission{
		ApiVersion: apiVersion,
		Kind:       kind,
		Group:      gVer.Group,
		Version:    gVer.Version,
		Resource:   resourceName,
		Name:       name,
		Namespace:  namespace,
		Rules:      rules,
		Verbs:      map[string]bool{},
	}
}

func (p Permission) isRole() bool {
	return p.Kind == "role" || p.Kind == "clusterrole"
}

func (p Permission) String() string {
	return fmt.Sprintf("%s/%s/%s/%s/%s/%v", p.Group, p.Version, p.Resource, p.Namespace, p.Name, p.Verbs)
}

func (p Permission) RBACRuleFromRole() []string {
	if len(p.Rules) == 0 {
		panic("No rules parsed on :" + p.String())
	}
	var ns string
	if p.Kind == "role" && p.Namespace != "" {
		ns = fmt.Sprintf("namespace=%s,", p.Namespace)
	}

	var ok bool
	var rbacs []string
	for _, rule := range p.Rules {
		var ruleMap map[string]interface{}
		if ruleMap, ok = rule.(map[string]interface{}); !ok {
			panic("Could not parse rule")
		}
		apiGroups, ok := ruleMap["apiGroups"].([]interface{})
		if !ok {
			panic("Could not parse apiGroups")
		}
		resourcesList, ok := ruleMap["resources"].([]interface{})
		if !ok {
			panic("Could not parse resourcesList")
		}
		verbsList, ok := ruleMap["verbs"].([]interface{})
		if !ok {
			panic("Could not parse verbsList")
		}
		for _, group := range apiGroups {
			groupStr, _ := group.(string)
			if groupStr == "" {
				groupStr = "\"\""
			}
			for _, res := range resourcesList {
				resStr, _ := res.(string)
				var verbsArray []string
				for _, v := range verbsList {
					if verb, ok := v.(string); ok {
						verbsArray = AddStringUnique(verbsArray, verb)
					} else {
						panic("We could not parse a verb as string")
					}
				}
				sort.Strings(verbsArray)
				rbac := fmt.Sprintf("//+kubebuilder:rbac:groups=%s,%sresources=%s,verbs=%s", groupStr, ns, resStr, strings.Join(verbsArray, ";"))
				rbacs = append(rbacs, rbac)
			}
		}
	}
	return rbacs
}

func (p Permission) RBACRule() []string {
	if p.isRole() {
		return p.RBACRuleFromRole()
	}
	verbs := []string{}
	verbs = append(verbs, defaultVerbs...)
	sort.Strings(verbs)
	var ns string
	if p.Namespace != "" {
		ns = fmt.Sprintf("namespace=%s,", p.Namespace)

	}
	var groupStr string
	if p.Group == "" {
		groupStr = "\"\""
	} else {
		groupStr = p.Group
	}
	s := fmt.Sprintf("//+kubebuilder:rbac:groups=%s,%sresources=%s,verbs=%s",
		groupStr, ns, p.Resource, strings.Join(verbs, ";"))

	return strings.Fields(s)
}

func ExtractRBACRules(yamlContent []byte) ([]Permission, error) {
	// We create a map of GroupVersionResource to verbs where the verbs is a map just to simulate a set really
	resources := []Permission{}

	decoder := yaml.NewDecoder(bytes.NewReader(yamlContent))
	for {
		var raw map[string]interface{}
		if err := decoder.Decode(&raw); err != nil { // End of file
			break
		}

		p := NewPermission(raw)
		if p == nil { // Skip if empty piece of yaml
			continue
		}

		resources = append(resources, *p)
	}

	return resources, nil
}

func GenerateRBACMarkers(permissions []Permission) []string {
	uniques := []string{}
	for _, p := range permissions {
		rules := p.RBACRule()
		for _, s := range rules {
			uniques = AddStringUnique(uniques, s)
		}
	}
	return uniques
}
