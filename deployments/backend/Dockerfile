FROM golang:1.13 as builder

WORKDIR /source

ADD go.mod .
ADD go.sum .
RUN go mod download

ADD . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app main.go

FROM alpine

WORKDIR /apps

COPY --from=builder /source/app /apps/app
COPY data /apps/data
COPY web/app/build /apps/static

CMD ["/apps/app"]
