podca-api [![Build Status](https://travis-ci.org/te-th/podca-api.svg?branch=master)](https://travis-ci.org/te-th/podca-api)
===

A minimal repository serving Podcast Feeds as JSON Objects. Build on Google Cloud Platform.

## Getting started

Go-Gettable via
```bash
go get github.com/te-th/podca-api
```
Run by
```make
make serve
```
or
```
goapp serve ./app
```
## Usecase search

```
http://localhost:8080/podcasts/search?term=WDR&limit=20
```

## Usecase feeds:

```
http://localhost:8080/feeds
``` 

## License

Licensed under Apache License, Version 2.0.
