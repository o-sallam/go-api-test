// Prefetch and swap article HTML on prev/next navigation
function prefetchArticle(slug, cache) {
    if (!slug || slug === "#" || slug === "") return;
    if (cache[slug]) return; // already prefetched
    fetch(`/post-partial-html/${slug}`)
        .then(res => res.ok ? res.text() : null)
        .then(html => {
            if (html) cache[slug] = html;
        });
}

function showSpinner() {
    const spinner = document.getElementById('loading-spinner');
    if (spinner) spinner.style.display = 'flex';
}

function hideSpinner() {
    const spinner = document.getElementById('loading-spinner');
    if (spinner) spinner.style.display = 'none';
}

function smoothScrollToTop() {
    window.scrollTo({
        top: 0,
        behavior: 'smooth'
    });
}

document.addEventListener('DOMContentLoaded', function () {
    const cache = {};
    const main = document.querySelector('main');
    if (!main) return;
    
    // Prefetch only once per link on mouseenter
    document.body.addEventListener('mouseenter', function (e) {
        const link = e.target.closest('a[href^="/"]');
        if (link && link.getAttribute('href') && link.getAttribute('href') !== '#') {
            const slug = link.getAttribute('href').replace(/^\//, '');
            prefetchArticle(slug, cache);
        }
    }, true); // useCapture=true to catch events on bubbling
    
    // Handle clicks on article cards and prev/next buttons
    document.body.addEventListener('click', function (e) {
        const link = e.target.closest('a[href^="/"]');
        if (link && link.getAttribute('href') && link.getAttribute('href') !== '#') {
            e.preventDefault();
            const slug = link.getAttribute('href').replace(/^\//, '');
            if (cache[slug]) {
                main.innerHTML = cache[slug];
                window.history.pushState({}, '', '/' + slug);
                smoothScrollToTop();
            } else {
                showSpinner();
                fetch(`/post-partial-html/${slug}`)
                    .then(res => res.ok ? res.text() : null)
                    .then(html => {
                        hideSpinner();
                        if (html) {
                            cache[slug] = html;
                            main.innerHTML = html;
                            window.history.pushState({}, '', '/' + slug);
                            smoothScrollToTop();
                        }
                    })
                    .catch(() => {
                        hideSpinner();
                    });
            }
        }
    });
    
    // Handle browser back/forward
    window.addEventListener('popstate', function () {
        const slug = location.pathname.replace(/^\//, '');
        if (slug && cache[slug]) {
            main.innerHTML = cache[slug];
            smoothScrollToTop();
        } else if (slug) {
            showSpinner();
            fetch(`/post-partial-html/${slug}`)
                .then(res => res.ok ? res.text() : null)
                .then(html => {
                    hideSpinner();
                    if (html) {
                        cache[slug] = html;
                        main.innerHTML = html;
                        smoothScrollToTop();
                    }
                })
                .catch(() => {
                    hideSpinner();
                });
        }
    });
}); 