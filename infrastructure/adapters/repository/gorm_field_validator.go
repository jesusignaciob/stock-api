package repository

import (
	"reflect"
	"strings"
	"sync"
)

type GormFieldValidator struct {
	model     interface{}
	cache     map[string]bool
	cacheLock sync.RWMutex
}

func NewGormFieldValidator(model interface{}) *GormFieldValidator {
	return &GormFieldValidator{
		model: model,
		cache: make(map[string]bool),
	}
}

// IsValidField checks if the given field is valid by first looking it up in a cache.
// If the field's validity is not cached, it performs a validation check and updates the cache.
// This method is thread-safe as it uses read-write locks to manage concurrent access to the cache.
//
// Parameters:
//   - field: The name of the field to validate.
//
// Returns:
//   - bool: True if the field is valid, false otherwise.
func (v *GormFieldValidator) IsValidField(field string) bool {
	v.cacheLock.RLock()
	if isValid, exists := v.cache[field]; exists {
		v.cacheLock.RUnlock()
		return isValid
	}
	v.cacheLock.RUnlock()

	isValid := v.checkField(field)

	v.cacheLock.Lock()
	v.cache[field] = isValid
	v.cacheLock.Unlock()

	return isValid
}

// GetAllValidFields retrieves all valid fields from the cache.
// It acquires a read lock to ensure thread-safe access to the cache.
// If the cache is empty, it returns nil. Otherwise, it iterates through
// the cache and collects all fields marked as valid into a slice, which
// is then returned.
//
// Returns:
//
//	[]string - A slice containing all valid field names, or nil if the cache is empty.
func (v *GormFieldValidator) GetAllValidFields() []string {
	v.cacheLock.RLock()
	defer v.cacheLock.RUnlock()

	if len(v.cache) == 0 {
		return nil
	}

	validFields := make([]string, 0, len(v.cache))
	for field := range v.cache {
		if v.cache[field] {
			validFields = append(validFields, field)
		}
	}

	return validFields
}

// checkField checks if a given field exists in the model associated with the GormFieldValidator.
// It verifies the presence of the field by inspecting the struct's fields and their "gorm" tags.
//
// Parameters:
//   - field: The name of the field to check.
//
// Returns:
//   - bool: True if the field exists in the model, either as a struct field name
//     or as a column name specified in the "gorm" tag; otherwise, false.
func (v *GormFieldValidator) checkField(field string) bool {
	modelType := reflect.TypeOf(v.model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	for i := 0; i < modelType.NumField(); i++ {
		fieldType := modelType.Field(i)
		gormTag := fieldType.Tag.Get("gorm")

		if gormTag != "" {
			for _, part := range strings.Split(gormTag, ";") {
				if strings.HasPrefix(part, "column:") {
					colName := strings.TrimPrefix(part, "column:")
					if colName == field {
						return true
					}
				}
			}
		}

		if strings.EqualFold(fieldType.Name, field) {
			return true
		}
	}

	return false
}
