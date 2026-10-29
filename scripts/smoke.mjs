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
  assert.match(api.headers.get('content-type') ?? '', /application\/json/);
  assert.equal(api.headers.get('cache-control'), 'no-store');
  const missing = await api.json();
  assert.equal(missing.error.code, 'NOT_FOUND');
  assert.match(missing.error.request_id, /^req_/);
  assert.equal(missing.error.request_id, api.headers.get('x-request-id'));

  const unsafe = await fetch('http://127.0.0.1:8080/api/v1/tasks/synthetic', {
    method: 'POST',
    signal: AbortSignal.timeout(5_000),
    redirect: 'error',
  });
  assert.equal(unsafe.status, 403, 'Origin yang hilang harus ditolak.');
  assert.equal((await unsafe.json()).error.code, 'INVALID_ORIGIN');
  console.log('PASS: web SSR Indonesia; API JSON, request ID, no-store, dan penolakan Origin lokal.');
} catch (error) {
  console.error('FAIL: ' + error.message);
  console.error('Jalankan pnpm dev:web dan pnpm dev:api sesuai README.md.');
  process.exit(1);
}
