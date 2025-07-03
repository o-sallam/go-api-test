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

document.addEventListener('DOMContentLoaded', function () {
    const cache = {};
    const main = document.querySelector('main');
    if (!main) return;
    // Delegate for prev/next buttons
    document.body.addEventListener('mousemove', function (e) {
        const link = e.target.closest('.article-nav-card a');
        if (link && link.getAttribute('href') && link.getAttribute('href') !== '#') {
            const slug = link.getAttribute('href').replace(/^\//, '');
            prefetchArticle(slug, cache);
        }
    });
    document.body.addEventListener('click', function (e) {
        const link = e.target.closest('.article-nav-card a');
        if (link && link.getAttribute('href') && link.getAttribute('href') !== '#') {
            e.preventDefault();
            const slug = link.getAttribute('href').replace(/^\//, '');
            if (cache[slug]) {
                main.innerHTML = cache[slug];
                window.history.pushState({}, '', '/' + slug);
            } else {
                fetch(`/post-partial-html/${slug}`)
                    .then(res => res.ok ? res.text() : null)
                    .then(html => {
                        if (html) {
                            cache[slug] = html;
                            main.innerHTML = html;
                            window.history.pushState({}, '', '/' + slug);
                        }
                    });
            }
        }
    });
    // Handle browser back/forward
    window.addEventListener('popstate', function () {
        const slug = location.pathname.replace(/^\//, '');
        if (slug && cache[slug]) {
            main.innerHTML = cache[slug];
        } else if (slug) {
            fetch(`/post-partial-html/${slug}`)
                .then(res => res.ok ? res.text() : null)
                .then(html => {
                    if (html) {
                        cache[slug] = html;
                        main.innerHTML = html;
                    }
                });
        }
    });
}); 