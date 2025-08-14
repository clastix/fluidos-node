// Copyright 2022-2025 FLUIDOS Project
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package models

import (
	"encoding/json"

	"k8s.io/apimachinery/pkg/api/resource"
)

type FilterData interface {
	GetFilterType() FilterType
}

// FilterType represents the type of Filter.
type FilterType string

const (
	// MatchFilter is the identifier for a match filter.
	MatchFilter FilterType = "Match"
	// RangeFilter is the identifier for a range filter.
	RangeFilter FilterType = "Range"
)

// NumberFilter represents a filter for a numeric value.
type NumberFilter struct {
	Name FilterType      `scheme:"name"`
	Data json.RawMessage `scheme:"data"`
}

func (r NumberFilter) GetFilterType() FilterType {
	return r.Name
}

type NumberMatchFilter struct {
	Value float64 `scheme:"value"`
}

type NumberRangeFilter struct {
	Min *float64 `scheme:"min,omitempty"`
	Max *float64 `scheme:"max,omitempty"`
}

// ResourceQuantityFilter represents a filter for a resource quantity.
type ResourceQuantityFilter struct {
	Name FilterType      `scheme:"name"`
	Data json.RawMessage `scheme:"data"`
}

func (r ResourceQuantityFilter) GetFilterType() FilterType {
	return r.Name
}

// BooleanFilter is a filter that can be applied to a boolean value.
type BooleanFilter struct {
	Condition bool `scheme:"condition"`
}

// GetFilterType returns the type of the Filter.
func (BooleanFilter) GetFilterType() FilterType {
	return MatchFilter
}

// StringFilter represents a filter for a string.
type StringFilter struct {
	Name FilterType      `scheme:"name"`
	Data json.RawMessage `scheme:"data"`
}

func (s StringFilter) GetFilterType() FilterType {
	return s.Name
}

// ResourceQuantityMatchFilter represents a match filter for a resource quantity.
type ResourceQuantityMatchFilter struct {
	Value resource.Quantity `scheme:"value"`
}

// ResourceQuantityRangeFilter represents a range filter for a resource quantity.
type ResourceQuantityRangeFilter struct {
	Min *resource.Quantity `scheme:"min,omitempty"`
	Max *resource.Quantity `scheme:"max,omitempty"`
}

// StringMatchFilter represents a match filter for a string.
type StringMatchFilter struct {
	Value string `scheme:"value"`
}

// StringRangeFilter represents a range filter for a string.
type StringRangeFilter struct {
	Regex string `scheme:"regex"`
}
