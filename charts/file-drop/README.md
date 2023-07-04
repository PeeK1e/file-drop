# File Drop Helm

This Chart deploys the filedrop service.

All values and their usage is described in the `values.yaml` of the chart.

Running without an ingress is not supported at the moment but still possible. Feel free to do so.

## Deployment

### Requirements
#### Database
You'll need to have a running Postgres Database reachable from your Cluster.

If not already existent, create a Secret containing the Username and Password for the Database. And set the keys in `extraEnvSecret.database.usernameKey` and `extraEnvSecret.database.passwordKey`. Also set the hostname of the DB in `extraEnvSecret.database.hostname`.

Alternatively create a Secret Containing all your environment variables needed to launch the application and reference it in the `extraEnvSecret.name` field and set `extraEnvSecret.create = false`.

#### Storage
Your Storage <text style="color:red">**must**</text> support multiattach (**ReadWriteMany**) Volumes so . Since the api and cleaner both need to access the files.

### Install

Clone the repo and run  the helm install command from within the chart's directory

```sh
helm upgrade --create-namespace --install filedrop -n filedrop-1 -f values.yaml .
```

## Architecture

```mermaid
graph TD

subgraph "Ingress"
    ingress(file-drop-ingress)
end

subgraph "Service"
    svc-api(file-drop-api)
    svc-web(file-drop-web)
end

subgraph "Deployment"
    deployment-web(file-drop-web)
    deployment-api(file-drop-api)
    deployment-cleaner(file-drop-cleaner)
end

subgraph "Job"
    migrations-job(file-drop-migrations-lz7e)
end

subgraph "Secret"
    secret(file-drop-secret)
end

subgraph "Database/Statefulset"
    db(Postgres)
end

subgraph "PVC"
    pvc[data]
end

Ingress <--> Service

Service <--> deployment-api
Service <--> deployment-web

Deployment --> Secret
Job --> Secret

Job --> Database/Statefulset

deployment-cleaner --> PVC
deployment-api --> PVC

deployment-api --> Database/Statefulset

```