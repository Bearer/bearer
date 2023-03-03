const config = {};
const html = SanitizationLib.sanitize(dirty, config);
document.body.innerHTML = html;
