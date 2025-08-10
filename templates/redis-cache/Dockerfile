FROM redis:alpine

# Copy Redis configuration
COPY redis.conf /usr/local/etc/redis/redis.conf

# Expose default Redis port (informational)
EXPOSE "{{.Port}}"

# Health check (password-aware for local/dev)
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD {{if .Password}}redis-cli -a "{{.Password}}" ping{{else}}redis-cli ping{{end}} || exit 1

# Start Redis with custom config
CMD ["redis-server", "/usr/local/etc/redis/redis.conf"] 