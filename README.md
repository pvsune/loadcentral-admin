# LoadCentral Admin
A simple interface to let you topup mobile load to multiple phone numbers concurrently using Goroutines.

## Quickstart
1. Build the image:
```
$ docker build -t pvsune/loadcentral-admin .
```

2. Fill up required environment variables in `.env`.
  - `APP_AUTH_USERNAME` `APP_AUTH_PASSWORD` - User credentials is not necessary to save to DB yet.
  - `APP_AUTH_KEY` - JWT key used for signing the token.
  - `APP_AUTH_COOKIEDOMAIN` - Set to localhost for development.
  - `APP_LC_USERNAME` `APP_LC_PASSWORD` - LoadCentral credentials.
  - `APP_AUTH_SECURECOOKIE` - Avoid XSS; Set to false for development.
  
3. Start the app:
```
$ docker run --rm --env-file .env -p 8000:8000 pvsune/loadcentral-admin
```

4. Alternatively, you can run the app directly and export enironment variables.
```
$ go get .
$ go run main.go 
```

## How it works
### Authentication
There's no need to use DB yet; the user credentials are set by environment variables. When the user logins, the app generates a JWT token and save it to cookie. Cookie settings are set to avoid CSRF, XSS. The token will expire for one hour and the user needs to login again.

### Load topup
Each `(PhoneNumber, Pcode)` is posted to LoadCentral API using its own Goroutine. The response is sent back to main program using channels and finally sent to templates for display (duhh!).

## Notes
1. Maybe good idea to setup proper CSRF token and not rely on `Set-Cookie: <cookie-name>=<cookie-value>; SameSite=Strict` response header.
2. Send email for every topup request.
3. Use sentry free tier to send app logs.
