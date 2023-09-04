package main

import (
	"strings"
	"testing"
)

func TestingParser(t *testing.T, parser Parser, value interface{}, expectedError *ValidationErrorKind) {
	context := ParserContext{Path: []string{}}
	_, err := parser.Parse(value, context)
	if expectedError == nil && err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if expectedError != nil && err == nil {
		t.Errorf("Expected error kind %s, but no error returned", *expectedError)
	}
	if expectedError != nil && err != nil && !strings.Contains(err.Error(), string(*expectedError)) {
		t.Errorf("Expected error kind: %s, but got %v", *expectedError, err)
	}
}

func TestStrParserValidString(t *testing.T) {
	p := StrParser{Optional: false, AllowNone: false}
	TestingParser(t, p, "Valid string", nil)
}

func TestStrParserRaisesOnNone(t *testing.T) {
	p := StrParser{Optional: false, AllowNone: false}
	error_kind := INVALID_NONE
	TestingParser(t, p, nil, &error_kind)
}

func TestStrParserAllowedNone(t *testing.T) {
	p := StrParser{Optional: false, AllowNone: true}
	TestingParser(t, p, nil, nil)
}

func TestStrParserWrongType(t *testing.T) {
	p := StrParser{Optional: false, AllowNone: true}
	error_kind := WRONG_TYPE
	TestingParser(t, p, 2, &error_kind)
}

func TestStrParserTooShort(t *testing.T) {
	minLength := 6
	p := StrParser{Optional: false, AllowNone: false, MinLength: &minLength}
	error_kind := INVALID
	TestingParser(t, p, "short", &error_kind)
}

func TestStrParserLongEnough(t *testing.T) {
	minLength := 6
	p := StrParser{Optional: false, AllowNone: false, MinLength: &minLength}
	TestingParser(t, p, "long enough", nil)
}

func TestIntParserValidInt(t *testing.T) {
	p := IntParser{Optional: false, AllowNone: false}
	TestingParser(t, p, 56, nil)
}

func TestIntParserRaisesOnNone(t *testing.T) {
	p := IntParser{Optional: false, AllowNone: false}
	error_kind := INVALID_NONE
	TestingParser(t, p, nil, &error_kind)
}

func TestIntParserAllowedNone(t *testing.T) {
	p := IntParser{Optional: false, AllowNone: true}
	TestingParser(t, p, nil, nil)
}

func TestIntParserWrongType(t *testing.T) {
	p := StrParser{Optional: false, AllowNone: true}
	error_kind := WRONG_TYPE
	TestingParser(t, p, 2, &error_kind)
}

func TestIntParserMin(t *testing.T) {
	minValue := 6
	p := IntParser{Optional: false, AllowNone: false, MinValue: &minValue}
	error_kind := INVALID
	TestingParser(t, p, 5, &error_kind)
}

func TestIntParserMax(t *testing.T) {
	maxValue := 6
	p := IntParser{Optional: false, AllowNone: false, MaxValue: &maxValue}
	error_kind := INVALID
	TestingParser(t, p, 7, &error_kind)
}

func TestIntParserMinAllowed(t *testing.T) {
	minValue := 3
	p := IntParser{Optional: false, AllowNone: false, MinValue: &minValue}
	TestingParser(t, p, 5, nil)
}

func TestIntParserMaxAllowed(t *testing.T) {
	maxValue := 6
	p := IntParser{Optional: false, AllowNone: false, MaxValue: &maxValue}
	TestingParser(t, p, 5, nil)
}

func TestDictParserEmptyDict(t *testing.T) {
	schema := map[string]Parser{
	}
	p := DictParser{Schema: schema, Optional: false, AllowNone: false}
	TestingParser(t, p, map[string]int{}, nil)
}

func TestDictParserNoneError(t *testing.T) {
	p := DictParser{Optional: false, AllowNone: false}
	error_kind := INVALID_NONE
	TestingParser(t, p, nil, &error_kind)
}

func TestDictParserAllowedNone(t *testing.T) {
	p := DictParser{Optional: false, AllowNone: true}
	TestingParser(t, p, nil, nil)
}

func TestDictParserValidStringDict(t *testing.T) {
	schema := map[string]Parser{
		"attr":       StrParser{Optional: false, AllowNone: false},
	}
	p := DictParser{Schema: schema, Optional: false, AllowNone: false}
	val := map[string]interface{}{
		"attr": "value",
	}
	TestingParser(t, p, val, nil)
}