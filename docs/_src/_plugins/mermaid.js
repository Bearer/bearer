const htmlencode = require("htmlencode");
module.exports = (eleventyConfig, options) => {
  const highlighter = eleventyConfig.markdownHighlighter;
  const defaults = {
    startOnLoad: true,
    theme: "base",
    flowchart: {
      curve: "stepAfter",
    },
  };
  eleventyConfig.addMarkdownHighlighter((str, language) => {
    if (language === "mermaid") {
      return `<pre class="mermaid">${htmlencode.htmlEncode(str)}</pre>`;
    }
    if (highlighter) {
      return highlighter(str, language);
    }
    return `<pre class="${language}">${str}</pre>`;
  });
  eleventyConfig.addShortcode("mermaid_js", () => {
    const mermaid_config = {
      ...defaults,
      ...options,
    };
    let src = "https://cdn.jsdelivr.net/npm/mermaid@10.1.0/+esm";
    return `<script type="module" async>import mermaid from "${src}"; document.addEventListener('DOMContentLoaded', mermaid.initialize(${JSON.stringify(
      mermaid_config
    )}));</script>`;
  });
};
