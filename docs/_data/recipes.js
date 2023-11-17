const { readdir, readFile } = require("node:fs/promises");
const path = require("path");

async function fetchFile(location) {
  return readFile(location, { encoding: "utf8" }).then((file) =>
    JSON.parse(file)
  );
}

async function fetchData(dir) {
  const files = await readdir(dir);
  let result = await Promise.all(
    files.map(async (file) => {
      let data = await fetchFile(path.join(dir, file));
      return {
        ...data,
        id: path.basename(file, ".json"),
        source: `/internal/classification/db/recipes/${file}`,
      };
    })
  );
  return result;
}
module.exports = async function () {
  let recipes = await fetchData("../internal/classification/db/recipes/");
  return recipes;
};
