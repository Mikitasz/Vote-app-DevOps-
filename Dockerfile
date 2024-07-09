# STAGE 1:
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum .
RUN go mod download

COPY . .

RUN go build -o vote_app


#STAGE 2:

FROM scratch

WORKDIR /app

COPY --from=builder /app/vote_app .
COPY --from=builder /app/template ./template

# Define environment variables for the runtime stage
ENV DB_NAME=vote_app \
    DB_HOST=vote_db \
    DB_PASSWORD=mysecret \
    DB_USER=postgres \
    DB_PORT=5432

EXPOSE 8080

ENTRYPOINT ["./vote_app"]
