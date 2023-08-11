FROM quay.io/keycloak/keycloak:21.1.2

# Copy the realm data to the appropriate location inside the container
COPY /data /opt/keycloak/data/import

# Start Keycloak with uri_redirect option for logout and import realm settings
CMD ["start-dev", "--spi-login-protocol-openid-connect-legacy-logout-redirect-uri=true", "--import-realm"]
