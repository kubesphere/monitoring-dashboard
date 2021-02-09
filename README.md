# KubeSphere Monitoring Dashboard

The project is inspired by [Grafana](http://grafana.com/) but with significant difference in data persistence, multitenancy supports and dashboard template sharing to fit KubeSphere's context. It is not a replacement to Grafana. It requires KubeSphere backend and frontend to work.

This repo is aimed at KubeSphere developers who want to understand dashboard data model, concepts and usage and how to contribute, as the custom monitoring feature is introduced as of v3.0.  

## Table of contents

- [Get Started](#get-started)
  - [Prerequisites](#prerequisites)
  - [Quick Start](#quick-start)
- [Concept and Design](#concept-and-design)
  - [Data Model](#data-model)
  - [Data Source](#data-source)
  - [Multi-tenancy](#multi-tenancy)
  - [Dashboard Template](#dashboard-template)
- [Manual](#manual)
  - [Query](#query)
  - [Panels](#panels-1)
  - [Time Range](#time-range)
  - [Variables](#variables)
- [Development](#development)
  - [APIs](#apis)
  - [Backend](#backend)
  - [Frontend](#frontend)
- [Contributing](#contributing)
  - [Dashboard Gallery](#dashboard-gallery)

## Get Started

## Prerequisites

- Kubernetes v17.0+
- KubeSphere v3.0+

## Quick Start

TODO(@FeynmanZhou)

The stack making these possible includes KubeSphere backend, console and custom resources for dashboards. 

- [KubeSphere Backend](https://github.com/kubesphere/kubesphere): proxies metrics query, ensures data isolation over namespaces.
- [KubeSphere Console](https://github.com/kubesphere/console): renders the dashboard with data and charts.
- CustomResourceDefinition for Dashboard: defines dashboard data model.

## Concept and Design

### Data Model

Dashboards are backed by the custom resource definition (CRD) **`Dashboard`**. It is compromised of three components: metadata, panels and templatings. Below is an example: 

```yaml
apiVersion: monitoring.kubesphere.io/v1alpha1
kind: Dashboard
metadata:
  name: mysql-overview
  namespace: default
spec:
  title: MySQL Overview
  description: MySQL dashboard for the mysql exporter
  time:
    from: now-1h
    to: now
  datasource: prometheus
  panels:
  - type: singlestat
    title": Instance Up
    targets:
    - expr: mysql_up{release="$release"}
      instant: "true"
  - type: graph
    title: mysql  disk reads vs writes
    targets:
    - expr": irate(mysql_global_status_innodb_data_reads{release="$release"}[10m])
      legendFormat": reads
    - expr": irate(mysql_global_status_innodb_data_writes{release="$release"}[10m])
      legendFormat": write
  templatings:
  - name: release
    query: label_values(mysql_up,release)
    type: query
    sort: 0
```

#### Metadata

|Name|Desc|
|---|---|
|`spec.title`|dashboard title|
|`spec.description`|dashboard description|
|`spec.time`|time range for display. see [Time Range](#time-range) for more info|
|`spec.datasource`|data source to query, defaults to Prometheus|

#### Panels

The `spec.panels` defines a collection of panels. Panels are build blocks of a dashboard. Currently supported panels are row, singlestat, graph and table. See [Query](#query) and [Panels](#panels) for more info.

#### Templatings

The `spec.templatings` defines a collection of variables. It is convenient to use variables in query expressions. See [Variable](#variable) for more info.

> Note that the data model is heavily inspired by [Grafana JSON model](https://grafana.com/docs/grafana/latest/reference/dashboard/) to gain compatibility. However, to adapt it to KubeSphere's context, we may bring in new fields and break changes.

### Data Source

Note: we currently only support Prometheus as data source.

### Multi-tenancy

Metrics data should be isolated across namespaces, which means namespace members can only view metrics in the namespace they belong to. This is implemented in the phase of [querying](#query). Any user-written expression will be mutated to make sure no query outside the scope of the namespace. 

Take the datasource Prometheus for example, KubeSphere backend will add a namespace matcher, i.e. `<metric_name>{namespace=<namespace_name>}`, to each query.

### Dashboard Template

Dashboard is represented by a custom resource object. Select templates from [contrib/gallery](contrib/gallery), and run the following command to import:

```
kubectl apply --namespace <NAMESPACE> -f contrib/gallery/<TEMPLATE_YAML_FILE>
```

You can also open source your template and contribute to Dashboard Gallery. Templates in Dashboard Gallery will be shipped with KubeSphere.

## Manual

### Query

Except for the panel `row`, each panel accepts at least one data source query. It allows you to input query expressions and fetch metrics data.

You can use placeholders in queries. See [Variables](#variables) for more information.

### Panels

#### Chart

We currently support four types of panels

- Row
- Singlestat
- Graph

#### Legend

### Time Range

Time range specifies current dashboard time for display. The following are examples in use.

|Example|From|To|
|--|--|--|
|Last 5 minutes|now-5m|now|
|Today|now-1d|now|
|This week|now-1w|now|
|Last month|now-2M|now-1M|

### Variables

|Query|Desc|
|--|--|
|label_values(metric, label)|Returns a list of label values for the label in the specified metric.|

## converter

we support a converter tool which can be used for importing dashboards from Grafana dashboard templates.

### how to use
if we want to convert a dashboard json template to a k8s manifest, you can use make cmdline like this below:
```
	make converter -isClusterCrd=$(IS_CLUSTER_CRD) -namespace=$(NAMESPACE) -inputPath=$(INPUT) -outputPath=$(OUTPUT)
```
or:
```
	go run ./cmd/converter -isClusterCrd=$(IS_CLUSTER_CRD) -namespace=$(NAMESPACE) -inputPath=$(INPUT) -outputPath=$(OUTPUT) -name=$(Name)
```
### Usage of converter
```
Usage of converter:
  -inputPath string
        a input path for the converter to look for jobs (default "./manifests/inputs")
  -isClusterCrd
        a flag that defines whether build the cluster dashboard resource or not
  -name string
        name of the dashboard resource (default "your file name")
  -namespace string
        namespace of the dashboard resource (default "default")
  -outputPath string
        a output path for the converter to store manifests (default "./manifests/outputs")
```
## Development

### APIs

For dashboard APIs, see [docs/crd.md](docs/crd.md)

The project is built with kubebuilder v2.3.0.

### Backend

If you find some fields should be included in the CRD, edit [api package](https://github.com/kubesphere/monitoring-dashboard/tree/master/api/v1alpha1) and regenerate the project by `make`. Kubebuilder is the tool we are using.

If you find bugs or want to add new APIs, implement new datasources to KubeSphere, read KubeSphere [developer guides](https://github.com/kubesphere/community/tree/master/developer-guide/concepts-and-designs) for monitoring.

### Frontend

@TODO(justahole)

## Contributing

### Dashboard Gallery

If you want to share your dashboard templates, you can submit dashboard template yaml files to gallery folder with an elaborate readme. Outstanding templates will be selected to ship with KubeSphere future releases.