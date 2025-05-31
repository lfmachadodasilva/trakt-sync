# Use Node.js Alpine version 23 as the base image
FROM node:23-alpine

# Set the working directory
WORKDIR /app

# Install Bun and other dependencies
RUN apk add --no-cache curl bash && \
    curl -fsSL https://bun.sh/install | bash && \
    apk del curl

# Add Bun to PATH
ENV PATH="/root/.bun/bin:$PATH"

# Copy package files and install dependencies
COPY package.json bun.lock ./
RUN bun install

# Copy the rest of the application files
COPY . .

# Build the application
RUN bun run build

# Expose the application on port 3000
EXPOSE 3000

# Start the application
CMD ["bun", "run", "start"]