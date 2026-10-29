---
name: work-on-marketing-workspace
description: Jalankan workflow berulang untuk Unggulmart Marketing Workspace dengan stack SvelteKit/TypeScript dan Go/Chi/PostgreSQL. Gunakan saat mengimplementasikan task milestone, menguji perubahan, meninjau kode atau keamanan, memperbaiki bug, melakukan refactor, mengubah kontrak API atau migration, memeriksa kalender/reminder/PWA, menyiapkan release, atau memperbarui TODO.md. Jangan gunakan untuk project lain atau menganggap penggunaan skill sebagai izin deploy produksi.
---

# Marketing Workspace Workflows

Gunakan satu workflow utama sesuai permintaan. Jalankan workflow pendukung hanya jika diperlukan untuk menyelesaikan scope. Nama workflow di bawah adalah mode dalam skill ini, bukan skill terpisah yang otomatis terdaftar.

## Pilih workflow

| Workflow | Contoh permintaan | Hasil utama |
|---|---|---|
| implement-task | Kerjakan M2-03 | Implementasi, acceptance evidence, update tracker |
| verify-change | Test perubahan task ini | Checks relevan dan laporan PASS/FAIL/NOT RUN |
| review-code | Review diff terakhir | Temuan behavior, kontrak, dan architecture |
| review-security | Audit keamanan login/review/push | Temuan dengan bukti dan dampak |
| fix-bug | Perbaiki stale approval | Reproduksi, perbaikan akar masalah, regression check |
| refactor-safely | Rapikan service tasks | Struktur lebih jelas dengan behavior tetap |
| change-api-db | Tambah field ke task API | Kontrak, migration, service, client konsisten |
| verify-calendar-jobs | Periksa WITA dan reminder | Bukti correctness waktu, dedup, race/recovery |
| verify-pwa-push | Periksa offline dan push | Bukti cache privacy, lifecycle, delivery |
| prepare-release | Cek kesiapan rilis | Readiness dan gate operasional yang belum terpenuhi |
| update-progress | Update TODO dari hasil sesi | Tabel progres dengan bukti aktual |

Jika hanya diminta review, laporkan temuan tanpa langsung mengedit kode. Jika diminta memperbaiki atau mengimplementasikan, kerjakan perubahan dan verifikasi sesuai scope. Jangan menjalankan seluruh workflow untuk perubahan kecil.

## Muat konteks project

1. Temukan root repository yang dimaksud dari konteks pengguna atau workspace. Jangan mengasumsikan current directory adalah root dan jangan mencari repository lain tanpa kebutuhan.
2. Baca `AGENTS.md` serta instruksi lokal direktori yang disentuh.
3. Baca `WORKFLOW.md`, `PRD.md`, `ARCHITECTURE.md`, dan `TODO.md` pada versi repository terkini. Seluruh dokumen Markdown proyek berada di root, termasuk SKILL.md ini; gunakan nama/kapitalisasi aktual tanpa membuat salinan. Kontrak API berada di openapi.yaml saat dibuat.
4. Gunakan PRD untuk requirement, AGENTS untuk aturan kerja, architecture untuk ownership/dependency, TODO untuk urutan dan status. Baca bagian domain terkait, kode, OpenAPI, migrations, scripts, serta CI sebelum mengubahnya.
5. Jika dokumen penting tidak tersedia, laporkan bagian yang hilang. Lanjutkan pemeriksaan independen yang aman; jangan menebak acceptance atau mengklaim kesesuaian penuh.
6. Bedakan keputusan tetap, default implementasi, dan gate yang belum dibuktikan. Jelaskan konflik dokumen/kode sebelum perubahan yang akan mengubah behavior produk.

Jangan menyalin PRD ke skill atau memperbarui PRD hanya untuk membenarkan bug. Skill ini tidak mengganti instruksi pengguna atau memberi otorisasi perubahan produksi.

## Invariant yang selalu dijaga

- Pertahankan modular monolith: web untuk UI/API consumption; Go untuk business logic, authorization dan persistence; SQL migrations untuk schema.
- Pertahankan stack dan satu workspace kecil. Jangan menambah analytics, upload, recurring items, multi-tenancy, broker atau infrastruktur di luar PRD.
- Ambil actor dari sesi. Hak baca tim tidak berarti hak edit. Supervisor bukan reviewer wajib dan tidak memiliki force approve.
- Pertahankan empat status task; overdue turunan dan archive flag. Requires Review default OFF, review sukarela tetap boleh, dan ON memblok direct completion.
- Bedakan complete dengan approve. Validasi actor, round aktif, version dan content revision; cegah self-approval serta approval hasil yang sudah berubah.
- Commit mutasi, audit, inbox/outbox yang relevan dalam transaksi yang sama. Jalankan network push setelah commit melalui state persisten.
- Gunakan Asia/Makassar, date untuk date-only, timestamptz untuk instant, all-day end eksklusif. Gregorian tetap sumber tanggal; Hijriyah memiliki provenance atau label perkiraan.
- Pertahankan server-side session, Origin/CSRF, no-store/private-cache boundary dan environment separation. Jangan menyimpan auth di localStorage atau secrets di laporan.

## implement-task

1. Resolve ID dan acceptance dari TODO. Jika tidak disebut, pilih task siap pertama dalam milestone aktif; jelaskan pilihan singkat.
2. Periksa dependensi dan gate milestone sebelumnya. Jangan membuka milestone berikut sebelum gate lulus. Jika tugas pengguna bertentangan dengan urutan, jelaskan konflik sebelum mengubah scope.
3. Isi PIC aktual dan status IN_PROGRESS untuk pekerjaan yang benar-benar dimulai. Jangan menetapkan Sultan/Alvin sebagai developer secara otomatis.
4. Tentukan perubahan minimum: kontrak → migration jika perlu → policy/service/transaksi → handler → UI.
5. Reuse boundary existing. Letakkan HTTP frontend di lib/api, UI domain di features, dan business rules di Go. Jangan membuat circular service import atau transaksi independen untuk mutasi atomik.
6. Jalankan verify-change dan acceptance terkait. Perbarui docs yang terdampak serta tracker dengan hasil nyata.

## verify-change

1. Baca diff dan tentukan risiko/perilaku yang berubah. Periksa scripts yang tersedia; jangan mengarang `pnpm quality` atau command generator.
2. Pilih checks melalui tabel di bawah. Jalankan checks wajib CI sebelum merge/release; gunakan checks terfokus selama iterasi.
3. Jika gagal, bedakan regresi perubahan, kegagalan existing, dan kendala environment berdasarkan bukti. Perbaiki sesuai scope lalu ulangi checks terdampak.
4. Laporkan command persis, hasil dan kendala. Jangan mengubah test untuk menerima pelanggaran requirement.

| Perubahan | Pemeriksaan yang relevan |
|---|---|
| Frontend | Format/lint dan Svelte/TypeScript checks melalui scripts nyata; build dan test behavior terkait; Svelte autofixer bila tersedia |
| Go | Gofmt check, `go vet ./...`, `go test ./...` dari apps/api; tambah pemeriksaan terfokus sesuai risiko |
| SQL/transaksi | Migration dan integration test pada PostgreSQL sungguhan; rollback, constraint serta concurrency |
| API | OpenAPI validation, regenerasi TypeScript client, drift dan request/response contract tests |
| Auth/review | API permission negatif, stale session/version/round, self-approval, idempotency serta race |
| Workflow lintas UI/API | E2E dengan dua sesi role berbeda dan seed sintetis |
| Copy/style ringan | Pemeriksaan visual/build relevan; jangan membuat test yang sekadar menyalin implementasi |

## review-code

1. Baca diff beserta caller, policy, DTO, schema dan tests terkait; jangan menilai file secara terisolasi.
2. Cari bug behavior terlebih dahulu: kehilangan update, izin salah, transaksi parsial, stale UI, error handling, kontrak drift, dan invalid state.
3. Periksa ownership sesuai architecture: tidak ada SQL/backend kedua di SvelteKit, raw fetch tersebar, cache ganda, circular import, atau interface berlapis tanpa kebutuhan.
4. Laporkan temuan berdasarkan dampak dengan lokasi dan skenario konkret. Pisahkan saran perapian dari bug. Jika tidak ada temuan, sebutkan batas pemeriksaan dan checks yang belum dijalankan.

## review-security

1. Tentukan scope dan trust boundary dari input browser/proxy ke session, service, SQL serta push destination.
2. Tinjau auth/session: cookie, expiry/revoke/rotation, Origin/CSRF termasuk login, rate limit durable, temporary password dan guard Supervisor terakhir.
3. Tinjau authorization lewat API: creator/PIC/reviewer, actor dari sesi, admin override beralasan, ownership inbox/subscription dan larangan self-approval.
4. Tinjau input/output: unknown fields, SQL parameter, limits, URL scheme, raw HTML/XSS, safe redirects, headers dan error/log redaction.
5. Untuk push, tinjau provider allowlist, DNS/destination validation, SSRF, redirect, session ownership dan payload lockscreen. Untuk PWA, tinjau cache privat. Untuk deployment, tinjau secrets dan pairing environment.
6. Gunakan fixture/sesi uji; jangan mengirim push atau melakukan pengujian aktif terhadap produksi tanpa otorisasi yang sesuai.
7. Laporkan severity berdasarkan exploitability/dampak, bukti, skenario, perbaikan, dan cara memverifikasi. Bedakan temuan terkonfirmasi dari hipotesis; jangan mengklaim audit menyeluruh hanya dari scan dependency.

## fix-bug

1. Catat expected behavior dari PRD serta actual behavior dari reproduksi/log/test. Jangan menebak akar masalah dari pesan error saja.
2. Tambahkan atau tautkan BUG-ID di TODO dengan langkah reproduksi dan task terdampak. Buka kembali task/gate bila acceptance tidak lagi terpenuhi.
3. Buat reproduksi terfokus, lacak akar masalah, lalu lakukan perbaikan minimum. Pertahankan perubahan pengguna yang tidak terkait.
4. Jalankan regression check yang membuktikan bug tertutup dan test terkait. Buat test baru bila risikonya sebanding.
5. Catat apa yang diperbaiki dan batas verifikasi; gunakan update-progress.

## refactor-safely

1. Nyatakan behavior yang harus tetap: kontrak HTTP, permission, state transitions, schema, completion/history, waktu dan notification semantics yang terkait.
2. Periksa baseline checks dan konsumen kode. Tentukan refactor kecil yang menyelesaikan masalah nyata.
3. Pindahkan tanggung jawab ke boundary yang benar tanpa fitur/dependency baru atau formatting massal.
4. Pertahankan OpenAPI/DB compatibility. Jika perlu perubahan behavior atau schema, jelaskan dan tangani sebagai perubahan eksplisit, bukan menyembunyikannya sebagai refactor.
5. Jalankan checks relevan dan bandingkan behavior sebelum/sesudah. Laporkan struktur yang berubah dan bukti behavior tetap.

## change-api-db

1. Baca inventory PRD, OpenAPI, handler/service/repository, generated client, migrations dan semua konsumen terkait.
2. Tentukan DTO serta compatibility; jangan mengekspos field database sensitif secara otomatis.
3. Perbarui kontrak bersama implementasi. Untuk schema, buat migration baru yang immutable; gunakan expand/deploy/contract sesuai kebutuhan.
4. Pertahankan FK, CHECK, indexes relevan, append-only history, version dan idempotency. Jangan hard delete resource atau mengubah migration terpakai.
   Ikuti staging schema PRD §13.5: completion CHECK dasar sejak M2, FK review/event/reminder saat tabel tujuan tersedia. Uji upgrade dari milestone sebelumnya dan seluruh rantai dari database kosong.
5. Uji migration dan transaksi pada PostgreSQL, regenerate client dan validasi drift. Uji stale expected_version dan retry key yang sama bila mutasi terdampak.
6. Jalankan migration produksi hanya dalam scope yang diotorisasi melalui release runner/direct connection. Jangan migrasi pada cold start atau build frontend.

## verify-calendar-jobs

1. Gunakan clock dan dataset fixture deterministik; uji perangkat/host di timezone berbeda.
2. Uji due_date vs due_at, midnight WITA, akhir bulan/tahun, exclusive all-day end, overlap/range dan projection task tanpa duplikasi event.
3. Uji metadata/cakupan Hijriyah, fallback berlabel dan perubahan dataset tanpa pergeseran Gregorian.
4. Untuk reminder, uji H-N sebagai hari kalender lokal, default/custom/rules kosong, recipients dedup, skipped dan catch-up sesuai PRD.
5. Uji reschedule/archive/restore, revision, concurrent cron, event→reminder lock order, unique inbox dan cancellation delivery lama.
6. Uji cron secret/flags/batch/budget serta recovery sebelum/sesudah commit. Jangan mengandalkan scheduler tepat waktu atau background goroutine.

## verify-pwa-push

1. Periksa manifest, satu registrasi service worker, lifecycle cleanup serta prompt update yang aman untuk form.
2. Periksa cache: asset publik saja; API/halaman privat network-only/no-store; offline save diblok, tanpa mutation queue.
3. Periksa feature detection, permission setelah klik, denied tanpa prompt berulang, dan panduan install pada perangkat aktual.
4. Uji login/logout, session revoke, ownership transfer, preferensi dan tidak ada backlog perangkat baru.
5. Dengan fake transport, uji lease recovery, retry/TTL, 404/410, 429/5xx, config error dan crash setelah provider menerima. Pastikan inbox tetap satu dan tidak menjanjikan exactly-once perangkat.
6. Uji payload generik dan notificationclick/deep link aman setelah login. Laporkan perangkat yang benar-benar diuji serta yang belum tersedia.

## prepare-release

1. Baca gate M7, Definition of Done PRD, CI, deployment/runbook dan ADR terkini. Susun readiness berdasarkan bukti, bukan status TODO saja.
2. Verifikasi checks penuh, Go routing/proxy cookie, environment pairing, secrets, runtime limits, pool, data kalender, paket bisnis/cron dan budget.
3. Periksa urutan migration sekali dengan lock → backend → web → smoke → aktivasi job; pastikan rollback schema tidak dianggap otomatis.
4. Periksa restore rehearsal database terpisah, RPO/RTO, monitoring/cron alert, retensi aman dedup dan UAT kedua pengguna.
5. Jika pengguna hanya meminta pemeriksaan kesiapan, hasilkan laporan tanpa deploy atau perubahan akses. Jika deploy diotorisasi, ikuti gate dan runbook, lalu verifikasi status serta smoke aktual.
6. Tulis READY hanya jika seluruh gate wajib memiliki bukti. Daftarkan gate belum terpenuhi beserta tindakan konkret; jangan menjanjikan hosting gratis atau dukungan perangkat yang belum diuji.

## update-progress

1. Baca tabel TODO terkini, diff, hasil checks dan bukti sesi. Jangan menandai task selesai hanya dari rencana atau klaim sesi lama.
2. Perbarui baris ID yang sama: status, PIC aktual, acceptance evidence dan catatan bertanggal WITA. Pertahankan ID stabil serta perubahan pengguna.
3. Gunakan vocabulary TODO: TODO, IN_PROGRESS, VERIFY, BLOCKED, DONE, CANCELLED. VERIFY jika implementasi tersedia tetapi acceptance wajib belum lengkap; BLOCKED untuk hambatan konkret dengan aksi pembuka dan owner.
4. Beri DONE hanya setelah acceptance/checks wajib terpenuhi. Jangan membatalkan scope tanpa keputusan eksplisit.
5. Hitung ulang ringkasan dari leaf task aktif; keluarkan gate dari denominator, hindari double count parent/child. Persentase bukan estimasi usaha.
6. Perbarui fokus sesi berikut, issue terkait, gate dan log singkat. Jika bug membuka task DONE, buka kembali gate yang terdampak dan nilai dependensi milestone lanjutan.
7. Jangan mengubah status berdasarkan dokumen spesifikasi yang baru selesai ditulis. Bedakan kode lokal, preview dan production.
8. Periksa pemetaan OQ ke issue keputusan dan task/gate pada TODO. Default kerja bukan bukti konfirmasi; catat keputusan nyata sebelum tenggat PRD dan jangan menutup OQ karena implementasi default atau perubahan dokumen selesai.

## Format laporan

Gunakan bagian yang relevan saja; tulis Bahasa Indonesia dan utamakan hasil.

Untuk implementasi/fix/refactor, jelaskan perubahan dan alasan, lalu bukti:

| Pemeriksaan | Command / metode aktual | Hasil | Bukti / kendala |
|---|---|---|---|

Gunakan PASS, FAIL, NOT RUN; jangan isi PASS tanpa menjalankan pemeriksaan. Sertakan file/commit bila tersedia tanpa mengarang ID atau hasil.

Untuk review, gunakan:

| Severity | Lokasi | Temuan dan skenario | Dampak | Perbaikan / verifikasi |
|---|---|---|---|---|

Untuk progres/handoff, gunakan:

| Task ID | Status sebelum → sesudah | Bukti | Langkah berikut / blocker |
|---|---|---|---|

Jangan mencantumkan secrets atau data kerja sensitif. Jangan mengklaim seluruh quality gates lulus jika hanya subset yang dijalankan. Tutup laporan dengan langkah berikut yang konkret hanya jika masih ada pekerjaan tersisa.
