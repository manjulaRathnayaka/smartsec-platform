#!/bin/sh

# Replace environment variables in nginx config
envsubst '${BFF_SERVICE_URL}' < /etc/nginx/nginx.conf > /tmp/nginx.conf
mv /tmp/nginx.conf /etc/nginx/nginx.conf

# Replace environment variables in built app
find /usr/share/nginx/html -name "*.js" -exec sed -i "s|REACT_APP_API_BASE_URL_PLACEHOLDER|${REACT_APP_API_BASE_URL}|g" {} \;
find /usr/share/nginx/html -name "*.js" -exec sed -i "s|REACT_APP_TELEMETRY_URL_PLACEHOLDER|${REACT_APP_TELEMETRY_URL}|g" {} \;
find /usr/share/nginx/html -name "*.js" -exec sed -i "s|REACT_APP_MCP_URL_PLACEHOLDER|${REACT_APP_MCP_URL}|g" {} \;

# Start nginx
exec "$@"
