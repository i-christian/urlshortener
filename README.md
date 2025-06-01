# URL shortening service

## Description
A URL Shortener is a service that takes a long URL and generates a shorter, unique alias that redirects users to the original URL. This alias is often a fixed-length string of characters. The system should be able to handle millions of URLs, allowing users to create, store, and retrieve shortened URLs efficiently. Each shortened URL needs to be unique and persistent. Additionally, the service should be able to handle high traffic, with shortened URLs redirecting to the original links in near real-time. In some cases, the service may include analytics to track link usage, such as click counts and user locations.

## Key technologies used in the project include:
- **Backend**: Golang
- **Database**: Redis

## Usage:
- Send a POST request with a long `url` using Curl
```
   curl -vL -XPOST http://localhost:8080/ \
         -H "Content-Type: application/x-www-form-urlencoded" \
         -d "long_url=https://www.google.com"
```
- This should return a JSON response like this
```
  {
    "key":"8dQsYab",
    "long_url":"https://www.google.com",
    "short_url":"http://localhost/aT3s1Nq"
  }
```
- If a long url is already in the database, an error in returned.

- Redirect a client request for the shortened URL
```
   curl http://localhost:8080/aT3s1Nq -i


  HTTP/1.1 302 Found
  Content-Type: text/html; charset=utf-8
  Location: https://www.google.com
  Date: Sun, 01 Jun 2025 20:39:03 GMT
  Content-Length: 45

  <a href="https://www.google.com">Found</a>.
```

- When we use an invalid short URL we will get the following:
```
  curl http://localhost:8080/aT3s1N -i

  HTTP/1.1 404 Not Found
  Content-Type: text/plain; charset=utf-8
  Content-Length: 14

  URL not found
```

- If we use `-L` flag with curl we will be redirected to the long url.
```
    curl http://localhost:8080/aT3s1Nq -L
```

