const CACHE_NAME = 'lxc-terminal-v5'; // Fixed Terminal class availability
const urlsToCache = [
  '/',
  '/css/tmux-ui.css',
  '/js/terminal.js',
  '/js/session-tabs.js', // T3.1 session tabs
  '/js/window-tabs.js', // T3.2 window tabs
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