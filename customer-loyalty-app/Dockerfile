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

# Copy the WAR file to a location that won't be overwritten by tmpfs
COPY raincoat-loyalty-manager.war /opt/ROOT.war

# Set proper permissions for the WAR file
RUN chown choreo:choreo /opt/ROOT.war

# Switch to non-root user (Choreo requires UID between 10000-20000)
USER 10001

# Use environment variable for port (Choreo injects PORT env var)
ENV PORT=8080
EXPOSE ${PORT}

# Start Tomcat with Choreo-compatible setup using /tmp for writable directories
CMD ["sh", "-c", "mkdir -p /tmp/tomcat/webapps /tmp/tomcat/work /tmp/tomcat/temp /tmp/tomcat/logs /tmp/tomcat/conf /tmp/tomcat/bin && cp -r /usr/local/tomcat/conf/* /tmp/tomcat/conf/ 2>/dev/null || true && cp /usr/local/tomcat/bin/*.jar /tmp/tomcat/bin/ 2>/dev/null || true && cp /opt/ROOT.war /tmp/tomcat/webapps/ROOT.war && java -Djava.security.egd=file:/dev/./urandom -Dfile.encoding=UTF-8 -Duser.timezone=UTC -Djava.awt.headless=true -Djava.net.preferIPv4Stack=true -Xms128m -Xmx512m -Djava.util.logging.config.file=/tmp/tomcat/conf/logging.properties -Djava.util.logging.manager=org.apache.juli.ClassLoaderLogManager -Djdk.tls.ephemeralDHKeySize=2048 -Djava.protocol.handler.pkgs=org.apache.catalina.webresources -Dorg.apache.catalina.security.SecurityListener.UMASK=0027 -Dignore.endorsed.dirs= -Dcatalina.base=/tmp/tomcat -Dcatalina.home=/usr/local/tomcat -Djava.io.tmpdir=/tmp/tomcat/temp -cp /tmp/tomcat/bin/bootstrap.jar:/tmp/tomcat/bin/tomcat-juli.jar org.apache.catalina.startup.Bootstrap start"]
