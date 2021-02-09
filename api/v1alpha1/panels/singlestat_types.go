// +kubebuilder:object:generate=true

package panels

// SingleStat shows instant query result
type SingleStat struct {
	// Name of the signlestat panel
	Title string `json:"title,omitempty" yaml:"title,omitempty"`
	// Must be `singlestat`
	Type string `json:"type,omitempty" yaml:"type,omitempty"`
	// Panel ID
	Id int64 `json:"id,omitempty" yaml:"id,omitempty"`
	// A collection of queries
	Targets []Target `json:"targets,omitempty" yaml:"targets,omitempty"`
	// Limit the decimal numbers
	Decimals int64 `json:"decimals,omitempty" yaml:"decimals,omitempty"`
	// Display unit
	Format string `json:"format,omitempty" yaml:"format,omitempty"`
	// datasource
	Datasource string `json:"datasource,omitempty" yaml:"datasource,omitempty"`
	// colors
	Colors []string `json:"colors,omitempty" yaml:"colors,omitempty"`
	// color settings
	Color []string `json:"color,omitempty" yaml:"color,omitempty"`
	// spark line: full or bottom
	SparkLine string `json:"sparkline" yaml:"sparkline"`
	// height
	Height string `json:"height,omitempty" yaml:"height,omitempty"`
	// gauge
	Gauge Gauge `json:"gauge,omitempty" yaml:"gauge,omitempty"`
	// value name
	ValueName string `json:"valueName,omitempty" yaml:"valueName,omitempty"`
}

// Gauge for a stat
type Gauge struct {
	MaxValue         int64 `json:"maxValue,omitempty" yaml:"maxValue,omitempty"`
	MinValue         int64 `json:"minValue,omitempty" yaml:"minValue,omitempty"`
	Show             bool  `json:"show,omitempty" yaml:"show,omitempty"`
	ThresholdLabels  bool  `json:"thresholdLabels,omitempty" yaml:"thresholdLabels,omitempty"`
	ThresholdMarkers bool  `json:"thresholdMarkers,omitempty" json:"thresholdMarkers,omitempty"`
}
