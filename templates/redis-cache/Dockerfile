FROM redis:alpine

# Copy Redis configuration
COPY redis.conf /usr/local/etc/redis/redis.conf

# Expose port
EXPOSE {{.Port}}

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD redis-cli ping || exit 1

# Start Redis with custom config
CMD ["redis-server", "/usr/local/etc/redis/redis.conf"] 