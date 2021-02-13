FROM golang:alpine

RUN go version
ENV GOPATH=/

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod tidy
# Download all dependencies. Dependencies will be cached if the go.mod and the go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the working Directory inside the container
COPY ./ ./

RUN ls
# Build the Go app
RUN go build -o octoserver -v ./cmd

# Expose port 8000 to the outside world
EXPOSE 8000

#Command to run the executable
CMD ["./octoserver"]