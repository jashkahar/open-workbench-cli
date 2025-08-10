FROM node:18-alpine

WORKDIR /app

COPY package*.json ./

# Install full deps to build; fall back if no lockfile
RUN sh -lc 'if [ -f package-lock.json ]; then npm ci; else npm install; fi'

COPY . .

RUN npm run build

# Vite preview default port (overridable)
EXPOSE 4173

# Local/dev-friendly start
CMD ["npm", "run", "preview", "--", "--host", "0.0.0.0", "--port", "4173"]