FROM golang:bookworm
WORKDIR /app
COPY . .
RUN apt update -y && apt upgrade -y
CMD ["go", "run", "."]
