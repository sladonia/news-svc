# news-sv

News service

### usage

Build and launch service and db with docker-compose
```shell
make up
```

Build docker image using git tag or `dev`
```shell
make build
```

Run tests
```shell
make test
```

### http client
Preconfigured rest-client.http is provided

### database migrations

stored ander `migration` directory in pure sql for simplification

### api endpoints

Create post
```http request
POST /posts

{
  "title": "top news!",
  "content": "covid is over!"
}
```

Get post by id
```http request
GET /posts/{id}
```

Update post
```http request
PUT /posts/{id}

{
  "title": "updated title",
  "content": "updated content"
}
```

Delete post
```http request
DELETE /posts/{id}
```

Find posts
```http request
GET /posts
 ?limit=2
 &offset=2
 &from=2021-11-26T16:03:40.000Z
 &to=2021-11-26T16:03:40.000Z
```
