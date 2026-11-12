# Spesifikasi: Membuat workout

`POST /api/workouts` · revisi 2

## Prinsip pembagian tanggung jawab

Mobile menghitung semua nilai turunan (langkah, jarak, kalori, pace, heart rate, splits, jenis lari per km) dan menampilkan ringkasan sebelum data dikirim. Server tidak menghitung ulang algoritma tersebut. Server memeriksa **invarian**: syarat yang harus benar pada workout mana pun, apa pun algoritma yang dipakai mobile (tidak ada nilai mustahil, dan antar-field konsisten).

Server menjamin data _wajar dan konsisten_, bukan _dihitung dengan benar_. Perubahan algoritma di mobile tidak butuh perubahan server selama invarian tetap terpenuhi.

## Pemicu dan prasyarat

### Pemicu

`POST /api/workouts` dari mobile (Android) saat pengguna menekan "Submit" pada layar ringkasan.

### Prasyarat

- Token akses valid.
- Pengguna dikenali dari token.
- Mobile menyimpan data workout secara lokal sampai menerima respons 200 atau 201, supaya pengiriman ulang tetap mungkin.

## Input

### Request body

```json
{
  "client_session_id": "string",
  "device_name": "string",
  "device_mac_address": "string",
  "workout_type": "running | walking",
  "started_at": "2026-11-10T06:42:00+07:00",
  "active_duration_sec": 0,
  "summary": {
    "steps": 0,
    "distance_m": 0,
    "calories_kcal": 0,
    "avg_pace_sec_per_km": null,
    "avg_heart_rate_bpm": null,
    "max_heart_rate_bpm": null
  },
  "splits": [
    {
      "pace_sec_per_km": 0,
      "avg_heart_rate_bpm": null,
      "kind": "default | warm_up | fastest | cool_down"
    }
  ],
  "heart_rates": [
    {
      "bpm": 0,
      "recorded_at": "2026-11-10T06:42:05+07:00"
    }
  ]
}
```

### Aturan bentuk (pelanggaran = 400)

| Field                                                                                                                    | Aturan                                                                 |
| ------------------------------------------------------------------------------------------------------------------------ | ---------------------------------------------------------------------- |
| `client_session_id`                                                                                                      | wajib, UUID                                                            |
| `device_name`                                                                                                            | wajib, 1 sampai 100 karakter                                           |
| `device_mac_address`                                                                                                     | wajib, 6 oktet heksadesimal (pemisah `:` atau `-`, atau tanpa pemisah) |
| `workout_type`, `splits[].kind`                                                                                          | wajib, string (keanggotaan enum diperiksa sebagai invarian)            |
| `started_at`, `heart_rates[].active_ms`                                                                                  | wajib, RFC 3339 dengan zona waktu (pecahan detik diperbolehkan)        |
| `active_duration_sec`, `summary.steps`, `splits[].pace_sec_per_km`, `heart_rates[].bpm`                                  | wajib, bilangan bulat                                                  |
| `summary.distance_m`, `summary.calories_kcal`                                                                            | wajib, bilangan                                                        |
| `summary.avg_pace_sec_per_km`, `summary.avg_heart_rate_bpm`, `summary.max_heart_rate_bpm`, `splits[].avg_heart_rate_bpm` | bilangan bulat atau null                                               |
| `summary`, `splits`, `heart_rates`                                                                                       | wajib ada (array boleh kosong)                                         |

### Pemilik

Pemilik workout diambil dari token, tidak pernah dari body.

## Aturan bisnis

1. Pemilik workout adalah pengguna dari token.
2. `device_mac_address` dinormalisasi (huruf besar, pemisah `:`) sebelum dihitung sidik jarinya dan disimpan.
3. Workout dikenali lewat (`user_id`, `client_session_id`). Bila sudah ada, server membandingkan sidik jari payload (lihat bagian Idempotensi). Sama: kembalikan yang tersimpan tanpa membuat baru dan tanpa memeriksa ulang invarian. Berbeda: tolak dengan 409.
4. Isi diperiksa terhadap katalog invarian. Semua pelanggaran dikumpulkan dan dilaporkan sekaligus dengan kodenya, tidak berhenti di yang pertama.
5. `workout_type` dan `kind` hanya diperiksa keanggotaannya dalam enum. Aturan penentuan `kind` milik mobile.
6. Workout yang lolos disimpan bersama `splits` dan `heart_rates` secara atomik (satu transaksi).
7. Setiap penolakan (422 dan 409) dicatat di log dengan kode pelanggaran, `user_id`, `client_session_id`, dan User-Agent. Isi data kesehatan (bpm, series, ringkasan) tidak dicatat.

## Urutan pemrosesan

1. Autentikasi. Gagal: 401.
2. Batas ukuran body. Terlampaui: 413.
3. Parsing dan bentuk. Gagal: 400.
4. Normalisasi (MAC).
5. Cari workout (`user_id`, `client_session_id`). Ada dan sidik jari sama: 200. Ada dan sidik jari berbeda: 409.
6. Periksa invarian. Ada pelanggaran: 422 dengan semua kode sekaligus.
7. Simpan atomik. Bila gagal karena unique constraint, ambil yang tersimpan lalu bandingkan sidik jari (200 atau 409). Kegagalan lain: 500.
8. Catat log, kembalikan 201.

Batas 400 dan 422: tidak bisa di-parse atau tipenya salah adalah 400. Bisa di-parse tetapi di luar aturan domain (termasuk nilai enum yang tidak dikenal) adalah 422.

## Katalog invarian

Satu pelanggaran dilaporkan dengan satu kode. Nilai parameter ada di tabel sesudahnya.

| Kode                      | Aturan                                                                                                                                                                                            |
| ------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `NEGATIVE_VALUE`          | `steps` dan `calories_kcal` tidak negatif                                                                                                                                                         |
| `DURATION_OUT_OF_RANGE`   | `MinDurasiAktifSec`/60 di prod                                                                                                                                                                    |
| `DISTANCE_OUT_OF_RANGE`   | `distance_m` berada di \[0, MaxJarakM\]                                                                                                                                                           |
| `WORKOUT_IN_FUTURE`       | `started_at` tidak lebih dari ToleransiJamPerangkat setelah waktu server                                                                                                                          |
| `PACE_PRESENCE_INVALID`   | `avg_pace_sec_per_km` kosong bila dan hanya bila `distance_m` = 0                                                                                                                                 |
| `PACE_OUT_OF_RANGE`       | `avg_pace_sec_per_km` (bila ada) dan setiap `splits[].pace_sec_per_km` berada di \[PaceMin, PaceMax\]                                                                                             |
| `AVG_PACE_INCONSISTENT`   | bila `distance_m` > 0, selisih `avg_pace_sec_per_km` dengan `active_duration_sec ÷ (distance_m ÷ 1000)` tidak melebihi ToleransiPace                                                              |
| `SPLITS_COUNT_MISMATCH`   | jumlah `splits` sama dengan floor(`distance_m` ÷ 1000)                                                                                                                                            |
| `HR_PRESENCE_INVALID`     | bila `heart_rates` kosong: `avg_heart_rate_bpm`, `max_heart_rate_bpm`, dan semua `splits[].avg_heart_rate_bpm` kosong. Bila tidak kosong: `avg_heart_rate_bpm` dan `max_heart_rate_bpm` wajib ada |
| `BPM_OUT_OF_RANGE`        | setiap bpm (series, rata-rata, maksimum, dan rata-rata per split, bila ada) berada di \[BpmMin, BpmMax\]                                                                                          |
| `HR_SUMMARY_INCONSISTENT` | `max_heart_rate_bpm` ≥ `avg_heart_rate_bpm`; `max_heart_rate_bpm` ≥ bpm tertinggi di `heart_rates`; setiap `splits[].avg_heart_rate_bpm` ≤ `max_heart_rate_bpm`                                   |
| `HR_TIME_ORDER_INVALID`   | `heart_rates[].active_ms` tidak menurun menurut urutan array                                                                                                                                      |
| `HR_TIME_OUT_OF_RANGE`    | setiap `heart_rates[].active_ms` berada di \[`started_at`, `started_at` + MaxRentangSesi\]                                                                                                        |
| `TOO_MANY_HEART_RATES`    | jumlah `heart_rates` tidak melebihi MaxHeartRates                                                                                                                                                 |
| `TYPE_INVALID`            | `workout_type` adalah anggota enum (`running`, `walking`)                                                                                                                                         |
| `KIND_INVALID`            | setiap `splits[].kind` adalah anggota enum (`default`, `warm_up`, `fastest`, `cool_down`)                                                                                                         |

### Parameter (nilai awal)

| Parameter             | Nilai awal         | Catatan                                               |
| --------------------- | ------------------ | ----------------------------------------------------- |
| BpmMin, BpmMax        | 30, 250            | batas "mustahil", bukan batas "tidak biasa"           |
| PaceMin, PaceMax      | 120, 1800 detik/km | 2:00 sampai 30:00 per km, mencakup jalan santai       |
| ToleransiPace         | 5%                 | dari pace hasil hitung                                |
| ToleransiJamPerangkat | 5 menit            | selisih jam ponsel terhadap server                    |
| MaxDurasiAktifSec     | 86400              | 24 jam                                                |
| MaxJarakM             | 350000             | 350 km                                                |
| MaxRentangSesi        | 24 jam             | batas atas waktu sampel dari `started_at`             |
| MaxHeartRates         | 86400              | MaxDurasiAktifSec × 1 sampel per detik                |
| Batas body (413)      | TBD                | turunkan dari batas jumlah elemen × ukuran per elemen |

Nilai awal dikalibrasi dari data live testing. Pantau jumlah pelanggaran per kode di log, lalu ketatkan bertahap.

Bila waktu sampel berada pada linimasa aktif (jeda dikeluarkan), batas atas `HR_TIME_OUT_OF_RANGE` dapat diperketat menjadi `started_at` + `active_duration_sec` + toleransi (lihat Keputusan terbuka nomor 1).

## Idempotensi

Workout dikenali lewat (`user_id`, `client_session_id`). Dua situasi harus ditangani.

**1. ID sama, isi berbeda.** Bila `client_session_id` sudah ada tetapi isinya berbeda dari yang tersimpan, server tidak boleh diam-diam mengembalikan data lama. Server membandingkan sidik jari: hash dari representasi kanonik workout setelah normalisasi (bukan byte JSON mentah, supaya urutan field dan spasi tidak memicu konflik palsu). Sidik jari disimpan di kolom `payload_hash` saat workout dibuat. Sama: 200. Berbeda: 409 dengan kode `SESSION_ID_REUSED`.

**2. Dua request identik bersamaan.** Pengecekan "sudah ada?" hanyalah jalur cepat. Penjaga sebenarnya adalah unique constraint (`user_id`, `client_session_id`) di database. Bila insert gagal karena pelanggaran unique, ambil workout yang tersimpan, bandingkan sidik jari, lalu kembalikan 200 atau 409. Pelanggaran unique tidak pernah menjadi 500.

Unique index harus memperhitungkan baris yang dihapus lunak (soft delete): gunakan partial index (`WHERE deleted_at IS NULL`) atau sengaja memblokir ID lama. Putuskan saat menulis migrasi.

## Hasil

| Situasi             | Status | Isi `data`                                                  |
| ------------------- | ------ | ----------------------------------------------------------- |
| Dibuat              | 201    | `id`, `client_session_id`, `created_at`                     |
| Sudah ada, isi sama | 200    | `id`, `client_session_id`, `created_at` dari yang tersimpan |

Contoh respons 201:

```json
{
  "code": 201,
  "status": true,
  "data": {
    "id": "uuid",
    "client_session_id": "string",
    "created_at": "2026-11-10T06:50:00Z"
  },
  "message": "workout created"
}
```

## Kegagalan

| Situasi                                              | Status | Respons                                                    |
| ---------------------------------------------------- | ------ | ---------------------------------------------------------- |
| Bentuk salah                                         | 400    | galat validasi                                             |
| Token tidak valid                                    | 401    |                                                            |
| Ukuran melebihi batas                                | 413    | batas ditegakkan server                                    |
| `client_session_id` sudah dipakai dengan isi berbeda | 409    | `SESSION_ID_REUSED`                                        |
| Invarian dilanggar                                   | 422    | "workout tidak wajar" + daftar kode pada `data.violations` |
| Gagal menyimpan                                      | 500    | pesan generik, detail hanya di log                         |

Contoh respons 422:

```json
{
  "code": 422,
  "status": false,
  "data": { "violations": ["AVG_PACE_INCONSISTENT", "BPM_OUT_OF_RANGE"] },
  "message": "workout tidak wajar"
}
```

## Perilaku mobile terhadap respons

| Respons                      | Tindakan mobile                                                                                                 |
| ---------------------------- | --------------------------------------------------------------------------------------------------------------- |
| 200, 201                     | tandai tersinkron, hapus dari antrean lokal                                                                     |
| 401                          | segarkan token, kirim ulang dengan `client_session_id` yang sama                                                |
| 400, 409, 413, 422           | jangan kirim ulang otomatis (hasilnya akan sama); simpan data lokal, tandai gagal sinkron, tampilkan pesan umum |
| 500, timeout, jaringan putus | kirim ulang dengan jeda bertambah, `client_session_id` yang sama                                                |

## Keputusan terbuka

1. **Nama dan semantik `active_ms`.** Nilainya RFC 3339. Apakah itu waktu dinding atau linimasa aktif (jeda dikeluarkan)? Bila linimasa aktif, perketat batas atas `HR_TIME_OUT_OF_RANGE`.
2. **Sumber jarak** (GPS ponsel atau estimasi langkah). Menentukan apakah pemeriksaan jarak per langkah layak ditambahkan.
3. **Sesi dengan jarak sangat kecil tetapi lebih dari 0** menghasilkan pace ekstrem dan akan ditolak `PACE_OUT_OF_RANGE`. Putuskan: mobile membuang sesi itu, atau ambang "pace kosong" dinaikkan dari 0.
4. **Apakah `heart_rates` di-downsample sebelum dikirim?** Mempengaruhi pemeriksaan `HR_SUMMARY_INCONSISTENT`.
5. **Kalibrasi parameter** dari data live testing.
6. **MaxHeartRates dan batas body** diturunkan dari laju sampling maksimum mobile.
7. **Unique index dan soft delete** (lihat bagian Idempotensi).
