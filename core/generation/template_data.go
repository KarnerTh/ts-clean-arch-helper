package generation

import (
	"fmt"
	"strings"

	"github.com/KarnerTh/ts-clean-arch-helper/analyze"
)

type ObjectTemplateData struct {
	Suffix SuffixType
	Data   []analyze.ObjectDetail
}

type ConverterTemplateData struct {
	ClassName string
	Type      analyze.ObjectType
	Methods   []ConverterTemplateMethodData
	Variables []analyze.VariableDetail
}

type ConverterTemplateMethodData struct {
	Name          string
	ParameterName string
	FromObject    string
	ToObject      string
}

func getMethods(objectName string, suffix SuffixType) []ConverterTemplateMethodData {
	if suffix == Model {
		return []ConverterTemplateMethodData{
			{
				Name:          "toDomain",
				ParameterName: "model",
				FromObject:    fmt.Sprintf("%s%s", objectName, suffix),
				ToObject:      objectName,
			},
			{
				Name:          "toModel",
				ParameterName: "domain",
				FromObject:    objectName,
				ToObject:      fmt.Sprintf("%s%s", objectName, suffix),
			},
		}
	}

	return []ConverterTemplateMethodData{
		{
			Name:          "toDomain",
			ParameterName: "entity",
			FromObject:    fmt.Sprintf("%s%s", objectName, suffix),
			ToObject:      objectName,
		},
	}
}

func getVariables(variables []analyze.VariableDetail, objectType analyze.ObjectType, suffix SuffixType) []analyze.VariableDetail {
	result := make([]analyze.VariableDetail, len(variables))
	for index, variable := range variables {
		var value string
		var converterName string

		if objectType == analyze.ObjectTypeEnum {
			value = variable.Value
		} else {
			if isCustomType(variable.Value) {
				value = fmt.Sprintf("%s%s", variable.Value, suffix)
				converterName, _ = getConverterName(variable.Value, suffix)
				if strings.Contains(value, "[]") {
					value = strings.ReplaceAll(value, "[]", "") + "[]"
				}
			} else {
				value = variable.Value
			}
		}

		detail := analyze.VariableDetail{
			Name:          variable.Name,
			Value:         value,
			IsCustomType:  isCustomType(variable.Value),
			ConverterName: converterName,
			IsArray:       strings.Contains(value, "[]"),
		}

		result[index] = detail
	}
	return result
}

func getTsTypes() []string {
	return []string{
		"any",
		"boolean",
		"Date",
		"never",
		"null",
		"number",
		"string",
		"undefined",
		"unknown",
	}
}

func isCustomType(value string) bool {
	if len(value) == 0 {
		return false
	}

	for _, tsType := range getTsTypes() {
		if tsType == value {
			return false
		}
	}

	return true
}
