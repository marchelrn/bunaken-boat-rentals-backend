# API Documentation - Bunaken Boat Rentals Backend

## Base URL

```
Production: https://bunakencharter.up.railway.app/api
Local: http://localhost:8080/api
```

## Authentication

API menggunakan JWT (JSON Web Token) untuk autentikasi. Untuk endpoint yang dilindungi, sertakan token di header Authorization:

```
Authorization: Bearer <your_token>
```

Token diperoleh dari endpoint `/auth/login`.

---

## Response Format

Semua response mengikuti format standar:

### Success Response
```json
{
  "success": true,
  "message": "Pesan sukses",
  "data": { ... }
}
```

### Error Response
```json
{
  "success": false,
  "message": "Pesan error",
  "data": null
}
```

---

## Error Codes

| Status Code | Description |
|------------|-------------|
| 200 | Success |
| 400 | Bad Request (Format JSON salah, validasi gagal) |
| 401 | Unauthorized (Token tidak valid/tidak ditemukan) |
| 404 | Not Found (Resource tidak ditemukan) |
| 500 | Internal Server Error |

---

## Endpoints

### Authentication

#### 1. Register (Create Admin Account)

**Endpoint:** `POST /auth/register`

**Description:** Membuat akun admin baru. **Catatan:** Endpoint ini sebaiknya dinonaktifkan setelah admin pertama dibuat untuk keamanan.

**Request Body:**
```json
{
  "username": "admin",
  "password": "password123",
  "role": "admin"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Registrasi berhasil",
  "data": null
}
```

**Error Response (400 Bad Request):**
```json
{
  "success": false,
  "message": "Username mungkin sudah digunakan",
  "data": null
}
```

---

#### 2. Login

**Endpoint:** `POST /auth/login`

**Description:** Login untuk mendapatkan JWT token.

**Request Body:**
```json
{
  "username": "admin",
  "password": "password123"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Login berhasil",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "role": "admin"
  }
}
```

**Error Response (400 Bad Request):**
```json
{
  "success": false,
  "message": "Username atau password salah",
  "data": null
}
```

---

### Packages (Public)

#### 3. Get All Packages

**Endpoint:** `GET /packages`

**Description:** Mendapatkan semua paket wisata yang tersedia.

**Headers:** Tidak diperlukan (Public endpoint)

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Berhasil mengambil data Packages",
  "data": [
    {
      "ID": 1,
      "CreatedAt": "2024-01-01T00:00:00Z",
      "UpdatedAt": "2024-01-01T00:00:00Z",
      "DeletedAt": null,
      "name": "Kapal Speed",
      "description": "Kapal cepat untuk perjalanan singkat",
      "capacity": "10 orang",
      "duration": "4 jam",
      "is_popular": true,
      "image_url": "https://example.com/image.jpg",
      "routes": [
        {
          "name": "Bunaken",
          "price": "1200000"
        },
        {
          "name": "Bunaken - Siladen",
          "price": "1500000"
        }
      ],
      "features": [
        "Snorkeling equipment",
        "Life jacket",
        "Guide"
      ],
      "excludes": [
        "Makan siang",
        "Minuman"
      ]
    }
  ]
}
```

---

#### 4. Get Package by ID

**Endpoint:** `GET /packages/:id`

**Description:** Mendapatkan detail paket wisata berdasarkan ID.

**Parameters:**
- `id` (path parameter) - ID paket

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Berhasil mengambil detail Package",
  "data": {
    "ID": 1,
    "CreatedAt": "2024-01-01T00:00:00Z",
    "UpdatedAt": "2024-01-01T00:00:00Z",
    "DeletedAt": null,
    "name": "Kapal Speed",
    "description": "Kapal cepat untuk perjalanan singkat",
    "capacity": "10 orang",
    "duration": "4 jam",
    "is_popular": true,
    "image_url": "https://example.com/image.jpg",
    "routes": [
      {
        "name": "Bunaken",
        "price": "1200000"
      }
    ],
    "features": ["Snorkeling equipment"],
    "excludes": ["Makan siang"]
  }
}
```

**Error Response (404 Not Found):**
```json
{
  "success": false,
  "message": "Package tidak ditemukan",
  "data": null
}
```

---

### Packages (Protected - Admin Only)

#### 5. Create Package

**Endpoint:** `POST /admin/packages`

**Description:** Membuat paket wisata baru. **Memerlukan autentikasi.**

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "name": "Kapal Speed",
  "description": "Kapal cepat untuk perjalanan singkat",
  "capacity": "10 orang",
  "duration": "4 jam",
  "is_popular": true,
  "image_url": "https://example.com/image.jpg",
  "routes": [
    {
      "name": "Bunaken",
      "price": "1200000"
    },
    {
      "name": "Bunaken - Siladen",
      "price": "1500000"
    }
  ],
  "features": [
    "Snorkeling equipment",
    "Life jacket",
    "Guide"
  ],
  "excludes": [
    "Makan siang",
    "Minuman"
  ]
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Berhasil membuat package baru",
  "data": {
    "ID": 1,
    "name": "Kapal Speed",
    ...
  }
}
```

**Error Response (400 Bad Request):**
```json
{
  "success": false,
  "message": "Format JSON salah: ...",
  "data": null
}
```

**Error Response (401 Unauthorized):**
```json
{
  "success": false,
  "message": "Token tidak ditemukan",
  "data": null
}
```

---

#### 6. Update Package

**Endpoint:** `PUT /admin/packages/:id`

**Description:** Mengupdate paket wisata yang sudah ada. **Memerlukan autentikasi.**

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Parameters:**
- `id` (path parameter) - ID paket yang akan diupdate

**Request Body:**
```json
{
  "name": "Kapal Speed Updated",
  "description": "Deskripsi baru",
  "capacity": "12 orang",
  "duration": "5 jam",
  "is_popular": false,
  "image_url": "https://example.com/new-image.jpg",
  "routes": [
    {
      "name": "Bunaken",
      "price": "1300000"
    }
  ],
  "features": ["Snorkeling equipment"],
  "excludes": ["Makan siang"]
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Berhasil update package",
  "data": {
    "ID": 1,
    "name": "Kapal Speed Updated",
    ...
  }
}
```

**Error Response (404 Not Found):**
```json
{
  "success": false,
  "message": "Package tidak ditemukan",
  "data": null
}
```

---

#### 7. Delete Package

**Endpoint:** `DELETE /admin/packages/:id`

**Description:** Menghapus paket wisata. **Memerlukan autentikasi.**

**Headers:**
```
Authorization: Bearer <token>
```

**Parameters:**
- `id` (path parameter) - ID paket yang akan dihapus

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Berhasil menghapus package",
  "data": null
}
```

**Error Response (404 Not Found):**
```json
{
  "success": false,
  "message": "Package tidak ditemukan",
  "data": null
}
```

---

### Add-Ons (Public)

#### 8. Get All Add-Ons

**Endpoint:** `GET /addons`

**Description:** Mendapatkan semua layanan add-on yang tersedia.

**Headers:** Tidak diperlukan (Public endpoint)

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Berhasil mengambil data Add-Ons",
  "data": [
    {
      "ID": 1,
      "CreatedAt": "2024-01-01T00:00:00Z",
      "UpdatedAt": "2024-01-01T00:00:00Z",
      "DeletedAt": null,
      "name": "Kamera Underwater",
      "price": "150000",
      "description": "Sewa kamera underwater untuk dokumentasi"
    },
    {
      "ID": 2,
      "CreatedAt": "2024-01-01T00:00:00Z",
      "UpdatedAt": "2024-01-01T00:00:00Z",
      "DeletedAt": null,
      "name": "Makan Siang",
      "price": "100000",
      "description": "Paket makan siang"
    }
  ]
}
```

---

#### 9. Get Add-On by ID

**Endpoint:** `GET /addons/:id`

**Description:** Mendapatkan detail add-on berdasarkan ID.

**Parameters:**
- `id` (path parameter) - ID add-on

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Berhasil mengambil detail Add-On",
  "data": {
    "ID": 1,
    "CreatedAt": "2024-01-01T00:00:00Z",
    "UpdatedAt": "2024-01-01T00:00:00Z",
    "DeletedAt": null,
    "name": "Kamera Underwater",
    "price": "150000",
    "description": "Sewa kamera underwater untuk dokumentasi"
  }
}
```

**Error Response (404 Not Found):**
```json
{
  "success": false,
  "message": "Add-On tidak ditemukan",
  "data": null
}
```

---

### Add-Ons (Protected - Admin Only)

#### 10. Create Add-On

**Endpoint:** `POST /admin/addons`

**Description:** Membuat add-on baru. **Memerlukan autentikasi.**

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "name": "Kamera Underwater",
  "price": "150000",
  "description": "Sewa kamera underwater untuk dokumentasi"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Berhasil membuat add-on baru",
  "data": {
    "ID": 1,
    "name": "Kamera Underwater",
    "price": "150000",
    "description": "Sewa kamera underwater untuk dokumentasi"
  }
}
```

**Error Response (400 Bad Request):**
```json
{
  "success": false,
  "message": "Format JSON salah: ...",
  "data": null
}
```

---

#### 11. Update Add-On

**Endpoint:** `PUT /admin/addons/:id`

**Description:** Mengupdate add-on yang sudah ada. **Memerlukan autentikasi.**

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Parameters:**
- `id` (path parameter) - ID add-on yang akan diupdate

**Request Body:**
```json
{
  "name": "Kamera Underwater Pro",
  "price": "200000",
  "description": "Sewa kamera underwater profesional untuk dokumentasi"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Berhasil update add-on",
  "data": {
    "ID": 1,
    "name": "Kamera Underwater Pro",
    "price": "200000",
    "description": "Sewa kamera underwater profesional untuk dokumentasi"
  }
}
```

**Error Response (404 Not Found):**
```json
{
  "success": false,
  "message": "Add-On tidak ditemukan",
  "data": null
}
```

---

#### 12. Delete Add-On

**Endpoint:** `DELETE /admin/addons/:id`

**Description:** Menghapus add-on. **Memerlukan autentikasi.**

**Headers:**
```
Authorization: Bearer <token>
```

**Parameters:**
- `id` (path parameter) - ID add-on yang akan dihapus

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Berhasil menghapus add-on",
  "data": null
}
```

**Error Response (404 Not Found):**
```json
{
  "success": false,
  "message": "Add-On tidak ditemukan",
  "data": null
}
```

---

### Authentication (Protected - Admin Only)

#### 13. Change Password

**Endpoint:** `PUT /admin/change-password`

**Description:** Mengubah password admin. **Memerlukan autentikasi.**

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "old_password": "password123",
  "new_password": "newpassword456"
}
```

**Validation:**
- `old_password`: Required
- `new_password`: Required, minimum 6 karakter

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Password berhasil diubah",
  "data": null
}
```

**Error Response (400 Bad Request):**
```json
{
  "success": false,
  "message": "Password lama salah",
  "data": null
}
```

**Error Response (400 Bad Request - Validation):**
```json
{
  "success": false,
  "message": "Format JSON salah: Key: 'ChangePasswordInput.NewPassword' Error:Field validation for 'NewPassword' failed on the 'min' tag",
  "data": null
}
```

**Error Response (401 Unauthorized):**
```json
{
  "success": false,
  "message": "User ID tidak ditemukan",
  "data": null
}
```

---

## Data Models

### Package Model

```json
{
  "ID": 1,
  "CreatedAt": "2024-01-01T00:00:00Z",
  "UpdatedAt": "2024-01-01T00:00:00Z",
  "DeletedAt": null,
  "name": "string",
  "description": "string",
  "capacity": "string",
  "duration": "string",
  "is_popular": boolean,
  "image_url": "string",
  "routes": [
    {
      "name": "string",
      "price": "string"
    }
  ],
  "features": ["string"],
  "excludes": ["string"]
}
```

### Add-On Model

```json
{
  "ID": 1,
  "CreatedAt": "2024-01-01T00:00:00Z",
  "UpdatedAt": "2024-01-01T00:00:00Z",
  "DeletedAt": null,
  "name": "string",
  "price": "string",
  "description": "string"
}
```

### User Model

```json
{
  "ID": 1,
  "CreatedAt": "2024-01-01T00:00:00Z",
  "UpdatedAt": "2024-01-01T00:00:00Z",
  "DeletedAt": null,
  "username": "string",
  "password": "string (hashed)",
  "role": "string"
}
```

---

## Example Usage

### cURL Examples

#### Login
```bash
curl -X POST https://bunakencharter.up.railway.app/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "password123"
  }'
```

#### Get All Packages
```bash
curl -X GET https://bunakencharter.up.railway.app/api/packages
```

#### Create Package (Protected)
```bash
curl -X POST https://bunakencharter.up.railway.app/api/admin/packages \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "name": "Kapal Speed",
    "description": "Kapal cepat untuk perjalanan singkat",
    "capacity": "10 orang",
    "duration": "4 jam",
    "is_popular": true,
    "image_url": "https://example.com/image.jpg",
    "routes": [
      {
        "name": "Bunaken",
        "price": "1200000"
      }
    ],
    "features": ["Snorkeling equipment"],
    "excludes": ["Makan siang"]
  }'
```

#### Change Password (Protected)
```bash
curl -X PUT https://bunakencharter.up.railway.app/api/admin/change-password \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "old_password": "password123",
    "new_password": "newpassword456"
  }'
```

---

## Notes

1. **CORS**: API mengizinkan semua origin untuk development. Untuk production, sebaiknya dibatasi ke domain tertentu.

2. **Token Expiration**: Token JWT saat ini tidak memiliki expiration time. Untuk production, sebaiknya ditambahkan expiration time.

3. **Password Security**: Password disimpan dalam bentuk hash menggunakan bcrypt.

4. **Soft Delete**: Delete operations menggunakan soft delete (tidak benar-benar menghapus dari database).

5. **Register Endpoint**: Endpoint `/auth/register` sebaiknya dinonaktifkan setelah admin pertama dibuat untuk keamanan.

---

## Support

Untuk pertanyaan atau masalah, silakan hubungi tim development.

