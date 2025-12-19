# Stage 1: Build Frontend
FROM node:20 AS frontend-builder

WORKDIR /app/frontend

# Copy frontend package files for dependency caching
COPY frontend/package*.json ./

# Install all dependencies (including dev dependencies for build)
RUN npm ci

# Copy frontend source code
COPY frontend/ .

# Build the React application
RUN npm run build

# Stage 2: Build Backend and Create Production Image
FROM node:20 AS production

WORKDIR /app

# Copy backend package files
COPY backend/package*.json ./

# Install all dependencies (including dev dependencies for TypeScript build)
RUN npm ci

# Copy backend source code
COPY backend/ .

# Build TypeScript backend
RUN npm run build

# Remove dev dependencies and clean cache to reduce image size
RUN npm prune --production && npm cache clean --force

# Copy frontend build output from Stage 1
COPY --from=frontend-builder /app/frontend/dist ./frontend/dist

# # Create non-root user for security
# RUN addgroup -g 1001 -S nodejs && \
#     adduser -S nodejs -u 1001 && \
#     chown -R nodejs:nodejs /app

# USER nodejs

# Expose the port the backend runs on (default, can be overridden)
EXPOSE 3000

# Set build-time environment variables only
ENV NODE_ENV=production
ENV FRONTEND_DIST=/app/frontend/dist

# # Default runtime environment variables (can be overridden at runtime)
# ENV PORT=3000
# ENV HOST=0.0.0.0

# # Health check (uses PORT env var)
# HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
#     CMD node -e "const port = process.env.PORT || 3000; require('http').get('http://localhost:' + port + '/health', (r) => {process.exit(r.statusCode === 200 ? 0 : 1)})"

# Start the backend server
CMD ["npm", "start"]
