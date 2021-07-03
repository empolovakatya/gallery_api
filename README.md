# gallery_api

## Run the app

    go run ./cmd/web
    
    port: 8000
    
# REST API

The REST API to the example app is described below.

## Get list of Photos

### Request

`GET /photos`

### Response

    HTTP/1.1 200 OK
    Date: Sat, 03 Jul 2021 14:33:54 GMT
    Content-Type: application/json
    Content-Length: 275
    
    []

## Create a new Photo

### Request

`POST /photos`

### Response

    HTTP/1.1 200 OK 
    Date: Sat, 03 Jul 2021 14:33:54 GMT
    Content-Type: text/plain; charset=utf-8
    Content-Length: 36


## Get a specific Photo

### Request

`GET /photos/id`


### Response

    HTTP/1.1 200 OK
    Date: Sat, 03 Jul 2021 14:37:52 GMT
    Content-Type: application/json
    Content-Length: 136
    
    {"id":1,"image":"Foo","preview":"Bar"}
    
## Delete a Photo

### Request

`DELETE /photo/id`

### Response

    HTTP/1.1 200 OK
    Date: Sat, 03 Jul 2021 14:40:57 GMT
    Content-Length: 17
