package hooks

type Identity struct {
	Claims struct {
		Groups []string `json:"groups,omitempty"`
	} `json:"claims,omitempty"`
}

type TokenRequest struct {
	Data struct {
		Identity `json:"identity,omitempty"`
	} `json:"data,omitempty"`
}
