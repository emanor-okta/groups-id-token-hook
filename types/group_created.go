package hooks

type TargetItem struct {
	Id          string `json:"id,omitempty"`
	Type        string `json:"type,omitempty"`
	AlternateId string `json:"alternateId,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
}

type Event struct {
	EventType string `json:"eventType,omitempty"`
	Outcome   struct {
		Result string `json:"result,omitempty"`
	} `json:"outcome,omitempty"`
	Target []TargetItem `json:"target,omitempty"`
}

type GroupCreated struct {
	Data struct {
		Events []Event `json:"events,omitempty"`
	} `json:"data,omitempty"`
}
