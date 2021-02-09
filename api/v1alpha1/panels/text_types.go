// +kubebuilder:object:generate=true

package panels

// dashboard text type
// referers to https://pkg.go.dev/github.com/K-Phoen/grabana/decoder#DashboardText
type Text struct {
	Title       string `json:"title,omitempty" yaml:"title,omitempty"`
	Id          int64  `json:"id,omitempty" yaml:"id,omitempty"`
	Height      string `json:"height,omitempty" yaml:"height,omitempty"`
	Transparent bool   `json:"transparent,omitempty" yaml:"transparent,omitempty"`
	HTML        string `json:"html,omitempty" yaml:"html,omitempty"`
	Markdown    string `json:"markdown,omitempty" yaml:"markdown,omitempty"`
}
