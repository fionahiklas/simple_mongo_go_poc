## Overview

Following this [tutorial](https://www.codementor.io/codehakase/building-a-restful-api-with-golang-a6yivzqdo) to
produce a REST API.

Also using [go-logger](https://github.com/bestmethod/go-logger)

## Setup

### GoLang

Firstly you need GoLang - this may be available from the MDM self service (if you have one of those laptops) or from the language [site](https://golang.org) the Ubuntu 18.04 VM I'm using has golang version 1.10.4 in the main repo so it can simply be installed with `apt install`.

### GOPATH

Set this up with two entries

* `~/wd/gobase`
* `~/wd/simple_mongo_go_poc`

Set this up using the following command

```
export GOPATH=~/wd/gobase:~/wd/simple_mongo_go_poc
```

## Tutorial

### Build/install

```
go install tutorial
```

### Run

```
bin/tutorial
```
