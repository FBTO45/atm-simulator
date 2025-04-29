# 💳 Simulasi Mesin ATM - CLI Go + MySQL

📘 **Assignment Guidance - Digital Skill Fair 38.0**  
🧠 *Golang Fundamental & Database MySQL*

---

## 🎯 Objectives

1. 🛠 Membangun aplikasi CLI sederhana menggunakan **Go** dan `urfave/cli`
2. 💾 Menggunakan **SQL** untuk manipulasi data (CRUD)
3. 🗂 Merancang relational database yang mencerminkan proses bisnis

---

## 📝 Deskripsi Assignment

Kamu diminta membuat aplikasi **simulasi mesin ATM berbasis CLI (Command-Line Interface)** menggunakan bahasa **Go** dan database backend **MySQL**.

Aplikasi ini harus mampu menangani proses keuangan dasar:

- Registrasi akun
- Login
- Cek saldo
- Setor (deposit)
- Tarik tunai (withdraw)
- Transfer antar akun

Mahasiswa juga harus menyimpan **histori transaksi** dan menangani **error** atau **edge case** secara eksplisit.

---

## 🔧 Fungsi CLI yang Harus Diimplementasikan

| Perintah CLI     | Deskripsi |
|------------------|-----------|
| `register`       | Membuat akun baru (input: nama, PIN) |
| `login`          | Login ke akun dengan nomor akun & PIN |
| `check-balance`  | Menampilkan saldo akun saat ini |
| `deposit`        | Menambahkan saldo ke akun login |
| `withdraw`       | Mengurangi saldo akun login (tidak boleh negatif) |
| `transfer`       | Mentransfer saldo ke akun lain (input: nomor akun tujuan) |

---

## 🧱 Spesifikasi Teknis

### 🗃 Struktur Tabel MySQL

#### Tabel `accounts`

| Kolom      | Tipe       | Keterangan                             |
|------------|------------|----------------------------------------|
| `id`       | INT (PK)   | Nomor akun                             |
| `name`     | VARCHAR    | Nama pengguna                          |
| `pin`      | VARCHAR    | PIN (boleh dalam plain atau hash)      |
| `balance`  | DECIMAL    | Saldo rekening                         |
| `created_at`| TIMESTAMP | Tanggal pembuatan akun                 |

#### Tabel `transactions`

| Kolom       | Tipe       | Keterangan                                      |
|-------------|------------|-------------------------------------------------|
| `id`        | INT (PK)   | ID transaksi                                    |
| `account_id`| INT (FK)   | Akun pemilik transaksi                          |
| `type`      | ENUM       | `deposit`, `withdraw`, `transfer_in`, `transfer_out` |
| `amount`    | DECIMAL    | Jumlah transaksi                                |
| `target_id` | INT (NULL) | Akun tujuan transfer (jika `transfer`)          |
| `description` | VARCHAR(255) | Keterangan transaksi (opsional)             |
| `created_at`| TIMESTAMP  | Tanggal transaksi                               |

---

## ✅ Skenario & Validasi

1. 🔐 **Login** hanya berlaku selama satu sesi program (gunakan variabel global)
2. 🚫 **Transaksi tidak boleh dilakukan jika belum login**
3. 💸 **Withdraw & transfer harus menolak jika saldo tidak mencukupi**
4. 🔄 **Transfer** harus membuat **dua entri transaksi**:
   - `transfer_out` dari akun pengirim
   - `transfer_in` ke akun penerima

---

## 🧰 Tools

- 🖥️ [Visual Studio Code](https://code.visualstudio.com/)
- 🧠 [Goland](https://www.jetbrains.com/go/)
- 🛢️ [MySQL](https://www.mysql.com/)

---

## 🚀 Cara Menjalankan

```bash
# Clone repository
git clone https://github.com/username/atm-cli-go.git
cd atm-cli-go

# Jalankan aplikasi
go run main.go
