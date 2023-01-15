## Clone repo

`$ git clonehttps://github.com/jesseinvent/go-gin-mongo-jwt`

## Run containers

`$ docker-compose up -d --build`

## Endpoints

- SIGNUP

`POST /api/v1/auth/signup`

```
{
    "first_name": "",
    "last_name": "",
    "password": "",
    "email": "",
    "phone": "",
    "user_type": ""
}
```

- LOGIN

`POST /api/v1/auth/login`

```
{
    "email": "",
    "password": ""
}
```

- GET ALL USERS (admin)

`Headers:`

`Authorization: Bearer $token`

`GET /api/v1/users`

- GET USER (user & admin)

`Headers:`

`Authorization: Bearer $token`

`GET api/v1/users/:user_id`