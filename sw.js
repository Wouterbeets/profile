self.addEventListener('install', (event) => {
	event.waitUntil(
		caches.open('portfolio-v1').then((cache) => {
			return cache.addAll([
				'/',
				'/static/styles.css',
				'/static/app.js',
				'/manifest.json'
			]);
		})
	);
});

self.addEventListener('fetch', (event) => {
	event.respondWith(
		caches.match(event.request).then((response) => {
			return response || fetch(event.request);
		})
	);
});
