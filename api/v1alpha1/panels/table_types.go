// +kubebuilder:object:generate=true

package panels

// a table panel
type Table struct {
	Title                  string   `json:"title,omitempty" yaml:"title,omitempty"`
	Id                     int64    `json:"id,omitempty" yaml:"id,omitempty"`
	Height                 string   `json:"height,omitempty" yaml:"height,omitempty"`
	Transparent            bool     `json:"transparent,omitempty" yaml:"transparent,omitempty"`
	Datasource             string   `json:"datasource,omitempty" yaml:"datasource,omitempty"`
	Targets                []Target `json:"targets,omitempty" yaml:"targets,omitempty"`
	HiddenColumns          []string `json:"hidden_columns,flow" yaml:"hidden_columns,flow"`
	TimeSeriesAggregations []string `json:"time_series_aggregations" yaml:"time_series_aggregations"`
}
