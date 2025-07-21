#!/bin/bash
set -e

# Configure Tomcat to use the PORT environment variable
# Default to 8080 if PORT is not set
PORT=${PORT:-8080}

# Configure Tomcat connector port
if [ -f /usr/local/tomcat/conf/server.xml ]; then
    # Create a backup of the original server.xml
    cp /usr/local/tomcat/conf/server.xml /usr/local/tomcat/conf/server.xml.bak
    
    # Update the HTTP connector port
    sed -i "s/port=\"8080\"/port=\"$PORT\"/g" /usr/local/tomcat/conf/server.xml
    
    echo "Tomcat configured to run on port: $PORT"
fi

# Set JVM options for better performance and security
export CATALINA_OPTS="$CATALINA_OPTS -Djava.security.egd=file:/dev/./urandom"
export CATALINA_OPTS="$CATALINA_OPTS -Dfile.encoding=UTF-8"
export CATALINA_OPTS="$CATALINA_OPTS -Duser.timezone=UTC"

# Start Tomcat
echo "Starting Tomcat..."
exec catalina.sh run
