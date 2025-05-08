FROM golang:1.24 AS build

WORKDIR /go/src/app
COPY . .

RUN go mod download
RUN go vet -v ./src/cmd
RUN go test -v ./src/cmd

RUN CGO_ENABLED=0 go build -o /go/bin/app ./src/cmd

FROM gcr.io/distroless/static-debian12

COPY --from=build /go/bin/app /
CMD ["/app"]