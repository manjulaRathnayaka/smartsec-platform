FROM tomcat:8.5-jdk11-openjdk-slim

# Set fixed UID
ENV TOMCAT_UID=10001

# Create non-root user with no home directory
RUN useradd -u ${TOMCAT_UID} --no-create-home --shell /sbin/nologin tomcatuser

# Clean default webapps
RUN rm -rf /usr/local/tomcat/webapps/*

# Copy WAR file into webapps
COPY raincoat-loyalty-manager.war /usr/local/tomcat/webapps/ROOT.war

# Create writable temp directory for Tomcat (for emptyDir mount)
RUN mkdir -p /usr/local/tomcat/temp && \
    chown -R tomcatuser:tomcatuser /usr/local/tomcat/temp

# Set ownership for logs/work/webapps dirs that need to be writable
RUN chown -R tomcatuser:tomcatuser /usr/local/tomcat/logs /usr/local/tomcat/work /usr/local/tomcat/webapps

# Make specific directories writable for tomcat operations
RUN chmod u+w /usr/local/tomcat/logs /usr/local/tomcat/work /usr/local/tomcat/webapps

# Switch to non-root user
USER ${TOMCAT_UID}

# Set safe working directory
WORKDIR /usr/local/tomcat

# Declare temp as a writable volume (to be used with emptyDir)
VOLUME ["/usr/local/tomcat/temp"]

# Expose default Tomcat port
EXPOSE 8080

# Start Tomcat
CMD ["catalina.sh", "run"]
