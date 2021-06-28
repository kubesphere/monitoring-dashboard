package v1alpha2

import (
	"testing"

	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/util/json"
	"kubesphere.io/monitoring-dashboard/api/v1alpha2/panels"
)

var (
	datasource = "prometheus"

	js = `
	{
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
	`
	expected = DashboardSpec{
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
	}
)

func TestDashboardSpecSerde(t *testing.T) {

	var actual DashboardSpec

	err := json.Unmarshal([]byte(js), &actual)
	if err != nil {
		panic(err)
	}

	require.EqualValues(t, expected, actual)
}
