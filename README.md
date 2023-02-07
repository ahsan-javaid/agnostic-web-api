## About

Web rest api prototype as generic backend application

## Documentation
This project is generic enough to handle (GET|PUT|POST|DELETE) routes on any entity specified in url. 

## Examples
```
GET /<entity name>
PUT /<entity name>/id
POST /<entity name>
DELETE /<entity name>/id
```
Where entity name could be anything like users, books, .... so on (unstructured json)

## Requirements
- Golang
- MongoDB 

## Project Setup
```
go mod download
```

## Requirements
```
Mongodb running at localhost:27017
```

## Authentication 
Coming soon...

### Architecture
Here is a light architecture diagram describing the relationship between packages in agnostic-web-api:

```
    +--main--+       +---api--+   +------db-----+   
    |        |-url-> +  GET   |-->|   Perform   |   
    |  http  |       +  PUT   |   |  operation  |    
    | server |       |  POST  |   |  based on   | 
    |        |       | DELETE |   |entity in url|
    +---+----+       +--------+   +-------------+
     |              http listner              |
     |  Routing & db read/write based on url  |
     +----------------------------------------+
```