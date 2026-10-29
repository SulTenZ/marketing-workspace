import assert from 'node:assert/strict';

// Jalankan kedua server lokal lebih dahulu; tidak pernah mengakses environment lain.
try {
  const web = await fetch('http://127.0.0.1:5173/', {
    signal: AbortSignal.timeout(10_000),
    redirect: 'error',
  });
  assert.equal(web.status, 200, 'Web harus merespons 200.');
  assert.match(web.headers.get('content-type') ?? '', /text\/html/);
  const html = await web.text();
  assert.match(html, /<html lang="id">/, 'Dokumen harus berbahasa Indonesia.');
  assert.match(html, /<title>Unggulmart Marketing Workspace<\/title>/);
  assert.match(html, /<h1[^>]*>Unggulmart Marketing Workspace<\/h1>/);

  const api = await fetch('http://127.0.0.1:8080/', {
    signal: AbortSignal.timeout(5_000),
    redirect: 'error',
  });
  assert.equal(api.status, 404, 'Scaffold API belum memiliki endpoint produk.');
  console.log('PASS: web SSR 200 (Bahasa Indonesia) dan server Go/Chi lokal merespons.');
} catch (error) {
  console.error('FAIL: ' + error.message);
  console.error('Jalankan pnpm dev:web dan pnpm dev:api sesuai README.md.');
  process.exit(1);
}
