# Build Stage
FROM node:18-alpine AS builder

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci

COPY . .
RUN npm run build

# Run Stage
FROM nginx:alpine

COPY --from=builder /app/dist /usr/share/nginx/html
# Add nginx config to handle client-side routing if needed, but default is often okay for simple apps.
# For React Router, we need to redirect 404 to index.html.
RUN echo 'server { \
    listen 80; \
    location / { \
        root /usr/share/nginx/html; \
        index index.html index.htm; \
        try_files $uri $uri/ /index.html; \
    } \
}' > /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
