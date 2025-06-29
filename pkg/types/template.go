package types

// Template represents a project template
type Template struct {
	ID           string             `yaml:"id" json:"id"`
	Name         string             `yaml:"name" json:"name"`
	Description  string             `yaml:"description" json:"description"`
	Type         string             `yaml:"type" json:"type"`
	Architecture string             `yaml:"architecture" json:"architecture"`
	Version      string             `yaml:"version" json:"version"`
	Author       string             `yaml:"author" json:"author"`
	License      string             `yaml:"license" json:"license"`
	Include      *TemplateIncludes  `yaml:"include" json:"include"`
	Variables    []TemplateVariable `yaml:"variables" json:"variables"`
	Files        []TemplateFile     `yaml:"files" json:"files"`
	Dependencies []Dependency       `yaml:"dependencies" json:"dependencies"`
	PostHooks    []Hook             `yaml:"post_hooks" json:"post_hooks"`
	Features     []TemplateFeature  `yaml:"features" json:"features"`
	Validation   []ValidationRule   `yaml:"validation" json:"validation"`
	Metadata     map[string]any     `yaml:"metadata" json:"metadata"`
}

// TemplateVariable represents a configurable variable in a template
type TemplateVariable struct {
	Name        string   `yaml:"name" json:"name"`
	Type        string   `yaml:"type" json:"type"`
	Description string   `yaml:"description" json:"description"`
	Default     any      `yaml:"default" json:"default"`
	Required    bool     `yaml:"required" json:"required"`
	Choices     []string `yaml:"choices" json:"choices"`
	Validation  string   `yaml:"validation" json:"validation"`
}

// TemplateFile represents a file in a template
type TemplateFile struct {
	Source      string `yaml:"source" json:"source"`
	Destination string `yaml:"destination" json:"destination"`
	Condition   string `yaml:"condition" json:"condition"`
	Executable  bool   `yaml:"executable" json:"executable"`
}

// Dependency represents a Go module dependency
type Dependency struct {
	Module    string `yaml:"module" json:"module"`
	Version   string `yaml:"version" json:"version"`
	Condition string `yaml:"condition" json:"condition"`
}

// Hook represents a post-generation hook
type Hook struct {
	Name    string   `yaml:"name" json:"name"`
	Command string   `yaml:"command" json:"command"`
	Args    []string `yaml:"args" json:"args"`
	WorkDir string   `yaml:"work_dir" json:"work_dir"`
}

// TemplateIncludes represents external file includes
type TemplateIncludes struct {
	Variables    string `yaml:"variables" json:"variables"`
	Dependencies string `yaml:"dependencies" json:"dependencies"`
	Features     string `yaml:"features" json:"features"`
}

// TemplateFeature represents a feature provided by the template
type TemplateFeature struct {
	Name        string `yaml:"name" json:"name"`
	Description string `yaml:"description" json:"description"`
	EnabledWhen string `yaml:"enabled_when" json:"enabled_when"`
}

// ValidationRule represents a validation rule for the template
type ValidationRule struct {
	Name        string `yaml:"name" json:"name"`
	Description string `yaml:"description" json:"description"`
	Value       string `yaml:"value" json:"value"`
}

// TemplateMetadata represents metadata about available templates
type TemplateMetadata struct {
	Templates []TemplateInfo `json:"templates"`
}

// TemplateInfo represents basic information about a template
type TemplateInfo struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Type         string   `json:"type"`
	Architecture string   `json:"architecture"`
	Tags         []string `json:"tags"`
}
