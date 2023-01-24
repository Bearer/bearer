const { readFile, readdir } = require("node:fs/promises");
const { statSync } = require("fs");
const path = require("path");
const yaml = require("js-yaml");
const rulesPath = "../pkg/commands/process/settings/rules/";
// Target file structure
// lang
//  -name
//  -children [] (ditch this unless needed. Seems unnecessary)
//    - framework
//      - name
//      - children []
//        - yml file
//          - metadata

function isDirectory(dir) {
  const result = statSync(dir);
  return result.isDirectory();
}

async function fetchData(location) {
  let rules = [];
  let groupedRules = [];
  try {
    const dirs = await readdir(location);
    // ex: looping through [ruby, gitleaks, sql]
    dirs.forEach(async (dir) => {
      const dirPath = path.join(rulesPath, dir);
      if (isDirectory(dirPath)) {
        const dirData = {
          name: dir,
          children: [],
        };

        const subDirs = await readdir(dirPath);
        // ex. looping through rules/ruby [lang, rails]
        subDirs.forEach(async (subDir) => {
          // const child = {
          //   name: subDir,
          //   // children: [],
          // };
          const subDirPath = path.join(dirPath, subDir);
          if (isDirectory(subDirPath)) {
            const files = await readdir(subDirPath);
            const children = await fetchAllFiles(subDirPath, files);
            // child.children = [...children];
            // dirData.children.push(child);

            dirData.children.push(...children);
            rules.push(...children);
          }
        });
        groupedRules.push(dirData);
      }
    });
    return { rules, groupedRules };
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
  if (path.extname(location) === ".yml") {
    return readFile(location, { encoding: "utf8" }).then((file) => {
      let out = yaml.load(file);

      return {
        name: path.basename(location, ".yml"),
        ...out,
      };
    });
  }
}

module.exports = async function () {
  return await fetchData("../pkg/commands/process/settings/rules/");

  // return await fetchFile("../pkg/commands/process/settings/policies.yml");
};
