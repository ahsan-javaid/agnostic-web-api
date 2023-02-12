## About

Web rest api prototype as generic backend application

## Documentation
This project is generic enough to handle (GET|PUT|POST|DELETE) routes on any entity specified in url. 

## Examples
Base url: http://localhost:5000
```
GET    http://localhost:5000/{entity_name}
PUT    http://localhost:5000/{entity_name}/id 
       Body: { json data here } 
POST   http://localhost:5000/{entity_name}
       Body: { json data here } 
DELETE http://localhost:5000/{entity_name}/id
```
Where entity name could be anything like users, books, .... so on (unstructured json)

## Requirements
- Golang
- MongoDB 

## Project Setup
Enviroment variables loaded from .env file
```
PORT=5000
DB_URI=mongodb://localhost:27017
DB_NAME=agnostic-web-api
```
Run following command in the root of the project
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