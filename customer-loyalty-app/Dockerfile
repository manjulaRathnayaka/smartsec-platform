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

# Create writable temp directory (will be backed by emptyDir at runtime)
RUN mkdir -p /usr/local/tomcat/temp && \
    chown -R tomcatuser:tomcatuser /usr/local/tomcat/temp

# Fix ownership of writable directories for non-root execution
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

# Declare writable volume (emptyDir) for /temp
VOLUME ["/usr/local/tomcat/temp"]

# Expose default Tomcat port
EXPOSE 8080

# Use full path to catalina.sh for non-root user
CMD ["/usr/local/tomcat/bin/catalina.sh", "run"]
