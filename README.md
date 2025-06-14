# 🛒 Shopping Cart Service – Backend Assignment (Golang)

## 📌 Objective
Design and implement a backend service in Golang that allows users to manage shopping carts. This service should support typical cart operations such as adding/removing items, viewing the cart, and calculating totals.

---

## 📦 Requirements

### 1. Entities

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

---

## 💾 Tech Stack

- **Language**: Go
- **Framework**: Any (Gin, Echo, Chi, or `net/http`)
- **Database**: PostgreSQL (via GORM, SQLX, or `database/sql`)
- *Optional*: Use Devbox or Docker if it helps

---

## 🔒 Constraints

- Assume user authentication is done; mock `user_id=1` for all endpoints.
- Don’t build full product CRUD — just seed a few sample products.

---

## 🧪 Bonus (Optional)

If you want to show off:

- [ ] Add concurrency safety (prevent race conditions on cart updates)
- [ ] Handle inventory limits (e.g., prevent adding more than stock)
- [ ] Add rate limiting per user
- [ ] Add unit + integration tests
- [ ] Emit metrics or logs for key events

---

## 📤 Submission

Please include:

- `README.md` with setup instructions
- Sample `cURL` commands or a Postman collection
- SQL schema/migration files
- Clear folder structure and idiomatic Go code
