FROM docker.io/golang:1.19-alpine AS base

WORKDIR /src
COPY . /src

RUN go build -o /bin/mock cmd/mock/main.go

FROM alpine
COPY --from=base /bin/mock /bin/mock
EXPOSE 8070 8071 8072
CMD /bin/mock