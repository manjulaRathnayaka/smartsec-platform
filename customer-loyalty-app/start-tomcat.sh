#!/bin/bash
set -e

# SECURITY NOTICE: CVE-2024-52316 - ACCEPTED RISK
# Tomcat 8.5 authentication bypass vulnerability is known but skipped
# due to application version constraints. Ensure Jakarta Authentication API is not used.

#!/bin/bash
set -e

# SECURITY NOTICE: CVE-2024-52316 - ACCEPTED RISK
# Tomcat 8.5 authentication bypass vulnerability is known but skipped
# due to application version constraints. Ensure Jakarta Authentication API is not used.

# Tomcat default port is 8080, which matches Choreo's expected port
echo "Starting Tomcat 8 on default port 8080"

# Set JVM options for better performance and security (Tomcat 8 compatible)
export CATALINA_OPTS="$CATALINA_OPTS -Djava.security.egd=file:/dev/./urandom"
export CATALINA_OPTS="$CATALINA_OPTS -Dfile.encoding=UTF-8"
export CATALINA_OPTS="$CATALINA_OPTS -Duser.timezone=UTC"
export CATALINA_OPTS="$CATALINA_OPTS -Djava.awt.headless=true"
export CATALINA_OPTS="$CATALINA_OPTS -Djava.net.preferIPv4Stack=true"

# Memory settings for container environment
export CATALINA_OPTS="$CATALINA_OPTS -Xms128m -Xmx512m"

# Start Tomcat 8
echo "Starting Tomcat 8 with optimized settings..."
exec catalina.sh run

# Use system properties to configure Tomcat port (avoids server.xml modification)
export CATALINA_OPTS="$CATALINA_OPTS -Dorg.apache.catalina.connector.http.port=$PORT"

echo "Tomcat 8 configured to run on port: $PORT using system properties"

# Set JVM options for better performance and security (Tomcat 8 compatible)
export CATALINA_OPTS="$CATALINA_OPTS -Djava.security.egd=file:/dev/./urandom"
export CATALINA_OPTS="$CATALINA_OPTS -Dfile.encoding=UTF-8"
export CATALINA_OPTS="$CATALINA_OPTS -Duser.timezone=UTC"
export CATALINA_OPTS="$CATALINA_OPTS -Djava.awt.headless=true"
export CATALINA_OPTS="$CATALINA_OPTS -Djava.net.preferIPv4Stack=true"

# Memory settings for container environment
export CATALINA_OPTS="$CATALINA_OPTS -Xms128m -Xmx512m"

# Security-related JVM options - Java Security Manager disabled due to application compatibility issues
# Note: Java Security Manager can cause issues with legacy applications and third-party libraries
# Relying on container-level security and Choreo platform security instead
# export CATALINA_OPTS="$CATALINA_OPTS -Djava.security.manager"
# export CATALINA_OPTS="$CATALINA_OPTS -Djava.security.policy=/usr/local/tomcat/conf/catalina.policy"

# Start Tomcat 8
echo "Starting Tomcat 8 with security configurations..."
exec catalina.sh run