# 🔗 URL Shortener Backend Assignment (1–2 hours)

## 📌 Objective
Build a minimal RESTful backend service in Go that shortens long URLs and supports basic redirection. This exercise tests clean API design, data modeling, and minimal persistence logic.


## 📦 Requirements

### 1. Core Functionality

Implement a service that allows users to:
- Create a short URL from a long one
- Redirect users who visit the short URL to the original one

**Note: You are free to use gRPC or HTTP, your choice.  The provided template offers flexibility to use either.**

### 2. RESTful Endpoints

- `POST /api/shortlinks`
  Creates a new short link
  **Request:**
  ```json
  {
    "original_url": "https://example.com"
  }
````

**Response:**

```json
{
  "short_url": "http://localhost:8080/r/abc123"
}
```

* `GET /r/{short_code}`
  Redirects to the original URL
  **Response:** 302 Found → `Location: original_url`


## 💾 Tech Stack

* **Language**: Go
* **Framework**: Any (Chi, Gin, Echo, or `net/http`)
* **Storage**: Use an in-memory map or SQLite/Postgres (your choice)


## ✅ Constraints

* Short code (`short_code`) should be 6–8 alphanumeric characters
* Must validate that `original_url` is a valid URL
* No authentication, no UI
* Assume `http://localhost:8080` as base URL


## 🧪 Bonus (Optional, if you have time)

* [ ] Handle duplicate original URLs by returning the existing short code
* [ ] Add basic unit tests


## 📤 Submission

- Please submit your project as a private repository on GitHub, with an invitation to `joel at elishah dot biz` Username: `dashyonah`.
- Send an email to `joel at elishah dot biz` with the link to the private repository.

- Please include in your README.md:
  - Setup instructions
  - Sample `cURL` commands or a Postman collection link to download and test

- In your code, you should have:
  - SQL schema/migration files where needed
  - Clear folder structure and idiomatic Go code
  - Tests


## Setting up the project

You are not required to use this template.

Should you choose to, this project will require you to install the following dependencies:

- Git
- Docker
- Devbox
- Make

> Do note that Go, Python etc. dependencies are managed by Devbox and installed in a virtual environment.  You can use the `devbox shell` command to enter the devbox shell, and then use the `make` command to run the project.

### Setting up the environment

Please run the following command to setup the project:

```bash
make
```

### Running the project

Please run the following command to run the project:

```bash
make up
```

### Tearing down the project

Please run the following command to tear down the project:

```bash
make down
```
