# k3d-gitea

Local development environment using a k3d cluster for testing/deploying gitea.

# Instructions

## Helm Template Generation
To generate the deployment yaml, clone the gitea helm chart into the root
of this repository from [here](https://gitea.com/gitea/helm-chart/)

Generating the yaml is done by running the following command
within the helm-chart directory you just cloned:

```bash
helm template . --values ../sqlite-values.yaml
```


# Description

## [gitea-postgres-values.yaml](gitea-sqlite-values.yaml)
YAML file for deploying gitea with postgres enabled.

## [gitea-sqlite-values.yaml](gitea-sqlite-values.yaml)
YAML file for deploying gitea with sqlite enabled.

## [postgres-values.yaml](gitea-sqlite-values.yaml)
Helm values file for bitnami's postgres chart: https://github.com/bitnami/charts/tree/master/bitnami/postgresql

## [docker-compose.yaml](docker-compose.yml)
Docker compose for installing postgres locally using docker.

## [main.go](main.go)
Entrypoint for golang cli for configuring postgres for use with gitea. 
See [database preparation](https://docs.gitea.io/en-us/database-prep/#postgresql) for
examples of configuring databases with gitea.