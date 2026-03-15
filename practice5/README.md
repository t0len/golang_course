# Practice 5 — Go Pagination, Filtering & Friends API

## Prerequisites
- Go 1.21+
- PostgreSQL running locally

## Setup

### 1. Create database
```bash
psql -U postgres -c "CREATE DATABASE practice5;"
```

### 2. Run migrations
```bash
psql -U postgres -d practice5 -f migrations/001_create_tables.sql
psql -U postgres -d practice5 -f migrations/002_seed.sql
```

### 3. Install dependencies
```bash
go mod tidy
```

### 4. Run the server
```bash
go run cmd/main.go
```

Server starts on **:8080**

Environment variables (optional, defaults shown):
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=practice5
```

---

## API Endpoints

### GET /users
Paginated, filtered, and sorted list of users.

**Query params:**
| Param | Description | Example |
|-------|-------------|---------|
| page | Page number (default 1) | `page=1` |
| page_size | Items per page (default 10) | `page_size=5` |
| order_by | Column to sort by (default: id) | `order_by=name` |
| id | Filter by exact UUID | `id=a1000...` |
| name | Filter by name (ILIKE) | `name=alice` |
| email | Filter by email (ILIKE) | `email=example.com` |
| gender | Filter by gender | `gender=female` |
| birth_date | Filter by birth date (YYYY-MM-DD) | `birth_date=1995-03-15` |

**Examples:**
```
GET http://localhost:8080/users
GET http://localhost:8080/users?page=1&page_size=5&order_by=name
GET http://localhost:8080/users?name=alice
GET http://localhost:8080/users?gender=female&order_by=birth_date&page=1&page_size=3
GET http://localhost:8080/users?page=2&page_size=5&order_by=name&gender=male
```

---

### GET /users/common-friends
Returns common friends of two users.

**Query params:**
| Param | Description |
|-------|-------------|
| user1 | UUID of first user |
| user2 | UUID of second user |

**Example:**
```
GET http://localhost:8080/users/common-friends?user1=a1000000-0000-0000-0000-000000000001&user2=a1000000-0000-0000-0000-000000000002
```
→ Returns Carol, David, Eva (3 common friends)

---

## Demo Video Scenario (Postman)

1. **Basic pagination:**
   `GET /users?page=1&page_size=5`

2. **With order_by:**
   `GET /users?page=1&page_size=5&order_by=name`

3. **Filter by ID:**
   `GET /users?id=a1000000-0000-0000-0000-000000000001`

4. **Filter by Name:**
   `GET /users?name=alice`

5. **Filter by Email:**
   `GET /users?email=example.com&page=1&page_size=5&order_by=email`

6. **Filter by 3 fields + pagination + order_by:**
   `GET /users?gender=female&page=1&page_size=3&order_by=birth_date`

7. **GetCommonFriends:**
   `GET /users/common-friends?user1=a1000000-0000-0000-0000-000000000001&user2=a1000000-0000-0000-0000-000000000002`

8. **Show DB in psql/TablePlus:** `SELECT * FROM users;` and `SELECT * FROM user_friends;`
