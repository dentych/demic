# Build image
FROM golang:alpine AS build
WORKDIR /app
COPY . .
ARG API_BASE_URL
RUN go build -o demic main.go

# App
FROM alpine
WORKDIR /app
COPY --from=build /app/demic .
EXPOSE 8080
CMD ./demic
