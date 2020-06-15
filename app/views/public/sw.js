// Only set the service worker once the user has logged in.
// This prevents certain login bugs where the PWA is switched to during the login process but doesn't recieve auth cookies.
// It also only asks the user to install the app if they log in
isLoggedIn.then((loggedIn) => {
  if (loggedIn) {
    const CACHE_NAME = 'stereodose-cache';
    const urlsToCache = ['/manifest.json', '/'];

    self.addEventListener('install', function (event) {
      // Perform install steps
      event.waitUntil(
        caches.open(CACHE_NAME).then(function (cache) {
          console.log('Opened cache');
          return cache.addAll(urlsToCache);
        })
      );
    });

    self.addEventListener('fetch', function (event) {
      // https://jakearchibald.com/2014/offline-cookbook/#network-falling-back-to-cache
      event.respondWith(
        (async function () {
          try {
            return await fetch(event.request);
          } catch (err) {
            return caches.match(event.request);
          }
        })()
      );
    });
  }
});

async function isLoggedIn() {
  const response = await fetch('/api/users/me');
  if (response.status === 200) {
    return true;
  }
  return false;
}
