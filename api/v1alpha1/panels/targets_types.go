// +kubebuilder:object:generate=true

package panels

// Query editor options
// Referers to https://pkg.go.dev/github.com/grafana-tools/sdk#Target
type Target struct {
	// common fields
	// Reference ID
	RefID int64 `json:"refId,omitempty" yaml:"refId,omitempty"`
	// Set a datasource
	Datasource string `json:"datasource,omitempty" yaml:"datasource,omitempty"`
	// hide: bool
	Hide bool `json:"hide,omitempty" yaml:"hide,omitempty"`

	// only support prometheus,and the corresponding fields are as follows:
	// Input for fetching metrics.
	Expression string `json:"expr,omitempty" yaml:"expr,omitempty"`
	// Legend format for outputs. You can make a dynamic legend with templating variables.
	LegendFormat string `json:"legendFormat,omitempty" yaml:"legendFormat,omitempty"`
	// Set series time interval
	Step string `json:"step,omitempty" yaml:"step,omitempty"`
	// extra fields
	IntervalFactor int    `json:"intervalFactor,omitempty" yaml:"intervalFactor,omitempty"`
	Instant        bool   `json:"instant,omitempty" yaml:"instant,omitempty"`
	Format         string `json:"format,omitempty" yaml:"format,omitempty"`
	Interval       string `json:"interval,omitempty" yaml:"interval,omitempty"`
}
