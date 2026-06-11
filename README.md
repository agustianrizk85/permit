# Green Park — Legal Permit System

Sistem manajemen legal & perizinan lahan untuk developer properti, mencakup
seluruh alur **Proses A–I**, Master Data PT, Early Warning AI, OCR, watermark
Confidential, dan setting DACI.

- **Backend:** Go (Gin + GORM), arsitektur berlapis `model → repository → service → handler`. Mendukung **PostgreSQL** (produksi) & **SQLite** (run cepat) via `DB_DRIVER`.
- **Frontend:** React + TypeScript + Vite, dengan lapisan `models` dan `services`.

```
legalpermit/
├── backend/
│   ├── cmd/server/main.go
│   └── internal/
│       ├── config database model repository service handler middleware dto seed
│       ├── ocr/         # OCR AI pluggable (mock default; Claude vision dll bisa di-plug)
│       ├── storage/     # helper simpan file (dipakai semua uploader)
│       └── watermark/   # Confidential + hitam-putih untuk gambar
└── frontend/
    └── src/
        ├── models/ services/ context/ pages/ components/ lib/
```

## Menjalankan

### Opsi A — SQLite (langsung jalan, tanpa setup)
```bash
cd backend
# .env sudah di-set DB_DRIVER=sqlite
go run ./cmd/server            # API :8080, file DB: backend/legalpermit.db
```

### Opsi B — PostgreSQL
```bash
docker compose up -d db        # atau pakai Postgres lokal
cd backend
# di .env: DB_DRIVER=postgres + isi DB_PASSWORD
go run ./cmd/server
```

### Frontend
```bash
cd frontend && npm install && npm run dev   # http://localhost:5173 (atau 5174)
```

Akun seed otomatis (password dari env, ganti setelah login):

| Email | Role | Password default |
|---|---|---|
| dirops@greenpark.id | dirops | `dirops123` |
| kadep@greenpark.id | kadep | `kadep123` |

## Cakupan fitur

| Modul | Status |
|---|---|
| **Proses A** Pra-Akad (A1–A8) | ✅ penuh, aturan Harga Fix/SPK ditegakkan |
| **Proses B** Akad Lahan (B1–B4) | ✅ |
| **Proses C** Permit (C1–C12) | ✅ (C10 Siteplan flag notif lintas dep; C12 IMB/PBG confidential) |
| **Proses D** Legal (D1–D2) | ✅ confidential output |
| **Proses E** Master Data PT | ✅ entity terpisah + upload dokumen |
| **Proses F** Master Data PKS Bank (F1–F2) | ✅ |
| **Proses H** Flow PKS Bank (H1–H5) | ✅ |
| **Proses I** SPK Legal Permit | ✅ (registrasi SPK) |
| **SLA / deadline** per step (dinamis) | ✅ default per template, override via API |
| **Early Warning System (AI)** | ✅ rule-based (overdue / due-soon / input kurang) |
| **Search All Dokumen** | ✅ |
| **DACI setting (KADEP/DIROPS)** | ✅ driver/approver/consulting/informed |
| **OCR AI** (KTP, KK, NPWP, SHM, PBB) | ✅ scaffold pluggable, provider **mock** default |
| **Watermark Confidential + hitam-putih** | ✅ untuk gambar (`?watermark=1`) |
| Notifikasi WA / Audit AI chatbot | ⚙️ konfigurasi tersedia; pengiriman WA perlu di-wire ke provider |
| Watermark untuk PDF | ⚙️ saat ini gambar (JPG/PNG); PDF perlu lib tambahan |

Total **34 langkah** ter-seed otomatis saat membuat satu lahan (A=8, B=4, C=12, D=2, F=2, H=5, I=1).

## API utama

| Method | Endpoint | Keterangan |
|---|---|---|
| POST | `/api/auth/login` · GET `/api/auth/me` | Auth JWT |
| GET | `/api/meta/categories` | Label kategori A–I |
| GET/POST | `/api/projects` | List / buat lahan (auto-seed A–I) |
| GET | `/api/projects/:id` · `/api/projects/:id/progress` | Detail + progres |
| GET/PUT | `/api/steps/:id` | Detail / update (status, harga, SPK, metadata, **SLA**) |
| POST | `/api/steps/:id/documents` | Upload dokumen |
| GET | `/api/documents/:id/download` `?watermark=1` | Unduh (opsi confidential) |
| GET/POST | `/api/pt` · GET `/api/pt/:id` | Master Data PT (E) |
| POST | `/api/pt/:id/documents` · GET `/api/pt-documents/:id/download` | Dokumen PT |
| GET | `/api/dashboard/warnings` | Early Warning AI |
| GET | `/api/dashboard/documents?q=` | Search dokumen |
| POST | `/api/ocr/extract` | OCR AI (multipart) |
| GET/PUT | `/api/settings/daci` · `/api/settings/notification` | Setting (PUT: KADEP/DIROPS) |

## Cara menambah proses baru
Cukup tambahkan `StepTemplate` di
[`backend/internal/service/catalog.go`](backend/internal/service/catalog.go)
dan (opsional) hint UI di
[`frontend/src/lib/processCatalog.ts`](frontend/src/lib/processCatalog.ts).
Seluruh engine (status, harga, SPK, SLA, dokumen, warning) sudah generik.

## Menyambungkan OCR / WA sungguhan
- **OCR:** implementasikan `ocr.Provider` (mis. Claude vision) lalu ganti
  `ocr.NewMockProvider()` di [router.go](backend/internal/handler/router.go).
- **WhatsApp:** baca `NotificationConfig` (URL & toggle dari setting) dan kirim
  via provider WA pada job reminder harian.
