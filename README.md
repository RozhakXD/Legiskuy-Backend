# LegisKuy - Backend API Pemilu

**LegisKuy** adalah backend RESTful API untuk sistem pemilu digital yang dibangun menggunakan **Go (Golang)**, Fiber framework, dan database **SQLite**. Proyek ini merupakan pengembangan lanjutan dari tugas akademik yang kemudian dikembangkan secara mandiri sebagai bagian dari portofolio profesional. Dirancang dengan pendekatan arsitektur multilayer yang bersih, aman, dan efisien.

## ✨ Fitur Utama

- **Otentikasi & Otorisasi Berbasis Peran:**
  - Sistem login yang aman menggunakan **JWT (JSON Web Token)**.
  - Pemisahan hak akses yang jelas antara `petugas` (manajemen penuh) dan pemilih (hanya bisa memilih dan melihat data).
- **Manajemen Data (CRUD):**
  - Pengelolaan data **Calon Legislatif** (tambah, lihat, ubah, hapus).
  - Pengelolaan data **Pemilih** (registrasi, lihat, ubah, hapus).
- **Proses Pemilu yang Aman:**
  - Endpoint khusus untuk melakukan voting (`POST /api/v1/votes`)
  - Validasi untuk memastikan setiap pemilih hanya bisa memberikan suara satu kali.
  - Penggunaan **transaksi database** untuk menjamin integritas data saat proses pemilihan.
- **Pencarian & Pengurutan Data:**
  - Pencarian calon berdasarkan nama atau partai.
  - Pengurutan data calon berdasarkan nama (menggunakan _Selection Sort_), partai, dan jumlah suara (menggunakan _Insertion Sort_) secara `ascending` maupun `descending`.
- **Pengaturan Dinamis Pemilu:**
  - Petugas dapat mengatur **jadwal (waktu mulai & selesai)** pemilu. Proses _voting_ hanya bisa dilakukan dalam rentang waktu tersebut.
  - Petugas dapat menentukan **ambang batas (threshold)** suara minimum bagi calon untuk dapat terpilih.

## 🏛️ Arsitektur & Teknologi

Proyek ini dibangun di atas arsitektur _N-Tier_ untuk memastikan kode yang _modular_, _scalable_, dan _maintainable_.

- **Presentation Layer (Handler):** Menangani request & response HTTP.
- **Business Logic Layer (Service):** Berisi semua logika bisnis dan validasi.
- **Data Access Layer (Repository):** Bertanggung jawab untuk interaksi dengan database.

### Tumpukan Teknologi (Tech Stack):

- **Bahasa:** Go (Golang)
- **Framework:** Fiber
- **Database:** SQLite 3
- **Keamanan:** JWT (JSON Web Token) & Bcrypt
- **Dokumentasi API:** Swagger (OpenAPI)

## 📚 Dokumentasi API

Dokumentasi API yang lengkap dan interaktif tersedia melalui Swagger. Setelah menjalankan server, akses URL berikut:

[http://localhost:3000/swagger/index.html](http://localhost:3000/swagger/index.html)

## 🚀 Instalasi & Menjalankan Proyek

1. **Prasyarat**
   
   - Pastikan Go versi 1.22 atau lebih baru sudah terinstal.
   - Pastikan Git sudah terinstal.

2. **Clone Repositori**
   
   ```bash
   git clone https://github.com/RozhakXD/legiskuy-backend.git
   cd legiskuy-backend
   ```

3. **Instalasi Dependensi**
   
   ```bash
   go mod tidy
   ```

4. **Konfigurasi Lingkungan**
   
   Proyek ini menggunakan variabel lingkungan untuk JWT Secret. Anda bisa membuat file .env di root direktori (opsional).
   
   ```sh
      JWT_SECRET=kunci_rahasia_yang_sangat_aman
   ```
   
   Jika file `.env` tidak ditemukan, aplikasi akan menggunakan `fallback` secret bawaan.

5. **Jalankan Aplikasi**
   
   ```bash
   go run cmd/api/main.go
   ```
   
    Server akan berjalan di `http://localhost:3000`.

## 📂 Struktur Proyek

```text
/legiskuy-backend
|-- /cmd/api/main.go        # Titik masuk aplikasi & registrasi rute
|-- /docs                   # File dokumentasi Swagger
|-- /internal               # Logika inti aplikasi
|   |-- /auth               # Modul otentikasi & otorisasi
|   |-- /candidate          # Modul manajemen calon
|   |-- /election           # Modul proses pemilu
|   |-- /voter              # Modul manajemen pemilih
|-- /pkg                    # Paket pendukung
|   |-- /database           # Koneksi & inisialisasi DB
|   |-- /middleware         # Middleware untuk otentikasi
|-- go.mod
|-- go.sum
|-- README.md
```

## 👤 Kontributor

LegisKuy bukan sekadar proyek backend — ini karya dari malam-malam panjang, kopi yang gak habis-habis, dan rasa yang gak bisa disampaikan langsung.

- 💻 **Rozhak** — Backend Developer, System Architect  
  _"Didedikasikan diam-diam untuk seseorang yang pernah jadi semangat paling tulus... walau hanya bisa mencintai dari jauh, setidaknya aku bisa mengubah rasa itu jadi baris-baris kode."_  

Kalau kamu baca ini, terima kasih… karena tanpamu, mungkin proyek ini nggak pernah selesai.

**Proyek ini dilisensikan di bawah `MIT License`.**
