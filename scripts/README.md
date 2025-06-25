# 🛠️ Domain Check Scripts

Bu dizin, Domain Check uygulamasını kolayca yönetmeniz için hazırlanmış shell scriptlerini içerir.

## 📋 Mevcut Scriptler

### 🚀 **Development & Çalıştırma**

#### `dev.sh` - Tam Geliştirme Ortamı
```bash
./scripts/dev.sh
```
- API ve Frontend'i birlikte başlatır
- Bağımlılıkları otomatik yükler
- Canlı monitoring ile çalışır
- Ctrl+C ile tüm servisleri durdurur

#### `start-api.sh` - Sadece API
```bash
./scripts/start-api.sh
```
- Sadece Go API sunucusunu başlatır (port 8080)

#### `start-frontend.sh` - Sadece Frontend
```bash
./scripts/start-frontend.sh
```
- Sadece Vue.js dev server'ı başlatır (port 3000)

### 🏗️ **Build & Dağıtım**

#### `build.sh` - Production Build
```bash
./scripts/build.sh
```
- Frontend'i production için build eder
- Go binary oluşturur (tüm platformlar için)
- Dağıtım paketi hazırlar
- `dist/` klasörüne her şeyi kopyalar

#### `build-frontend.sh` - Sadece Frontend Build
```bash
./scripts/build-frontend.sh
```
- Sadece Vue.js frontend'ini build eder

### 🛑 **Yönetim**

#### `stop.sh` - Servisleri Durdur
```bash
./scripts/stop.sh           # Sadece durdur
./scripts/stop.sh --clean   # Durdur + build dosyalarını temizle
```

#### `logs.sh` - Log ve Monitoring
```bash
./scripts/logs.sh           # Genel durum
./scripts/logs.sh api       # API logları
./scripts/logs.sh frontend  # Frontend logları
./scripts/logs.sh live      # Canlı monitoring
./scripts/logs.sh system    # Sistem bilgileri
./scripts/logs.sh ports     # Port kullanımı
./scripts/logs.sh processes # Çalışan süreçler
./scripts/logs.sh clear     # Log dosyalarını temizle
```

#### `test-api.sh` - API Test
```bash
./scripts/test-api.sh
```
- Tüm API endpoint'lerini test eder
- Performance testleri çalıştırır
- Detaylı sonuç raporu verir

## 🎯 **Kullanım Örnekleri**

### Geliştirme Başlatma
```bash
# Tam geliştirme ortamı
./scripts/dev.sh

# Veya ayrı terminallerde
./scripts/start-api.sh      # Terminal 1
./scripts/start-frontend.sh # Terminal 2
```

### Production Build
```bash
# Tam build
./scripts/build.sh

# Sadece frontend
./scripts/build-frontend.sh
```

### Test ve Monitoring
```bash
# API'yi test et
./scripts/test-api.sh

# Canlı monitoring
./scripts/logs.sh live

# Sistem durumu
./scripts/logs.sh system
```

### Temizlik
```bash
# Servisleri durdur ve temizle
./scripts/stop.sh --clean

# Log dosyalarını temizle
./scripts/logs.sh clear
```

## 🔧 **Gereksinimler**

### Geliştirme için:
- **Go** 1.19+
- **Node.js** 16+
- **npm** 8+

### Çalışma için:
- **curl** (test scriptleri için)
- **lsof** (port kontrolü için)
- **jq** (JSON formatting için, opsiyonel)

## 📊 **Script Özellikleri**

### ✅ **Özellikler**
- 🎨 Renkli output
- 🔍 Dependency kontrolü
- 🧹 Otomatik cleanup
- 📋 Detaylı logging
- ⚡ Performance monitoring
- 🛡️ Error handling
- 🚀 Cross-platform support

### 🎯 **Akıllı Özellikler**
- Port çakışmalarını otomatik çözer
- Bağımlılıkları kontrol eder
- Build durumunu doğrular
- Process ID'leri takip eder
- Graceful shutdown yapar

## 🆘 **Sorun Giderme**

### Port Çakışması
```bash
./scripts/stop.sh    # Tüm servisleri durdur
./scripts/logs.sh ports  # Port kullanımını kontrol et
```

### Build Hataları
```bash
./scripts/logs.sh system  # Sistem gereksinimlerini kontrol et
./scripts/build-frontend.sh  # Sadece frontend'i yeniden build et
```

### Performance Problemleri
```bash
./scripts/logs.sh processes  # Çalışan süreçleri kontrol et
./scripts/test-api.sh        # API performance'ını test et
```

## 💡 **İpuçları**

1. **jq yükleyin** - JSON output'u daha güzel görmek için:
   ```bash
   brew install jq  # macOS
   ```

2. **Alias oluşturun** - Kolay erişim için:
   ```bash
   alias dc-dev='./scripts/dev.sh'
   alias dc-build='./scripts/build.sh'
   alias dc-test='./scripts/test-api.sh'
   ```

3. **IDE Integration** - VS Code'da task olarak ekleyin

4. **CI/CD** - Bu scriptleri CI/CD pipeline'ınızda kullanabilirsiniz

## 🔗 **Linkler**

- **API**: http://localhost:8080
- **Frontend**: http://localhost:3000
- **Health Check**: http://localhost:8080/api/health
- **API Docs**: http://localhost:8080

---

**💬 Yardım gerekirse**: Tüm scriptler `--help` parametresini destekler veya parametre olmadan çalıştırıldığında usage bilgisini gösterir.
