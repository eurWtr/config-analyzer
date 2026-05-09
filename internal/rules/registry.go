package rules

import (
	"config-analyzer/internal/models"
	"fmt"
)

// Rule defines the interface for a validation rule.
// To add a new rule implement this interface and register it.
type Rule interface {
	Name() string

	Check(config map[string]interface{}, filePath string) []models.Issue
}

// Registry holds registered rules.
type Registry struct {
	rules []Rule
}

// NewRegistry creates a registry with a set of default rules.
func NewRegistry() *Registry {
	r := &Registry{}

	r.Register(&DebugModeRule{})
	r.Register(&PlainPasswordRule{})
	r.Register(&BindAddressRule{})
	r.Register(&TLSDisabledRule{})
	r.Register(&WeakAlgorithmRule{})
	r.Register(&FilePermissionsRule{})

	return r
}

// Register adds a rule to the registry.
func (r *Registry) Register(rule Rule) {
	r.rules = append(r.rules, rule)
}

// CheckAll applies all registered rules to the configuration.
func (r *Registry) CheckAll(config map[string]interface{}, filePath string) []models.Issue {
	var allIssues []models.Issue
	for _, rule := range r.rules {
		issues := rule.Check(config, filePath)
		allIssues = append(allIssues, issues...)
	}
	return allIssues
}

// Rules returns the list of registered rules.
func (r *Registry) Rules() []Rule {
	return r.rules
}

// flatten recursively walks a nested map and returns a flat map
// with keys like "parent.child.key".
func flatten(prefix string, m map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range m {
		fullKey := k
		if prefix != "" {
			fullKey = prefix + "." + k
		}
		switch val := v.(type) {
		case map[string]interface{}:
			for fk, fv := range flatten(fullKey, val) {
				result[fk] = fv
			}
		case map[interface{}]interface{}:
			converted := make(map[string]interface{})
			for mk, mv := range val {
				converted[fmt.Sprintf("%v", mk)] = mv
			}
			for fk, fv := range flatten(fullKey, converted) {
				result[fk] = fv
			}
		default:
			result[fullKey] = v
		}
	}
	return result
}
