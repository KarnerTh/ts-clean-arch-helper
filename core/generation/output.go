package generation

type Output struct {
	DomainObjects   string `json:"domainObjects"`
	ModelObjects    string `json:"modelObjects"`
	EntityConverter string `json:"entityConverter"`
	ModelConverter  string `json:"modelConverter"`
}
