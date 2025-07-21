# Use official Tomcat runtime as the base image
FROM tomcat:9.0-jdk11-openjdk

# Install curl for health checks and create non-root user with specific UID for Choreo
RUN apt-get update && \
    apt-get install -y curl && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* && \
    groupadd -r -g 10001 choreo && \
    useradd -r -u 10001 -g choreo choreo

# Remove default webapps for security
RUN rm -rf /usr/local/tomcat/webapps/*

# Copy the WAR file and startup script
COPY raincoat-loyalty-manager.war /usr/local/tomcat/webapps/ROOT.war
COPY start-tomcat.sh /usr/local/bin/start-tomcat.sh

# Set proper permissions
RUN chown -R choreo:choreo /usr/local/tomcat/ && \
    chmod +x /usr/local/bin/start-tomcat.sh && \
    chmod -R 755 /usr/local/tomcat/

# Switch to non-root user (Choreo requires UID between 10000-20000)
USER 10001

# Use environment variable for port (Choreo injects PORT env var)
ENV PORT=8080
EXPOSE ${PORT}

# Health check for better container management
HEALTHCHECK --interval=30s --timeout=3s --start-period=30s --retries=3 \
    CMD curl -f http://localhost:${PORT}/ || exit 1

# Start Tomcat with our custom script
CMD ["/usr/local/bin/start-tomcat.sh"]
