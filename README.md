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
  }
}```