FROM golang:alpine AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

# flags for decrise binary size, by removing debug stuff
RUN go build -ldflags="-s -w" -o /app/lalachka ./cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=build /app/lalachka .

RUN mkdir templates


COPY --from=build /app/templates/home.html ./templates/home.html
COPY --from=build /app/templates/login.html ./templates/login.html
COPY --from=build /app/templates/startup.html ./templates/startup.html


CMD ["./lalachka"]