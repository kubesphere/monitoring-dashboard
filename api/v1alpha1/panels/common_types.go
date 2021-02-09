// +kubebuilder:object:generate=true

package panels

// Query editor options
type CommonPanel struct {
	// Name of the  panel
	Title string `json:"title,omitempty" yaml:"title,omitempty"`
	// Type of the  panel
	Type string `json:"type,omitempty" yaml:"type,omitempty"`
	// Panel ID
	Id int64 `json:"id,omitempty" yaml:"id,omitempty"`
}
