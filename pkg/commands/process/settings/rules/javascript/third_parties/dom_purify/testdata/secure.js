const config = { RETURN_DOM_FRAGMENT: true };
let html = DOMPurify.sanitize(dirty);
document.body.innerHTML = html;
