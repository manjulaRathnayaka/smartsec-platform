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
    rm -rf /var/lib/apt/lists/*

# Remove default webapps for security and remove unnecessary files
RUN rm -rf /usr/local/tomcat/webapps/* && \
    rm -rf /usr/local/tomcat/webapps.dist && \
    rm -rf /usr/local/tomcat/work/* && \
    rm -rf /usr/local/tomcat/temp/*

# Copy the WAR file
COPY raincoat-loyalty-manager.war /usr/local/tomcat/webapps/ROOT.war

# Create non-root user for Choreo (UID between 10000-20000)
RUN groupadd -r -g 10001 choreo && \
    useradd -r -u 10001 -g choreo choreo && \
    chown -R choreo:choreo /usr/local/tomcat

# Switch to non-root user
USER 10001

# Use environment variable for port (Choreo injects PORT env var)
ENV PORT=8080
EXPOSE ${PORT}

# Start Tomcat with standard /usr/local/tomcat directories
CMD ["sh", "-c", "mkdir -p /usr/local/tomcat/work /usr/local/tomcat/temp /usr/local/tomcat/logs && java -Djava.security.egd=file:/dev/./urandom -Dfile.encoding=UTF-8 -Duser.timezone=UTC -Djava.awt.headless=true -Djava.net.preferIPv4Stack=true -Xms128m -Xmx512m -Djava.util.logging.config.file=/usr/local/tomcat/conf/logging.properties -Djava.util.logging.manager=org.apache.juli.ClassLoaderLogManager -Djdk.tls.ephemeralDHKeySize=2048 -Djava.protocol.handler.pkgs=org.apache.catalina.webresources -Dorg.apache.catalina.security.SecurityListener.UMASK=0027 -Dignore.endorsed.dirs= -Dcatalina.base=/usr/local/tomcat -Dcatalina.home=/usr/local/tomcat -Djava.io.tmpdir=/usr/local/tomcat/temp -cp /usr/local/tomcat/bin/bootstrap.jar:/usr/local/tomcat/bin/tomcat-juli.jar org.apache.catalina.startup.Bootstrap start"]
