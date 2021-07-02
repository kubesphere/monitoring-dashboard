package converter

import (
	"bytes"
	"testing"

	"github.com/grafana-tools/sdk"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	v1alpha2 "kubesphere.io/monitoring-dashboard/api/v1alpha2"
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
	err := converter.ConvertToKubsphereDashboardManifests(bytes.NewBufferString(""), bytes.NewBufferString(""), false, "default", "test-dashboard")

	req.Error(err)
}

func TestConvertValidJSONToKubsphereClusterManifest(t *testing.T) {
	req := require.New(t)

	converter := NewConverter(zap.NewNop())
	err := converter.ConvertToKubsphereDashboardManifests(bytes.NewBufferString("{}"), bytes.NewBufferString(""), true, "default", "test-cluster-dashboard")

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

	dashboard := &v1alpha2.DashboardSpec{}
	converter.convertGeneralSettings(board, dashboard)

	req.Equal("title", dashboard.Title)
	req.Equal("5s", dashboard.AutoRefresh)
	req.Equal([]string{"tag", "other"}, dashboard.Tags)
	req.True(dashboard.Editable)
	req.True(dashboard.SharedCrosshair)
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

	dashboard := &v1alpha2.DashboardSpec{}
	converter.convertVariables([]sdk.TemplateVar{variable}, dashboard)

	req.Len(dashboard.Templatings, 1)
	req.NotNil(dashboard.Templatings[0])

	query := dashboard.Templatings[0]

	req.Equal("var_query", query.Name)
	req.Equal("query", query.Type)
	req.Equal("Query", query.Label)
	req.Equal(datasource, *query.Datasource)
	req.Equal("prom_query", query.Query)
	req.True(query.IncludeAll)

}

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
	dashboard := &v1alpha2.DashboardSpec{}

	NewConverter(zap.NewNop()).convertAnnotations([]sdk.Annotation{annotation}, dashboard)

	req.Len(dashboard.Annotations, 0)
}

func TestConvertTagAnnotationIgnoresUnknownTypes(t *testing.T) {
	req := require.New(t)

	annotation := sdk.Annotation{Name: "Will be ignored", Type: "dashboard"}
	dashboard := &v1alpha2.DashboardSpec{}

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
	dashboard := &v1alpha2.DashboardSpec{}

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
	var dashboard *v1alpha2.DashboardSpec
	var panels []*sdk.Panel
	converter.convertPanels(panels, dashboard, false)
	req.Zero(dashboard)
}

func TestConvertBarGaugePanel(t *testing.T) {
	datasource := "a bar guage datasource"
	bargaugePanel := &sdk.Panel{
		CommonPanel: sdk.CommonPanel{
			Title:      "a common panel for test",
			Type:       "bargauge",
			Datasource: &datasource,
		},
	}

	panels := []*sdk.Panel{
		bargaugePanel,
	}

	req := require.New(t)
	converter := NewConverter(zap.NewNop())
	dashboard := &v1alpha2.DashboardSpec{}
	converter.convertPanels(panels, dashboard, false)
	req.Equal(*dashboard.Panels[0].CommonPanel.Datasource, datasource)
}

func TestConvertGraphPanel(t *testing.T) {
	datasource := "a graph datasource"
	graphPanel := &sdk.Panel{
		CommonPanel: sdk.CommonPanel{
			Title:      "a common panel for test",
			Type:       "graph",
			Datasource: &datasource,
		},
	}

	panels := []*sdk.Panel{
		graphPanel,
	}

	req := require.New(t)
	converter := NewConverter(zap.NewNop())
	dashboard := &v1alpha2.DashboardSpec{}
	converter.convertPanels(panels, dashboard, false)
	req.Equal(*dashboard.Panels[0].CommonPanel.Datasource, datasource)
}

func TestConvertSinglestatPanel(t *testing.T) {
	datasource := "a single stat datasource"
	singlestatPanel := &sdk.Panel{
		CommonPanel: sdk.CommonPanel{
			Title:      "a common panel for test",
			Type:       "singlestat",
			Datasource: &datasource,
		},
	}

	panels := []*sdk.Panel{
		singlestatPanel,
	}

	req := require.New(t)
	converter := NewConverter(zap.NewNop())
	dashboard := &v1alpha2.DashboardSpec{}
	converter.convertPanels(panels, dashboard, false)
	req.Equal(*dashboard.Panels[0].CommonPanel.Datasource, datasource)
}

func TestConvertTablePanel(t *testing.T) {
	datasource := "a table datasource"
	tablePanel := &sdk.Panel{
		CommonPanel: sdk.CommonPanel{
			Title:      "a common panel for test",
			Type:       "table",
			Datasource: &datasource,
		},
	}

	panels := []*sdk.Panel{
		tablePanel,
	}

	req := require.New(t)
	converter := NewConverter(zap.NewNop())
	dashboard := &v1alpha2.DashboardSpec{}
	converter.convertPanels(panels, dashboard, false)
	req.Equal(*dashboard.Panels[0].CommonPanel.Datasource, datasource)
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
	dashboard := &v1alpha2.DashboardSpec{}
	converter.convertPanels(panels, dashboard, false)
	req.Equal(dashboard.Panels[0].TextPanel.Content, "a markdown content for test")
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
