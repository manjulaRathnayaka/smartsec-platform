# Security Risk Assessment - Customer Loyalty App

## Accepted Security Risks

### CVE-2024-52316 - Apache Tomcat Authentication Bypass
- **Severity**: CRITICAL
- **Status**: ACCEPTED RISK
- **Reason**: Application constrained to Tomcat 8.x for compatibility
- **CVSS Score**: 9.8
- **Description**: Authentication bypass when using Jakarta Authentication API

#### Risk Details:
- **Affected Versions**: Tomcat 9.0.0-M1 to 9.0.95, 10.1.0-M1 to 10.1.29, 11.0.0-M1 to 11.0.0
- **Note**: Tomcat 8.x is not directly affected by this specific CVE as it doesn't support Jakarta Authentication API

#### Mitigation Measures:
1. **Network Security**: Deploy behind Choreo's security layer
2. **Access Control**: Implement application-level authentication
3. **Monitoring**: Enable comprehensive logging and monitoring
4. **Input Validation**: Ensure all inputs are properly validated
5. **Regular Updates**: Apply all available Tomcat 8.5.x security patches

#### Compensating Controls:
- Container isolation with non-root user (UID 10001)
- Java Security Manager enabled
- Server information disclosure disabled
- Choreo platform security features
- Regular security scanning

#### Future Actions:
- **Priority**: Plan migration to Tomcat 9.0.96+ when application compatibility allows
- **Timeline**: Evaluate upgrade path within next development cycle
- **Testing**: Verify WAR file compatibility with newer Tomcat versions

---
**Risk Owner**: Development Team  
**Review Date**: January 2025  
**Next Review**: March 2025  
**Approved By**: Security Team (Pending)
