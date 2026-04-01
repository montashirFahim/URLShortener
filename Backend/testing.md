# API Testing Guide - URL Shortener Backend (Normalized)

This guide provides real-world `curl` commands for all 12 implemented API endpoints, updated for the **normalized schema**.

> [!NOTE]
> Replace variables like `<access_token>`, `<id>`, `<short_url_id>`, and `<code>` with actual values from your session.

## 1. Authentication & Registration

### 1.1 User Registration
**Endpoint:** `POST /api/v1/auth/register`
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
     -H "Content-Type: application/json" \
     -d '{
       "name": "Jane Doe",
       "username": "janedoe",
       "password": "securepassword",
       "email": "jane@example.com",
       "phone": "555-0199"
     }'
```

### 1.2 Get Access Token (Password Grant)
**Endpoint:** `POST /oauth/token`
```bash
curl -X POST http://localhost:8080/oauth/token \
     -H "Content-Type: application/json" \
     -d '{
       "grant_type": "password",
       "username": "janedoe",
       "password": "securepassword",
       "client_id": "client123",
       "client_secret": "secret123"
     }'
```

---

## 2. URL Management (Authenticated)

### 2.1 Create Short Link
**Endpoint:** `POST /api/v1/users/{id}/urls`
```bash
curl -X POST http://localhost:8080/api/v1/users/<id>/urls \
     -H "Authorization: Bearer <access_token>" \
     -H "Content-Type: application/json" \
     -d '{
       "url": "https://www.github.com",
       "custom_code": "mygithub"
     }'
```

### 2.2 List User Links
**Endpoint:** `GET /api/v1/users/{id}/urls?limit=10&offset=0`
```bash
curl -G http://localhost:8080/api/v1/users/<id>/urls \
     -H "Authorization: Bearer <access_token>" \
     -d "limit=5" \
     -d "offset=0"
```

### 2.3 Get URL Details
**Endpoint:** `GET /api/v1/users/{id}/urls/{short_url_id}`
```bash
curl http://localhost:8080/api/v1/users/<id>/urls/<short_url_id> \
     -H "Authorization: Bearer <access_token>"
```

### 2.4 Delete Link
**Endpoint:** `DELETE /api/v1/users/{id}/urls/{short_url_id}`
```bash
curl -X DELETE http://localhost:8080/api/v1/users/<id>/urls/<short_url_id> \
     -H "Authorization: Bearer <access_token>"
```

---

## 3. Guest & Redirection

### 3.1 Generate Guest URL
**Endpoint:** `POST /api/v1/gen`
```bash
curl -X POST http://localhost:8080/api/v1/gen \
     -H "Authorization: Bearer <access_token>" \
     -H "Content-Type: application/json" \
     -d '{"url": "https://www.google.com"}'
```

### 3.2 Get Guest URL Info
**Endpoint:** `GET /api/v1/gen/{code}`
```bash
curl http://localhost:8080/api/v1/gen/<code> \
     -H "Authorization: Bearer <access_token>"
```

### 3.3 Redirect (Public)
**Endpoint:** `GET /{code}`
```bash
curl -i http://localhost:8080/<code>
```
*Expected: 302 Found with Location header.*

---

## 4. Analytics

### 4.1 Get URL Stats
**Endpoint:** `GET /api/v1/users/{id}/urls/{short_url_id}/analytics`
```bash
curl http://localhost:8080/api/v1/users/<id>/urls/<short_url_id>/analytics \
     -H "Authorization: Bearer <access_token>"
```
*Returns click counts and last access time from the `url_stats` table.*
