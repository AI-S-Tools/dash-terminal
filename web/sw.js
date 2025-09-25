const CACHE_NAME = 'lxc-terminal-v2'; // Updated for session-tabs.js and responsive fixes
const urlsToCache = [
  '/',
  '/css/tmux-ui.css',
  '/js/terminal.js',
  '/js/session-tabs.js', // Added session tabs
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