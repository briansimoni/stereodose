const CACHE_NAME = 'stereodose-cache';
const urlsToCache = [
  '/manifest.json',
  '/sw.js'
];

self.addEventListener('install', function(event) {
  // Perform install steps
  event.waitUntil(
    caches.open(CACHE_NAME)
      .then(function(cache) {
        console.log('Opened cache');
        return cache.addAll(urlsToCache);
      })
  );
});


self.addEventListener('fetch', function(event) {
  let path = new URL(event.request.url).pathname;
  if (path === '/' && !navigator.onLine) {
    event.respondWith(new Promise( function(resolve, reject) {
      const body = new Blob(['You need to be online for this app to work'], {type : 'text/html'});
      const res = new Response(body, {status: 200, statusText: 'OK'});
      resolve(res);
    }));
    return;
  }

  event.respondWith(
    caches.match(event.request)
      .then(function(response) {
        // Cache hit - return response
        if (response) {
          return response;
        } 
        return fetch(event.request);
      }
    )
  );
});