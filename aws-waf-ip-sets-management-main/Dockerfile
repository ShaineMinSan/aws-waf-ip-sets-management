# Use the minimal base image
FROM alpine:latest

# Set the current working directory inside the container
WORKDIR /app

# Copy the compiled Go binary from the local file system
COPY backend/backend /app/backend

# Copy the frontend files
COPY frontend /app/frontend

# # Copy the environment configuration file
# COPY config.env /app/config.env

# Expose the port the app runs on
EXPOSE 3000

# Command to run the executable
CMD ["/app/backend"]
