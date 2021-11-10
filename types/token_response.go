package hooks

type Value struct {
	Op    string   `json:"op,omitempty"`
	Path  string   `json:"path,omitempty"`
	Value []string `json:"value,omitempty"`
}

type Values struct {
	Type  string   `json:"type,omitempty"`
	Value []*Value `json:"value,omitempty"`
}

type TokenResponse struct {
	Commands []*Values `json:"commands,omitempty"`
}
