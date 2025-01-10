package http

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/goravel/framework/contracts/foundation"
	contractstesting "github.com/goravel/framework/contracts/testing"
	"github.com/goravel/framework/support/maps"
)

type AssertableJson struct {
	t       *testing.T
	json    foundation.Json
	jsonStr string
	decoded map[string]any
}

func NewAssertableJSON(t *testing.T, json foundation.Json, jsonStr string) (contractstesting.AssertableJSON, error) {
	var decoded map[string]any
	err := json.Unmarshal([]byte(jsonStr), &decoded)
	if err != nil {
		return nil, err
	}

	return &AssertableJson{
		t:       t,
		json:    json,
		jsonStr: jsonStr,
		decoded: decoded,
	}, nil
}

func (r *AssertableJson) Json() map[string]any {
	return r.getDecoded()
}

func (r *AssertableJson) Count(key string, length int) contractstesting.AssertableJSON {
	actual := maps.Get(r.getDecoded(), key)
	assert.Len(r.t, actual, length, fmt.Sprintf("Property [%s] does not have the expected size.", key))

	return r
}

func (r *AssertableJson) Has(key string) contractstesting.AssertableJSON {
	exists := maps.Has(r.getDecoded(), key)
	assert.True(r.t, exists, fmt.Sprintf("Property [%s] does not exist.", key))

	return r
}

func (r *AssertableJson) HasAll(keys []string) contractstesting.AssertableJSON {
	for _, key := range keys {
		r.Has(key)
	}

	return r
}

func (r *AssertableJson) HasAny(keys []string) contractstesting.AssertableJSON {
	assert.True(r.t, maps.HasAny(r.getDecoded(), keys...), fmt.Sprintf("None of properties %v exist.", keys))

	return r
}

func (r *AssertableJson) Missing(key string) contractstesting.AssertableJSON {
	assert.False(r.t, maps.Has(r.getDecoded(), key), fmt.Sprintf("Property [%s] was found while it was expected to be missing.", key))

	return r
}

func (r *AssertableJson) MissingAll(keys []string) contractstesting.AssertableJSON {
	for _, key := range keys {
		r.Missing(key)
	}

	return r
}

func (r *AssertableJson) Where(key string, value any) contractstesting.AssertableJSON {
	r.Has(key)

	actual := maps.Get(r.getDecoded(), key)

	assert.Equal(r.t, value, actual, fmt.Sprintf("Expected property [%s] to have value [%v], but got [%v].", key, value, actual))

	return r
}

func (r *AssertableJson) WhereNot(key string, value any) contractstesting.AssertableJSON {
	r.Has(key)

	actual := maps.Get(r.getDecoded(), key)
	assert.NotEqual(r.t, value, actual, fmt.Sprintf("Expected property [%s] to not have value [%v], but it did.", key, value))
	return r
}

func (r *AssertableJson) First(key string, callback func(contractstesting.AssertableJSON)) contractstesting.AssertableJSON {
	value, exists := r.getDecoded()[key]
	if !assert.True(r.t, exists, fmt.Sprintf("Property [%s] does not exist.", key)) {
		return r
	}

	array, ok := value.([]any)
	if !assert.True(r.t, ok, fmt.Sprintf("Expected key [%s] to hold an array, but got %T", key, value)) {
		return r
	}

	if len(array) == 0 {
		assert.Fail(r.t, fmt.Sprintf("Expected a non-empty array for key [%s].", key))
		return r
	}

	firstItem := array[0]
	itemJson, err := r.json.Marshal(firstItem)
	if assert.NoError(r.t, err, "Failed to marshal the first item") {
		newJson, err := NewAssertableJSON(r.t, r.json, string(itemJson))
		if assert.NoError(r.t, err, "Failed to create AssertableJSON for first item") {
			callback(newJson)
		}
	}

	return r
}

func (r *AssertableJson) HasWithScope(key string, length int, callback func(contractstesting.AssertableJSON)) contractstesting.AssertableJSON {
	value, exists := r.getDecoded()[key]
	if !assert.True(r.t, exists, fmt.Sprintf("Property [%s] does not exist.", key)) {
		return r
	}

	array, ok := value.([]any)
	if !assert.True(r.t, ok, fmt.Sprintf("Expected key [%s] to hold an array, but got %T", key, value)) {
		return r
	}

	if !assert.Len(r.t, array, length, fmt.Sprintf("Property [%s] does not have the expected length of %d.", key, length)) {
		return r
	}

	if len(array) > 0 {
		itemJson, err := r.json.Marshal(array[0])
		if !assert.NoError(r.t, err, "Failed to marshal the first item of array") {
			return r
		}

		newJson, err := NewAssertableJSON(r.t, r.json, string(itemJson))
		if !assert.NoError(r.t, err, "Failed to create AssertableJSON for first item in scoped array") {
			return r
		}

		callback(newJson)
	}

	return r
}

func (r *AssertableJson) Each(key string, callback func(contractstesting.AssertableJSON)) contractstesting.AssertableJSON {
	value, exists := r.getDecoded()[key]
	if !assert.True(r.t, exists, fmt.Sprintf("Property [%s] does not exist.", key)) {
		return r
	}

	array, ok := value.([]any)
	if !assert.True(r.t, ok, fmt.Sprintf("Expected key [%s] to hold an array, but got %T", key, value)) {
		return r
	}

	for _, item := range array {
		itemJson, err := r.json.Marshal(item)
		if !assert.NoError(r.t, err) {
			continue
		}

		newJson, err := NewAssertableJSON(r.t, r.json, string(itemJson))
		if !assert.NoError(r.t, err) {
			continue
		}

		callback(newJson)
	}

	return r
}

func (r *AssertableJson) getDecoded() map[string]any {
	return r.decoded
}