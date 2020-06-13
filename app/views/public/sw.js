
// Only set the service worker once the user has logged in.
// This prevents certain login bugs where the PWA is switched to during the login process but doesn't recieve auth cookies.
// It also only asks the user to install the app if they log in
if (userLoggedIn()) {
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

// userLoggedIn returns true if the user is logged in, false otherwise
function userLoggedIn() {
  // stolen from Stack Overflow
  function getCookie(name) {
    var dc = document.cookie;
    var prefix = name + '=';
    var begin = dc.indexOf('; ' + prefix);
    if (begin === -1) {
      begin = dc.indexOf(prefix);
      if (begin !== 0) return null;
    } else {
      begin += 2;
      var end = document.cookie.indexOf(';', begin);
      if (end === -1) {
        end = dc.length;
      }
    }
    // because unescape has been deprecated, replaced with decodeURI
    //return unescape(dc.substring(begin + prefix.length, end));
    return decodeURI(dc.substring(begin + prefix.length, end));
  }

  let cookie = getCookie('stereodose_session');
  if (!cookie) {
    return false;
  }

  return true;
}
