#!/bin/bash
# GitHub'a Push Etme Scripti

echo "🚀 GitHub'a Yükleme Scripti"
echo "============================"
echo ""

# GitHub kullanıcı adını sor
read -p "GitHub kullanıcı adınızı girin: " GITHUB_USER

if [ -z "$GITHUB_USER" ]; then
    echo "❌ Kullanıcı adı boş olamaz!"
    exit 1
fi

REPO_NAME="WAF-DRAFT"
REPO_URL="https://github.com/${GITHUB_USER}/${REPO_NAME}.git"

echo ""
echo "📋 Repository bilgileri:"
echo "   Kullanıcı: $GITHUB_USER"
echo "   Repository: $REPO_NAME"
echo "   URL: $REPO_URL"
echo ""

read -p "Devam etmek istiyor musunuz? (y/n): " CONFIRM

if [ "$CONFIRM" != "y" ] && [ "$CONFIRM" != "Y" ]; then
    echo "❌ İptal edildi"
    exit 0
fi

echo ""
echo "🔄 Remote repository ekleniyor..."
git remote add origin $REPO_URL 2>/dev/null || git remote set-url origin $REPO_URL

echo "🔄 Branch main olarak ayarlanıyor..."
git branch -M main

echo ""
echo "⚠️  ÖNEMLİ: Önce GitHub'da repository oluşturmanız gerekiyor!"
echo ""
echo "1. https://github.com/new adresine git"
echo "2. Repository adı: $REPO_NAME"
echo "3. Public veya Private seç"
echo "4. 'Initialize with README' seçme!"
echo "5. 'Create repository' tıkla"
echo ""
read -p "Repository'yi oluşturdunuz mu? (y/n): " REPO_CREATED

if [ "$REPO_CREATED" != "y" ] && [ "$REPO_CREATED" != "Y" ]; then
    echo "❌ Önce repository'yi oluşturun, sonra tekrar çalıştırın"
    exit 1
fi

echo ""
echo "🚀 GitHub'a push ediliyor..."
git push -u origin main

if [ $? -eq 0 ]; then
    echo ""
    echo "✅ Başarılı! Repository GitHub'da:"
    echo "   https://github.com/${GITHUB_USER}/${REPO_NAME}"
else
    echo ""
    echo "❌ Push başarısız oldu. Hata mesajını kontrol edin."
    echo "   Eğer authentication hatası alıyorsanız:"
    echo "   - GitHub token kullanmanız gerekebilir"
    echo "   - veya SSH key kullanın"
fi

