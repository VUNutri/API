# Nutri APP API

### Introduction
This api is based on REST architecture. The whole communication process is being made with HTTP JSON request and responses. 
 ___

### Product

**Product creation**
POST request url: `<host>/v1/api/products/create`
JSON Response: 
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
JSON Response:
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
**Fetching menu**
GET request url: `<host>/v1/api/recipes/getMenu/{daysCount}/{mealsCount}/{caloriesCount}`
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

