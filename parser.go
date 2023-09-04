package main

import (
	"fmt"
	"reflect"
	"strings"
)

type ValidationErrorKind string

const (
	INVALID      ValidationErrorKind = "invalid"
	WRONG_TYPE   ValidationErrorKind = "wrong_type"
	MISSING_ATTR ValidationErrorKind = "missing_attr"
	INVALID_NONE ValidationErrorKind = "invalid_none"
)

type ParserContext struct {
	Path []string
}

type ValidationError struct {
	Kind ValidationErrorKind
	Path []string
}

// Basically the toString of an error
func (e ValidationError) Error() string {
	return fmt.Sprintf("Error: %s, Path: %v", e.Kind, strings.Join(e.Path, ", "))
}

type Parser interface {
	Parse(val interface{}, context ParserContext) (interface{}, error)
	/* This getter is needed because interfaces cannot have variables and we need this to be
	in the interface in order to check for it in DictParser where the exact type is not known */
	IsOptional() bool
}

type StrParser struct {
	Optional  bool
	AllowNone bool
	// In order for MinLength to be optional, we define it as a pointer
	MinLength *int
}

func (p StrParser) IsOptional() bool {
	return p.Optional
}

func (p StrParser) Parse(val interface{}, context ParserContext) (interface{}, error) {
	if val == nil {
		if p.AllowNone {
			return val, nil
		}
		return nil, ValidationError{Kind: INVALID_NONE, Path: context.Path}
	}

	strVal, ok := val.(string)
	if !ok {
		return nil, ValidationError{Kind: WRONG_TYPE, Path: context.Path}
	}

	if p.MinLength != nil && len(strVal) < *p.MinLength {
		return nil, ValidationError{Kind: INVALID, Path: context.Path}
	}

	return strVal, nil
}

type IntParser struct {
	Optional  bool
	AllowNone bool
	MinValue  *int
	MaxValue  *int
}

func (p IntParser) IsOptional() bool {
	return p.Optional
}

func (p IntParser) Parse(val interface{}, context ParserContext) (interface{}, error) {
	if val == nil {
		if p.AllowNone {
			return val, nil
		}
		return nil, ValidationError{Kind: INVALID_NONE, Path: context.Path}
	}

	intVal, ok := val.(int)
	if !ok {
		return nil, ValidationError{Kind: WRONG_TYPE, Path: context.Path}
	}

	if p.MinValue != nil && intVal < *p.MinValue {
		return nil, ValidationError{Kind: INVALID, Path: context.Path}
	}

	if p.MaxValue != nil && intVal > *p.MaxValue {
		return nil, ValidationError{Kind: INVALID, Path: context.Path}
	}

	return intVal, nil
}

type DictParser struct {
	Schema    map[string]Parser
	Optional  bool
	AllowNone bool
}

func (p DictParser) Parse(val interface{}, context ParserContext) (interface{}, error) {
	if val == nil {
		if p.AllowNone {
			return val, nil
		}
		return nil, ValidationError{Kind: INVALID_NONE, Path: context.Path}
	}

	mapVal, _ := val.(map[string]interface{})
	// With empty maps, the type assertion does not work, but the following does.
	if reflect.ValueOf(val).Kind() != reflect.Map {
		return nil, ValidationError{Kind: WRONG_TYPE, Path: context.Path}
	}

	output := make(map[string]interface{})
	for key, parser := range p.Schema {
		subVal, exists := mapVal[key]
		if !exists && parser.IsOptional() {
			return nil, ValidationError{Kind: MISSING_ATTR, Path: append(context.Path, key)}
		}
		subContext := ParserContext{Path: append(context.Path, key)}
		parsed, err := parser.Parse(subVal, subContext)
		if err != nil {
			return nil, err
		}
		output[key] = parsed
	}

	return output, nil
}

func (p DictParser) IsOptional() bool {
	return p.Optional
}
