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
	for i, variable := range variables {
		if objectType == analyze.ObjectTypeEnum {
			result[i] = getEnumVariables(variable)
		} else {
			result[i] = getInterfaceVariables(variable, suffix)
		}
	}
	return result
}

func getInterfaceVariables(variable analyze.VariableDetail, suffix SuffixType) analyze.VariableDetail {
	value := []string{}
	var converterName string
	includesCustomType := false
	includesArray := false

	for i, ty := range variable.Types {
		if !isCustomType(ty) {
			value = append(value, ty)
			continue
		}
		includesCustomType = true
		converterName, _ = getConverterName(ty, suffix)
		value = append(value, fmt.Sprintf("%s%s", ty, suffix))
		if strings.Contains(value[i], "[]") {
			includesArray = true
			value[i] = strings.ReplaceAll(value[i], "[]", "") + "[]"
		}
	}

	return analyze.VariableDetail{
		Name:          variable.Name,
		Types:         variable.Types,
		IsCustomType:  includesCustomType,
		ConverterName: converterName,
		IsArray:       includesArray,
	}
}

func getEnumVariables(v analyze.VariableDetail) analyze.VariableDetail {
	return v
}

func getTsTypes() []string {
	return []string{
		"any",
		"boolean",
		"Boolean",
		"Date",
		"never",
		"null",
		"number",
		"Number",
		"string",
		"String",
		"undefined",
		"unknown",
		"void",
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
