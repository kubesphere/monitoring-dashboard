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
`
	var expected = DashboardSpec{
		Title:      "MySQL Overview",
		DataSource: "prometheus",
		Panels: []panels.Panel{
			{
				Id:    1,
				Type:  "singlestat",
				Title: "Current QPS",
				Targets: []panels.Target{{
					Expression: "vector(1)",
				}},
			},
			{

				Id:    2,
				Title: "Connections and Threads",
				Type:  "row",
			},
			{

				Id:    3,
				Title: "Connections",
				Type:  "graph",
				Targets: []panels.Target{{
					Expression: "vector(1)",
				}},
			}},
	}

	var actual DashboardSpec

	err := json.Unmarshal([]byte(jsonString), &actual)
	if err != nil {
		panic(err)
	}

	require.Equal(t, expected, actual)
}
