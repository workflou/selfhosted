document.addEventListener('htmx:beforeSwap', function (evt) {
    if (evt.detail.xhr.status === 400) {
        evt.detail.shouldSwap = true;
        evt.detail.isError = false;
    }

    if (evt.detail.xhr.status === 429) {
        evt.detail.shouldSwap = true;
        evt.detail.isError = false;
    }
});

document.addEventListener('htmx:afterSwap', function (event) {
    _hyperscript.processNode(event.target);
});

document.addEventListener('redirect', function (event) {
    htmx.ajax('GET', event.detail.value, 'body')
});