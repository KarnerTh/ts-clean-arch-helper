package analyze

import (
	"bufio"
	"regexp"
	"strings"
)

type analyzer struct{}

type Analyzer interface {
	Analyze(input string) ([]ObjectDetail, error)
}

func NewAnalyzer() Analyzer {
	return analyzer{}
}

var regObject = regexp.MustCompile(`(?:export\s)?(?P<type>interface|class|enum)\s(?P<name>.*\S)\s?\{`)
var regVariable = regexp.MustCompile(`\s*(?P<name>.*):\s(?P<type>.*);|\n`)
var regEnumWithValue = regexp.MustCompile(`\s*(?P<name>.*\S)\s*=\s*(?P<value>.*)`)
var regEnumWithoutValue = regexp.MustCompile(`\s*(?P<name>.*\S),`)

func (a analyzer) Analyze(input string) ([]ObjectDetail, error) {
	var objects []ObjectDetail

	scanner := bufio.NewScanner(strings.NewReader(input))
	var curObject *ObjectDetail
	for scanner.Scan() {
		line := scanner.Text()
		if regObject.MatchString(line) {
			if curObject != nil {
				objects = append(objects, *curObject)
			}

			groups := regObject.FindStringSubmatch(line)
			name := groups[regObject.SubexpIndex("name")]
			var objectType ObjectType
			if groups[regObject.SubexpIndex("type")] == string(ObjectTypeEnum) {
				objectType = ObjectTypeEnum
			} else {
				objectType = ObjectTypeInterface
			}

			curObject = &ObjectDetail{
				Name: sanitizeObjectName(name),
				Type: objectType,
			}
		}

		if curObject == nil {
			continue
		}

		if curObject.Type == ObjectTypeInterface && regVariable.MatchString(line) {
			curObject.Variables = append(curObject.Variables, variableDetailFromInterface(line))
		} else if curObject.Type == ObjectTypeEnum {
			enumDetail := variableDetailFromEnum(line)
			if enumDetail != nil {
				curObject.Variables = append(curObject.Variables, *enumDetail)
			}
		}
	}

	if curObject != nil {
		objects = append(objects, *curObject)
	}

	return objects, nil
}

func variableDetailFromInterface(line string) VariableDetail {
	groups := regVariable.FindStringSubmatch(line)
	nameIndex := regVariable.SubexpIndex("name")
	tyIndex := regVariable.SubexpIndex("type")

	tyStr := sanitizeVariableTypeStr(groups[tyIndex])
	return VariableDetail{
		Name:  sanitizeVariableName(groups[nameIndex]),
		Types: parseTypes(tyStr),
	}
}

func variableDetailFromEnum(line string) *VariableDetail {
	if regEnumWithValue.MatchString(line) {
		groups := regEnumWithValue.FindStringSubmatch(line)
		nameIndex := regEnumWithValue.SubexpIndex("name")
		valueIndex := regEnumWithValue.SubexpIndex("value")

		return &VariableDetail{
			Name:  sanitizeVariableName(groups[nameIndex]),
			Value: sanitizeEnumValueStr(groups[valueIndex]),
		}
	} else if regEnumWithoutValue.MatchString(line) {
		groups := regEnumWithoutValue.FindStringSubmatch(line)
		nameIndex := regEnumWithoutValue.SubexpIndex("name")

		return &VariableDetail{
			Name:  sanitizeVariableName(groups[nameIndex]),
			Value: "",
		}
	}

	return nil
}

func sanitizeVariableName(name string) string {
	name = strings.ReplaceAll(name, "!", "")
	name = strings.Trim(name, " ")
	return name
}

func sanitizeVariableTypeStr(tyStr string) string {
	return strings.ReplaceAll(tyStr, "Entity", "")
}

func sanitizeEnumValueStr(value string) string {
	return strings.ReplaceAll(value, ",", "")
}

func sanitizeObjectName(object string) string {
	return strings.ReplaceAll(object, "Entity", "")
}

func parseTypes(tyStr string) []string {
	return strings.Split(strings.ReplaceAll(tyStr, " ", ""), "|")
}
