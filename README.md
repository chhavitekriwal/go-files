# go-files
A simple file management system in Golang

## Run Locally

Clone the project

```bash
  git clone https://github.com/chhavitekriwal/go-files.git
  cd go-files
```
Build and start the server
```
  go build .
```

## API Reference

Base URL for local: http://127.0.0.1:8080

### Register

```shell
  POST /register
```

<details>
<summary>Request</summary>
<pre>
{
    "username": "userabc",
    "password": "abc1234"
}
</pre>
</details>

<details>
<summary>Response</summary>
<pre>
{
    "Message": "User registered and logged in",
    "Token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InVzZXJhYmMiLCJleHAiOjE3MDI5MDc0OTYsImlhdCI6MTcwMjkwNTY5Nn0.CMGQDyfoe5Oe1nICoTFpdJoMHvblV6pHdLKoTM7FM6Y"
}
</pre>
</details>

### Login 

```shell
  POST /login
```

<details>
<summary>Request</summary>
<pre>
{
    "username": "userabc",
    "password": "abc1234"
}
</pre>
</details>

<details>
<summary>Response</summary>
<pre>
{
    "Message": "Logged in",
    "Token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InVzZXJhYmMiLCJleHAiOjE3MDI5MDc1NjMsImlhdCI6MTcwMjkwNTc2M30.w2ta9UHGor8haAOcDRSzAzjr7O7NUnVE_lfFqSDtVto"
}
</pre>
</details>

### Upload

```shell
  POST /upload
```

<details>
<summary>Request</summary>
<pre>
![Alt text](assets/image-1.png)
</pre>
</details>

<details>
<summary>Response</summary>
<pre>
{
    "message": "Successfully uploaded",
    "filename": "xRsuEua.jpg",
    "size_in_bytes": 759085,
    "type": "image/jpeg"
}
</pre>
</details>

### Delete

```shell
  DELETE /delete?filename=F3000TEK.JPG
```

<details>
<summary>Response</summary>
<pre>
{
    "message": "Successfully deleted"
}
</pre>
</details>

```shell
  GET /files
```
### List all files

<details>
<summary>Response</summary>
<pre>
[
    {
        "filename": "Why_Phishing_works.pdf",
        "size_in_bytes": 1419266,
        "type": "application/pdf",
        "modified": "2023-12-18 04:31:43"
    },
    {
        "filename": "xRsuEua.jpg",
        "size_in_bytes": 759085,
        "type": "image/jpeg",
        "modified": "2023-12-18 19:04:08"
    }
]
</pre>
</details>

### Get file transactions

```shell
  GET /transactions
```

<details>
<summary>Response</summary>
<pre>
[
    {
        "id": 23,
        "timestamp": "2023-12-18T19:01:18.900674321+05:30",
        "filename": "F0003TEK.JPG",
        "transaction_type": "DOWNLOAD",
        "username": "newuser"
    },
    {
        "id": 22,
        "timestamp": "2023-12-18T18:55:28.166416481+05:30",
        "filename": "F0003TEK.JPG",
        "transaction_type": "UPLOAD",
        "username": "newuser"
    },
    {
        "id": 21,
        "timestamp": "2023-12-18T18:34:08.089294576+05:30",
        "filename": "Social_Phishing.pdf",
        "transaction_type": "DELETE",
        "username": "newuser"
    }
]
</pre>