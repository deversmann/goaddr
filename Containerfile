FROM registry.access.redhat.com/ubi8/go-toolset:latest AS build
USER root
WORKDIR /work
COPY go.mod go.sum main.go ./
ADD log/ ./log
RUN go mod tidy
RUN go build

FROM registry.access.redhat.com/ubi8/ubi-minimal:latest
COPY --from=build /work/goaddr .
CMD ["./goaddr"]