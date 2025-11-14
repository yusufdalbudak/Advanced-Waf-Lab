# GitHub'a Yükleme Rehberi

## 🚀 Adım Adım GitHub'a Yükleme

### 1. GitHub'da Repository Oluştur

1. GitHub.com'a git: https://github.com
2. Sağ üstteki "+" butonuna tıkla
3. "New repository" seç
4. Repository adı: `WAF-DRAFT` (veya istediğin isim)
5. Description: "Production-ready Web Application Firewall (WAF) in Go"
6. Public veya Private seç
7. **"Initialize with README" seçme!** (zaten README var)
8. "Create repository" butonuna tıkla

### 2. Git Yapılandırması (İlk kez ise)

```bash
git config --global user.name "Your Name"
git config --global user.email "your.email@example.com"
```

### 3. Dosyaları Commit Et

```bash
cd /Users/yusufdalbudak/Documents/github/WAF-DRAFT

# Tüm dosyaları ekle
git add .

# Commit oluştur
git commit -m "Initial commit: Production-ready WAF with attack detection, logging, and dashboard"
```

### 4. GitHub Repository'ye Bağla ve Push Et

GitHub'da oluşturduğun repository'nin URL'sini kullan:

```bash
# Remote ekle (YOUR_USERNAME'i değiştir)
git remote add origin https://github.com/YOUR_USERNAME/WAF-DRAFT.git

# Veya SSH kullanıyorsan:
# git remote add origin git@github.com:YOUR_USERNAME/WAF-DRAFT.git

# Branch'i main olarak ayarla
git branch -M main

# GitHub'a push et
git push -u origin main
```

### 5. Alternatif: GitHub CLI ile (Daha Kolay)

Eğer GitHub CLI kuruluysa:

```bash
# GitHub CLI ile oturum aç
gh auth login

# Repository oluştur ve push et
gh repo create WAF-DRAFT --public --source=. --remote=origin --push
```

## 📝 Commit Mesajı Örnekleri

```bash
# İlk commit
git commit -m "Initial commit: Production-ready WAF MVP"

# Feature eklerken
git commit -m "feat: Add attack logging and dashboard"

# Bug fix
git commit -m "fix: Fix path traversal detection"

# Documentation
git commit -m "docs: Add professional roadmap"
```

## 🔐 GitHub Token (Gerekirse)

Eğer push sırasında authentication hatası alırsan:

1. GitHub Settings > Developer settings > Personal access tokens
2. "Generate new token" (classic)
3. `repo` scope'u seç
4. Token'ı kopyala
5. Push yaparken password yerine token kullan

## ✅ Kontrol

Push'tan sonra GitHub'da repository'ni kontrol et:
- Tüm dosyalar görünüyor mu?
- README.md düzgün görünüyor mu?
- .gitignore çalışıyor mu?

## 🎯 Sonraki Adımlar

1. **README.md'yi güncelle** - GitHub repository linkini ekle
2. **Topics ekle** - GitHub'da repository'ye topics ekle: `waf`, `security`, `golang`, `web-application-firewall`
3. **License ekle** - LICENSE dosyası ekle (MIT, Apache, vs.)
4. **Badges ekle** - CI/CD, coverage, license badge'leri

## 📚 Faydalı Komutlar

```bash
# Değişiklikleri kontrol et
git status

# Commit geçmişi
git log --oneline

# Remote repository'yi kontrol et
git remote -v

# Son commit'i değiştir
git commit --amend -m "New message"

# Belirli dosyaları commit et
git add file1.go file2.go
git commit -m "Update specific files"
```

## 🚨 Önemli Notlar

- **.gitignore** dosyası zaten var - hassas bilgiler commit edilmeyecek
- **waf.log** dosyası ignore edilecek
- **Binary dosyalar** ignore edilecek
- **IDE dosyaları** ignore edilecek

## 💡 İpuçları

1. İlk push'tan önce `make clean` çalıştır (binary'leri temizle)
2. README.md'yi kontrol et - güzel görünüyor mu?
3. LICENSE dosyası ekle
4. GitHub Actions CI/CD otomatik çalışacak

