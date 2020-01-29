# API documentation

## Auth
to be fulfilled


## Routes
The API is currently only composed of one endpoint, reachable with POST, GET, PUT and DELETE.
It allows you to perform CRUD operations on the rules database.

| route       | verb   | description   |
|-------------|--------|---------------|
| /rules      | GET    | retrieve rule |
| /rules      | POST   | create rule   |
| /rules      | PUT    | update rule   |
| /rules      | DELETE | remove rule   |
| /rules/list | GET    | list rules    |


## Swagger
The golang code has swagger tags comments. You can generate a swagger spec file with the [goswagger](https://goswagger.io/) tool:
```shell script
swagger generate spec -o ./swagger.json
```

## Database
The database engine used is SQLite 3.x.
 