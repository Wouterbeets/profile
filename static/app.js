// Load theme
if (localStorage.getItem('theme') === 'dark' || (!localStorage.getItem('theme') && window.matchMedia('(prefers-color-scheme: dark)').matches)) {
	document.documentElement.classList.add('dark');
}

// Language switcher
function setLanguage(lang) {
	document.cookie = `language=${lang}; path=/; max-age=31536000`; // 1 year
	window.location.reload();
}
