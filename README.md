# k3d-gitea

Local development environment using a k3d cluster for testing/deploying gitea using argocd.

# Instructions
To setup everything, run:
```bash
task install
```
# Description
## [gitea-postgres-values.yaml](./helm/gitea-sqlite-values.yaml)
YAML file for deploying gitea with postgres enabled.

## [gitea-sqlite-values.yaml](./helm/gitea-sqlite-values.yaml)
YAML file for deploying gitea with sqlite enabled.

## [main.go](app/main.go)
Entrypoint for golang cli for configuring postgres for use with gitea. 
See [database preparation](https://docs.gitea.io/en-us/database-prep/#postgresql) for
examples of configuring external database with gitea.

# Optional Instructions
## Clone Gitea Helm Chart
```bash
task gitea:clone
```
This command will clone the gitea helm chart and build it's dependencies.