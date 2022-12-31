FROM golang:1.19-alpine AS build

WORKDIR /app

COPY ./ ./

# Install dependencies
RUN go mod download && \
  # Build the app
  GOOS=linux GOARCH=amd64 go build -o main && \
  # Make the final output executable
  chmod +x ./main

FROM alpine:latest

# Install os packages
WORKDIR /app

COPY --from=build /app/ .

CMD ["./main"]