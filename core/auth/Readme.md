# Auth module

OpenId Connect Implementation to login against a configured SSO

## Configuration example

```
auth:
  server: '%%ENV:AUTH_SERVER%%'
  secret: '%%ENV:AUTH_CLIENT_SECRET%%'
  clientid: '%%ENV:AUTH_CLIENT_ID%%'
  myhost: '%%ENV:FLAMINGO_HOSTNAME%%'
  disableOfflineToken: true
```

# Specific scopes

By default, email and profile are added into scopes list (openid scope is
attached always to the list, so it'' not necessary to add it).

```
auth:
  ...
  scopes:
  - email
  - profile
  - address
```

# Specific claims

As openid connect standard it's possible to require claims in auth request.
By default, claims are empty, but it's possible to define a list of voluntaries
claims as a list named "claims".

```
auth:
  ...
  claims:
  - someName
  - someEmail
  - someSalutation
```

# Specific mapping

If it's necessary, fields from id_token and userinfo can be mapped to 
actual user entity. By default, only "sub", "name" and "email" fields
are mapped. 
To map to specific field, use top-level attribute mapping (in example, fields, like 
"someEmail" or "someName" from "id_token", would be mapped to desired fields
"email" and "name" in User entity).

```
auth:
  ...
  mapping.idToken:
    sub: someSub
    name: someName
    email: someEmail
    salutation: someSalutation
    firstName: someFirstName
    lastName: someLastName
    street: someStreet
    zipCode: someZipCode
    city: someCity
    dateOfBirth: someDateOfBirth
    country: someCountry
    customFields:
    - someField1
    - someField2
```
