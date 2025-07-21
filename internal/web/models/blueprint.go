package models

// Blueprint represents a project template
type Blueprint struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Type         string   `json:"type"`
	Architecture string   `json:"architecture,omitempty"`
	Complexity   string   `json:"complexity"`
	FileCount    int      `json:"file_count"`
	Dependencies []string `json:"dependencies,omitempty"`
	Features     []string `json:"features,omitempty"`
}

// BlueprintFile represents a file in a blueprint
type BlueprintFile struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Condition   string `json:"condition,omitempty"`
}

// BlueprintListResponse is the response for listing blueprints
type BlueprintListResponse struct {
	Blueprints []Blueprint `json:"blueprints"`
}

// BlueprintDetailResponse is the response for getting blueprint details
type BlueprintDetailResponse struct {
	Blueprint Blueprint              `json:"blueprint"`
	Files     []BlueprintFile        `json:"files"`
	Variables map[string]interface{} `json:"variables"`
}