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

func convertToPlural(kind string) string {
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
			panic(fmt.Sprintf("Value '%s' already exists in the slice", value))
		}
	}
	return append(slice, value)
}

type Permission struct {
	apiVersion string
	kind       string
	group      string
	version    string
	resource   string
	name       string
	namespace  string
	rules      []interface{}
	verbs      map[string]bool
}

func NewPermission(raw map[string]interface{}) *Permission {
	obj := &unstructured.Unstructured{Object: raw}
	apiVersion := obj.GetAPIVersion()
	kind := obj.GetKind()
	name := obj.GetName()
	namespace := obj.GetNamespace()
	if kind == "" && apiVersion == "" {
		return nil
	}
	gVer, err := schema.ParseGroupVersion(apiVersion)
	if err != nil {
		panic("Could not parse group version")
	}
	resourceName := convertToPlural(kind)
	rules, _, err := unstructured.NestedSlice(raw, "rules")
	if err != nil {
		panic("Could not parse rules")
	}

	return &Permission{
		apiVersion: apiVersion,
		kind:       kind,
		group:      gVer.Group,
		version:    gVer.Version,
		resource:   resourceName,
		name:       name,
		namespace:  namespace,
		rules:      rules,
		verbs:      map[string]bool{},
	}
}

func (p Permission) isRole() bool {
	return p.kind == "Role" || p.kind == "ClusterRole"
}

func (p Permission) SetDefaultVerbs() {
	for _, i := range defaultVerbs {
		p.verbs[i] = true
	}
}

func (p Permission) String() string {
	return fmt.Sprintf("%s/%s/%s/%s/%s/%v", p.group, p.version, p.resource, p.namespace, p.name, p.verbs)
}

func (p Permission) RBACRuleFromRole() []string {
	if len(p.rules) == 0 {
		panic("No rules parsed on :" + p.String())
	}
	var ns string
	if p.kind == "Role" && p.namespace != "" {
		ns = fmt.Sprintf("namespace=%s,", p.namespace)
	}

	var ok bool
	var rbacs []string
	for _, rule := range p.rules {
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
	p.SetDefaultVerbs()
	verbs := []string{}
	for v := range p.verbs {
		verbs = append(verbs, v)
	}
	sort.Strings(verbs)
	var ns string
	if p.namespace != "" {
		ns = fmt.Sprintf("namespace=%s,", p.namespace)

	}
	var groupStr string
	if p.group == "" {
		groupStr = "\"\""
	} else {
		groupStr = p.group
	}
	s := fmt.Sprintf("//+kubebuilder:rbac:groups=%s,resources=%s,%sverbs=%s",
		groupStr, p.resource, ns, strings.Join(verbs, ";"))

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

func GenerateRBACMarkers(permissions []Permission) {
	uniques := map[string]bool{}
	for _, p := range permissions {
		rules := p.RBACRule()
		for _, s := range rules {
			if _, exists := uniques[s]; exists {
				continue
			}
			uniques[s] = true
			fmt.Println(s)
		}
	}
}
