# API Docs
This Document documents the types introduced by the Monitoring Dashboard to be consumed by users.
> Note this document is generated from code comments. When contributing a change to this document please do so by changing the code comments.
## Table of Contents
* [Dashboard](#dashboard)
* [DashboardList](#dashboardlist)
* [DashboardSpec](#dashboardspec)
* [Templating](#templating)
* [Time](#time)
## Dashboard

Dashboard is the Schema for the dashboards API


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | metav1.ObjectMeta |
| spec |  | DashboardSpec |

[Back to TOC](#table-of-contents)
## DashboardList

DashboardList contains a list of Dashboard


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | metav1.ListMeta |
| items |  | []Dashboard |

[Back to TOC](#table-of-contents)
## DashboardSpec

DashboardSpec defines the desired state of Dashboard


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| title | Dashboard title | string |
| description | Dashboard description | string |
| datasource | Dashboard datasource | string |
| time | Time range for display | Time |
| panels | Collection of panels. Panel is one of [Row](row.md), [Singlestat](#singlestat.md) or [Graph](graph.md) | []Panel |
| templating | Templating variables | []Templating |

[Back to TOC](#table-of-contents)
## Templating

Templating defines a variable, which can be used as a placeholder in query


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| name | Variable name | string |
| query | Set variable values to be the return result of the query | string |

[Back to TOC](#table-of-contents)
## Time

Time ranges of the metrics for display


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| from | Start time in the format of `^now([+-][0-9]+[smhdwMy])?$`, eg. `now-1M`. It denotes the end time is set to the last month since now. | string |
| to | End time in the format of `^now([+-][0-9]+[smhdwMy])?$`, eg. `now-1M`. It denotes the start time is set to the last month since now. | string |

[Back to TOC](#table-of-contents)
