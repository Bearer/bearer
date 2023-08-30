const { readFile, readdir } = require("node:fs/promises");
const { statSync } = require("fs");
const EleventyFetch = require("@11ty/eleventy-fetch");
const path = require("path");
const yaml = require("js-yaml");
const cweList = require("./cweList.json");
const gitly = require("gitly");
const source = "bearer/bearer-rules";
const rulesPath = "_tmp/rules-data";
const excludeDirectories = [".github", "scripts"];

let counts = {
  languages: {},
};

function updateCounts(lang, framework = null, id = null) {
  if (framework !== null) {
    if (counts.languages[lang].frameworks[framework]) {
      counts.languages[lang].frameworks[framework].count++;
      counts.languages[lang].frameworks[framework].rules.push(id);
    } else if (framework !== "lang") {
      newFramework(lang, framework);
    } else if (framework === "lang") {
      counts.languages[lang].baseRules.push(id);
    }
    counts.languages[lang].count++;
  } else {
    newLang(lang);
  }
}

function newFramework(lang, name) {
  counts.languages[lang].frameworks[name] = {
    name,
    count: 1,
    rules: [],
  };
}

function newLang(name) {
  counts.languages[name] = {
    name,
    count: 0,
    baseRules: [],
    frameworks: {},
  };
}

function isDirectory(dir) {
  const result = statSync(dir);
  return result.isDirectory();
}

async function fetchRelease() {
  let latest = {};
  try {
    latest = await EleventyFetch(
      `https://api.github.com/repos/${source}/releases/latest`,
      {
        duration: "60m",
        type: "json",
      }
    );
  } catch (e) {
    console.error(e);
  }
  try {
    let src = await gitly.download(`${source}#${latest.tag_name}`);
    await gitly.extract(src, rulesPath);
  } catch (e) {
    throw console.error(e);
  }
}

async function fetchData(location) {
  let rules = [];
  try {
    const dirs = await readdir(location);
    // ex: looping through rules [ruby, gitleaks, sql]
    dirs.forEach(async (dir) => {
      const dirPath = path.join(location, dir);
      if (isDirectory(dirPath) && !excludeDirectories.includes(dir)) {
        const subDirs = await readdir(dirPath);
        updateCounts(dir);
        // ex. looping through rules/ruby [lang, rails]
        subDirs.forEach(async (subDir) => {
          const subDirPath = path.join(dirPath, subDir);
          if (isDirectory(subDirPath) && subDir !== "shared") {
            const files = await readdir(subDirPath);
            const children = await fetchAllFiles(
              subDirPath,
              path.join(dir, subDir),
              files
            );
            rules.push(...children);
          }
        });
      }
    });
    return { counts, rules };
  } catch (err) {
    throw err;
  }
}

async function fetchAllFiles(directory, breadcrumb, files) {
  let result = await Promise.all(
    files.reduce((all, file) => {
      const location = path.join(directory, file);
      if (path.extname(location) === ".yml") {
        return [...all, fetchFile(path.join(directory, file), breadcrumb)];
      } else {
        return all;
      }
    }, [])
  );
  return result;
}

async function fetchFile(location, breadcrumb) {
  return readFile(location, { encoding: "utf8" }).then((file) => {
    let out = yaml.load(file);
    let owasps = new Set();
    let subdir = breadcrumb.split("/");
    let framework = subdir[subdir.length - 1];
    updateCounts(subdir[subdir.length - 2], framework, out.metadata.id);
    if (out.metadata.cwe_id) {
      out.metadata.cwe_id.forEach((i) => {
        if (cweList[i].owasp) {
          owasps.add(cweList[i].owasp.id);
        }
      });
    }
    return {
      name: path.basename(location, ".yml"),
      location: path.join(breadcrumb, path.basename(location, ".yml")),
      owasp_ids: [...owasps].sort(),
      framework,
      ...out,
    };
  });
}

module.exports = async function () {
  await fetchRelease();
  return await fetchData(path.join(rulesPath, "rules"));
};
