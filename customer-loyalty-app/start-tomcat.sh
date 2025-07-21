#!/bin/bash
set -e

# SECURITY NOTICE: CVE-2024-52316 - ACCEPTED RISK
# Tomcat 8.5 authentication bypass vulnerability is known but skipped
# due to application version constraints. Ensure Jakarta Authentication API is not used.

# Configure Tomcat to use the PORT environment variable
# Default to 8080 if PORT is not set
PORT=${PORT:-8080}

# Configure Tomcat connector port
if [ -f /usr/local/tomcat/conf/server.xml ]; then
    # Create a backup of the original server.xml
    cp /usr/local/tomcat/conf/server.xml /usr/local/tomcat/conf/server.xml.bak

    # Update the HTTP connector port
    sed -i "s/port=\"8080\"/port=\"$PORT\"/g" /usr/local/tomcat/conf/server.xml

    # Security hardening for Tomcat 8
    # Disable server information disclosure
    sed -i 's/server="Apache-Tomcat\/8.5.[0-9]*"/server=""/g' /usr/local/tomcat/conf/server.xml

    echo "Tomcat 8 configured to run on port: $PORT with security hardening"
fi

# Set JVM options for better performance and security (Tomcat 8 compatible)
export CATALINA_OPTS="$CATALINA_OPTS -Djava.security.egd=file:/dev/./urandom"
export CATALINA_OPTS="$CATALINA_OPTS -Dfile.encoding=UTF-8"
export CATALINA_OPTS="$CATALINA_OPTS -Duser.timezone=UTC"
export CATALINA_OPTS="$CATALINA_OPTS -Djava.awt.headless=true"
export CATALINA_OPTS="$CATALINA_OPTS -Djava.net.preferIPv4Stack=true"

# Memory settings for container environment
export CATALINA_OPTS="$CATALINA_OPTS -Xms128m -Xmx512m"

# Security-related JVM options
export CATALINA_OPTS="$CATALINA_OPTS -Djava.security.manager"
export CATALINA_OPTS="$CATALINA_OPTS -Djava.security.policy=/usr/local/tomcat/conf/catalina.policy"

# Start Tomcat 8
echo "Starting Tomcat 8 with security configurations..."
exec catalina.sh run