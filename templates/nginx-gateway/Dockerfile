FROM nginx:alpine

# Install curl for healthcheck in local/dev
RUN apk add --no-cache curl

# Copy nginx configuration
COPY nginx.conf /etc/nginx/nginx.conf

# Copy default site configuration
COPY default.conf /etc/nginx/conf.d/default.conf

# Expose default HTTP port (informational)
EXPOSE 80

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost/health || exit 1

CMD ["nginx", "-g", "daemon off;"]