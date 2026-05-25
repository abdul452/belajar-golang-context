# Belajar Golang Context & Request Lifecycle Control 🚀

Repositori ini berisi kumpulan catatan eksperimen, analisis kode, dan implementasi fungsional mengenai package `context` di bahasa Go. Fokus utama dari modul ini adalah memahami bagaimana mengelola siklus hidup permintaan (*request lifecycle*), menghindari kebocoran memori pada Goroutine (*Goroutine leak*), serta mendistribusikan data lintas batas fungsi secara aman dan efisien.

---

## 📌 Modul Latihan
* `01-context_test`: Berisi implementasi dasar, pewarisan nilai (*Context Value*), mekanisme pembatalan manual (*WithCancel*), dan pembatalan otomatis berbasis durasi (*WithTimeout*).

---

## 🔍 Analisis Mendalam & Filosofi Context di Go

Berbeda dengan `Context` pada bahasa lain seperti Java/Android yang bertindak sebagai objek berat untuk mengakses *resource* sistem operasi, `context.Context` di Go backend adalah sebuah *interface* super ringan yang berfungsi sebagai **Surat Perintah Jalan & Batas Waktu** yang dioper sebagai parameter pertama di setiap fungsi.

### 1. Struktur Hirarki Pewarisan Nilai (`context.WithValue`)
Fungsi `context.WithValue` bersifat *Immutable* (tidak dapat diubah). Setiap kali kita menambahkan data *key-value*, Go akan membuat objek context anak (*child*) baru yang membentuk struktur data **Pohon Keluarga (Context Tree)**.

#### Prinsip Pencarian Data (Bottom-Up)
* **Anak bisa melihat data Bapak:** Ketika memanggil `.Value("key")`, Go akan mencari di context dirinya sendiri. Jika tidak ditemukan, ia akan naik mencari ke bapaknya, kakeknya, hingga mentok di akar (`context.Background()`).
* **Bapak TIDAK BISA melihat data Anak:** Context akar atau *parent* tidak memiliki akses ke data yang dititipkan pada context *child*.
* *Real Case:* Digunakan untuk membawa data global tipis per request (*Request-Scoped Data*) seperti **Trace ID / Request ID** untuk kebutuhan *logging logging tracking*, atau data klaim token JWT (`user_id`).

---

### 2. Penanggulangan Kebocoran Memori (`context.WithCancel`)
Saat memicu Goroutine asinkron yang berjalan dalam perulangan tanpa henti (*infinite loop*), menghentikan proses pembacaan di program utama menggunakan keyword `break` saja tidak cukup. Goroutine di latar belakang akan terkunci selamanya (*blocking*) saat mencoba mengirim data ke channel yang sudah tidak ada pembacanya. Fenomena ini disebut **Goroutine Leak**.

#### Solusi Sinyal Pembatalan
Dengan memanfaatkan `context.WithCancel`, kita memantau channel `<-ctx.Done()` di dalam blok `select case` Goroutine. Begitu fungsi `cancel()` dipanggil dari program utama, channel `Done()` akan tertutup, memicu Goroutine untuk melakukan `return` secara terhormat dan membersihkan memorinya.

```go
select {
case <-ctx.Done():
    return // 🎯 Membunuh Goroutine secara aman dari kebocoran memori
default:
    destination <- counter
    counter++
}
```
### Sistem Pertahanan API Otomatis (context.WithTimeout)
`context.WithTimeout` berfungsi seperti alarm bom waktu otomatis. Fitur ini sangat krusial di industri backend untuk menangani operasi I/O Bound (seperti query database atau menembak API payment gateway pihak ketiga).

## ⚠️ Aturan Emas: Wajib Menggunakan defer cancel()
Meskipun `WithTimeout` dapat memicu pembatalan otomatis saat durasi habis (misal 5 detik), kita tetap wajib menuliskan `defer cancel()`.

Alasannya: Jika proses kodingan ternyata selesai lebih cepat (misal hanya butuh 1 detik), memanggil `cancel()` secara manual lewat `defer` akan langsung menghancurkan timer internal context saat itu juga tanpa perlu menunggu sisa waktu 4 detiknya habis. Ini sangat menghemat alokasi memori CPU server.

```go
// 🎯 Standar Penulisan Industri Backend
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel() // Memastikan resource langsung bersih setelah fungsi selesai
```

## 🛠️ Cara Menjalankan Pengujian
Masuk ke direktori folder modul lalu jalankan perintah pengujian dengan flag verbose (`-v`) untuk melihat output log secara detail:
``` bash
cd 01-context_test
go test -v
```
