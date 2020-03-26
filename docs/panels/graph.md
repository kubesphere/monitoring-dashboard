# Graph

Graph visualizes range query results into a linear graph


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| title | Name of the graph panel | string |
| type | Must be `graph` | string |
| id | Panel ID | int64 |
| description | Panel description | string |
| targets | A collection of queries | [][Target](common.md) |
| bars | Display as a bar chart | bool |
| colors | Set series color | []string |
| lines | Display as a line chart | bool |
| stack | Display as a stacked chart | bool |
| yaxes | Y-axis options | []Yaxis |
# Yaxis




| Field | Description | Scheme |
| ----- | ----------- | ------ |
| decimals | Limit the decimal numbers | int64 |
| format | Display unit | string |
