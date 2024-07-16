# Vote app

## Table of Contents

- [TODO](#todo)
- [Environment Variables](#environment-variables)
- [Implementations of Deployment](#implementations-of-deployment)
  - [Docker-compose](#docker-compose)
  - [Kubernetes](#kubernetes)
- [Preview web page](#preview)

## TODO

- [ ] Implement volumes in Kubernetes
- [ ] Figure out how to create automatic tags in actions
- [ ] Maybe add some simple tests
- [ ] Implement end-to-end CI/CD pipeline in GitHub Actions
- [ ] Deploy web app on AWS on EKS using Terraform
- [ ] Implement monitoring, like Grafana
- [x] Implement ingrees-nginx load

## Environmetn variables

    DB_NAME --- name of database
    DB_PASSWORD --- password for user postgres
    DB_USER --- username, by default `postgres`
    DB_HOST --- host name of database
    DB_PORT --- port where the database starts

## Implementations of Deployment ( actual stage )

I have implemented two deployment options for my simple web app:

- Kubernetes (localy on minikube)
- Docker-compose (localy)

I chose **Go** to develop the web app because it's a compiled language, which allows for clearer implementation in future CI/CD pipelines.

For containerization, I've used a multi-stage [Dockerfile](https://github.com/Mikitasz/Vote-app-DevOps-/blob/main/Dockerfile) to build a Docker image starting from a `scratch` base. The total size of the resulting image is just 8.49MB. You can find the Docker image on [Docker Hub](https://hub.docker.com/repository/docker/mikitasz/golang-vote-app/general).

To automate building and pushing the Docker image, I've set up [GitHub Actions](https://github.com/Mikitasz/Vote-app-DevOps-/blob/main/.github/workflows/publish.yaml). This workflow simplifies the process of building and pushing the image to the Docker registry.

### Docker-compose (Localy)

This is the [docker-compose.yaml](https://github.com/Mikitasz/Vote-app-DevOps-/blob/main/docker-compose.yaml) file.

To modify variables, edit the `environment` section in the `web` service. By default, it is configured as follows:

      DB_HOST: db
      DB_NAME: vote_app
      DB_PORT: 5432
      DB_USER: postgres
      DB_PASSWORD: mysecret

To run application:
`docker-compose up -d`

To down:
`docker-compose down`

### Kubernetes

#### Localy using minikube

This is the [kubernetes](https://github.com/Mikitasz/Vote-app-DevOps-/tree/main/kubernetes) folder, containing 4 files:

- **secrets.yaml**: Specifies secret variables.
- **configmap.yaml**: Specifies non-secret variables.
- **webdeployment.yaml**: Includes Deployment and Service definitions to deploy the Golang web app.
- **postgresql.yaml**: Includes definitions to deploy the PostgreSQL database.

To modify variables, edit the `secret.yaml` and `configmap.yaml`.

To run the application:
`kubectl apply -f kubernetes`
