# Flow

### 1. Read the bytes from source(file, network connection,...).
### 2. Accept the TCP connection.
### 3. HTTP body
  - GET
  ```txt
    read: GET /coffee HTTP/1.1
    read: Host: localhost:42069
    read: User-Agent: curl/8.9.1
    read: Accept: */*
    read: 
  ```
  - POST
  ```txt
    read: POST /coffee HTTP/1.1
    read: Host: localhost:42069
    read: User-Agent: curl/8.9.1
    read: Accept: */*
    read: Content-Type: application/json
    read: Content-Length: 14
    read: 
    read: {"name":"aum"}
  ```

  ### NOTE: 
      1. each line is seperated by \r\n
      2. there is one last \r\n before the reqBody starts
      3. Generally there is Content-Length given to know the size of the reqBody
      4. There is none in GET cause there is not
