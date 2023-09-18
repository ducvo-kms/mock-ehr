package r4

type Bundle struct {
	ResourceType string  `json:"resourceType"`
	Type         string  `json:"type"`
	Entry        []Entry `json:"entry"`
}

type Entry struct {
	Resource Resource `json:"resource"`
}
