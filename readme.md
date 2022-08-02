# Phoenix [![Go Reference](https://pkg.go.dev/badge/github.com/marius004/phoenix.svg)](https://pkg.go.dev/github.com/marius004/phoenix) 

Phoenix is a work in progress platform for learning data structures and algoritms. 

## Description

Being a platform for solving algorithms and data structures problems, Phoenix features 
a complex code execution environment that allows users to run their source code against 
custom made tests. After each execution, the user receives a score from 0 to 100 based 
on the correctitude and runtime performance of the submission. 

## Getting Started

### Dependencies
* PostgreSQL
* Linux 
* Isolate
* GCC
* Go 

### 1. Setting up isolate 
```
sudo apt-get update
sudo apt-get install libcap-dev asciidoc-base
cd ~/Downloads
git clone https://github.com/marius004/isolate
cd isolate
make
sudo make install
```

### 2. Setting up GCC
```
sudo apt-get update
sudo apt install build-essential
```

### 3. Clone the github repository

### 4. Configure up the enviroment variables 
1. create the .env file in the root of the project
2. set the following environment variables in your .env file
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_NAME=project
DB_PASSWORD=?

JWT_SECRET=?
COOKIE_LIFETIME=7

SERVER_HOST=localhost
SERVER_PORT=8080
```

### 5. Running the server
```
go mod download
go run .
```