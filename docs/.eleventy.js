const syntaxHighlight = require("@11ty/eleventy-plugin-syntaxhighlight");
const { EleventyHtmlBasePlugin } = require("@11ty/eleventy");
const yaml = require("js-yaml");
const markdownIt = require("markdown-it");
const markdownItEmoji = require("markdown-it-emoji");
const now = String(Date.now());

const pathPrefix = (module.exports = function (eleventyConfig) {
  eleventyConfig.addWatchTarget("./_src/styles/tailwind.config.js");
  eleventyConfig.addWatchTarget("./_src/styles/tailwind.css");
  eleventyConfig.addPassthroughCopy("assets/img");
  eleventyConfig.addPassthroughCopy("assets/fonts");
  eleventyConfig.addPassthroughCopy({ "./_tmp/style.css": "./style.css" });
  eleventyConfig.addDataExtension("yaml", (contents) => yaml.load(contents));
  eleventyConfig.addShortcode("version", function () {
    return now;
  });

  eleventyConfig.addPlugin(EleventyHtmlBasePlugin, {
    baseHref: process.env.ELEVENTY_PRODUCTION ? "/curio/" : "/",
  });
  eleventyConfig.addPlugin(syntaxHighlight);
  eleventyConfig.setLibrary(
    "md",
    markdownIt({ html: true }).use(markdownItEmoji)
  );

  return {
    dir: {
      includes: "_src/_includes",
      output: "_site",
    },
  };
});
