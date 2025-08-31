# JobGen API — Frontend Reference

**Base URL:** `{{HOST}}/api/v1` *(replace `{{HOST}}` with the environment host, e.g. `https://api.jobgen.io` or `http://localhost:8080`)*

**Overview**
This document is a living Markdown reference for the frontend team describing all currently implemented API routes, request/response shapes, authentication, and useful examples (cURL). I'll update this file as new routes are implemented.

---

## Table of contents

* [Authentication / Security](#authentication--security)
* [Standard response format](#standard-response-format)
* [Admin endpoints](#admin-endpoints)

  * [GET /admin/users](#get-adminusers)
  * [DELETE /admin/users/{user\_id}](#delete-adminusersuser_id)
  * [PUT /admin/users/{user\_id}/role](#put-adminusersuser_idrole)
  * [PUT /admin/users/{user\_id}/toggle-status](#put-adminusersuser_idtoggle-status)
* [Authentication endpoints](#authentication-endpoints)

  * [POST /auth/register](#post-authregister)
  * [POST /auth/login](#post-authlogin)
  * [POST /auth/logout](#post-authlogout)
  * [POST /auth/refresh](#post-authrefresh)
  * [POST /auth/verify-email](#post-authverify-email)
  * [POST /auth/resend-otp](#post-authresend-otp)
  * [POST /auth/forgot-password](#post-authforgot-password)
  * [POST /auth/reset-password](#post-authreset-password)
  * [POST /auth/change-password](#post-authchange-password)
* [User profile endpoints](#user-profile-endpoints)

  * [GET /users/profile](#get-usersprofile)
  * [PUT /users/profile](#put-usersprofile)
  * [DELETE /users/account](#delete-usersaccount)

---

## Authentication & Security

* Most protected routes use **Bearer JWT** authentication.
* Send header: `Authorization: Bearer <access_token>`
* All requests using JSON must include header: `Content-Type: application/json`
* When a route requires auth, it is documented under the route as **Auth: Required**.

## Standard response format

All endpoints return the project's standard wrapper type `controllers.StandardResponse`:

```json
{
  "success": true,
  "message": "Human readable message",
  "data": { /* object or array depending on endpoint */ },
  "error": {
    "code": "optional_error_code",
    "message": "error message",
    "details": {}
  }
}
```

* `success` boolean indicates request-level success.
* `message` is a short human readable message suitable for UI to display.
* `data` contains the payload on success — its shape varies by endpoint.
* `error` contains details when `success` is `false`.

> **Note:** The exact contents of `data` differ per endpoint (e.g., login returns tokens). Where the API definition did not include explicit `data` fields, this document includes suggested example payloads; please confirm with the backend if you rely on specific fields.

---

# Admin endpoints

### GET /admin/users

**Description:** Get a paginated list of users (admin only)

* **URL:** `/admin/users`
* **Method:** `GET`
* **Auth:** Required (Bearer JWT, admin role)
* **Query Parameters**

  * `page` (integer, default: 1) — page number
  * `limit` (integer, default: 10) — items per page
  * `role` (string: `user` | `admin`) — filter by role
  * `active` (boolean) — filter by active status
  * `search` (string) — search in email, username, or full\_name
  * `sort_by` (string, default: `created_at`) — field to sort by
  * `sort_order` (string: `asc` | `desc`, default: `desc`)

**Success response (200)** — `data` will include a paginated list. Suggested shape:

```json
{
  "success": true,
  "message": "Users retrieved",
  "data": {
    "items": [
      {
        "id": "user-id",
        "username": "jdoe",
        "email": "jdoe@example.com",
        "full_name": "John Doe",
        "role": "user",
        "active": true,
        "created_at": "2025-08-01T12:34:56Z"
      }
    ],
    "page": 1,
    "limit": 10,
    "total": 123
  }
}
```

**Errors:** `401 Unauthorized`, `403 Forbidden` (if not admin)

**cURL example:**

```bash
curl -H "Authorization: Bearer $ACCESS" \
  "{{HOST}}/api/v1/admin/users?page=1&limit=20&sort_by=created_at&sort_order=desc"
```

---

### DELETE /admin/users/{user\_id}

**Description:** Delete a user account permanently (admin only)

* **URL:** `/admin/users/{user_id}`
* **Method:** `DELETE`
* **Auth:** Required (admin)
* **Path Params:** `user_id` (string)

**Success response (200):** standard response with success message.
**Errors:** `400 Bad Request`, `401 Unauthorized`, `403 Forbidden`.

**cURL example:**

```bash
curl -X DELETE -H "Authorization: Bearer $ADMIN_TOKEN" \
  "{{HOST}}/api/v1/admin/users/63f2a3..."
```

---

### PUT /admin/users/{user\_id}/role

**Description:** Update the role of a user (admin only)

* **URL:** `/admin/users/{user_id}/role`
* **Method:** `PUT`
* **Auth:** Required (admin)
* **Path Params:** `user_id` (string)
* **Body (JSON)**

  * `role` (string) — one of: `user`, `admin` (required)

**Example request body:**

```json
{ "role": "admin" }
```

**Success (200):** updated user role in `data` or confirmation message.
**Errors:** `400`, `401`, `403`.

**cURL example:**

```bash
curl -X PUT -H "Content-Type: application/json" -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{"role":"admin"}' \
  "{{HOST}}/api/v1/admin/users/63f2a3.../role"
```

---

### PUT /admin/users/{user\_id}/toggle-status

**Description:** Activate or deactivate a user account (admin only)

* **URL:** `/admin/users/{user_id}/toggle-status`
* **Method:** `PUT`
* **Auth:** Required (admin)
* **Path Params:** `user_id` (string)

**Success (200):** toggled user `active` status; response `data` may include updated user info.
**Errors:** `400`, `401`, `403`.

**cURL example:**

```bash
curl -X PUT -H "Authorization: Bearer $ADMIN_TOKEN" \
  "{{HOST}}/api/v1/admin/users/63f2a3.../toggle-status"
```

---

# Authentication endpoints

> Note: Request/validation schemas are defined under `definitions` in the swagger. Required fields are documented per endpoint.

### POST /auth/register

**Description:** Register a new user account (email verification/OTP flow)

* **URL:** `/auth/register`
* **Method:** `POST`
* **Auth:** Public
* **Body (JSON)** — `controllers.RegisterRequest` (required fields shown):

  * `email` (string) **required**
  * `full_name` (string) **required**
  * `username` (string, min 3, max 30) **required**
  * `password` (string, minLength 8) **required**
  * optional: `bio`, `location`, `phone_number`, `experience_years` (integer), `skills` (array of string)

**Success (201):** user created; likely returns message and user `id` or a partial user in `data`.
**Errors:** `400 Bad Request`, `409 Conflict` when user exists.

**Example request:**

```json
{
  "email": "jdoe@example.com",
  "full_name": "John Doe",
  "username": "jdoe",
  "password": "S3cur3Pass!",
  "skills": ["go","react"]
}
```

**cURL:**

```bash
curl -X POST -H "Content-Type: application/json" \
  -d '{"email":"jdoe@example.com","full_name":"John Doe","username":"jdoe","password":"S3cur3Pass!"}' \
  "{{HOST}}/api/v1/auth/register"
```

---

### POST /auth/login

**Description:** Authenticate and receive tokens

* **URL:** `/auth/login`
* **Method:** `POST`
* **Auth:** Public
* **Body (JSON)** — `controllers.LoginRequest` (required):

  * `email` (string)
  * `password` (string)

**Success (200):** returns tokens in `data`. Typical suggested `data` shape (confirm with backend):

```json
"data": {
  "access_token": "eyJ...",
  "refresh_token": "rftok...",
  "expires_in": 3600
}
```

**Errors:** `400 Bad Request`, `401 Invalid credentials`.

**cURL example:**

```bash
curl -X POST -H "Content-Type: application/json" \
  -d '{"email":"jdoe@example.com","password":"S3cur3Pass!"}' \
  "{{HOST}}/api/v1/auth/login"
```

---

### POST /auth/logout

**Description:** Logout logged-in user and invalidate tokens

* **URL:** `/auth/logout`
* **Method:** `POST`
* **Auth:** Required (Bearer)

**Success (200):** confirmation message.
**Errors:** `401 Unauthorized`.

**cURL:**

```bash
curl -X POST -H "Authorization: Bearer $ACCESS_TOKEN" \
  "{{HOST}}/api/v1/auth/logout"
```

---

### POST /auth/refresh

**Description:** Exchange a refresh token for a new access token (and possibly a new refresh token)

* **URL:** `/auth/refresh`
* **Method:** `POST`
* **Auth:** Public (sends refresh token in body)
* **Body (JSON)** — `controllers.RefreshTokenRequest` (required):

  * `refresh_token` (string)

**Success (200):** `data` will contain `access_token` (and possibly `refresh_token`).
**Errors:** `400 Bad Request`, `401 Invalid or expired refresh token`.

**cURL:**

```bash
curl -X POST -H "Content-Type: application/json" \
  -d '{"refresh_token":"rftok..."}' \
  "{{HOST}}/api/v1/auth/refresh"
```

---

### POST /auth/verify-email

**Description:** Verify user email using OTP code

* **URL:** `/auth/verify-email`
* **Method:** `POST`
* **Auth:** Public
* **Body (JSON)** — `controllers.VerifyEmailRequest` (required):

  * `email` (string)
  * `otp` (string)

**Success (200):** email verified; user account becomes verified.
**Errors:** `400 Bad Request` (invalid or expired OTP).

---

### POST /auth/resend-otp

**Description:** Resend OTP verification to email

* **URL:** `/auth/resend-otp`
* **Method:** `POST`
* **Auth:** Public
* **Body (JSON)** — `controllers.ResendOTPRequest` (required): `email` (string)

**Success (200):** OTP resent. Errors: `400`, `404` if user not found.

---

### POST /auth/forgot-password

**Description:** Send password reset email with token

* **URL:** `/auth/forgot-password`
* **Method:** `POST`
* **Auth:** Public
* **Body (JSON)** — `controllers.RequestPasswordResetRequest` (required): `email` (string)

**Success (200):** confirmation that reset email was sent.
**Errors:** `400`.

---

### POST /auth/reset-password

**Description:** Reset password using a token sent via email

* **URL:** `/auth/reset-password`
* **Method:** `POST`
* **Auth:** Public
* **Body (JSON)** — `controllers.ResetPasswordRequest` (required):

  * `token` (string)
  * `new_password` (string, minLength: 8)

**Success (200):** password reset confirmation.
**Errors:** `400`.

---

### POST /auth/change-password

**Description:** Change password while logged in

* **URL:** `/auth/change-password`
* **Method:** `POST`
* **Auth:** Required (Bearer)
* **Body (JSON)** — `controllers.ChangePasswordRequest` (required):

  * `old_password` (string)
  * `new_password` (string, minLength: 8)

**Success (200):** password updated.
**Errors:** `400`, `401`.

---

# User Profile endpoints

### GET /users/profile

**Description:** Get the current authenticated user's profile

* **URL:** `/users/profile`
* **Method:** `GET`
* **Auth:** Required (Bearer)

**Success (200):** `data` will contain the user profile. Suggested shape:

```json
{
  "id": "user-id",
  "username": "jdoe",
  "email": "jdoe@example.com",
  "full_name": "John Doe",
  "bio":"...",
  "skills":["go","react"],
  "experience_years": 3,
  "location":"Addis Ababa",
  "profile_picture":"https://.../avatar.jpg",
  "active": true
}
```

**Errors:** `401`, `404`.

**cURL:**

```bash
curl -H "Authorization: Bearer $ACCESS_TOKEN" \
  "{{HOST}}/api/v1/users/profile"
```

---

### PUT /users/profile

**Description:** Update the current authenticated user's profile

* **URL:** `/users/profile`
* **Method:** `PUT`
* **Auth:** Required
* **Body (JSON)** — `controllers.UpdateProfileRequest` fields are optional:

  * `full_name`, `bio`, `location`, `phone_number`, `profile_picture`, `experience_years`, `skills` (array)

**Success (200):** updated profile in `data`.
**Errors:** `400`, `401`.

**Example request:**

```json
{
  "full_name": "John Q. Doe",
  "bio": "5+ years building Go services",
  "skills": ["go","docker"]
}
```

---

### DELETE /users/account

**Description:** Permanently delete current user's account

* **URL:** `/users/account`
* **Method:** `DELETE`
* **Auth:** Required (Bearer)

**Success (200):** account deleted confirmation.
**Errors:** `401`.

**cURL:**

```bash
curl -X DELETE -H "Authorization: Bearer $ACCESS_TOKEN" \
  "{{HOST}}/api/v1/users/account"
```

---

## Common notes for frontend integration

* **Content-Type:** Always set `Content-Type: application/json` unless uploading files (not covered here).
* **Auth header:** `Authorization: Bearer <access_token>` — the frontend should attach the token to every protected request.
* **Token storage:** Use secure storage (HttpOnly cookies recommended for web apps; localStorage can be used if you handle XSS carefully).
* **Pagination:** `/admin/users` returns `page`, `limit`, `total`, and `items`. If you need cursor-style pagination later, backend will add it.
* **Error handling:** If `success: false`, show `error.message` (or `message`) and optionally show `error.details` for dev/debug.
* **Validation:** Backend returns `400` for validation errors. Show field-level messages where available.

---

## Things I couldn't infer precisely (please confirm with backend)

1. **Exact `data` payloads** for some endpoints (e.g., login token property names or user object fields) were not explicitly defined in the swagger snippet. I added sensible example shapes in this doc — ask the backend to confirm exact field names (especially token fields like `access_token`, `refresh_token` and `expires_in`).
2. **Status codes & error codes** beyond the ones listed (some endpoints may return additional codes).

If you want, I can:

* Add concrete JSON response examples captured from the running server (you can paste sample responses here), or
* Add TypeScript interfaces for the frontend based on the definitions in this doc.

---

*Document generated from the current swagger docs. I will update this as new routes are implemented.*
