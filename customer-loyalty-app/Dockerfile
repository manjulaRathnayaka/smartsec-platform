FROM tomcat:8.5-jdk11-openjdk-slim

# Create non-root user with fixed UID and no home directory
RUN useradd -u 10001 --no-create-home --shell /sbin/nologin tomcatuser

# Clean default webapps
RUN rm -rf /usr/local/tomcat/webapps/*

# Copy WAR and extract it to ROOT/
COPY raincoat-loyalty-manager.war /tmp/ROOT.war
RUN mkdir -p /usr/local/tomcat/webapps/ROOT && \
    cd /usr/local/tomcat/webapps/ROOT && \
    jar -xf /tmp/ROOT.war && \
    rm /tmp/ROOT.war

# Copy custom Tomcat configuration files
COPY conf/tomcat-users.xml /usr/local/tomcat/conf/tomcat-users.xml
COPY conf/context.xml /usr/local/tomcat/conf/context.xml

# Ensure config files are readable by non-root user
RUN chmod 644 /usr/local/tomcat/conf/tomcat-users.xml /usr/local/tomcat/conf/context.xml

# Create writable temp directory (backed by emptyDir at runtime)
RUN mkdir -p /usr/local/tomcat/temp && \
    chown -R tomcatuser:tomcatuser /usr/local/tomcat/temp

# Fix ownership of writable directories
RUN chown -R tomcatuser:tomcatuser /usr/local/tomcat/logs \
    /usr/local/tomcat/work \
    /usr/local/tomcat/webapps

# Ensure writable permissions for necessary Tomcat dirs
RUN chmod u+w /usr/local/tomcat/logs \
    /usr/local/tomcat/work \
    /usr/local/tomcat/webapps

# Switch to non-root user
USER 10001

# Set safe working directory
WORKDIR /usr/local/tomcat

# Declare writable temp volume (used with emptyDir)
VOLUME ["/usr/local/tomcat/temp"]

# Expose default Tomcat port
EXPOSE 8080

# Start Tomcat
CMD ["catalina.sh", "run"]
