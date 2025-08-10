# Use Node.js 18 Alpine as base image
FROM node:18-alpine

WORKDIR /app

COPY package*.json ./
RUN sh -lc 'if [ -f package-lock.json ]; then npm ci; else npm install; fi'

COPY . .

RUN npm run build

EXPOSE 3000

# Local/dev-friendly Nuxt 3 preview server
CMD ["npm", "run", "preview", "--", "--host", "0.0.0.0", "--port", "3000"]