# Multi-stage build untuk aplikasi React Vite production
FROM node:22.15-alpine3.21 AS builder

# Set working directory
WORKDIR /app

# Install dependencies untuk build process
RUN apk add --no-cache git

# Copy package.json dan package-lock.json (jika ada)
COPY package*.json ./

# Install dependencies
RUN npm ci --only=production --silent

# Copy source code
COPY . .

# Build aplikasi untuk production
RUN npm run build

# Production stage dengan Nginx
FROM nginx:1.28.0-alpine3.21

# Install additional packages untuk production environment
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    && update-ca-certificates

# Remove default nginx configuration
RUN rm /etc/nginx/conf.d/default.conf

# Copy custom nginx configuration
COPY nginx/nginx.conf /etc/nginx/conf.d/

# Copy built application dari builder stage
COPY --from=builder /app/dist /usr/share/nginx/html

# Create nginx user dan set proper permissions
RUN chown -R nginx:nginx /usr/share/nginx/html \
    && chmod -R 755 /usr/share/nginx/html

# Create directory untuk nginx logs
RUN mkdir -p /var/log/nginx \
    && chown -R nginx:nginx /var/log/nginx

# Expose port 80 dan 443
EXPOSE 80 443

# Health check untuk container monitoring
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost/ || exit 1

# Start nginx
CMD ["nginx", "-g", "daemon off;"]