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

# Use fakes

For testing purposes it's possible to use fakes. In this case, login/logout process
is simulated and it doesn't use any real SSO service. Still, all login and logout
links are valid and clickable, and user data provided from UserService is still
present, after "login". Whole process simply redirects to internal pages and handle
session user data.
To specify using of fake services and user data, check configuration bellow.
Attribute names used for fakeUserData are the same ones used for id_token mapping.

```
auth:
  ...
  useFake: true
  fakeLoginTemplate: "fake/login"
  fakeUserData:
    sub: ID123456
    email: email@domain.com
    name: "Mr. Flamingo"
    ...
```

It's possible to provide fake login page. In this case, template for fake login page
would be shown. Expected behaviour would be to have login button that points
to "auth.callback" handler, so it can finish login process. Fake user data is stored
in session anyway, but with "fakeLoginTemplate" parameter it's allowed to
add dummy login page in the middle of fake auth process.

```
html
  ...
  a(href=url("auth.callback")) Login
```
