const { readFile, readdir } = require("node:fs/promises");
const { statSync } = require("fs");
const fetch = require("node-fetch");
const path = require("path");
const yaml = require("js-yaml");
const cweList = require("./cweList.json");
const gitly = require("gitly");
const source = "bearer/bearer-rules";
const rulesPath = "_tmp/rules-data";
const languageDirectories = ["ruby", "javascript"];

function isDirectory(dir) {
  const result = statSync(dir);
  return result.isDirectory();
}

async function fetchRelease() {
  const latest = await fetch(
    `https://api.github.com/repos/${source}/releases/latest`
  )
    .then((res) => res.json())
    .then((data) => {
      return data.tag_name;
    });
  try {
    let src = await gitly.download(`${source}#${latest}`);
    await gitly.extract(src, rulesPath);
  } catch (e) {
    throw console.error(e);
  }
}

async function fetchData(location) {
  let rules = [];
  let groupedRules = [];
  try {
    const dirs = await readdir(location);
    // ex: looping through rules [ruby, gitleaks, sql]
    dirs.forEach(async (dir) => {
      const dirPath = path.join(rulesPath, dir);
      if (isDirectory(dirPath) && languageDirectories.includes(dir)) {
        const dirData = {
          name: dir,
          children: [],
        };

        const subDirs = await readdir(dirPath);
        // ex. looping through rules/ruby [lang, rails]
        subDirs.forEach(async (subDir) => {
          const subDirPath = path.join(dirPath, subDir);
          if (isDirectory(subDirPath)) {
            const files = await readdir(subDirPath);
            const children = await fetchAllFiles(subDirPath, files);
            dirData.children.push(...children);
            rules.push(...children);
          }
        });
        groupedRules.push(dirData);
      }
    });
    return rules;
  } catch (err) {
    throw err;
  }
}

async function fetchAllFiles(directory, files) {
  let result = await Promise.all(
    files.reduce((all, file) => {
      const location = path.join(directory, file);
      if (path.extname(location) === ".yml") {
        return [...all, fetchFile(path.join(directory, file))];
      } else {
        return all;
      }
    }, [])
  );
  return result;
}

async function fetchFile(location) {
  return readFile(location, { encoding: "utf8" }).then((file) => {
    let out = yaml.load(file);
    let owasps = new Set();
    if (out.metadata.cwe_id) {
      out.metadata.cwe_id.forEach((i) => {
        if (cweList[i].owasp) {
          owasps.add(cweList[i].owasp.id);
        }
      });
    }
    return {
      name: path.basename(location, ".yml"),
      location: location.substring(2),
      owasp_ids: [...owasps].sort(),
      ...out,
    };
  });
}

module.exports = async function () {
  await fetchRelease();
  return await fetchData(rulesPath);
};
