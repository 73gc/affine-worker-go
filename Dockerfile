# Step 1: Use the official golang:1.22.10 image as the base image
FROM golang:1.22.10

# Step 2: Set the working directory inside the container
WORKDIR /app

# Step 3: Copy the files of the application code
COPY . .

# Step 4: Install dependencies
RUN go mod tidy

# Step 5: Build the golang project
RUN go build -o affine-worker-go

# Step 6: Expose the port the app runs on
EXPOSE 8888

# Step 7: Define the command to run the app
CMD ["./affine-worker-go"]