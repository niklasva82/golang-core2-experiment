package main

import (
	"fmt"
)

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
