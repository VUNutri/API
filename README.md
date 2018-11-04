# Nutri APP API

### Introduction
This api is based on REST architecture. The whole communication process is being made with HTTP json request and responses.
___

### Product

**Product creation**
Url is: `<host>/v1/api/products/create`
. JSON object:
```
{
  "title": string,
  "calories": int,
  "carbs": int,
  "proteins": int
}
```
**Fetching all products**
Url is: `<host>/v1/api/products/getAll`
. JSON object:
```
{
  "id: int,
  "title": string,
  "calories": int,
  "carbs": int,
  "proteins": int
}
```

### Recipe
**Recipe creation**
Url is: `<host>/v1/api/recipes/create`
. JSON object:
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
Url is: `<host>/v1/api/recipes/getAll`
. JSON object:
```
{
  "id": int,
  "title": string,
  "category": int,
  "time": int,
  "image": string,
  "instructions": string,
  "products": [
    {
      "title": string,
      "value": int,
      "calories": int,
      "carbs": int,
      "proteins": int
    }
  ]
}
```
