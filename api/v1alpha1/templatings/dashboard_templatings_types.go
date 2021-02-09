// +kubebuilder:object:generate=true

package templatings

//referers to https://pkg.go.dev/github.com/K-Phoen/grabana/decoder#DashboardVariable
type Templatings struct {
	Interval   *VariableInterval   `json:"interval,omitempty" yaml:"interval,omitempty"`
	Custom     *VariableCustom     `json:"custom,omitempty" yaml:"custom,omitempty"`
	Query      *VariableQuery      `json:"query,omitempty" yaml:"query,omitempty"`
	Const      *VariableConst      `json:"const,omitempty" yaml:"const,omitempty"`
	Datasource *VariableDatasource `json:"datasource,omitempty" yaml:"datasource,omitempty"`
}

type VariableInterval struct {
	Name    string   `json:"name,omitempty" yaml:"name,omitempty"`
	Label   string   `json:"label,omitempty" yaml:"label,omitempty"`
	Default string   `json:"default,omitempty" yaml:"default,omitempty"`
	Values  []string `json:"values,omitempty,flow" yaml:"values,omitempty,flow"`
}

type VariableCustom struct {
	Name       string            `json:"name,omitempty" yaml:"name,omitempty"`
	Label      string            `json:"label,omitempty" yaml:"label,omitempty"`
	Default    string            `json:"default,omitempty" yaml:"default,omitempty"`
	ValuesMap  map[string]string `json:"values_map,omitempty" yaml:",omitempty"`
	IncludeAll bool              `json:"include_all,omitempty" yaml:"include_all,omitempty"`
	AllValue   string            `json:"all_value,omitempty" yaml:"all_value,omitempty"`
}

type VariableConst struct {
	Name      string            `json:"name,omitempty" yaml:"name,omitempty"`
	Label     string            `json:"label,omitempty" yaml:"label,omitempty"`
	Default   string            `json:"default,omitempty" yaml:"default,omitempty"`
	ValuesMap map[string]string `json:"values_map,omitempty" yaml:"values_map,omitempty"`
}

type VariableQuery struct {
	Name       string `json:"name,omitempty" yaml:"name,omitempty"`
	Label      string `json:"label,omitempty" yaml:"label,omitempty"`
	Datasource string `json:"datasource,omitempty" yaml:"datasource,omitempty"`
	Request    string `json:"request,omitempty" yaml:"request,omitempty"`
	Regex      string `json:"regex,omitempty" yaml:"regex,omitempty"`
	IncludeAll bool   `json:"include_all,omitempty" yaml:"include_all,omitempty"`
	DefaultAll bool   `json:"default_all,omitempty" yaml:"default_all,omitempty"`
	AllValue   string `json:"all_value,omitempty" yaml:"all_value,omitempty"`
}

type VariableDatasource struct {
	Name       string `json:"name,omitempty" yaml:"name,omitempty"`
	Label      string `json:"label,omitempty" yaml:"label,omitempty"`
	Type       string `json:"type,omitempty" yaml:"type,omitempty"`
	Regex      string `json:"regex,omitempty" yaml:"regex,omitempty"`
	IncludeAll bool   `json:"include_all,omitempty" yaml:"include_all,omitempty"`
}
