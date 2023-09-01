package main

import (
	"fmt"
	"strings"
)

type ValidationErrorKind string

const (
	INVALID      ValidationErrorKind = "invalid"
	WRONG_TYPE                       = "wrong_type"
	MISSING_ATTR                     = "missing_attr"
	INVALID_NONE                     = "invalid_none"
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

	mapVal, ok := val.(map[string]interface{})
	if !ok {
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

func main() {
	// In order for MinLength, MinValue and MaxValue to be optional, they need to be passed as pointers, so we need to decalre them first
	titleMinLength := 2

	organizationParser := DictParser{
		Schema: map[string]Parser{
			"id":    IntParser{Optional: false, AllowNone: false},
			"title": StrParser{Optional: false, AllowNone: false, MinLength: &titleMinLength},
		},
		Optional:  true,
		AllowNone: true,
	}

	folderMinValue := 1
	folderMaxValue := 10
	mainParser := DictParser{
		Schema: map[string]Parser{
			"title":        StrParser{Optional: false, AllowNone: false, MinLength: &titleMinLength},
			"folder_id":    IntParser{Optional: true, AllowNone: true, MinValue: &folderMinValue, MaxValue: &folderMaxValue},
			"owner_id":     IntParser{Optional: true, AllowNone: true},
			"description":  StrParser{Optional: true, AllowNone: true},
			"organization": organizationParser,
		},
	}

	data := map[string]interface{}{
		"title":       "Title",
		"owner_id":    3,
		"folder_id":   1,
		"description": "",
		"organization": map[string]interface{}{
			"id":    3,
			"title": "My Organization",
		},
	}

	context := ParserContext{Path: []string{}}
	parsedData, err := mainParser.Parse(data, context)
	if err != nil {
		fmt.Println("Error parsing data:", err)
		return
	}

	fmt.Println("Parsed data:", parsedData)
}
