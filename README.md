# Domain Check

Go ile yazılmış basit bir domain kontrol API'si ve Vue.js frontend arayüzü.

## Özellikler

- Domain availability kontrolü
- DNS lookup ile IP adresi öğrenme
- Modern Vue.js arayüzü
- RESTful API
- CORS desteği
- Geçmiş kontrolleri görüntüleme

## 🚀 Hızlı Başlangıç

### ⚡ **Tek Komutla Başlat** (Önerilen)
```bash
./scripts/dev.sh
```
Bu komut otomatik olarak:
- Bağımlılıkları yükler
- API sunucusunu başlatır (port 8080)
- Frontend dev server'ı başlatır (port 3000)
- Monitoring sağlar

### 🔧 **Manuel Kurulum**

#### Backend (Go API)
```bash
go mod tidy
./scripts/start-api.sh
```

#### Frontend (Vue.js)
```bash
./scripts/start-frontend.sh
```

#### Production Build
```bash
./scripts/build.sh
```

### 📋 **Tüm Scriptler**
- `./scripts/dev.sh` - Tam geliştirme ortamı
- `./scripts/build.sh` - Production build
- `./scripts/test-api.sh` - API testleri
- `./scripts/stop.sh` - Servisleri durdur
- `./scripts/logs.sh` - Log ve monitoring

Detaylar için: [scripts/README.md](scripts/README.md)

## API Endpoints

- `GET /api/health` - API sağlık kontrolü
- `POST /api/check-domain` - Domain kontrolü
- `GET /api/domains` - Kontrol geçmişi

## Kullanım

1. Backend ve frontend sunucularını başlat
2. `http://localhost:3000` adresine git
3. Domain adı gir ve kontrol et
4. Sonuçları geçmişte görüntüle

## Teknolojiler

- **Backend**: Go, Gin framework
- **Frontend**: Vue.js 3, Axios, Tailwind CSS
- **CORS**: Frontend-backend iletişimi için
