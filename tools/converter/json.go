package converter

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"

	"github.com/grafana-tools/sdk"
	"github.com/mitchellh/mapstructure"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
	v1alpha1 "kubesphere.io/monitoring-dashboard/api/v1alpha1"
	ansModel "kubesphere.io/monitoring-dashboard/api/v1alpha1/annotations"
	panelsModel "kubesphere.io/monitoring-dashboard/api/v1alpha1/panels"
	templatingsModel "kubesphere.io/monitoring-dashboard/api/v1alpha1/templatings"
)

type k8sDashboard struct {
	APIVersion string                  `json:"apiVersion" yaml:"apiVersion"`
	Kind       string                  `json:"kind" yaml:"kind"`
	Metadata   map[string]string       `json:"metadata" yaml:"metadata"`
	Spec       *v1alpha1.DashboardSpec `json:"spec" yaml:"spec"`
}

// JSON struct: this struct has a log property, so other newly added methods can access this log
type JSON struct {
	logger *zap.Logger
}

// NewJSON: new a JSON struct object with a logger object
func NewJSON(logger *zap.Logger) *JSON {
	return &JSON{
		logger: logger,
	}
}

// ToYAML: this method converts to a yaml file
func (converter *JSON) ToYAML(input io.Reader, output io.Writer, isClusterCrd bool) error {
	dashboard, err := converter.parseInput(input, isClusterCrd)
	if err != nil {
		converter.logger.Error("could parse input", zap.Error(err))
		return err
	}

	converted, err := yaml.Marshal(dashboard)
	if err != nil {
		converter.logger.Error("could marshall dashboard to yaml", zap.Error(err))
		return err
	}

	_, err = output.Write(converted)

	return err
}

// ToK8SManifest: this method converts to a k8s mainfest file
func (converter *JSON) ToK8SManifest(input io.Reader, output io.Writer, isClusterCrd bool, ns string, name string) error {
	dashboard, err := converter.parseInput(input, isClusterCrd)
	if err != nil {
		converter.logger.Error("could parse input", zap.Error(err))
		return err
	}

	apiVersion := v1alpha1.GroupVersion.Group + "/" + v1alpha1.GroupVersion.Version

	kind := "Dashboard"
	if isClusterCrd {
		kind = "ClusterDashboard"
	}

	if ns == "" {
		ns = "default"
	}
	manifest := k8sDashboard{
		APIVersion: apiVersion,
		Kind:       kind,
		Metadata:   map[string]string{"name": name, "namespace": ns},
		Spec:       dashboard,
	}

	converted, err := yaml.Marshal(manifest)
	if err != nil {
		converter.logger.Error("could marshall dashboard to yaml", zap.Error(err))
		return err
	}

	_, err = output.Write(converted)

	return err
}

// parseInput: this method reads a input json file, then extract needed fields to the yaml model
func (converter *JSON) parseInput(input io.Reader, isClusterCrd bool) (*v1alpha1.DashboardSpec, error) {
	content, err := ioutil.ReadAll(input)
	if err != nil {
		converter.logger.Error("could not read input", zap.Error(err))
		return nil, err
	}

	board := &sdk.Board{}
	if err := json.Unmarshal(content, board); err != nil {
		converter.logger.Error("could not unmarshall dashboard", zap.Error(err))
		return nil, err
	}

	// a yaml model
	dashboard := &v1alpha1.DashboardSpec{}

	// starts to convert general settings
	converter.convertGeneralSettings(board, dashboard)
	// starts to convert templating variables
	converter.convertVariables(board.Templating.List, dashboard)
	// starts to convert annotations
	converter.convertAnnotations(board.Annotations.List, dashboard)
	// starts to convert pannels
	converter.convertPanels(board.Panels, dashboard, isClusterCrd)
	// starts to convert rows
	// some dashboards only include rows
	converter.convertRows(board.Rows, dashboard, isClusterCrd)

	return dashboard, nil
}

// convert GeneralSettings
func (converter *JSON) convertGeneralSettings(board *sdk.Board, dashboard *v1alpha1.DashboardSpec) {
	dashboard.Title = board.Title
	dashboard.Editable = board.Editable
	dashboard.SharedCrosshair = board.SharedCrosshair
	dashboard.Tags = board.Tags
	dashboard.Time.From = board.Time.From
	dashboard.Time.To = board.Time.To
	dashboard.Timezone = board.Timezone

	if board.Refresh != nil {
		dashboard.AutoRefresh = board.Refresh.Value
	}
}

// convert Annotations
func (converter *JSON) convertAnnotations(annotations []sdk.Annotation, dashboard *v1alpha1.DashboardSpec) {
	for _, annotation := range annotations {
		// grafana-sdk doesn't expose the "builtIn" field, so we work around that by skipping
		// the annotation we know to be built-in by its name
		if annotation.Name == "Annotations & Alerts" {
			continue
		}

		if annotation.Type != "tags" {
			converter.logger.Warn("unhandled annotation type: skipped", zap.String("type", annotation.Type), zap.String("name", annotation.Name))
			continue
		}

		datasource := ""
		if annotation.Datasource != nil {
			datasource = *annotation.Datasource
		}

		dashboard.Annotations = append(dashboard.Annotations, ansModel.Annotation{
			Name:        annotation.Name,
			Datasource:  datasource,
			IconColor:   annotation.IconColor,
			Tags:        annotation.Tags,
			ShowLine:    annotation.ShowLine,
			LineColor:   annotation.LineColor,
			IconSize:    annotation.IconSize,
			Enable:      annotation.Enable,
			Query:       annotation.Query,
			Expr:        annotation.Expr,
			Step:        annotation.Step,
			TextField:   annotation.TextField,
			TextFormat:  annotation.TextFormat,
			TitleFormat: annotation.TitleFormat,
			TagsField:   annotation.TagsField,
			TagKeys:     annotation.TagKeys,
			Type:        annotation.Type,
		})
	}
}

// convert diferent variables
func (converter *JSON) convertVariables(variables []sdk.TemplateVar, dashboard *v1alpha1.DashboardSpec) {
	for _, variable := range variables {
		switch variable.Type {
		case "interval":
			converter.convertIntervalVar(variable, dashboard)
		case "custom":
			converter.convertCustomVar(variable, dashboard)
		case "query":
			converter.convertQueryVar(variable, dashboard)
		case "const":
			converter.convertConstVar(variable, dashboard)
		case "datasource":
			converter.convertDatasourceVar(variable, dashboard)
		default:
			converter.logger.Warn("unhandled variable type found: skipped", zap.String("type", variable.Type), zap.String("name", variable.Name))
		}
	}
}

// convert interval variables
func (converter *JSON) convertIntervalVar(variable sdk.TemplateVar, dashboard *v1alpha1.DashboardSpec) {
	interval := templatingsModel.TemplateVar{
		Name:    variable.Name,
		Type:    variable.Type,
		Label:   variable.Label,
		Default: defaultOption(variable.Current),
		Values:  make([]string, 0, len(variable.Options)),
	}

	for _, opt := range variable.Options {
		interval.Values = append(interval.Values, opt.Value)
	}

	dashboard.Templatings = append(dashboard.Templatings, interval)
}

// convert custom variables
func (converter *JSON) convertCustomVar(variable sdk.TemplateVar, dashboard *v1alpha1.DashboardSpec) {
	custom := templatingsModel.TemplateVar{
		Name:       variable.Name,
		Label:      variable.Label,
		Type:       variable.Type,
		Default:    defaultOption(variable.Current),
		ValuesMap:  make(map[string]string, len(variable.Options)),
		AllValue:   variable.AllValue,
		IncludeAll: variable.IncludeAll,
	}

	for _, opt := range variable.Options {
		custom.ValuesMap[opt.Text] = opt.Value
	}

	dashboard.Templatings = append(dashboard.Templatings, custom)
}

// convert query variables
func (converter *JSON) convertQueryVar(variable sdk.TemplateVar, dashboard *v1alpha1.DashboardSpec) {
	datasource := ""
	if variable.Datasource != nil {
		datasource = *variable.Datasource
	}

	query := templatingsModel.TemplateVar{
		Name:       variable.Name,
		Label:      variable.Label,
		Type:       variable.Type,
		Datasource: datasource,
		Request:    variable.Query,
		Regex:      variable.Regex,
		IncludeAll: variable.IncludeAll,
		DefaultAll: variable.Current.Value == "$__all",
		AllValue:   variable.AllValue,
	}

	dashboard.Templatings = append(dashboard.Templatings, query)
}

// convert datasource variables
func (converter *JSON) convertDatasourceVar(variable sdk.TemplateVar, dashboard *v1alpha1.DashboardSpec) {
	datasource := templatingsModel.TemplateVar{
		Name:       variable.Name,
		Label:      variable.Label,
		Type:       variable.Type,
		Query:      variable.Query,
		Regex:      variable.Regex,
		IncludeAll: variable.IncludeAll,
	}

	dashboard.Templatings = append(dashboard.Templatings, datasource)
}

// convert const variables
func (converter *JSON) convertConstVar(variable sdk.TemplateVar, dashboard *v1alpha1.DashboardSpec) {
	text, _ := variable.Current.Text.(string)
	constant := templatingsModel.TemplateVar{
		Name:      variable.Name,
		Label:     variable.Label,
		Type:      variable.Type,
		Default:   text,
		ValuesMap: make(map[string]string, len(variable.Options)),
	}

	for _, opt := range variable.Options {
		constant.ValuesMap[opt.Text] = opt.Value
	}

	dashboard.Templatings = append(dashboard.Templatings, constant)
}

//convert rows
func (converter *JSON) convertPanels(panels []*sdk.Panel, dashboard *v1alpha1.DashboardSpec, isClusterCrd bool) {

	for _, panel := range panels {
		if panel.Type == "row" {
			for _, rowPanel := range panel.Panels {
				convertedPanel, ok := converter.convertDataPanel(rowPanel, isClusterCrd)
				if ok {
					dashboard.Panels = append(dashboard.Panels, convertedPanel)
				}
			}
		} else {
			convertedPanel, ok := converter.convertDataPanel(*panel, isClusterCrd)
			if ok {
				dashboard.Panels = append(dashboard.Panels, convertedPanel)
			}
		}
	}

}

//convert rows
func (converter *JSON) convertRows(rows []*sdk.Row, dashboard *v1alpha1.DashboardSpec, isClusterCrd bool) {

	for _, row := range rows {
		if row == nil {
			continue
		}
		panels := row.Panels
		if panels == nil || len(rows) == 0 {
			continue
		}
		for _, pl := range panels {
			convertedPanel, ok := converter.convertDataPanel(pl, isClusterCrd)
			if ok {
				dashboard.Panels = append(dashboard.Panels, convertedPanel)
			}
		}
	}
}

// convert different types of the given panel
func (converter *JSON) convertDataPanel(panel sdk.Panel, isClusterCrd bool) (panelsModel.Panel, bool) {
	switch panel.Type {
	case "graph":
		return converter.convertGraph(panel, isClusterCrd), true
	case "singlestat":
		return converter.convertSingleStat(panel, isClusterCrd), true
	case "gauge":
		return converter.convertCustom(panel, isClusterCrd), true
	case "bargauge":
		return converter.convertBarGauge(panel, isClusterCrd), true
	case "table":
		return converter.convertTable(panel, isClusterCrd), true
	case "text":
		return converter.convertText(panel), true
	default:
		converter.logger.Warn("unhandled panel type: skipped", zap.String("type", panel.Type), zap.String("title", panel.Title))
	}
	return panelsModel.Panel{}, false
}

// a graph panel
func (converter *JSON) convertGraph(panel sdk.Panel, isClusterCrd bool) panelsModel.Panel {
	// filled with values of the given fields
	var decimals int64
	if panel.GraphPanel.Decimals != nil {
		decimals = int64(*panel.GraphPanel.Decimals)
	}
	graph := &panelsModel.Panel{
		Title:       panel.Title,
		Type:        panel.Type,
		Decimals:    decimals,
		Colors:      defaultColors(),
		Description: panelDescription(panel.CommonPanel.Description),
		Id:          panelSpan(panel),
		Bars:        panel.GraphPanel.Bars,
		Lines:       panel.GraphPanel.Lines,
		Stack:       panel.GraphPanel.Stack,
		Legend:      converter.convertLegend(panel.GraphPanel.Legend),
	}

	if panel.Height != nil {
		graph.Height = *panel.Height
	}

	if panel.Datasource != nil {
		graph.Datasource = *panel.Datasource
	}

	// converts target
	if panel.GraphPanel.Targets != nil && len(panel.GraphPanel.Targets) > 0 {
		for index, target := range panel.GraphPanel.Targets {
			graphTarget := converter.convertTarget(target, isClusterCrd, index, panel.Type)
			if graphTarget == nil {
				continue
			}

			graph.Targets = append(graph.Targets, *graphTarget)
		}

	}

	// converts yaxes
	for _, yaxis := range panel.GraphPanel.Yaxes {
		d := int64(yaxis.Decimals)
		if d == 0 {
			d = 3
		}
		f := handleFormat(yaxis.Format)
		y := &panelsModel.Yaxis{
			Decimals: d,
			Format:   f,
		}
		graph.Yaxes = append(graph.Yaxes, *y)
		break
	}

	return *graph
}

func (converter *JSON) convertLegend(sdkLegend sdk.Legend) []string {
	var legend []string

	if !sdkLegend.Show {
		legend = append(legend, "hide")
	}
	if sdkLegend.AlignAsTable {
		legend = append(legend, "as_table")
	}
	if sdkLegend.RightSide {
		legend = append(legend, "to_the_right")
	}
	if sdkLegend.Min {
		legend = append(legend, "min")
	}
	if sdkLegend.Max {
		legend = append(legend, "max")
	}
	if sdkLegend.Avg {
		legend = append(legend, "avg")
	}
	if sdkLegend.Current {
		legend = append(legend, "current")
	}
	if sdkLegend.Total {
		legend = append(legend, "total")
	}
	if sdkLegend.HideEmpty {
		legend = append(legend, "no_null_series")
	}
	if sdkLegend.HideZero {
		legend = append(legend, "no_zero_series")
	}

	return legend
}

// singlestat panel
func (converter *JSON) convertSingleStat(panel sdk.Panel, isClusterCrd bool) panelsModel.Panel {
	singleStat := &panelsModel.Panel{
		Title:       panel.Title,
		Id:          panelSpan(panel),
		Type:        panel.Type,
		Description: panelDescription(panel.CommonPanel.Description),
		Format:      panel.SinglestatPanel.Format,
		Decimals:    int64(panel.SinglestatPanel.Decimals),
		ValueName:   panel.SinglestatPanel.ValueName,
	}

	if panel.Height != nil {
		singleStat.Height = *panel.Height
	}

	if panel.Datasource != nil {
		singleStat.Datasource = *panel.Datasource
	}

	if len(panel.SinglestatPanel.Colors) == 3 {
		singleStat.Colors = []string{
			panel.SinglestatPanel.Colors[0],
			panel.SinglestatPanel.Colors[1],
			panel.SinglestatPanel.Colors[2],
		}
	} else {
		singleStat.Colors = defaultColors()
	}

	var colorOpts []string
	if panel.SinglestatPanel.ColorBackground {
		colorOpts = append(colorOpts, "background")
	}
	if panel.SinglestatPanel.ColorValue {
		colorOpts = append(colorOpts, "value")
	}
	if len(colorOpts) != 0 {
		singleStat.Color = colorOpts
	}

	if panel.SinglestatPanel.SparkLine.Show && panel.SinglestatPanel.SparkLine.Full {
		singleStat.SparkLine = "full"
	}
	if panel.SinglestatPanel.SparkLine.Show && !panel.SinglestatPanel.SparkLine.Full {
		singleStat.SparkLine = "bottom"
	}

	// handles targets
	if panel.SinglestatPanel.Targets != nil && len(panel.SinglestatPanel.Targets) > 0 {
		for index, target := range panel.SinglestatPanel.Targets {
			graphTarget := converter.convertTarget(target, isClusterCrd, index, panel.Type)
			if graphTarget == nil {
				continue
			}

			singleStat.Targets = append(singleStat.Targets, *graphTarget)
		}

	}

	// handles gauge
	gauge := panelsModel.Gauge{
		MaxValue:         int64(panel.SinglestatPanel.Gauge.MaxValue),
		MinValue:         int64(panel.SinglestatPanel.Gauge.MinValue),
		Show:             panel.SinglestatPanel.Gauge.Show,
		ThresholdLabels:  panel.SinglestatPanel.Gauge.ThresholdLabels,
		ThresholdMarkers: panel.SinglestatPanel.Gauge.ThresholdMarkers,
	}

	singleStat.Gauge = gauge

	return *singleStat
}

// gauge
func (converter *JSON) convertCustom(panel sdk.Panel, isClusterCrd bool) panelsModel.Panel {
	// set options
	customPanel := &panelsModel.Panel{
		Decimals: 0,
		Title:    panel.Title,
		// Type:     panel.Type,
		Type:        "singlestat",
		Colors:      defaultColors(),
		Description: panelDescription(panel.CommonPanel.Description),
		Id:          panelSpan(panel),
	}

	// handles targets
	custom := *panel.CustomPanel
	if custom == nil {
		return *customPanel
	}

	var targets []sdk.Target

	if err := mapstructure.Decode(custom["targets"], &targets); err != nil {
		return *customPanel
	}

	for index, target := range targets {
		t := converter.convertTarget(target, isClusterCrd, index, panel.Type)
		customPanel.Targets = append(customPanel.Targets, *t)
	}

	return *customPanel
}

// bar gauge
func (converter *JSON) convertBarGauge(panel sdk.Panel, isClusterCrd bool) panelsModel.Panel {
	// set options
	barGaugePanel := &panelsModel.Panel{
		Options: panelsModel.BarGaugeOptions{
			Orientation: panel.BarGaugePanel.Options.Orientation,
			TextMode:    panel.BarGaugePanel.Options.TextMode,
			ColorMode:   panel.BarGaugePanel.Options.ColorMode,
			GraphMode:   panel.BarGaugePanel.Options.GraphMode,
			JustifyMode: panel.BarGaugePanel.Options.JustifyMode,
			DisplayMode: panel.BarGaugePanel.Options.DisplayMode,
			Content:     panel.BarGaugePanel.Options.Content,
			Mode:        panel.BarGaugePanel.Options.Mode,
		},
		Decimals:    0,
		Title:       panel.Title,
		Type:        panel.Type,
		Colors:      defaultColors(),
		Description: panelDescription(panel.CommonPanel.Description),
		Id:          panelSpan(panel),
	}

	// handles targets
	if panel.BarGaugePanel.Targets != nil && len(panel.BarGaugePanel.Targets) > 0 {
		for index, target := range panel.BarGaugePanel.Targets {
			barGaugeTarget := converter.convertTarget(target, isClusterCrd, index, panel.Type)
			if barGaugeTarget == nil {
				continue
			}

			barGaugePanel.Targets = append(barGaugePanel.Targets, *barGaugeTarget)
		}

	}

	return *barGaugePanel
}

// converts a table panel
func (converter *JSON) convertTable(panel sdk.Panel, isClusterCrd bool) panelsModel.Panel {
	tablePanel := &panelsModel.Panel{
		Title:       panel.Title,
		Id:          panelSpan(panel),
		Type:        panel.Type,
		Colors:      defaultColors(),
		Description: panelDescription(panel.CommonPanel.Description),
		Transparent: panel.Transparent,
		Decimals:    0,
	}

	if panel.Height != nil {
		tablePanel.Height = *panel.Height
	}

	if panel.Datasource != nil {
		tablePanel.Datasource = *panel.Datasource
	}

	if panel.TablePanel.Targets != nil && len(panel.TablePanel.Targets) > 0 {
		for index, target := range panel.TablePanel.Targets {
			graphTarget := converter.convertTarget(target, isClusterCrd, index, panel.Type)
			if graphTarget == nil {
				continue
			}

			tablePanel.Targets = append(tablePanel.Targets, *graphTarget)
		}
	}

	// hidden columns
	for _, columnStyle := range panel.TablePanel.Styles {
		if columnStyle.Type != "hidden" {
			continue
		}
		tablePanel.HiddenColumns = append(tablePanel.HiddenColumns, columnStyle.Pattern)
	}

	return *tablePanel
}

// converts a text panel
func (converter *JSON) convertText(panel sdk.Panel) panelsModel.Panel {
	textPanel := &panelsModel.Panel{
		Title:       panel.Title,
		Id:          panelSpan(panel),
		Type:        panel.Type,
		Colors:      defaultColors(),
		Description: panelDescription(panel.CommonPanel.Description),
		Transparent: panel.Transparent,
		Decimals:    0,
	}

	if panel.Height != nil {
		textPanel.Height = *panel.Height
	}

	if panel.TextPanel.Mode == "markdown" {
		textPanel.Markdown = panel.TextPanel.Content
	} else {
		textPanel.HTML = panel.TextPanel.Content
	}

	return *textPanel
}

func (converter *JSON) convertTarget(target sdk.Target, isClusterCrd bool, index int, tp string) *panelsModel.Target {
	// looks like a prometheus target
	converter.logger.Info("Only supported target type: prometheus", zap.Any("target", target))
	return converter.convertPrometheusTarget(target, isClusterCrd, index, tp)
}

func (converter *JSON) convertPrometheusTarget(target sdk.Target, isClusterCrd bool, index int, tp string) *panelsModel.Target {
	t := &panelsModel.Target{
		// RefID: target.RefID,
		RefID:          int64(index) + 1,
		Datasource:     target.Datasource,
		Hide:           target.Hide,
		LegendFormat:   handleLegendFormat(target.LegendFormat),
		Instant:        target.Instant,
		Format:         target.Format,
		Interval:       target.Interval,
		IntervalFactor: target.IntervalFactor,
	}

	// adjusts the query expression to adapt to the ks cluster
	transfered := transferExpr(target.Expr, isClusterCrd, tp)
	if transfered == "" {
		t.Expression = target.Expr
		return t
	}

	t.Expression = fmt.Sprintf("%s", transfered)
	t.Step = toString(target.Step)
	return t

}

func panelSpan(panel sdk.Panel) int64 {
	return int64(panel.ID)
}

func defaultOption(opt sdk.Current) string {
	if opt.Value == nil {
		return ""
	}

	return opt.Value.(string)
}

func handleZimu(z string) int64 {
	n, _ := strconv.Atoi(z)
	return int64(n) + 1
}

func toString(step int) string {
	// number := int(step / 60)
	// if number == 0 {
	// 	return strconv.Itoa(step) + "s"
	// }
	// if number > 60 {
	// 	return strconv.Itoa(number/60) + "h"
	// }
	// return strconv.Itoa(number) + "m"
	return "1m"
}

func handleFormat(f string) string {

	if f == "bytes" || f == "Bps" {
		f = "Byte"
	} else if f == "percent" || f == "percentunit" {
		f = "percent (0.0-1.0)"
	} else {
		f = "none"
	}
	return f
}

func handleLegendFormat(l string) string {
	badPat := regexp.MustCompile(`\{(\s+\w+\s+)\}`)
	if match := badPat.Match([]byte(l)); match {
		f := func(s string) string {
			stripReg := regexp.MustCompile(`\s+`)
			return stripReg.ReplaceAllString(s, "")
		}
		return badPat.ReplaceAllStringFunc(l, f)
	}
	return l
}

func panelDescription(des *string) string {
	d := ""
	if des != nil {
		d = *des
	}
	return d
}

func defaultColors() []string {
	return []string{"#60acfc", "#23c2db", "#64d5b2", "#d5ec5a", "#ffb64e", "#fb816d", "#d15c7f"}
}

func transferExpr(expr string, isClusterCrd bool, tp string) string {
	// example: sum (elasticsearch_jvm_memory_used_bytes{cluster="$cluster",name=~"$name"})/ sum (elasticsearch_jvm_memory_max_bytes{cluster="$cluster",name=~"$name"})
	// transfer to: sum by(cluster,name)  (elasticsearch_jvm_memory_used_bytes) / sum by(cluster,name)  (elasticsearch_jvm_memory_max_bytes) * 100

	// free the door if don't match a `[{}]` regex style
	pat := regexp.MustCompile(`[\{\}]`)
	if !pat.Match([]byte(expr)) {
		return ""
	}
	// handles $interval or $__interval
	pat1 := regexp.MustCompile(`\$_{0,2}interval`)
	if pat1.Match([]byte(expr)) {
		expr = pat1.ReplaceAllString(expr, "3m")
	}
	// if contains sum func, it can be divided by a flag(`isCluterCrd`) to distinguish with what kind of the resource
	pat2 := regexp.MustCompile(`sum\s+(\(\w+\{.*?\}\))`)
	if matchSum := pat2.Match([]byte(expr)); matchSum {
		return transferSumFunc(strings.TrimPrefix(expr, " "), isClusterCrd, tp)
	}

	// if contains irate/rate/count func, just removes `\{.*\}`
	pat3 := regexp.MustCompile(`\{.*?\}`)
	if matchCommon := pat3.Match([]byte(expr)); matchCommon {
		expr = pat3.ReplaceAllString(expr, "")
	}
	// if contains count, removes `>\d+`
	pat4 := regexp.MustCompile(`>\d+`)
	if matchCount := pat4.Match([]byte(expr)); matchCount {
		expr = pat4.ReplaceAllString(expr, "")
	}
	return expr

}

// convert sum func query to normal query which can be visualized
func transferSumFunc(expr string, isClusterCrd bool, tp string) string {
	pat := regexp.MustCompile(`\s+(\(\w+\{.*?\}\))`)
	f := func(s string) string {
		var trueByNames string
		params := regexp.MustCompile(`[\{\}]`).Split(s, -1)
		variable, byNameRaw := params[0], params[1]
		byNames := regexp.MustCompile(`(\w+)=`).FindAllString(byNameRaw, -1)
		if isClusterCrd || tp == "singlestat" {
			trueByNames = "cluster"
		} else {
			trueByNames = strings.ReplaceAll(strings.Join(byNames, ","), "=", "")
		}
		return fmt.Sprintf(" by(%s) %s)", trueByNames, variable)
	}

	return pat.ReplaceAllStringFunc(expr, f)

}
