# Use Node.js 18 Alpine as base image
FROM node:18-alpine

# Set working directory
WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm ci --only=production

# Copy application code
COPY . .

# Build the application (if there's a build script)
RUN npm run build || true

# Expose port
EXPOSE 3000

# Start the application
CMD ["npm", "start"] 