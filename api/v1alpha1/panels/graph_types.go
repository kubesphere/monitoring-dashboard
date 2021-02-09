// +kubebuilder:object:generate=true

package panels

// Graph visualizes range query results into a linear graph
type Graph struct {
	// Name of the graph panel
	Title string `json:"title,omitempty" yaml:"title,omitempty"`
	// Must be `graph`
	Type string `json:"type,omitempty" yaml:"type,omitempty"`
	// Panel ID
	Id int64 `json:"id,omitempty" yaml:"id,omitempty"`
	// // Panel description
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// A collection of queries
	Targets []Target `json:"targets,omitempty" yanl:"targets,omitempty"`
	// Display as a bar chart
	Bars bool `json:"bars,omitempty" yaml:"bars,omitempty"`
	// Set series color
	Colors []string `json:"colors,omitempty" yaml:"colors,omitempty"`
	// Display as a line chart
	Lines bool `json:"lines,omitempty" yaml:"lines,omitempty"`
	// Display as a stacked chart
	Stack bool `json:"stack,omitempty" yaml:"stack,omitempty"`
	// Y-axis options
	Yaxes []Yaxis `json:"yaxes,omitempty" yaml:"yaxes,omitempty"`
	// datasource
	Datasource string `json:"datasource,omitempty" yaml:"datasource,omitempty"`
	// legend
	Legend []string `json:"legend,omitempty,flow" yaml:"legend,omitempty,flow"`
	Height string   `json:"height,omitempty" yaml:"height,omitempty"`
}

type Yaxis struct {
	// Limit the decimal numbers
	Decimals int64 `json:"decimals" yaml:"decimals"`
	// Display unit
	Format string `json:"format,omitempty" yaml:"format,omitempty"`
}
