const CACHE_NAME = 'stereodose-cache';
const urlsToCache = [
  '/manifest.json',
  '/'
];

self.addEventListener('install', function (event) {
  // Perform install steps
  event.waitUntil(
    caches.open(CACHE_NAME)
      .then(function (cache) {
        console.log('Opened cache');
        return cache.addAll(urlsToCache);
      })
  );
});


self.addEventListener('fetch', function (event) {
  // https://jakearchibald.com/2014/offline-cookbook/#network-falling-back-to-cache
  event.respondWith(async function () {
    try {
      return await fetch(event.request);
    } catch (err) {
      return caches.match(event.request);
    }
  }());
});