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

const mdSetup = markdownIt({ html: true })
  .use(markdownItEmoji)
  .use(markdownItAnchor);

mdSetup.renderer.rules.code_inline = (tokens, idx, { langPrefix = "" }) => {
  const token = tokens[idx];
  return `<code class="code-inline ${langPrefix}">${token.content}</code>`;
};

module.exports = function (eleventyConfig) {
  eleventyConfig.addWatchTarget("./_src/styles/tailwind.config.js");
  eleventyConfig.addWatchTarget("./_src/styles/tailwind.css");
  eleventyConfig.addPassthroughCopy("assets/img");
  eleventyConfig.addPassthroughCopy("assets/fonts");
  eleventyConfig.addPassthroughCopy({
    "./_src/styles/prism-theme.css": "./prism-theme.css",
  });
  eleventyConfig.addPassthroughCopy({ "./_tmp/style.css": "./style.css" });
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

  return {
    dir: {
      includes: "_src/_includes",
      output: "_site",
    },
  };
};
