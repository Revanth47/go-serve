# go-serve

Command-Line Static File Server

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


#### Not a fan of directory listing? Just run
```bash
$ go-serve --disable-dir
```


#### Want to serve files through a different public path? (say localhost:8000/files/)?
```bash
$ go-serve --path /files/
```


## Installation
```bash
$ go get -u github.com/Revanth47/go-serve
```
# Contribute
Feel free to send Pull Requests to add/improve upon features or for fix bugs.

# LICENSE
[MIT](https://github.com/Revanth47/go-serve/blob/master/LICENSE)