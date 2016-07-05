# go-serve
Command-Line Static File Server

![Go Report Card](https://goreportcard.com/badge/github.com/Revanth47/go-serve)
![Travis](https://api.travis-ci.org/Revanth47/go-serve.svg?branch=master)
## Usage

#### How to use it? Just run (Default: serves current working directory on localhost:8000)
```bash
$ go-serve
``` 

#### Dont like the port? Then change it
```bash
$ go-serve -p 3000
```


#### Want to serve only a specific sub directory?
```bash
$ go-serve -d assets
```


#### Want to serve files through a different public path? (say localhost:8000/files/)?
```bash
$ go-serve --path /files/
```

#### Other options
```bash
$ go-serve --disable-dir // To disable directory listing
$ go-serve --read 10s    // Maximum request read time
$ go-serve --write 10s   // Maximum response write time 
```

## Installation
```bash
$ go get -u github.com/Revanth47/go-serve
```
# Contribute
Feel free to send Pull Requests to add/improve upon features or to fix bugs.

# LICENSE
[MIT](https://github.com/Revanth47/go-serve/blob/master/LICENSE)