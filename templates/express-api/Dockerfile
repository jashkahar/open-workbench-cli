###### Local/dev-friendly multi-stage build to ensure TypeScript compilation and native deps reliability

# Stage 1: build with devDependencies
FROM node:18-bullseye-slim AS builder
WORKDIR /app
COPY package*.json ./
RUN sh -lc 'if [ -f package-lock.json ]; then npm ci; else npm install; fi'
COPY . .
RUN npm run build

# Stage 2: production runtime with only prod deps
FROM node:18-bullseye-slim AS runner
WORKDIR /app
ENV NODE_ENV=development \
    # Defaults for local Docker Desktop to reach host DBs
    DB_HOST=host.docker.internal \
    MONGO_URI=mongodb://host.docker.internal:27017/express_api \
    PORT=3001

COPY package*.json ./
RUN sh -lc 'if [ -f package-lock.json ]; then npm ci --omit=dev; else npm install --omit=dev; fi'

# Copy compiled app
COPY --from=builder /app/dist ./dist

# Expose API port (unique for local/dev)
EXPOSE 3001

CMD ["node", "dist/index.js"]