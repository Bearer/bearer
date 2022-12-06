const { readFile } = require("node:fs/promises");
const yaml = require("js-yaml");

async function fetchFile(location) {
  return readFile(location, { encoding: "utf8" }).then((file) =>
    yaml.load(file)
  );
}

module.exports = async function () {
  return await fetchFile("../pkg/commands/process/settings/policies.yml");
};
