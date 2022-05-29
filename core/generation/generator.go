package generation

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"text/template"

	"github.com/KarnerTh/ts-clean-arch-helper/analyze"
)

type generator struct{}

type Generator interface {
	Generate(data []analyze.ObjectDetail) (*Output, error)
}

func NewGenerator() Generator {
	return generator{}
}

var objectTemplate = `
{{ range .Data }}
{{- if eq .Type "interface" -}}
export interface {{ .Name }}{{ $.Suffix }} {
  {{- range .Variables }}
  {{ .Name }}: {{ join .Types " | " }};
  {{- end }}
}
{{- else -}}
export enum {{ .Name }}{{ $.Suffix }} {
  {{- range .Variables }}
  {{ .Name }}{{if ne .Value ""}} = {{ .Value }}{{end}},
  {{- end }}
}
{{- end }}

{{ end }}
`

var converterTemplate = `
{{ range . }}
{{- $variables := .Variables }}
{{- $type := .Type }}
export class {{ .ClassName }} {
  {{- range .Methods}}
  {{- $parameterName := .ParameterName}}
  {{- $fromObject := .FromObject }}
  {{- $toObject := .ToObject }}
  {{- $methodName := .Name }}
  public static {{ .Name }}({{ .ParameterName }}: {{ .FromObject }}): {{ .ToObject }} {
	{{ if eq $type "interface" -}}
    return {
        {{- range $variables }}
        {{ .Name }}: {{ if and (.IsCustomType) (.IsArray) -}}
					 {{ $parameterName }}.{{ .Name }}.map({{ .ConverterName }}.{{ $methodName }}),
					 {{- else if .IsCustomType -}}
					 {{ .ConverterName }}.{{ $methodName }}({{ $parameterName }}.{{ .Name }}),
					 {{- else -}}
					 {{ $parameterName }}.{{ .Name }},
					 {{- end }}
        {{- end }}
      }
      {{- else -}}
      switch({{ .ParameterName }}) {
        {{- range $variables }}
        case {{ $fromObject }}.{{ .Name }}:
          return {{ $toObject }}.{{ .Name }};
        {{- end }}
      }
      {{- end }}
  }
  {{- end }}
}
{{ end }}
`

func (g generator) Generate(data []analyze.ObjectDetail) (*Output, error) {
	domainObjects, err := generateObjects(data, Domain)
	if err != nil {
		return nil, err
	}

	modelObjects, err := generateObjects(data, Model)
	if err != nil {
		return nil, err
	}

	entityConverter, err := generateConverters(data, Entity)
	if err != nil {
		return nil, err
	}

	modelConverter, err := generateConverters(data, Model)
	if err != nil {
		return nil, err
	}

	return &Output{
		DomainObjects:   domainObjects,
		ModelObjects:    modelObjects,
		EntityConverter: entityConverter,
		ModelConverter:  modelConverter,
	}, nil
}

func generateObjects(data []analyze.ObjectDetail, suffix SuffixType) (string, error) {
	preparedData := make([]analyze.ObjectDetail, len(data))
	for index, item := range data {
		detail := analyze.ObjectDetail{
			Name:      item.Name,
			Type:      item.Type,
			Variables: getVariables(item.Variables, item.Type, suffix),
		}
		preparedData[index] = detail
	}

	templateData := ObjectTemplateData{Suffix: suffix, Data: preparedData}
	tmplObject, err := template.New("objectTemplate").Funcs(template.FuncMap{
		"join": func(s []string, d string) string {
			return strings.Join(s, d)
		},
	}).Parse(objectTemplate)
	if err != nil {
		return "", err
	}

	var output bytes.Buffer
	err = tmplObject.Execute(&output, templateData)
	if err != nil {
		return "", err
	}

	return output.String(), nil
}

func generateConverters(data []analyze.ObjectDetail, suffix SuffixType) (string, error) {
	var converters []ConverterTemplateData
	for _, d := range data {
		className, err := getConverterName(d.Name, suffix)
		if err != nil {
			return "", err
		}

		methods := getMethods(d.Name, suffix)
		variables := getVariables(d.Variables, d.Type, suffix)

		converterData := ConverterTemplateData{
			ClassName: className,
			Type:      d.Type,
			Methods:   methods,
			Variables: variables,
		}

		converters = append(converters, converterData)
	}

	tmplConverter, err := template.New("converterTemplate").Parse(converterTemplate)
	if err != nil {
		return "", err
	}

	var output bytes.Buffer
	err = tmplConverter.Execute(&output, converters)
	if err != nil {
		return "", err
	}

	return output.String(), nil
}

func getConverterName(objectName string, suffix SuffixType) (string, error) {
	objectName = strings.ReplaceAll(objectName, "[]", "")

	switch suffix {
	case Entity:
		return fmt.Sprintf("%s%s%s", objectName, suffix, "Converter"), nil
	case Model:
		return fmt.Sprintf("%s%s", objectName, "Converter"), nil
	default:
		return "", errors.New("suffix not supported")
	}
}
