package v1alpha1

import (
	"testing"

	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/util/json"
	"kubesphere.io/monitoring-dashboard/api/v1alpha1/panels"
)

func TestDashboardSpecSerde(t *testing.T) {
	var jsonString = `
{
	"datasource": "prometheus",
	"panels": [
		{
			"id": 1,
			"title": "Current QPS",
			"type": "singlestat"			
		},
		{
			"id": 2,
			"title": "Connections and Threads",
			"type": "row"
		},
		{
			"id": 3,
			"title": "Connections",
			"type": "graph"
		}
	],
	"title": "MySQL Overview"
}
`
	var expected = DashboardSpec{
		Title:      "MySQL Overview",
		DataSource: "prometheus",
		Panels: []Panel{{
			SingleStat: &panels.SingleStat{
				Id:    1,
				Type:  "singlestat",
				Title: "Current QPS",
			},
		}, {
			Row: &panels.Row{
				Id:    2,
				Title: "Connections and Threads",
				Type:  "row",
			},
		}, {
			Graph: &panels.Graph{
				Id:    3,
				Title: "Connections",
				Type:  "graph",
			},
		}},
	}

	var actual DashboardSpec

	err := json.Unmarshal([]byte(jsonString), &actual)
	if err != nil {
		panic(err)
	}

	require.Equal(t, expected, actual)
}
