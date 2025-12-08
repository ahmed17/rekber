#!/bin/bash
echo "🔍 Verifikasi struktur project..."

echo "1. Cek file repository.go..."
if [ -f "internal/repositories/repository.go" ]; then
    echo "✅ internal/repositories/repository.go ditemukan"
else
    echo "❌ internal/repositories/repository.go tidak ditemukan"
    exit 1
fi

echo "2. Cek isi repository.go..."
grep -q "InitRepository" internal/repositories/repository.go && echo "✅ InitRepository ditemukan" || echo "❌ InitRepository tidak ditemukan"
grep -q "GetRepository" internal/repositories/repository.go && echo "✅ GetRepository ditemukan" || echo "❌ GetRepository tidak ditemukan"

echo "3. Build project..."
if go build ./...; then
    echo "✅ Build berhasil!"
    
    echo "4. Test server..."
    timeout 5s go run cmd/api/main.go &
    SERVER_PID=$!
    sleep 2
    
    if curl -s http://localhost:8080/health | grep -q "healthy"; then
        echo "✅ Server berjalan dengan baik"
        kill $SERVER_PID 2>/dev/null
        wait $SERVER_PID 2>/dev/null
        echo "🎉 Semua test berhasil!"
    else
        echo "❌ Server tidak merespon"
        kill $SERVER_PID 2>/dev/null
        exit 1
    fi
else
    echo "❌ Build gagal"
    exit 1
fi