// +kubebuilder:object:generate=true
package panels

type BarGaugeOptions struct {
	Orientation string `json:"orientation,omitempty" yaml:"orientation,omitempty"`
	TextMode    string `json:"textMode,omitempty" yaml:"textMode,omitempty"`
	ColorMode   string `json:"colorMode,omitempty" yaml:"colorMode,omitempty"`
	GraphMode   string `json:"graphMode,omitempty" yaml:"graphMode,omitempty"`
	JustifyMode string `json:"justifyMode,omitempty" yaml:"justifyMode,omitempty"`
	DisplayMode string `json:"displayMode,omitempty" yaml:"displayMode,omitempty"`
	Content     string `json:"content,omitempty" yaml:"content,omitempty"`
	Mode        string `json:"mode,omitempty" yaml:"mode,omitempty"`
}

// refers to https://pkg.go.dev/github.com/grafana-tools/sdk#BarGaugePanel
type BarGaugePanel struct {
	Options *BarGaugeOptions `json:"options,omitempty" yaml:"options,omitempty"`
	// FieldConfig FieldConfig `json:"fieldConfig,omitempty" yaml:"fieldConfig,omitempty"`
}
