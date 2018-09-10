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
