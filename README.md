# picshare

A photosharing application written in go.

.config file template :

```{
  "port": 3000,
  "env": "",
  "pepper": "",
  "hmac_key": "",
  "database": {
    "host": "localhost",
    "port": 5432,
    "user": "postgres",
    "password": "",
    "name": ""
  },
  "mailgun": {
    "api_key": "",
    "public_api_key": "",
    "domain": ""
  },
  "dropbox": {
    "id": "",
    "secret": "",
    "auth_url": "https://www.dropbox.com/oauth2/authorize",
    "token_url": "https://api.dropboxapi.com/oauth2/token",
    "redirect_url": ""
  }
}
```