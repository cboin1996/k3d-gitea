# Fetch gitea chart and dependencies                                                                                                                                                                                       
FROM alpine/helm as helm

RUN helm repo add gitea-charts https://dl.gitea.io/charts/ \
    && helm repo update

ARG GITEA_VERSION=5.0.3
RUN helm pull gitea-charts/gitea \
      --version $GITEA_VERSION \
      --untar \
    && helm dependency update gitea


# Build with the golang image
FROM golang:1.18-alpine AS build

# Add git
RUN apk add git

# Set workdir
WORKDIR /src

# Add dependencies
COPY app/go.mod .
COPY app/go.sum .
RUN go mod download

# Build
COPY app .
RUN CGO_ENABLED=0 go install .

# Generate final image
FROM alpine:3.15
RUN apk --update --no-cache add ca-certificates
COPY --from=build /go/bin/k3d-gitea /usr/local/bin/k3d-gitea
RUN mkdir -p /charts
COPY --from=helm /apps/gitea /charts
ENTRYPOINT [ "/usr/local/bin/k3d-gitea" ]
