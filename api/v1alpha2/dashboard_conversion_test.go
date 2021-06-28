package v1alpha2

import (
	"testing"

	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/util/json"
	"kubesphere.io/monitoring-dashboard/api/v1alpha1"
	v1alpha1panels "kubesphere.io/monitoring-dashboard/api/v1alpha1/panels"
	"kubesphere.io/monitoring-dashboard/api/v1alpha2/panels"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

var (
	v1alpha1DashboardString = `
{
	"spec": {
		"datasource": "prometheus",
		"panels": [
			{
				"id": 1,
				"title": "Current QPS",
				"type": "singlestat",
				"targets": [{
					"expr": "vector(1)"
				}]
			},
			{
				"id": 2,
				"title": "Connections and Threads",
				"type": "row"
			},
			{
				"id": 3,
				"title": "Connections",
				"type": "graph",
				"targets": [{
					"expr": "vector(1)"
				}]
			}
		],
		"title": "MySQL Overview"
	}
}  
`

	v1alpha2DashboardString = `
{
	"spec":{
		"panels": [
		{
			"id": 1,
			"title": "Current QPS",
			"type": "singlestat",
			"datasource": "prometheus",
			"targets": [{
				"expr": "vector(1)"
			}]
		},
		{
			"id": 2,
			"title": "Connections and Threads",
			"type": "row",
			"datasource": "prometheus"
		},
		{
			"id": 3,
			"title": "Connections",
			"type": "graph",
			"targets": [{
				"expr": "vector(1)"
			}],
			"datasource": "prometheus"
		}
	],
	"title": "MySQL Overview"
	}	
}
`
	expectedV1alpha1Dashboard = v1alpha1.Dashboard{
		Spec: v1alpha1.DashboardSpec{
			Title:      "MySQL Overview",
			DataSource: "prometheus",
			Panels: []v1alpha1.Panel{{
				PanelMeta: v1alpha1.PanelMeta{
					Id:    1,
					Type:  "singlestat",
					Title: "Current QPS",
				},
				Targets: []v1alpha1panels.Target{{
					Expression: "vector(1)",
				}},
			}, {
				PanelMeta: v1alpha1.PanelMeta{
					Id:    2,
					Title: "Connections and Threads",
					Type:  "row",
				},
			}, {
				PanelMeta: v1alpha1.PanelMeta{
					Id:    3,
					Title: "Connections",
					Type:  "graph",
				},
				Targets: []v1alpha1panels.Target{{
					Expression: "vector(1)",
				}},
			}},
		},
	}

	expectedV1alpha2Dashboard = Dashboard{
		Spec: DashboardSpec{
			Title: "MySQL Overview",
			Panels: []*panels.Panel{
				{
					CommonPanel: panels.CommonPanel{
						Id:         1,
						Type:       "singlestat",
						Title:      "Current QPS",
						Datasource: &datasource,
						Targets: []panels.Target{{
							Expression: "vector(1)",
						}},
					},
				},
				{
					CommonPanel: panels.CommonPanel{
						Id:         2,
						Type:       "row",
						Title:      "Connections and Threads",
						Datasource: &datasource,
					},
					// RowPanel: &panels.RowPanel{},
				},
				{
					CommonPanel: panels.CommonPanel{
						Id:         3,
						Title:      "Connections",
						Type:       "graph",
						Datasource: &datasource,
						Targets: []panels.Target{{
							Expression: "vector(1)",
						}},
					},
				}},
		},
	}
)

func TestV1alpha1ToV1alpha2(t *testing.T) {

	var v1alpha1ActualDashboard v1alpha1.Dashboard
	var v1alpha2ActualDashboard Dashboard

	err := json.Unmarshal([]byte(v1alpha1DashboardString), &v1alpha1ActualDashboard)
	if err != nil {
		panic(err)
	}

	v1alpha2ActualDashboard.ConvertFrom(conversion.Hub(&v1alpha1ActualDashboard))

	require.EqualValues(t, expectedV1alpha2Dashboard, v1alpha2ActualDashboard)
}

func TestV1alpha2ToV1alpha1(t *testing.T) {

	var v1alpha1ActualDashboard v1alpha1.Dashboard
	var v1alpha2ActualDashboard Dashboard

	err := json.Unmarshal([]byte(v1alpha2DashboardString), &v1alpha2ActualDashboard)
	if err != nil {
		panic(err)
	}

	v1alpha2ActualDashboard.ConvertTo(conversion.Hub(&v1alpha1ActualDashboard))

	require.EqualValues(t, expectedV1alpha1Dashboard, v1alpha1ActualDashboard)
}
