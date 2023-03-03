const config = {};
const html = DOMPurify.sanitize(dirty, config);
document.body.innerHTML = html;
