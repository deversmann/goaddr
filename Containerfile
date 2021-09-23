FROM registry.access.redhat.com/ubi8/go-toolset:latest AS build
USER root
WORKDIR /work
COPY . ./
RUN go mod tidy
RUN go build

FROM registry.access.redhat.com/ubi8/ubi-minimal:latest
COPY --from=build /work/goaddr .
CMD ["./goaddr"]