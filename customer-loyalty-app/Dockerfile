# Use Tomcat 8 with latest available security patches
# Using 8.5-jdk11 for compatibility while maintaining security
# NOTE: Tomcat 8.5 has known vulnerabilities that are fixed in newer versions
# Consider upgrading to Tomcat 9+ when application compatibility allows
FROM tomcat:8.5-jdk11-openjdk-slim

# Update base packages to patch security vulnerabilities
RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y --no-install-recommends curl && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* && \
    groupadd -r -g 10001 choreo && \
    useradd -r -u 10001 -g choreo choreo

# Remove default webapps for security and remove unnecessary files
RUN rm -rf /usr/local/tomcat/webapps/* && \
    rm -rf /usr/local/tomcat/webapps.dist && \
    rm -rf /usr/local/tomcat/work/* && \
    rm -rf /usr/local/tomcat/temp/*

# Copy the WAR file and startup script
COPY raincoat-loyalty-manager.war /usr/local/tomcat/webapps/ROOT.war
COPY start-tomcat.sh /usr/local/tomcat/bin/start-tomcat.sh

# Set proper permissions
RUN chmod +x /usr/local/tomcat/bin/start-tomcat.sh && \
    chown -R choreo:choreo /usr/local/tomcat/webapps /usr/local/tomcat/work /usr/local/tomcat/temp /usr/local/tomcat/logs /usr/local/tomcat/conf

# Switch to non-root user (Choreo requires UID between 10000-20000)
USER 10001

# Use environment variable for port (Choreo injects PORT env var)
ENV PORT=8080
EXPOSE ${PORT}

# Health check for better container management
HEALTHCHECK --interval=30s --timeout=3s --start-period=30s --retries=3 \
    CMD curl -f http://localhost:${PORT}/ || exit 1

# Start Tomcat with our custom script
CMD ["/usr/local/tomcat/bin/start-tomcat.sh"]
