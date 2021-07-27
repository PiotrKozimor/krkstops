FROM docker.io/golang:1.16-alpine AS base

WORKDIR /src
COPY . /src

RUN go build -o /bin/mock cmd/mock/main.go

FROM alpine
COPY --from=base /bin/mock /bin/mock
EXPOSE 8080 8070 8090
CMD /bin/mock