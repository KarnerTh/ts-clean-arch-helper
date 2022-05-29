package analyze

type ObjectType string

var (
	ObjectTypeInterface ObjectType = "interface"
	ObjectTypeEnum      ObjectType = "enum"
)

type ObjectDetail struct {
	Name      string
	Type      ObjectType
	Variables []VariableDetail
}

type VariableDetail struct {
	Name          string
	Value         string
	Types         []string
	IsCustomType  bool
	ConverterName string
	IsArray       bool
}
