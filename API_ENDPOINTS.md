# 🌐 API Endpoints Guide

Panduan lengkap URL untuk mengakses backend API dari berbagai environment.

## 📍 Base URL berdasarkan Environment

### 1. Backend berjalan di Docker (Recommended)

```javascript
// Frontend di browser (localhost atau network lain)
const API_BASE_URL = "http://localhost:8080"

// Contoh fetch
fetch("http://localhost:8080/health")
fetch("http://localhost:8080/api/v1/menu")
```

**Penjelasan:**
- API berjalan di Docker container
- Di-expose ke host di port `8080`
- Frontend (browser) akses via `localhost:8080`

---

### 2. Backend berjalan di Docker, Frontend berjalan di Docker (same network)

```javascript
// Jika frontend juga di Docker dan dalam network yang sama
const API_BASE_URL = "http://api:8080"

// Atau jika beda network, gunakan host.docker.internal (Windows/Mac)
const API_BASE_URL = "http://host.docker.internal:8080"
```

---

### 3. Backend berjalan Lokal (tanpa Docker)

```javascript
// Backend running dengan: go run ./cmd/api
const API_BASE_URL = "http://localhost:8080"
```

---

## 🎯 URL Endpoints yang Tersedia

Berdasarkan struktur POS FnB Backend:

### Public Endpoints (Tanpa Auth)

```javascript
// Health Check
GET http://localhost:8080/health

// Get Menu by Restaurant
GET http://localhost:8080/api/v1/menu?restaurant_id={uuid}

// Get Table Info (scan QR)
GET http://localhost:8080/api/v1/table/{qr_token}

// Create Order (customer)
POST http://localhost:8080/api/v1/orders
```

### Staff Endpoints (Perlu JWT Token)

```javascript
// Create Table
POST http://localhost:8080/api/v1/staff/tables

// Scan Order QR (kasir)
POST http://localhost:8080/api/v1/staff/orders/scan-qr

// Get All Orders
GET http://localhost:8080/api/v1/staff/orders

// Update Order Status
PATCH http://localhost:8080/api/v1/staff/orders/{id}/status
```

---

## 💻 Contoh Implementasi Frontend

### React / Next.js

```typescript
// config/api.ts
export const API_CONFIG = {
  BASE_URL: process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080',
  TIMEOUT: 10000,
}

// lib/api.ts
import axios from 'axios'
import { API_CONFIG } from '@/config/api'

const apiClient = axios.create({
  baseURL: API_CONFIG.BASE_URL,
  timeout: API_CONFIG.TIMEOUT,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Interceptor untuk JWT token
apiClient.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

export default apiClient

// services/menu.service.ts
export const getMenu = async (restaurantId: string) => {
  const response = await apiClient.get('/api/v1/menu', {
    params: { restaurant_id: restaurantId }
  })
  return response.data
}

export const createOrder = async (orderData: any) => {
  const response = await apiClient.post('/api/v1/orders', orderData)
  return response.data
}
```

**Environment Variables (.env.local):**
```bash
NEXT_PUBLIC_API_URL=http://localhost:8080
```

---

### Vue.js / Nuxt.js

```typescript
// plugins/axios.ts
import axios from 'axios'

const apiClient = axios.create({
  baseURL: import.meta.env.VITE_API_URL || 'http://localhost:8080',
  timeout: 10000,
})

// Request interceptor
apiClient.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

export default apiClient

// composables/useMenu.ts
export const useMenu = () => {
  const getMenu = async (restaurantId: string) => {
    try {
      const { data } = await apiClient.get('/api/v1/menu', {
        params: { restaurant_id: restaurantId }
      })
      return data
    } catch (error) {
      console.error('Error fetching menu:', error)
      throw error
    }
  }

  return { getMenu }
}
```

**Environment Variables (.env):**
```bash
VITE_API_URL=http://localhost:8080
```

---

### Vanilla JavaScript / HTML

```html
<!DOCTYPE html>
<html>
<head>
  <title>POS FnB</title>
</head>
<body>
  <script>
    const API_BASE_URL = 'http://localhost:8080'

    // Health Check
    async function checkHealth() {
      try {
        const response = await fetch(`${API_BASE_URL}/health`)
        const data = await response.json()
        console.log('Health:', data)
      } catch (error) {
        console.error('Error:', error)
      }
    }

    // Get Menu
    async function getMenu(restaurantId) {
      try {
        const response = await fetch(
          `${API_BASE_URL}/api/v1/menu?restaurant_id=${restaurantId}`
        )
        const data = await response.json()
        console.log('Menu:', data)
        return data
      } catch (error) {
        console.error('Error:', error)
      }
    }

    // Create Order
    async function createOrder(orderData) {
      try {
        const response = await fetch(`${API_BASE_URL}/api/v1/orders`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(orderData)
        })
        const data = await response.json()
        console.log('Order created:', data)
        return data
      } catch (error) {
        console.error('Error:', error)
      }
    }

    // With JWT Token
    async function createTable(tableData, token) {
      try {
        const response = await fetch(`${API_BASE_URL}/api/v1/staff/tables`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`
          },
          body: JSON.stringify(tableData)
        })
        const data = await response.json()
        return data
      } catch (error) {
        console.error('Error:', error)
      }
    }
  </script>
</body>
</html>
```

---

## 🌍 Production / Deployment Scenarios

### 1. Frontend di Vercel/Netlify, Backend di Server/Cloud

```javascript
// Frontend .env.production
NEXT_PUBLIC_API_URL=https://api.yourposfnb.com

// Atau dengan subdomain
NEXT_PUBLIC_API_URL=https://backend.yourposfnb.com

// Atau dengan path
NEXT_PUBLIC_API_URL=https://yourposfnb.com/api
```

### 2. Same Domain dengan Reverse Proxy (Nginx)

```javascript
// Frontend di: https://yourposfnb.com
// Backend di: https://yourposfnb.com/api

const API_BASE_URL = '/api'  // Relative path

// Nginx config:
// location /api/ {
//   proxy_pass http://backend:8080/;
// }
```

### 3. Microservices dengan API Gateway

```javascript
// API Gateway di: https://api.yourposfnb.com
// Backend services di belakang gateway

const API_BASE_URL = 'https://api.yourposfnb.com'

// Gateway routing:
// /menu -> menu-service:8080
// /orders -> order-service:8080
// /staff -> staff-service:8080
```

---

## 🔧 Environment-based Configuration

### Best Practice untuk Multi-Environment

```typescript
// config/environment.ts
export const environments = {
  development: {
    apiUrl: 'http://localhost:8080',
    wsUrl: 'ws://localhost:8080',
  },
  staging: {
    apiUrl: 'https://staging-api.yourposfnb.com',
    wsUrl: 'wss://staging-api.yourposfnb.com',
  },
  production: {
    apiUrl: 'https://api.yourposfnb.com',
    wsUrl: 'wss://api.yourposfnb.com',
  },
}

const env = process.env.NODE_ENV || 'development'
export const config = environments[env as keyof typeof environments]

// Usage
import { config } from '@/config/environment'
const response = await fetch(`${config.apiUrl}/api/v1/menu`)
```

---

## 📱 Mobile App (React Native / Flutter)

### React Native

```typescript
// config/api.config.ts
import { Platform } from 'react-native'

const getBaseUrl = () => {
  if (__DEV__) {
    // Development
    if (Platform.OS === 'ios') {
      return 'http://localhost:8080'
    } else {
      // Android emulator
      return 'http://10.0.2.2:8080'
    }
  }
  // Production
  return 'https://api.yourposfnb.com'
}

export const API_BASE_URL = getBaseUrl()

// services/api.service.ts
import axios from 'axios'
import { API_BASE_URL } from '@/config/api.config'

const apiClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
})

export default apiClient
```

**Note untuk Android Emulator:**
- `localhost` tidak bisa diakses dari emulator
- Gunakan `10.0.2.2` untuk akses host machine
- Atau gunakan IP address actual computer (e.g., `192.168.1.100:8080`)

### Flutter

```dart
// lib/config/api_config.dart
class ApiConfig {
  static const String baseUrl = String.fromEnvironment(
    'API_BASE_URL',
    defaultValue: 'http://localhost:8080',
  );
}

// lib/services/api_service.dart
import 'package:http/http.dart' as http;
import 'dart:convert';
import '../config/api_config.dart';

class ApiService {
  static Future<Map<String, dynamic>> getMenu(String restaurantId) async {
    final response = await http.get(
      Uri.parse('${ApiConfig.baseUrl}/api/v1/menu?restaurant_id=$restaurantId'),
    );
    
    if (response.statusCode == 200) {
      return json.decode(response.body);
    } else {
      throw Exception('Failed to load menu');
    }
  }
}
```

---

## 🧪 Testing Endpoints

### cURL
```bash
# Health Check
curl http://localhost:8080/health

# Get Menu
curl "http://localhost:8080/api/v1/menu?restaurant_id=11111111-1111-1111-1111-111111111111"

# Create Order
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "session-123",
    "table_id": "table-uuid",
    "customer_name": "John Doe",
    "payment_method": "cashier",
    "items": [{"menu_item_id": "item-uuid", "quantity": 2}]
  }'

# With JWT Token (Staff endpoint)
curl -X POST http://localhost:8080/api/v1/staff/tables \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "restaurant_id": "restaurant-uuid",
    "table_number": "A1"
  }'
```

### Postman / Insomnia
```
Base URL: http://localhost:8080
```

Import collection dengan endpoints yang sudah ada di README.md

---

## ⚙️ Port Configuration

Jika port `8080` sudah digunakan, edit file `.env`:

```bash
APP_PORT=8081
```

Kemudian restart Docker:
```bash
docker compose down
docker compose up -d
```

URL menjadi: `http://localhost:8081`

---

## 🔒 CORS Configuration

Backend sudah include CORS middleware. Jika ada masalah CORS:

1. **Pastikan Origin diizinkan** - Check `internal/middleware/cors.go`
2. **Development:** Biasanya `*` diizinkan
3. **Production:** Set specific origins

```go
// internal/middleware/cors.go
AllowedOrigins: []string{
  "http://localhost:3000",  // Next.js
  "http://localhost:5173",  // Vite
  "https://yourfrontend.com",
}
```

---

## 📊 Summary Table

| Scenario | Backend URL | Frontend Uses |
|----------|-------------|---------------|
| **Docker (default)** | Container port 8080 | `http://localhost:8080` |
| **Local Go** | Host port 8080 | `http://localhost:8080` |
| **Custom Port** | Port in `.env` | `http://localhost:{APP_PORT}` |
| **Production** | Domain/IP | `https://api.yourdomain.com` |
| **Same Docker Network** | Container name | `http://api:8080` |
| **Android Emulator** | Host network | `http://10.0.2.2:8080` |

---

## 🎯 Quick Answer

**Untuk kebanyakan kasus development:**

```javascript
// Frontend (React/Vue/etc) di browser
const API_URL = "http://localhost:8080"

// Contoh penggunaan
fetch("http://localhost:8080/api/v1/menu?restaurant_id=xxx")
```

**Backend di Docker = Frontend akses via `localhost:8080`** ✅

---

## 📞 Need Help?

- Port issues: [DOCKER_CHEATSHEET.md](./DOCKER_CHEATSHEET.md) - Troubleshooting
- API endpoints: [README.md](./README.md) - Testing section
- CORS issues: Check middleware configuration

---

**Last Updated:** September 16, 2024
