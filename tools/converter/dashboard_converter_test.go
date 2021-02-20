package converter

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/grafana-tools/sdk"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	v1alpha1 "kubesphere.io/monitoring-dashboard/api/v1alpha1"
)

func defaultVar(varType string) sdk.TemplateVar {
	return sdk.TemplateVar{
		Type:  varType,
		Name:  "var",
		Label: "Label",
	}
}

func TestConvertInvalidJSONToKubsphereDashboard(t *testing.T) {
	req := require.New(t)

	converter := NewConverter(zap.NewNop())
	err := converter.ConvertKubsphereDashboard(bytes.NewBufferString(""), bytes.NewBufferString(""), false, "default", "test-dashboard")

	req.Error(err)
}

func TestConvertValidJSONToKubsphereClusterManifest(t *testing.T) {
	req := require.New(t)

	converter := NewConverter(zap.NewNop())
	err := converter.ConvertKubsphereDashboard(bytes.NewBufferString("{}"), bytes.NewBufferString(""), true, "default", "test-cluster-dashboard")

	req.NoError(err)
}

func TestConvertGeneralSettings(t *testing.T) {
	req := require.New(t)

	board := &sdk.Board{}
	board.Title = "title"
	board.SharedCrosshair = true
	board.Editable = true
	board.Tags = []string{"tag", "other"}
	board.Refresh = &sdk.BoolString{
		Flag:  true,
		Value: "5s",
	}

	converter := NewConverter(zap.NewNop())

	dashboard := &v1alpha1.DashboardSpec{}
	converter.convertGeneralSettings(board, dashboard)

	req.Equal("title", dashboard.Title)
	req.Equal("5s", dashboard.AutoRefresh)
	req.Equal([]string{"tag", "other"}, dashboard.Tags)
	req.True(dashboard.Editable)
	req.True(dashboard.SharedCrosshair)
}

func TestConvertUnknownVar(t *testing.T) {
	req := require.New(t)

	variable := defaultVar("unknown")

	converter := NewConverter(zap.NewNop())

	dashboard := &v1alpha1.DashboardSpec{}
	converter.convertVariables([]sdk.TemplateVar{variable}, dashboard)

	req.Len(dashboard.Templatings, 0)
}

func TestConvertIntervalVar(t *testing.T) {
	req := require.New(t)

	variable := defaultVar("interval")
	variable.Name = "var_interval"
	variable.Label = "Label interval"
	variable.Current = sdk.Current{Text: "30sec", Value: "30s"}
	variable.Options = []sdk.Option{
		{Text: "10sec", Value: "10s"},
		{Text: "30sec", Value: "30s"},
		{Text: "1min", Value: "1m"},
	}

	converter := NewConverter(zap.NewNop())

	dashboard := &v1alpha1.DashboardSpec{}
	converter.convertVariables([]sdk.TemplateVar{variable}, dashboard)

	req.Len(dashboard.Templatings, 1)
	req.NotNil(dashboard.Templatings[0])

	interval := dashboard.Templatings[0]

	req.Equal("var_interval", interval.Name)
	req.Equal("interval", interval.Type)
	req.Equal("Label interval", interval.Label)
	req.Equal("30s", interval.Default)
	req.ElementsMatch([]string{"10s", "30s", "1m"}, interval.Values)
}

func TestConvertCustomVar(t *testing.T) {
	req := require.New(t)

	variable := defaultVar("custom")
	variable.Name = "var_custom"
	variable.Label = "Label custom"
	variable.Current = sdk.Current{Text: "85th", Value: "85"}
	variable.Options = []sdk.Option{
		{Text: "50th", Value: "50"},
		{Text: "85th", Value: "85"},
		{Text: "99th", Value: "99"},
	}

	converter := NewConverter(zap.NewNop())

	dashboard := &v1alpha1.DashboardSpec{}
	converter.convertVariables([]sdk.TemplateVar{variable}, dashboard)

	req.Len(dashboard.Templatings, 1)
	req.NotNil(dashboard.Templatings[0])

	custom := dashboard.Templatings[0]

	req.Equal("custom", custom.Type)
	req.Equal("var_custom", custom.Name)
	req.Equal("Label custom", custom.Label)
	req.Equal("85", custom.Default)
	req.True(reflect.DeepEqual(custom.ValuesMap, map[string]string{
		"50th": "50",
		"85th": "85",
		"99th": "99",
	}))
}

func TestConvertConstVar(t *testing.T) {
	req := require.New(t)

	variable := defaultVar("const")
	variable.Name = "var_const"
	variable.Label = "Label const"
	variable.Current = sdk.Current{Text: "85th", Value: "85"}
	variable.Options = []sdk.Option{
		{Text: "85th", Value: "85"},
		{Text: "99th", Value: "99"},
	}

	converter := NewConverter(zap.NewNop())

	dashboard := &v1alpha1.DashboardSpec{}
	converter.convertVariables([]sdk.TemplateVar{variable}, dashboard)

	req.Len(dashboard.Templatings, 1)
	req.NotNil(dashboard.Templatings[0])

	constant := dashboard.Templatings[0]

	req.Equal("var_const", constant.Name)
	req.Equal("const", constant.Type)
	req.Equal("Label const", constant.Label)
	req.Equal("85th", constant.Default)
	req.True(reflect.DeepEqual(constant.ValuesMap, map[string]string{
		"85th": "85",
		"99th": "99",
	}))
}

func TestConvertQueryVar(t *testing.T) {
	req := require.New(t)
	datasource := "prometheus-default"

	variable := defaultVar("query")
	variable.Name = "var_query"
	variable.Label = "Query"
	variable.IncludeAll = true
	variable.Current = sdk.Current{Value: "$__all"}
	variable.Datasource = &datasource
	variable.Query = "prom_query"

	converter := NewConverter(zap.NewNop())

	dashboard := &v1alpha1.DashboardSpec{}
	converter.convertVariables([]sdk.TemplateVar{variable}, dashboard)

	req.Len(dashboard.Templatings, 1)
	req.NotNil(dashboard.Templatings[0])

	query := dashboard.Templatings[0]

	req.Equal("var_query", query.Name)
	req.Equal("query", query.Type)
	req.Equal("Query", query.Label)
	req.Equal(datasource, query.Datasource)
	req.Equal("prom_query", query.Request)
	req.True(query.IncludeAll)
	req.True(query.DefaultAll)

}

// func TestConvertTargetFailsIfNoValidTargetIsGiven(t *testing.T) {
// 	req := require.New(t)
// 	converter := NewConverter(zap.NewNop())
// 	var target sdk.Target
// 	convertedTarget := converter.convertTarget(target, false, 0)
// 	req.Zero(*convertedTarget)
// }

func TestConvertTargetWithPrometheusTarget(t *testing.T) {
	req := require.New(t)

	converter := NewConverter(zap.NewNop())

	target := sdk.Target{
		Expr:         "prometheus_query",
		LegendFormat: "{{ field }}",
		RefID:        "A",
	}

	convertedTarget := converter.convertTarget(target, 0)

	req.NotNil(convertedTarget)
	req.Equal("prometheus_query", convertedTarget.Expression)
	req.Equal("{{field}}", convertedTarget.LegendFormat)
	req.Equal(int64(1), convertedTarget.RefID)
}

func TestConvertTagAnnotationIgnoresBuiltIn(t *testing.T) {
	req := require.New(t)

	annotation := sdk.Annotation{Name: "Annotations & Alerts"}
	dashboard := &v1alpha1.DashboardSpec{}

	NewConverter(zap.NewNop()).convertAnnotations([]sdk.Annotation{annotation}, dashboard)

	req.Len(dashboard.Annotations, 0)
}

func TestConvertTagAnnotationIgnoresUnknownTypes(t *testing.T) {
	req := require.New(t)

	annotation := sdk.Annotation{Name: "Will be ignored", Type: "dashboard"}
	dashboard := &v1alpha1.DashboardSpec{}

	NewConverter(zap.NewNop()).convertAnnotations([]sdk.Annotation{annotation}, dashboard)

	req.Len(dashboard.Annotations, 0)
}

func TestConvertTagAnnotation(t *testing.T) {
	req := require.New(t)

	converter := NewConverter(zap.NewNop())

	datasource := "-- Grafana --"
	annotation := sdk.Annotation{
		Type:       "tags",
		Datasource: &datasource,
		IconColor:  "#5794F2",
		Name:       "Deployments",
		Tags:       []string{"deploy"},
	}
	dashboard := &v1alpha1.DashboardSpec{}

	converter.convertAnnotations([]sdk.Annotation{annotation}, dashboard)

	req.Len(dashboard.Annotations, 1)
	req.Equal("Deployments", dashboard.Annotations[0].Name)
	req.ElementsMatch([]string{"deploy"}, dashboard.Annotations[0].Tags)
	req.Equal("#5794F2", dashboard.Annotations[0].IconColor)
	req.Equal(datasource, dashboard.Annotations[0].Datasource)
}

func TestConvertLegend(t *testing.T) {
	req := require.New(t)

	converter := NewConverter(zap.NewNop())

	rawLegend := sdk.Legend{
		AlignAsTable: true,
		Avg:          true,
		Current:      true,
		HideEmpty:    true,
		HideZero:     true,
		Max:          true,
		Min:          true,
		RightSide:    true,
		Show:         true,
		Total:        true,
	}

	legend := converter.convertLegend(rawLegend)

	req.ElementsMatch(
		[]string{"as_table", "to_the_right", "min", "max", "avg", "current", "total", "no_null_series", "no_zero_series"},
		legend,
	)
}

func TestConvertCanHideLegend(t *testing.T) {
	req := require.New(t)
	converter := NewConverter(zap.NewNop())

	legend := converter.convertLegend(sdk.Legend{Show: false})
	req.ElementsMatch([]string{"hide"}, legend)
}

func TestConvertPanelsAreEmpty(t *testing.T) {
	req := require.New(t)
	converter := NewConverter(zap.NewNop())
	var dashboard *v1alpha1.DashboardSpec
	var panels []*sdk.Panel
	converter.convertPanels(panels, dashboard, false)
	req.Zero(dashboard)
}

func TestConvertBarGaugePanel(t *testing.T) {
	bargaugePanel := &sdk.Panel{
		CommonPanel: sdk.CommonPanel{
			Title: "a common panel for test",
			Type:  "bargauge",
		},
		BarGaugePanel: &sdk.BarGaugePanel{
			Targets: []sdk.Target{
				sdk.Target{
					Datasource: "a bar guage datasource",
				},
			},
		},
	}

	panels := []*sdk.Panel{
		bargaugePanel,
	}

	req := require.New(t)
	converter := NewConverter(zap.NewNop())
	dashboard := &v1alpha1.DashboardSpec{}
	converter.convertPanels(panels, dashboard, false)
	req.Equal(dashboard.Panels[0].Targets[0].Datasource, "a bar guage datasource")
}

func TestConvertGraphPanel(t *testing.T) {
	graphPanel := &sdk.Panel{
		CommonPanel: sdk.CommonPanel{
			Title: "a common panel for test",
			Type:  "graph",
		},
		GraphPanel: &sdk.GraphPanel{
			Targets: []sdk.Target{
				sdk.Target{
					Datasource: "a graph datasource",
				},
			},
		},
	}

	panels := []*sdk.Panel{
		graphPanel,
	}

	req := require.New(t)
	converter := NewConverter(zap.NewNop())
	dashboard := &v1alpha1.DashboardSpec{}
	converter.convertPanels(panels, dashboard, false)
	req.Equal(dashboard.Panels[0].Targets[0].Datasource, "a graph datasource")
}

func TestConvertSinglestatPanel(t *testing.T) {
	singlestatPanel := &sdk.Panel{
		CommonPanel: sdk.CommonPanel{
			Title: "a common panel for test",
			Type:  "singlestat",
		},
		SinglestatPanel: &sdk.SinglestatPanel{
			Targets: []sdk.Target{
				sdk.Target{
					Datasource: "a single stat datasource",
				},
			},
		},
	}

	panels := []*sdk.Panel{
		singlestatPanel,
	}

	req := require.New(t)
	converter := NewConverter(zap.NewNop())
	dashboard := &v1alpha1.DashboardSpec{}
	converter.convertPanels(panels, dashboard, false)
	req.Equal(dashboard.Panels[0].Targets[0].Datasource, "a single stat datasource")
}

func TestConvertTablePanel(t *testing.T) {
	tablePanel := &sdk.Panel{
		CommonPanel: sdk.CommonPanel{
			Title: "a common panel for test",
			Type:  "table",
		},
		TablePanel: &sdk.TablePanel{
			Targets: []sdk.Target{
				sdk.Target{
					Datasource: "a table datasource",
				},
			},
		},
	}

	panels := []*sdk.Panel{
		tablePanel,
	}

	req := require.New(t)
	converter := NewConverter(zap.NewNop())
	dashboard := &v1alpha1.DashboardSpec{}
	converter.convertPanels(panels, dashboard, false)
	req.Equal(dashboard.Panels[0].Targets[0].Datasource, "a table datasource")
}

func TestConvertTextPanel(t *testing.T) {
	textPanel := &sdk.Panel{
		CommonPanel: sdk.CommonPanel{
			Title: "a common panel for test",
			Type:  "text",
		},
		TextPanel: &sdk.TextPanel{
			Content: "a markdown content for test",
			Mode:    "markdown",
		},
	}

	panels := []*sdk.Panel{
		textPanel,
	}

	req := require.New(t)
	converter := NewConverter(zap.NewNop())
	dashboard := &v1alpha1.DashboardSpec{}
	converter.convertPanels(panels, dashboard, false)
	req.Equal(dashboard.Panels[0].Markdown, "a markdown content for test")
}

func TestConvertExpr(t *testing.T) {
	req := require.New(t)
	testCase := []string{
		"sum (elasticsearch_jvm_memory_used_bytes{cluster=\"$cluster\",name=~\"$name\"}) / sum (elasticsearch_jvm_memory_max_bytes{cluster=\"$cluster\",name=~\"$name\"}) * 100",
		"elasticsearch_cluster_health_number_of_pending_tasks{cluster=\"$cluster\"}",
	}
	newExpr := convertExpr(testCase[0])
	req.Equal(newExpr, "sum (elasticsearch_jvm_memory_used_bytes) / sum (elasticsearch_jvm_memory_max_bytes) * 100")
	newExpr2 := convertExpr(testCase[1])
	req.Equal(newExpr2, "elasticsearch_cluster_health_number_of_pending_tasks")
}

func TestHandleLegendFormat(t *testing.T) {
	req := require.New(t)
	testCase := "'{{name}}: {{ type }}'"
	newExpr := handleLegendFormat(testCase)
	req.Equal(newExpr, "'{{name}}: {{type}}'")
}
