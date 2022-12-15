const { readFile } = require("node:fs/promises");
const yaml = require("js-yaml");

function convertAndSort(obj) {
  let list = [];
  for (key in obj) {
    if (obj.hasOwnProperty(key)) {
      list.push({
        ...obj[key],
      });
    }
  }
  return list.sort((a, b) => (a.display_id > b.display_id ? 1 : -1));
}

async function fetchFile(location) {
  return readFile(location, { encoding: "utf8" }).then((file) => {
    let out = yaml.load(file);
    return convertAndSort(out);
  });
}

module.exports = async function () {
  return await fetchFile("../pkg/commands/process/settings/policies.yml");
};
