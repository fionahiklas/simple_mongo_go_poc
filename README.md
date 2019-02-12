## Overview

Following this [tutorial](https://www.codementor.io/codehakase/building-a-restful-api-with-golang-a6yivzqdo) to
produce a REST API.

Also using [go-logger](https://github.com/bestmethod/go-logger)

Also referring to this example of a [mongo go application](https://labix.org/mgo)

## Setup

### Docker Mongo

For testing purposes running Mongo in a container is sufficient/relatively simple.

```
docker run -d --name devmongo -e MONGO_INITDB_ROOT_USERNAME=mongoadmin -e MONGO_INITDB_ROOT_PASSWORD=password -v $PWD/mongodb/devmongo:/data/db -p 27017:27017 mongo
```

This assumes a local directory called `mongodb/devmongo` (though docker seems to create missing parts of the path)

You can then connect to MongoDB using the command line tools

```
mongo -u mongoadmin -p password --port 27017 --host 127.0.0.1 admin
```

**NOTE:** You must connect to the `admin` database initially as that is where those credentials are setup

Create the database and a new user from the command line as follows

```
use persondb
db.createUser(
  {
    user: "persondb",
    pwd: "personpass",
    roles: [
       { role: "readWrite", db: "persondb" }
    ],
    passwordDigestor:"server"
  }
)
```

You can also update an existing user with new roles

```
use persondb
db.updateUser("persondb",
  {
    roles: [
       { role: "readWrite", db: "persondb" }
    ]
  }
)

```

You can now exit and logon using the new credentials

```
mongo -u persondb -p personpass --port 27017 --host 127.0.0.1 persondb
```

Adding test data to the DB

```
db.person.insert({
  id: "00001",
  firstname: "Angua",
  lastname: "Von Uberwald",
    address: {
      city: "Ankh Morpork",
      state: "The Shades"
    }
})
```

And another

```
db.person.insert({
  id: "00002",
  firstname: "Agnes",
  lastname: "Nitt",
    address: {
      city: "Mad Stoat",
      state: "Lancre"
    }
})
```

And one more

```
db.person.insert({
  id: "00003",
  firstname: "Polly",
  lastname: "Perks",
    address: {
      city: "Munz",
      state: "Borogravia"
    }
})
```

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

This is copied almost verbatim from the tutorial link above.  The only additions
have been the logging statements.

### Build/install

```
go install tutorial
```

### Run

```
bin/tutorial
```


## Tutorial Mongo

Again based on the tutorial code but attempts to connect to a MongoDB for
data rather than using hard-coded values.

### Build/install

```
go install tutorial_mongo
```

### Run

Pass the connection details for Mongo

```
bin/tutorial_mongo localhost 27017 persondb personpass
```
