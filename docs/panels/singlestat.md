# SingleStat

SingleStat shows instant query result


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| title | Name of the signlestat panel | string |
| type | Must be `singlestat` | string |
| id | Panel ID | int64 |
| targets | A collection of queries | [][Target](common.md) |
| decimals | Limit the decimal numbers | int64 |
| format | Display unit | string |
