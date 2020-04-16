# covid19kalteng
Local covid19 stats API develop with Golang (echo, gorm and goose)

## Usage
First running dependency docker for this project
```sh
$ docker-compose -f deploy/docker-compose.yaml up -d
```
running main project
```sh
$ docker-compose up --build
```

move in to docker container
```sh
$ docker exec -it covid19kalteng /bin/sh
covid19kalteng$
```
run migration inside container
```sh
covid19kalteng$ covid19kalteng migrate up
covid19kalteng$ covid19kalteng seed
```
open your Postman and set your base_url API env to ***localhost:8000***



License
----
MIT


**100% Free royadv/zoed**

[//]: # (These are reference links used in the body of this note and get stripped out when the markdown processor does its job. There is no need to format nicely because it shouldn't be seen. Thanks SO - http://stackoverflow.com/questions/4823468/store-comments-in-markdown-syntax)