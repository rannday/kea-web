/// <reference lib="webworker" />

const CACHE_NAME = 'cag-cache-v1';
const CACHE_TIMEOUT_SECONDS = 30 * 24 * 60 * 60; // 30 days in seconds
const UNSUPPORTED_SCHEMES = [
  'chrome-extension:',
  'moz-extension:',
  'ms-browser-extension:',
  'about:',
  'data:',
];

/**
 * Validates whether a request is safe and cacheable.
 * @param {Request} request - The intercepted request.
 * @param {Response} response - The fetched response.
 * @returns {boolean} True if the request is cacheable; otherwise, false.
 */
function isCacheableRequest(request, response) {
  const url = new URL(request.url);

  if (UNSUPPORTED_SCHEMES.includes(url.protocol)) return false; // Skip unsupported schemes
  if (!response || response.status !== 200 || response.type !== 'basic') return false; // Only cache successful, basic responses

  return url.pathname.endsWith('.css') || url.pathname.endsWith('.js'); // Only cache CSS/JS
}

/**
 * Finds the latest bundled CSS/JS file dynamically.
 * @param {string} path - The directory path to check.
 * @param {string} ext - The file extension (css/js).
 * @returns {Promise<string|null>} - The latest bundle filename or null if not found.
 */
async function findLatestBundle(path, ext) {
  try {
    const response = await fetch(path, { cache: 'no-store' });
    if (!response.ok) throw new Error(`Failed to list directory: ${path}`);

    const text = await response.text();
    const regex = new RegExp(`bundle\\.[a-f0-9]+\\.${ext}`, 'i');
    const match = text.match(regex);

    if (match) {
      console.log(`[SW] Found latest bundle: ${match[0]}`);
      return path + match[0]; // Return full path like "/static/css/bundle.abc123.css"
    }
  } catch (err) {
    console.error(`[SW] Error finding bundle in ${path}:`, err);
  }
  return null;
}

/**
 * Clears old cache when a new version is deployed.
 */
self.addEventListener('activate', (event) => {
  console.log('[SW] Activating...');
  event.waitUntil(
    caches.keys().then((cacheNames) => {
      return Promise.all(
        cacheNames.filter(cache => cache !== CACHE_NAME).map(cache => caches.delete(cache))
      );
    })
  );
});

/**
 * Install the service worker and pre-cache dynamic assets.
 */
self.addEventListener('install', (event) => {
  console.log('[SW] Installing...');
  event.waitUntil(
    (async () => {
      const cache = await caches.open(CACHE_NAME);

      // Dynamically fetch latest bundled files
      const cssBundle = await findLatestBundle('/static/css/', 'css');
      const jsBundle = await findLatestBundle('/static/js/', 'js');

      // List of assets to cache (only add if they exist)
      const assetsToCache = [
        cssBundle,
        jsBundle,
      ].filter(Boolean);

      console.log(`[SW] Caching assets: ${assetsToCache}`);

      // Cache assets one by one (skip failures)
      await Promise.all(
        assetsToCache.map(async (url) => {
          try {
            const response = await fetch(url, { cache: 'no-store' });
            if (!response.ok) throw new Error(`Failed to fetch ${url}, status: ${response.status}`);
            await cache.put(url, response);
          } catch (err) {
            console.error(`[SW] Error caching ${url}:`, err);
          }
        })
      );

      self.skipWaiting();
    })()
  );
});

/**
 * Fetch event listener to serve assets from cache.
 */
self.addEventListener('fetch', (event) => {
  if (event.request.method !== 'GET') {
    return; // Only cache GET requests
  }

  event.respondWith(
    (async () => {
      const cache = await caches.open(CACHE_NAME);
      const normalizedRequest = new Request(event.request.url.split('?')[0], {
        method: event.request.method,
        headers: event.request.headers,
      });

      // Check cache first
      const cachedResponse = await cache.match(normalizedRequest);
      if (cachedResponse) {
        console.log(`[SW] Serving from cache: ${event.request.url}`);
        return cachedResponse;
      }

      // Fetch and update cache
      try {
        const response = await fetch(event.request);
        if (isCacheableRequest(event.request, response)) {
          cache.put(normalizedRequest, response.clone());
        }
        return response;
      } catch (err) {
        console.error(`[SW] Fetch failed for ${event.request.url}:`, err);
        return new Response('Network error occurred.', { status: 503 });
      }
    })()
  );
});
