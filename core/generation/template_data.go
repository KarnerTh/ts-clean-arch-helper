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
		value := []string{}
		var converterName string
		includesCustomType := false
		includesArray := false

		if objectType == analyze.ObjectTypeEnum {
			value = append(value, variable.Value)
		} else {
			for i, v := range variable.Values {
				if isCustomType(v) {
					includesCustomType = true
					value = append(value, fmt.Sprintf("%s%s", v, suffix))
					if strings.Contains(value[i], "[]") {
						includesArray = true
						value[i] = strings.ReplaceAll(value[i], "[]", "") + "[]"
					}
				} else {
					value = append(value, v)
				}
			}
		}

		if includesCustomType {
			converterName, _ = getConverterName(variable.Value, suffix)
		}

		detail := analyze.VariableDetail{
			Name:          variable.Name,
			Value:         strings.Join(value, " | "),
			IsCustomType:  includesCustomType,
			ConverterName: converterName,
			IsArray:       includesArray,
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
