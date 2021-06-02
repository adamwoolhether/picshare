# picshare

A photosharing application written in go.

### Why did I make this?
To serve as proof of concept demonstrating my ability to create web applications with the [Go programming language](https://www.golang.org/).   
I may keep adding new features as I learn new tricks  
Go ahead and click [Sign Up](https://www.adamwoolhether.com/signup) on the top right and start uploading pictures!

### What is this?
🧐 This is back end development project.  
📸 Picture sharing web app written 💯 in the [Go programming language](https://www.golang.org/).  
📤 [Dropbox](https://www.dropbox.com/home) integration via persistent OAuth tokens.  
🦍 [Gorilla/mux](https://github.com/gorilla/mux) used for HTTP routing.  
📋 Underlying database: [PostgreSQL](https://postgresql.com/), implemented with [Gorm](https://gorm.io/).  
🔎 Data validation for DB consistency.  
👨 Served with [Caddy](https://caddyserver.com/) on a VM running CentOS 8 (rip) on [DigitalOcean](https://www.digitalocean.com/).  
📃 Frontend created with templatized .gohtml files and [Bootstrap](https://getbootstrap.com/).  
🍪 Session cookies with protection against [XSS](https://owasp.org/www-community/attacks/xss/) and [CSRF](https://owasp.org/www-community/attacks/csrf).  
🔐 Passwords are encrypted and hashed.  
❓ Reset password functionality.  
📧 Sign Up email and password reset implemented with [Mailgun-Go](https://github.com/mailgun/mailgun-go) (You may have to check your junkmail folder, as my domain is not widely recognized).

### How did you make this site?
I made this with the [GoLand IDE](https://www.jetbrains.com/go/) developed by Jet Brains.
Many many thanks to [Jon Calhoun](https://www.calhoun.io/)!

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

## TODO
- Fix Album not deleting images on disk
- Top nav-bar shouldn't ask to connect to dropbox if already connected
- Fix dropbox javascipt id to be dynamically set to config setting
- Fix logging level flag is correctly implemented
- Implement option to make galleries private

- Redirect to Gallery view after image upload (?)
- Reset password request should redirect to another place.