# Kea Web
Web server for configuring and monitoring ISC Kea  
## Dev Setup
### Install tdewolff/minify Library  
For minnifying CSS & JS assets  
https://github.com/tdewolff/minify  
`go get -u github.com/tdewolff/minify/v2`
### Install air  
For live reloading  
https://github.com/air-verse/air  
`go install github.com/air-verse/air@latest`
### Setup .env File  
```env
KEA_API_IP=192.68.0.1
KEA_API_USERNAME=test
KEA_API_PASSWORD=xxx
KEA_DB_USER=kea
KEA_DB_PASSWORD=xxx
KEA_DB_NAME=kea
```
### Run
`air`