# 🛒 Shopping Cart Service – Backend Assignment (Golang)

## 📌 Objective
Design and implement a backend service in Golang that allows users to manage shopping carts. This service should support typical cart operations such as adding/removing items, viewing the cart, and calculating totals.

## Setting up the project

This project will require you to install the following dependencies:

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

## 📦 Requirements

1. Implement the following endpoints with gRPC.  The endpoints should be implemented in the file `packages/proto/assets/cart/cart.proto`.

### 1. Entities to implement

- **User** (can be mocked or assumed `user_id=1`)
- **Product**
  - `id`: UUID
  - `name`: string
  - `price`: decimal (2dp)
- **CartItem**
  - `id`: UUID
  - `user_id`: UUID
  - `product_id`: UUID
  - `quantity`: int

### 2. Endpoints

- `POST /cart/items`
  Add item to cart (or increment quantity)

- `DELETE /cart/items/:product_id`
  Remove item from cart

- `PATCH /cart/items/:product_id`
  Update quantity of an item

- `GET /cart`
  Get current cart with item details and total cost

### 3. Business Logic

- If a product is added multiple times, increment the quantity.
- Return accurate **cart totals** (subtotal per item, overall total).
- Prevent quantity from being less than 1.
- Assume products already exist in the database.

2. Ensure the endpoints are documented with OpenAPI annotations and Swagger.  The `make up` command to generate the swagger.json file, and it will hot reload the swagger.json file when you make changes to the proto file.
3. Ensure all routes have tests written.  You should ensure > 80% coverage for the routes.
4. Please ensure interfaces are used where appropriate, including for test mocking.
5. You are free to edit the framework where you see fit.  You can use any framework you want, but please ensure the code is idiomatic and well-written.

---

## 💾 Tech Stack

- **Language**: Go
- **Framework**: Chi-router, gRPC, SQLX
- **Database**: MySQL (via GORM, SQLX, or `database/sql`.  The repository sets up a framework using SQLX.)

---

## 🔒 Constraints

- Assume user authentication is done; mock `user_id=1` for all endpoints.
- Don’t build full product CRUD — just seed a few sample products.

---

## 🧪 Bonus (Optional)

If you want to show off, you can implement the following:

- [ ] Add concurrency safety (prevent race conditions on cart updates)
- [ ] Handle inventory limits (e.g., prevent adding more than stock)
- [ ] Add rate limiting per user
- [ ] Add unit + integration tests
- [ ] Emit metrics or logs for key events.  The framework uses OpenTelemetry and Grafana stack.

---

## 📤 Submission

- Please submit your project as a private repository on GitHub, with an invitation to `joel at elishah dot biz` Username: `dashyonah`.
- Send an email to `joel at elishah dot biz` with the link to the private repository.

- Please include in your README.md:
  - Setup instructions
  - A Loom video showing the project setup and running (https://www.loom.com/looms/videos)
  - Sample `cURL` commands or a Postman collection link to download and test

- In your code, you should:
  - SQL schema/migration files
  - Clear folder structure and idiomatic Go code
  - Tests
