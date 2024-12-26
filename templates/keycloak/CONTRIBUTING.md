# Contribute to Keycloak version detection

Here are some tips to create new Verdex rules and variables for Keycloak:

- Start a test session:
```bash
go run . -product keycloak -test -test-version 24.0.4 -test-session -verbose
```

- Use official repo to compare versions, for example:
https://github.com/keycloak/keycloak/compare/24.0.3...24.0.4

- Here are nice files to create rules on:
    - Messages translations
       repo file: js/apps/admin-ui/maven-resources/theme/keycloak.v2/admin/messages/messages_en.properties
      exposed at: /resources/master/admin/en

    - OpenID connect step 1
       repo file: services/src/main/resources/org/keycloak/protocol/oidc/endpoints/3p-cookies-step1.html
      exposed at: /realms/master/protocol/openid-connect/3p-cookies/step1.html
