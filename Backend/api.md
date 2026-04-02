# API Documentation - URL Shortener Backend

## 1. Authentication

### Register User
**Method:** `POST`  
**Endpoint:** `/api/v1/auth/register`  
**Description:** Creates a new user account.  

**Request:**
```json
{
  "name"     : "string",
  "username" : "string",
  "password" : "string",
  "phone"    : "string",
  "email"    : "string"
}
```

**Response (201 Created):**
```json
{
  "id"        : "uuid",
  "username"  : "string",
  "email"     : "string",
  "created_at": "timestamp"
}
```

---

### Login User (OAuth2)
**Method:** `POST`  
**Endpoint:** `/oauth/token`  
**Description:** authenticates the user and returns access/refresh tokens.  

**Request:**
```json
{
  "grant_type"    : "password",
  "username"      : "string",
  "password"      : "string",
  "client_id"     : "string",
  "client_secret" : "string"
}
```

**Response (200 OK):**
```json
{
  "access_token"  : "string",
  "token_type"    : "Bearer",
  "expires_in"    : 3600,
  "refresh_token" : "string"
}
```

---

### Logout/Revoke Token
**Method:** `POST`  
**Endpoint:** `/oauth/revoke`  
**Headers:** `Authorization: Bearer <access_token>`  

**Request:**
```json
{
  "token": "string"
}
```

**Response (200 OK):**
```json
{
  "message": "logged out successfully"
}
```

---

## 2. URL Management

### Create Custom Short URL (Authenticated)
**Method:** `POST`  
**Endpoint:** `/api/v1/users/{id}/urls`  
**Headers:** `Authorization: Bearer <access_token>`  

**Request:**
```json
{
  "url"         : "https://example.com",
  "custom_code" : "short123"
}
```

**Response (201 Created):**
```json
{
  "id"           : "uint64",
  "code"         : "short123",
  "long_url"     : "https://example.com",
  "created_at"   : "timestamp"
}
```

---

### Remove URL
**Method:** `DELETE`  
**Endpoint:** `/api/v1/users/{id}/urls/{short_url_id}`  
**Headers:** `Authorization: Bearer <access_token>`  

**Response (24 No Content):**
*(No body)*

---

### Get URL Details
**Method:** `GET`  
**Endpoint:** `/api/v1/users/{id}/urls/{short_url_id}`  
**Headers:** `Authorization: Bearer <access_token>`  

**Response (200 OK):**
```json
{
  "id"           : "uint64",
  "code"         : "short123",
  "long_url"     : "https://example.com",
  "created_at"   : "timestamp",
  "expires_at"   : "timestamp"
}
```

---

### List User URLs
**Method:** `GET`  
**Endpoint:** `/api/v1/users/{id}/urls?limit=10&offset=0`  
**Headers:** `Authorization: Bearer <access_token>`  

**Response (200 OK):**
```json
{
  "data" : [
    {
      "id"           : "uint64",
      "code"         : "short123",
      "long_url"     : "https://example.com",
      "created_at"   : "timestamp"
    }
  ],
  "pagination": {
    "limit" : 10,
    "offset": 0,
    "total" : 100
  }
}
```

---

## 3. Public & Redirection

### Public Redirection
**Method:** `GET`  
**Endpoint:** `/{code}`  
**Description:** Redirects users to the original long URL stored in the cache/database.  

**Response (302 Found):**
*(Redirects to Location header)*

---

### Generate Guest Short URL
**Method:** `POST`  
**Endpoint:** `/api/v1/gen`  
**Headers:** `Authorization: Bearer <access_token>`  

**Request:**
```json
{
  "url": "https://example.com"
}
```

**Response (201 Created):**
```json
{
  "code"      : "short123",
  "expires_at": "timestamp"
}
```

---

### Get Guest URL Info
**Method:** `GET`  
**Endpoint:** `/api/v1/gen/{code}`  
**Headers:** `Authorization: Bearer <access_token>`  

**Response (200 OK):**
```json
{
  "long_url": "https://example.com",
  "code"    : "short123"
}
```

---

## 4. Analytics

### Get URL Analytics
**Method:** `GET`  
**Endpoint:** `/api/v1/users/{id}/urls/{short_url_id}/analytics`  
**Headers:** `Authorization: Bearer <access_token>`  

**Response (200 OK):**
```json
{
  "clicks"      : 10,
  "last_access" : "timestamp"
}
```
