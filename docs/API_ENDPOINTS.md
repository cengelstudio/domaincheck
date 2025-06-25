# 🚀 Domain Check API Endpoints

Bu dokümantasyon, Domain Check API'nin tüm endpoint'lerini, parametrelerini ve örnek kullanımlarını içerir.

## 📋 İçindekiler

- [Base URL](#base-url)
- [Authentication](#authentication)
- [Response Format](#response-format)
- [Health Check](#health-check)
- [Domain Operations](#domain-operations)
- [Extensions Management](#extensions-management)
- [Error Handling](#error-handling)
- [Rate Limiting](#rate-limiting)

## 🌐 Base URL

```
http://localhost:8080
```

## 🔐 Authentication

Bu API şu anda authentication gerektirmez. Tüm endpoint'ler herkese açıktır.

## 📝 Response Format

Tüm API yanıtları standart JSON formatında döner:

```json
{
  "success": true,
  "data": {},
  "message": "İşlem başarılı",
  "error": null,
  "meta": {
    "total": 0,
    "page": 1,
    "per_page": 20,
    "total_pages": 0,
    "request_id": "12345",
    "process_time_ms": 150
  }
}
```

---

## 💚 Health Check

### GET `/api/health`

API'nin sağlık durumunu kontrol eder.

#### Request
```http
GET /api/health
```

#### Response
```json
{
  "success": true,
  "data": {
    "status": "healthy",
    "version": "1.0.0",
    "timestamp": "2023-12-01T10:30:00Z",
    "uptime": "2h15m30s",
    "environment": "development"
  },
  "message": "Service is healthy"
}
```

### GET `/api/v1/health`

V1 API sağlık kontrolü (yukarıyla aynı).

---

## 🌍 Domain Operations

### 🎯 POST `/api/check-all-extensions` - **MAIN FEATURE**

**En önemli endpoint!** Bir domain adını tüm mevcut uzantılarla kontrol eder (228+ uzantı).

#### Request
```http
POST /api/check-all-extensions
Content-Type: application/json

{
  "domain_name": "metehansaral"
}
```

#### Request Parameters
| Parameter   | Type   | Required | Description                                    |
|-------------|--------|----------|------------------------------------------------|
| domain_name | string | Yes      | Kontrol edilecek domain adı (uzantısız, örn: "metehansaral") |

#### Response
```json
{
  "success": true,
  "data": {
    "domain_name": "metehansaral",
    "total_extensions": 228,
    "available_count": 195,
    "unavailable_count": 30,
    "error_count": 3,
    "checked_at": "2023-12-01T10:30:00Z",
    "total_time_ms": 4500,
    "available_domains": [
      {
        "domain": {
          "id": 1,
          "name": "metehansaral.com",
          "extension": ".com",
          "available": true,
          "dns_resolved": false,
          "checked_at": "2023-12-01T10:30:00Z",
          "response_time_ms": 120
        },
        "is_valid_tld": true,
        "supported_tld": true
      }
    ],
    "unavailable_domains": [...],
    "error_domains": [...],
    "all_results": [...],
    "summary": {
      "popular_available": ["metehansaral.com", "metehansaral.net", "metehansaral.org"],
      "recommended_domains": ["metehansaral.com", "metehansaral.org", "metehansaral.co"],
      "alternative_suggestions": ["metehansaralapp.com", "metehansaralpro.com"],
      "fastest_response": {
        "id": 5,
        "name": "metehansaral.net",
        "extension": ".net",
        "available": true,
        "response_time_ms": 45
      },
      "slowest_response": {
        "id": 6,
        "name": "metehansaral.museum",
        "extension": ".museum",
        "available": true,
        "response_time_ms": 1250
      }
    }
  },
  "message": "Domain extensions check completed successfully",
  "meta": {
    "total": 228,
    "process_time_ms": 4500,
    "request_id": "req_12345"
  }
}
```

#### Usage Examples

**cURL:**
```bash
# Check all extensions for "metehansaral"
curl -X POST http://localhost:8080/api/check-all-extensions \
  -H "Content-Type: application/json" \
  -d '{"domain_name": "metehansaral"}'
```

**JavaScript (Fetch):**
```javascript
const checkAllExtensions = async (domainName) => {
  const response = await fetch('/api/check-all-extensions', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ domain_name: domainName })
  });

  const result = await response.json();

  console.log(`Domain: ${result.data.domain_name}`);
  console.log(`Available: ${result.data.available_count}/${result.data.total_extensions}`);
  console.log(`Recommended:`, result.data.summary.recommended_domains);

  return result;
};

// Usage
checkAllExtensions('metehansaral');
```

**JavaScript (Axios):**
```javascript
const checkAllExtensions = async (domainName) => {
  try {
    const response = await axios.post('/api/check-all-extensions', {
      domain_name: domainName
    });

    const data = response.data.data;
    console.log(`✅ ${data.available_count} domains available out of ${data.total_extensions}`);
    console.log(`🎯 Recommended:`, data.summary.recommended_domains);

    return data;
  } catch (error) {
    console.error('❌ Error:', error.response?.data?.message);
  }
};
```

### 🆔 POST `/api/v1/domains/check-all-extensions`

V1 API versiyonu - yukarıyla aynı fonksiyonalite.

---

### POST `/api/check-domain` (Legacy)

Tek bir domain'in availability durumunu kontrol eder.

#### Request
```http
POST /api/check-domain
Content-Type: application/json

{
  "domain": "google.com"
}
```

#### Request Parameters
| Parameter | Type   | Required | Description          |
|-----------|--------|----------|----------------------|
| domain    | string | Yes      | Kontrol edilecek domain |

#### Response
```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "google.com",
    "extension": ".com",
    "available": false,
    "ip": "142.250.191.14",
    "dns_resolved": true,
    "checked_at": "2023-12-01T10:30:00Z",
    "response_time_ms": 245,
    "error": null
  },
  "message": "Domain check completed successfully",
  "meta": {
    "process_time_ms": 250,
    "request_id": "req_12345"
  }
}
```

### POST `/api/v1/domains/check`

V1 API - Tek domain kontrolü (yukarıyla aynı format).

### POST `/api/v1/domains/check-multiple`

Birden fazla domain'i aynı anda kontrol eder.

#### Request
```http
POST /api/v1/domains/check-multiple
Content-Type: application/json

{
  "domains": ["google.com", "github.com", "nonexistent-domain-12345.com"]
}
```

#### Request Parameters
| Parameter | Type     | Required | Description                    |
|-----------|----------|----------|--------------------------------|
| domains   | string[] | Yes      | Kontrol edilecek domain listesi (max 50) |

#### Response
```json
{
  "success": true,
  "data": [
    {
      "domain": {
        "id": 2,
        "name": "google.com",
        "extension": ".com",
        "available": false,
        "ip": "142.250.191.14",
        "dns_resolved": true,
        "checked_at": "2023-12-01T10:30:00Z",
        "response_time_ms": 120
      },
      "is_valid_tld": true,
      "supported_tld": true
    },
    {
      "domain": {
        "id": 3,
        "name": "github.com",
        "extension": ".com",
        "available": false,
        "ip": "140.82.112.3",
        "dns_resolved": true,
        "checked_at": "2023-12-01T10:30:00Z",
        "response_time_ms": 95
      },
      "is_valid_tld": true,
      "supported_tld": true
    },
    {
      "domain": {
        "id": 4,
        "name": "nonexistent-domain-12345.com",
        "extension": ".com",
        "available": true,
        "dns_resolved": false,
        "checked_at": "2023-12-01T10:30:00Z",
        "response_time_ms": 2500,
        "error": "no such host"
      },
      "is_valid_tld": true,
      "supported_tld": true
    }
  ],
  "message": "Domain checks completed successfully",
  "meta": {
    "total": 3,
    "process_time_ms": 2800,
    "request_id": "req_12346"
  }
}
```

---

## 📊 Domain History

### GET `/api/domains` (Legacy)

Domain kontrol geçmişini getirir.

#### Request
```http
GET /api/domains
```

#### Response
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "google.com",
      "extension": ".com",
      "available": false,
      "ip": "142.250.191.14",
      "dns_resolved": true,
      "checked_at": "2023-12-01T10:30:00Z",
      "response_time_ms": 245
    }
  ],
  "message": "Domains retrieved successfully"
}
```

### GET `/api/v1/domains/history`

V1 API - Sayfalama destekli domain geçmişi.

#### Request
```http
GET /api/v1/domains/history?page=1&per_page=10
```

#### Query Parameters
| Parameter | Type | Default | Description                    |
|-----------|------|---------|--------------------------------|
| page      | int  | 1       | Sayfa numarası                 |
| per_page  | int  | 20      | Sayfa başına kayıt (max 100)  |

#### Response
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "google.com",
      "extension": ".com",
      "available": false,
      "ip": "142.250.191.14",
      "dns_resolved": true,
      "checked_at": "2023-12-01T10:30:00Z",
      "response_time_ms": 245
    }
  ],
  "message": "Domain history retrieved successfully",
  "meta": {
    "total": 50,
    "page": 1,
    "per_page": 10,
    "total_pages": 5,
    "request_id": "req_12347"
  }
}
```

### DELETE `/api/v1/domains/history`

Domain kontrol geçmişini temizler.

#### Request
```http
DELETE /api/v1/domains/history
```

#### Response
```json
{
  "success": true,
  "message": "Domain history cleared successfully"
}
```

---

## 🔧 Extensions Management

### GET `/api/v1/extensions`

Desteklenen domain uzantıları listesini getirir.

#### Request
```http
GET /api/v1/extensions
```

#### Response
```json
{
  "success": true,
  "data": [
    ".com",
    ".net",
    ".org",
    ".edu",
    ".gov",
    ".tr",
    ".uk",
    ".de"
  ],
  "message": "Valid extensions retrieved successfully",
  "meta": {
    "total": 200,
    "request_id": "req_12348"
  }
}
```

### POST `/api/v1/extensions/reload`

Domain uzantıları dosyasını yeniden yükler.

#### Request
```http
POST /api/v1/extensions/reload
```

#### Response
```json
{
  "success": true,
  "message": "Extensions reloaded successfully"
}
```

---

## ❌ Error Handling

### Error Response Format

Hata durumlarında API aşağıdaki formatı kullanır:

```json
{
  "success": false,
  "data": null,
  "message": "Domain check failed",
  "error": "invalid domain format: invalid..domain"
}
```

### HTTP Status Codes

| Status Code | Description                    |
|-------------|--------------------------------|
| 200         | Success                        |
| 400         | Bad Request (validation error) |
| 404         | Not Found                      |
| 500         | Internal Server Error          |

### Common Error Types

#### 400 - Validation Errors
```json
{
  "success": false,
  "message": "Invalid request format",
  "error": "Domain parameter is required"
}
```

#### 400 - Invalid Domain Format
```json
{
  "success": false,
  "message": "Domain check failed",
  "error": "invalid domain format: invalid..domain"
}
```

#### 400 - Empty Domain
```json
{
  "success": false,
  "message": "Domain check failed",
  "error": "domain must have an extension"
}
```

#### 500 - Internal Server Error
```json
{
  "success": false,
  "message": "Failed to reload extensions",
  "error": "failed to open extensions file: no such file or directory"
}
```

---

## 🚦 Rate Limiting

Şu anda rate limiting implementasyonu yoktur, ancak gelecekte eklenebilir.

### Önerilen Limitler
- Domain check: 100 istek/dakika
- Multiple domain check: 10 istek/dakika (max 50 domain per request)
- History ve extensions: 1000 istek/dakika

---

## 📚 Usage Examples

### cURL Examples

#### Single Domain Check
```bash
curl -X POST http://localhost:8080/api/check-domain \
  -H "Content-Type: application/json" \
  -d '{"domain": "google.com"}'
```

#### Multiple Domain Check
```bash
curl -X POST http://localhost:8080/api/v1/domains/check-multiple \
  -H "Content-Type: application/json" \
  -d '{"domains": ["google.com", "github.com", "stackoverflow.com"]}'
```

#### Get History with Pagination
```bash
curl "http://localhost:8080/api/v1/domains/history?page=2&per_page=5"
```

#### Get Valid Extensions
```bash
curl http://localhost:8080/api/v1/extensions
```

#### Health Check
```bash
curl http://localhost:8080/api/health
```

### JavaScript Examples

#### Fetch API
```javascript
// Single domain check
const checkDomain = async (domain) => {
  const response = await fetch('http://localhost:8080/api/check-domain', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ domain })
  });

  return await response.json();
};

// Multiple domain check
const checkMultipleDomains = async (domains) => {
  const response = await fetch('http://localhost:8080/api/v1/domains/check-multiple', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ domains })
  });

  return await response.json();
};

// Get history
const getHistory = async (page = 1, perPage = 20) => {
  const response = await fetch(
    `http://localhost:8080/api/v1/domains/history?page=${page}&per_page=${perPage}`
  );

  return await response.json();
};
```

#### Axios Examples
```javascript
import axios from 'axios';

const api = axios.create({
  baseURL: 'http://localhost:8080/api',
  headers: {
    'Content-Type': 'application/json'
  }
});

// Single domain check
const checkDomain = (domain) => api.post('/check-domain', { domain });

// Multiple domain check
const checkMultipleDomains = (domains) =>
  api.post('/v1/domains/check-multiple', { domains });

// Get history
const getHistory = (page = 1, perPage = 20) =>
  api.get(`/v1/domains/history?page=${page}&per_page=${perPage}`);

// Get extensions
const getExtensions = () => api.get('/v1/extensions');

// Health check
const healthCheck = () => api.get('/health');
```

---

## 🔍 Testing

API'yi test etmek için proje içindeki test script'ini kullanabilirsiniz:

```bash
./scripts/test-api.sh
```

Bu script tüm endpoint'leri test eder ve detaylı rapor verir.

---

## 📝 Notes

### Domain Availability Logic

- **Available (true)**: Domain DNS çözümlenemedi (muhtemelen available)
- **Available (false)**: Domain DNS çözümlendi (registered/active)

### Performance

- Single domain check: ~100-500ms
- Multiple domain check: Concurrent olarak çalışır (max 10 concurrent)
- Timeout: 5 saniye (configurable)

### Supported Formats

Domain girişleri için desteklenen formatlar:
- `google.com`
- `www.google.com` (www. prefix'i otomatik kaldırılır)
- `http://google.com` (protocol prefix'i otomatik kaldırılır)
- `https://www.google.com` (tüm prefix'ler kaldırılır)

### Extensions File

Domain uzantıları `./data/domain_extensions.txt` dosyasından okunur. Bu dosyayı düzenleyebilir ve `/api/v1/extensions/reload` endpoint'ini kullanarak yeniden yükleyebilirsiniz.

---

## 🔗 Related Documentation

- [Project README](../README.md)
- [Scripts Documentation](../scripts/README.md)
- [Configuration Guide](../configs/config.yaml)

---

**📞 Support**: API ile ilgili sorularınız için GitHub Issues kullanabilirsiniz.
