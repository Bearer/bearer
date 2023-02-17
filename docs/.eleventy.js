const syntaxHighlight = require("@11ty/eleventy-plugin-syntaxhighlight");
const {
  EleventyHtmlBasePlugin,
  EleventyRenderPlugin,
} = require("@11ty/eleventy");
const yaml = require("js-yaml");
const markdownIt = require("markdown-it");
const markdownItEmoji = require("markdown-it-emoji");
const markdownItAnchor = require("markdown-it-anchor");
const pluginTOC = require("eleventy-plugin-toc");
const now = String(Date.now());
const path = require("path");

const mdSetup = markdownIt({ html: true })
  .use(markdownItEmoji)
  .use(markdownItAnchor);

mdSetup.renderer.rules.code_inline = (tokens, idx, { langPrefix = "" }) => {
  const token = tokens[idx];
  return `<code class="${langPrefix}">${mdSetup.utils.escapeHtml(
    token.content
  )}</code>`;
};

module.exports = function (eleventyConfig) {
  eleventyConfig.addWatchTarget("./_src/styles/tailwind.config.js");
  eleventyConfig.addWatchTarget("./_src/styles/tailwind.css");
  eleventyConfig.addWatchTarget("./_src/js/*.js");
  eleventyConfig.addPassthroughCopy("assets/img");
  eleventyConfig.addPassthroughCopy("assets/fonts");
  eleventyConfig.addPassthroughCopy({
    "./_src/styles/prism-theme.css": "./prism-theme.css",
  });
  eleventyConfig.addPassthroughCopy({ "./_src/js/app.js": "./app.js" });
  eleventyConfig.addPassthroughCopy({ "./_tmp/style.css": "./style.css" });
  eleventyConfig.addPassthroughCopy({ "./robots.txt": "./robots.txt" });
  eleventyConfig.addDataExtension("yaml", (contents) => yaml.load(contents));
  eleventyConfig.addShortcode("version", function () {
    return now;
  });
  eleventyConfig.setLibrary("md", mdSetup);
  eleventyConfig.addPlugin(EleventyHtmlBasePlugin, {
    baseHref: "/",
  });

  eleventyConfig.addPlugin(EleventyRenderPlugin);
  eleventyConfig.addPlugin(syntaxHighlight);
  eleventyConfig.addPlugin(pluginTOC, {
    wrapper: "nav",
  });

  eleventyConfig.addFilter("sortById", (arr) => {
    arr.sort((a, b) => (a.metadata.id > b.metadata.id ? 1 : -1));
    return arr;
  });
  eleventyConfig.addFilter("setAttribute", (obj, key, value) => {
    obj[key] = value;
    return obj;
  });
  eleventyConfig.addFilter("deduplicate", (arr) => {
    const result = arr.filter(
      (value, index, self) => index === self.findIndex((t) => t.id === value.id)
    );
    return result;
  });

  eleventyConfig.addNunjucksGlobal("navHighlight", (parent, child) => {
    const target = parent.split(path.sep).slice(1, -1);
    const check = child.split(path.sep).slice(1, -1);
    // handles individual rule pages highlighting "rule" in side nav
    const isRule =
      target.includes("rules") && check[check.length - 2] === "rules";
    if (child === parent || isRule) {
      return true;
    } else {
      return false;
    }
  });

  return {
    dir: {
      includes: "_src/_includes",
      output: "_site",
    },
  };
};
