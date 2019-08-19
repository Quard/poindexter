FROM golang:latest as build
WORKDIR /app
COPY . .
ENV GO111MODULE on
ENV CGO_ENABLED=0
ENV GOOS=linux
RUN go build -v -o /poindexter cmd/poindexter/main.go

FROM alpine
RUN apk add --no-cache ca-certificates
COPY --from=build /poindexter /poindexter
CMD /poindexter
