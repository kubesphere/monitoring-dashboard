// +kubebuilder:object:generate=true

package panels

// DashboardPanel referers to https://pkg.go.dev/github.com/K-Phoen/grabana/decoder#DashboardPanel
// type DashboardPanel struct {
// 	Name   string  `json:"name,omitempty" yaml:"name,omitempty"`
// 	Repeat string  `json:"repeat,omitempty" yaml:"repeat,omitempty"`
// 	Panels []Panel `json:"panels,omitempty" yaml:"panels,omitempty"`
// }

// Supported panel type
type Panel struct {
	// private property of the bargauge panel
	Options BarGaugeOptions `json:"options,omitempty,omitempty" yaml:"options,omitempty,omitempty"`

	// ****common properties start
	// Name pf the panel
	Title string `json:"title,omitempty" yaml:"title,omitempty"`
	// Must be `graph`
	Type string `json:"type,omitempty" yaml:"type,omitempty"`
	// Panel ID
	Id int64 `json:"id,omitempty" yaml:"id,omitempty"`
	// Panel description
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// datasource
	Datasource string `json:"datasource,omitempty" yaml:"datasource,omitempty"`
	// A collection of queries
	Targets []Target `json:"targets,omitempty,omitempty" yaml:"targets,omitempty,omitempty"`
	// Display as a bar chart
	Bars bool `json:"bars,omitempty" yaml:"bars,omitempty"`
	// Set series color
	Colors []string `json:"colors,omitempty" yaml:"colors,omitempty"`
	// color settings
	// todo graph
	Color []string `json:"color,omitempty" yaml:"color,omitempty"`
	// Display as a line chart
	Lines bool `json:"lines,omitempty" yaml:"lines,omitempty"`
	// Display as a stacked chart
	Stack bool `json:"stack,omitempty" yaml:"stack,omitempty"`
	// legend
	Legend []string `json:"legend,omitempty,flow" yaml:"legend,omitempty,flow"`
	// height
	Height string `json:"height,omitempty" yaml:"height,omitempty"`
	// transparent
	Transparent            bool     `json:"transparent,omitempty" yaml:"transparent,omitempty"`
	HiddenColumns          []string `json:"hidden_columns,omitempty,flow" yaml:"hidden_columns,omitempty,flow"`
	TimeSeriesAggregations []string `json:"time_series_aggregations,omitempty" yaml:"time_series_aggregations,omitempty"`
	// value name
	ValueName string `json:"valueName,omitempty" yaml:"valueName,omitempty"`
	// *****common properties end

	// private property of graph panel
	// Y-axis options
	Yaxes []Yaxis `json:"yaxes,omitempty" yaml:"yaxes,omitempty"`

	// private properties of singlestat panel
	// spark line: full or bottom
	SparkLine string `json:"sparkline,omitempty" yaml:"sparkline,omitempty"`
	// Limit the decimal numbers
	Decimals int64 `json:"decimals,omitempty" yaml:"decimals,omitempty"`
	// Display unit
	Format string `json:"format,omitempty" yaml:"format,omitempty"`
	// gauge
	Gauge Gauge `json:"gauge,omitempty" yaml:"gauge,omitempty"`

	// table panel has no private property

	// private properties of text panel
	HTML     string `json:"html,omitempty" yaml:"html,omitempty"`
	Markdown string `json:"markdown,omitempty" yaml:"markdown,omitempty"`
}

// type Panel struct {
// 	CommonPanel
// 	*Graph
// 	*Row
// 	*SingleStat
// 	*Table
// 	*Text
// }
