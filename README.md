# Nutri APP API

### Introduction
This api is based on REST architecture. The whole communication process is being made with HTTP JSON request and responses. To try this api locally, all you need to is to install [docker-compose](https://docs.docker.com/compose/), create .env file and run:  
```docker-compose up -d --build```
 ___

### Product

**Product creation**
POST request url: `<host>/v1/api/products/create`
JSON Request:
```
{
  "title": string,
  "calories": int,
  "carbs": int,
  "proteins": int
}
```
**Fetching all products**
GET request url: `<host>/v1/api/products/getAll`
JSON Request:
```
[
  {
    "id: int,
    "title": string,
    "calories": int,
    "carbs": int,
    "proteins": int
  }
]
```
**Fetching product by id**
GET request url: `<host>/v1/api/products/getById/{id}`
JSON Response:
```
{
  "id": int,
  "title": string,
  "calories": int,
  "carbs": int,
  "proteins": int
}
```

### Recipe
**Recipe creation**
POST request url: `<host>/v1/api/recipes/create`
JSON Response:
```
{
  "title": string,
  "category": int,
  "time": int,
  "image": string,
  "instructions": string,
  "products": [
    {
      "id": int,
      "value": int
    }
  ]
}
```
**Fetching all recipes**
GET request url: `<host>/v1/api/recipes/getAll`
JSON Response:
```
[
  {
    "id": int,
    "title": string,
    "category": int,
    "time": int,
    "image": string,
    "instructions": string,
    "calories": int,
    "carbs": int,
    "proteins": int
    "products": [
      {
        "id": int,
        "title": string,
        "value": int,
        "calories": int,
        "carbs": int,
        "proteins": int
      }
    ]
  }
]
```
**Fetching recipe by id**
GET request url: `<host>/v1/api/recipes/getById/{id}`
JSON Response:
```
{
  "id": int,
  "title": string,
  "category": int,
  "time": int,
  "image": string,
  "instructions": string,
  "calories": int,
  "carbs": int,
  "proteins": int
  "products": [
    {
      "id": int,
      "title": string,
      "value": int,
      "calories": int,
      "carbs": int,
      "proteins": int
    }
  ]
}
```
**Checking if title is valid**
GET request url: `<host>/v1/api/recipes/checkTitle/<title>`
JSON Response:
```
Title is valid/not valid
```

### Menu
**Fetching menu**
POST request url: `<host>/v1/api/menu/getMenu`
JSON Request:
```
{
  "days": int,
  "meals": int,
  "time": int,
  "calories": int,
  "blockedIngredients": int []
}
```
JSON Response:
```
[
  {
    "dayCount": int,
    "meals": [
      "id": int,
      "title": string,
      "category": int,
      "time": int,
      "image": string,
      "instructions": string,
      "calories": int,
      "carbs": int,
      "proteins": int
      "products": [
        {
          "id": int,
          "title": string,
          "value": int,
          "calories": int,
          "carbs": int,
          "proteins": int
        }
      ]
    ]
  }
]
```
**Fetching one day menu**
POST request url: `<host>/v1/api/menu/getDailyMenu`
JSON Request:
```
{
  "dayCount": int,
  "meals": int,
  "calories": int,
  "time": int,
  "blockedIngredients": int []
}
```
JSON Response:
```
{
  "dayCount": int,
  "meals": [
    "id": int,
    "title": string,
    "category": int,
    "time": int,
    "image": string,
    "instructions": string,
    "calories": int,
    "carbs": int,
    "proteins": int
    "products": [
      {
        "id": int,
        "title": string,
        "value": int,
        "calories": int,
        "carbs": int,
        "proteins": int
      }
    ]
  ]
}
```
**Fetching single day menu**
POST request url: `<host>/v1/api/menu/getDayOneMenu`
JSON Request:
```
{
  "category": int,
  "calories": int,
  "time": int,
  "blockedIngredients": int []
}
```
JSON Response:
```
{
  "id": int,
  "title": string,
  "category": int,
  "time": int,
  "image": string,
  "instructions": string,
  "calories": int,
  "carbs": int,
  "proteins": int
  "products": [
    {
      "id": int,
      "title": string,
      "value": int,
      "calories": int,
      "carbs": int,
      "proteins": int
    }
  ]
}
```
### Auth
**Admin creation**
POST request url: `<host>/v1/api/auth/create`
JSON Response:
```
{
  "username": string,
  "pass": string //hash instead of cleartext password
}
```
**Admin login**
POST request url: `<host>/v1/api/auth/login`
JSON Response:
```
{
  "username": string,
  "pass": string //hash instead of cleartext password
}
```
**Cookies**
TTL of cookies is 15 minutes
