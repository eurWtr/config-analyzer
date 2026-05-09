package rules

import (
	"config-analyzer/internal/models"
	"fmt"
)

// Rule определяет интерфейс правила проверки.
// Для добавления нового правила достаточно реализовать этот интерфейс и зарегистрировать его.
type Rule interface {
	Name() string

	Check(config map[string]interface{}, filePath string) []models.Issue
}

// Registry хранит зарегистрированные правила.
type Registry struct {
	rules []Rule
}

// NewRegistry создаёт реестр с набором правил по умолчанию.
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

// Register добавляет правило в реестр.
func (r *Registry) Register(rule Rule) {
	r.rules = append(r.rules, rule)
}

// CheckAll применяет все зарегистрированные правила к конфигурации.
func (r *Registry) CheckAll(config map[string]interface{}, filePath string) []models.Issue {
	var allIssues []models.Issue
	for _, rule := range r.rules {
		issues := rule.Check(config, filePath)
		allIssues = append(allIssues, issues...)
	}
	return allIssues
}

// Rules возвращает список зарегистрированных правил.
func (r *Registry) Rules() []Rule {
	return r.rules
}

// flatten рекурсивно обходит вложенную map и возвращает плоскую map
// с ключами вида "parent.child.key".
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
