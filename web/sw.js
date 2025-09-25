const CACHE_NAME = 'lxc-terminal-v1';
const urlsToCache = [
  '/',
  '/css/tmux-ui.css',
  '/js/terminal.js',
  '/assets/icons/icon-192x192.png'
];

self.addEventListener('install', event => {
  event.waitUntil(
    caches.open(CACHE_NAME)
      .then(cache => cache.addAll(urlsToCache))
  );
});

self.addEventListener('fetch', event => {
  event.respondWith(
    caches.match(event.request)
      .then(response => response || fetch(event.request))
  );
});