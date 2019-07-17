# OAuth module

OpenId connect implementation to login against a configured SSO

## Configuration example

```yaml
oauth:
  server: '%%ENV:AUTH_SERVER%%'
  secret: '%%ENV:AUTH_CLIENT_SECRET%%'
  clientid: '%%ENV:AUTH_CLIENT_ID%%'
  disableOfflineToken: true
```

# Specific scopes

By default, email and profile are added into scopes list (openid scope is
attached always to the list, so it's not necessary to add it).

```yaml
oauth:
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

```yaml
oauth:
  ...
  claims:
  - someName
  - someEmail
  - someSalutation
```

# Specific mapping

If it's necessary, fields from `id_token` and `userinfo` can be mapped to the
actual user entity. By default, only `sub`, `name` and `email` fields are mapped.
 
To map to a specific field, use top-level attribute mapping (in example, fields, like 
`someEmail` or `someName` from `id_token`, would be mapped to desired fields `email` and `name` in the user entity).

```yaml
oauth:
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
    groups: groupfield1;groupfield2
    customFields:
    - someField1
    - someField2
```

As you see above the mapping allows to specify multiple keys in the claim. So `groups: groupfield1;groupfield2` will map the group property of the user object from the claim `groupfield1` and if that is not present it will use `groupfield2`. 

# Use fakes

For testing purposes it's possible to use fakes. In this case, login/logout process
is simulated and it doesn't use any real SSO service. Still, all login and logout
links are valid and clickable, and user data provided from `UserService` is still
present, after "login". Whole process simply redirects to internal pages and handle
session user data.
To specify using of fake services and user data, check configuration below.
Attribute names used for fakeUserData are the same ones used for `id_token` mapping.

```yaml
oauth:
  ...
  useFake: true
  fakeLoginTemplate: "fake/login"
  fakeUserData:
    sub: ID123456
    email: email@domain.com
    name: "Mr. Flamingo"
    ...
```

It's possible to provide fake login page. In this case, the template for fake login page
would be shown. Expected behaviour would be to have a login button that points
to `auth.callback` handler, so it can finish the login process. Fake user data is stored
in session anyway, but with `fakeLoginTemplate` parameter it's allowed to
add a dummy login page in the middle of fake auth process.

```pug
html
  ...
  a(href=url("auth.callback")) Login
```


## Debugging

Start flamingo with the environment variable "OAUTHDEBUG" - to get raw dump of http request and responses to the configured oauth provider logged to stdout.