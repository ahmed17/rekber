Website Rekening Bersama (Rekber) dengan Go
1. Topik Utama Pembahasan
Membangun website Rekening Bersama (Rekber) menggunakan Go (Golang) dengan arsitektur berlapis (layered architecture). Website ini mendukung transaksi manual antara pembeli dan penjual dengan alur yang telah ditentukan.

2. Poin Kunci yang Sudah Dibahas
Step 1a: Project Initialization & Structure Setup

Inisialisasi project Go dengan struktur direktori yang terorganisir.

Setup konfigurasi dasar menggunakan environment variables.

Pembuatan Makefile untuk otomatisasi perintah.

Server dasar dengan Gin framework dan endpoint health check.

Step 1b: Database Setup & Migration

Setup koneksi database PostgreSQL dengan GORM.

Pembuatan migration SQL untuk tabel: users, transactions, testimonies.

Implementasi migration tool dengan golang-migrate.

Fix error SSL dengan menambahkan sslmode=disable untuk development.

Step 1c: Model Definitions (GORM Models)

Pembuatan model Go (struct) untuk: User, Transaction, Testimony.

Definisi tipe data dan relasi antar model.

Auto migration untuk development (namun di-nonaktifkan karena sudah menggunakan migrasi SQL).

Perbaikan error import dan struktur project.

Step 1d: Repository Layer

Pembuatan repository pattern untuk abstraksi akses data.

Interface dan implementasi untuk: UserRepository, TransactionRepository, TestimonyRepository.

Container repository untuk mengelola instance repository.

Testing repository dengan data dummy.

Step 1e: Service Layer

Pembuatan service layer untuk business logic.

Implementasi service: AuthService, UserService, TransactionService, TestimonyService.

Validasi, aturan bisnis, dan flow transaksi (status management, auto-complete dalam 2x24 jam).

JWT authentication dengan token generation dan validation.

Perbaikan error pada service layer (undefined types, missing methods, nil pointer).

3. Data Penting yang Telah Disebutkan
Teknologi Stack:

Backend: Go (Golang) dengan framework Gin.

Database: PostgreSQL (port 5433, nama database: db_rekber).

ORM: GORM.

Migration: golang-migrate.

Authentication: JWT (JSON Web Tokens).

Password hashing: bcrypt.

Alur Transaksi (Status Flow):

Waiting Payment → 2. Pending (setelah upload bukti bayar) → 3. Validated (admin memvalidasi) → 4. Processing → 5. Shipped (penjual mengirim data) → 6. Completed (otomatis setelah 2x24 jam atau manual oleh pembeli).

Data Sample (Seed):

User: admin, seller1, buyer1.

Transaction: REKBER-001, REKBER-002, REKBER-003 dengan berbagai status.

Testimony: untuk transaksi yang sudah selesai.

4. Step yang Masih Belum Terjawab / Belum Dikerjakan
Step 1f: Handler Layer (Controllers)

Pembuatan handler (controller) untuk menangani HTTP request dan response.

Routing untuk API endpoints (RESTful API).

Middleware untuk authentication, authorization, logging, dll.

Step 1g: Middleware & Authentication

Implementasi middleware untuk validasi JWT.

Middleware untuk role-based access control (admin, seller, buyer).

Step 1h: Frontend Integration (Go Templates)

Pembuatan template HTML dengan Go template engine.

Integrasi dengan CSS framework (Tailwind CSS) dan JavaScript (HTMX/Alpine.js).

Halaman-halaman untuk user interface.

Step 1i: Scheduler untuk Auto-Complete

Implementasi scheduler (cron job) untuk menangani auto-complete transaksi setelah 2x24 jam.

Step 1j: File Upload & Static File Serving

Handler untuk upload bukti pembayaran dan file lainnya.

Static file server untuk menyajikan file yang diupload.

Step 1k: Testing & Validation

Unit test dan integration test untuk layer-layer yang sudah dibuat.

Validasi input dengan validator.

Step 1l: Deployment Preparation

Dockerfile dan docker-compose untuk containerization.

Konfigurasi untuk production environment.

5. Arah Pembahasan Selanjutnya
Langkah selanjutnya adalah Step 1f: Handler Layer (Controllers), yang akan meliputi:

Pembuatan handler untuk setiap service (auth, user, transaction, testimony).

Definisi route (endpoint) API RESTful.

Binding dan validasi request.

Response formatting (JSON).

Error handling di level HTTP.

Setelah handler layer selesai, kita akan melanjutkan dengan middleware, frontend integration, dan fitur-fitur lainnya sesuai alur yang telah direncanakan.

Catatan: Semua kode yang dihasilkan sejauh ini dapat dijalankan dengan perintah make run dan sudah teruji (dengan beberapa penyesuaian untuk environment masing-masing). Pastikan database PostgreSQL berjalan di port 5433 dengan nama database db_rekber sebelum menjalankan aplikasi.

Sesi selanjutnya akan fokus pada implementasi handler layer yang menghubungkan service layer dengan HTTP requests.
