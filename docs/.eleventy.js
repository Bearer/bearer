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
const htmlencode = require("htmlencode");
const now = String(Date.now());
const path = require("path");
const mermaid = require("./_src/_plugins/mermaid");

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
  eleventyConfig.addPassthroughCopy({
    "./_src/styles/callout.css": "./callout.css",
  });
  eleventyConfig.addPassthroughCopy({ "./_src/js/app.js": "./app.js" });
  eleventyConfig.addPassthroughCopy({
    "./_src/js/rule-search.js": "./rule-search.js",
  });
  eleventyConfig.addPassthroughCopy({ "./_tmp/style.css": "./style.css" });
  eleventyConfig.addPassthroughCopy({ "./robots.txt": "./robots.txt" });
  eleventyConfig.addPassthroughCopy({ "./_redirects": "./_redirects" });
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
  
  // mermaid rendering
  eleventyConfig.addPlugin(mermaid, {
    themeVariables: {
      darkMode: true,
      primaryColor: "#d4bcf8",
      mainBkg: "#1E065F",
      clusterBorder: "#d4bcf8",
      lineColor: "#d4bcf8",
      titleColor: "#fff",
      edgeLabelBackground: "hsl(243,27%,35%)",
      fontSize: "1rem",
      clusterBkg: "transparent",
      secondaryColor: "hsl(243,27%,35%)",
    },
  });

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

  eleventyConfig.addFilter("keysToArr", (data) => {
    return Object.keys(data);
  });
  eleventyConfig.addFilter("rewriteFrameworks", (word) => {
    function updatePhrase(word) {
      const dictionary = {
        rails: "Ruby on Rails",
        javascript: "JavaScript / TypeScript",
        express: "ExpressJS",
        react: "React",
        third_parties: "Third party",
      };

      if (dictionary[word]) {
        return dictionary[word];
      }
      return word;
    }

    if (typeof word === "string") {
      return updatePhrase(word);
    } else if (Array.isArray(word)) {
      let cleaned = word.filter((w) => w !== "third_parties");
      return cleaned.map((w) => updatePhrase(w));
    } else if (typeof word === "object") {
      let cleaned = Object.keys(word).filter((w) => w !== "third_parties");
      return cleaned.map((w) => updatePhrase(w));
    }
    return word;
  });

  eleventyConfig.addNunjucksFilter("packageMap", (name, manager, group) => {
    switch (manager) {
      case "rubygems":
        return `https://rubygems.org/gems/${name}`;
      case "packagist":
        return `https://packagist.org/packages/${name}`;
      case "go":
        return `https://${name}`;
      case "npm":
        return `https://www.npmjs.com/package/${name}`;
      case "pypi":
        return `https://pypi.org/project/${name}`;
      case "maven":
        return `https://mvnrepository.com/artifact/${group}/${name}`;
      case "nuget":
        return `https://www.nuget.org/packages/${name}`;
      default:
        return "/";
    }
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
  
  eleventyConfig.addPairedShortcode("callout", function(content, level = "", format = "html", customLabel = "") {
		if( format === "md" ) {
			content = mdIt.renderInline(content);
		} else if( format === "md-block" ) {
			content = mdIt.render(content);
		}
		let label = "";
		if(customLabel) {
			label = customLabel;
		} else if(level === "info" || level === "error") {
			label = level.toUpperCase() + ":";
		} else if(level === "warn") {
			label = "WARNING:";
		}
		let labelHtml = label ? `<div class="elv-callout-label">${customLabel || label}</div>` : "";
		let contentHtml = (content || "").trim().length > 0 ? `<div class="elv-callout-c">${content}</div>` : "";

		return `<div class="elv-callout${level ? ` elv-callout-${level}` : ""}">${labelHtml}${contentHtml}</div>`;
	});

  return {
    dir: {
      includes: "_src/_includes",
      output: "_site",
    },
  };
};
