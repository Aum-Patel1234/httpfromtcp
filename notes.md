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

### 4. PARSE each Line such that if let say there are total of 10 bytes in a line and only 8 came the code should not break.

### 5. Parsing the HTTP Request

#### 1. Parse the Request Line

* Read until the first `\r\n`
* Extract:

  * `Method` (e.g., GET, POST)
  * `Request-Target` (e.g., `/`, `/api`)
  * `HTTP-Version` (e.g., `HTTP/1.1` → store `1.1`)
* Validate:

  * Must contain exactly 3 parts
  * Version must start with `"HTTP/"`

---

#### 2. Parse Headers

* Read line-by-line until an empty line (`\r\n`)
* Each header format:

  ```
  Key: Value
  ```
* Rules:

  * Header names are **case-insensitive** → store lowercase
  * Trim spaces around values
  * Duplicate headers → merge as CSV (`value1,value2`)
  * Validate header name (token rules from RFC 9110)

---

#### 3. Detect End of Headers

* Headers end when you encounter:

  ```
  \r\n
  ```

  (i.e., an empty line)

---

#### 4. Parse Body

* Check for `Content-Length` header:

  * If present:

    * Read exactly `Content-Length` bytes as body
  * If NOT present:

    * Assume **no body** (as per RFC 9110, unless using special encodings)

---

#### 5. Important Notes / Edge Cases

* `Content-Length: 0` → valid, empty body
* Do NOT over-read beyond `Content-Length`
* Data may arrive in chunks → support partial reads (streaming)
* Body can contain binary data → store as `[]byte`, not string

