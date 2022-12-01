const { readdir, readFile } = require("node:fs/promises");
const path = require("path");

async function fetchFile(location) {
  return readFile(location, { encoding: "utf8" }).then((file) =>
    JSON.parse(file)
  );
}

async function fetchData(dir) {
  try {
    const files = await readdir(dir);
    console.log(files.length);
    let result = await Promise.all(
      files.map(async (file) => {
        return await fetchFile(path.join(dir, file));
      })
    );
    return result;
  } catch (err) {
    throw err;
  }
}

function sortData(typesFile, catsFile, groupsFile) {
  let output = {};
  let counts = {
    types: typesFile.length,
  };

  // setup groups
  for (const key in groupsFile.groups) {
    output[key] = { uuid: key, name: groupsFile.groups[key], categories: {} };
  }

  // add categories to each group
  for (const key in groupsFile.category_mapping) {
    output[groupsFile.category_mapping[key].group_uuid].categories[key] = {
      types: [],
      uuid: key,
      ...groupsFile.category_mapping[key],
    };
  }

  // add types to each category
  // output[group uuid].categories[category uuid].types[]
  // note: inefficient, needs rewrite
  for (const key in output) {
    typesFile.forEach((item) => {
      if (output[key].categories[item.category_uuid]) {
        output[key].categories[item.category_uuid].types.push(item);
      }
    });
  }

  return { output, counts };
}
// example();
module.exports = async function () {
  let dataTypes = await fetchData("../pkg/classification/db/data_types/");
  let dataCats = await fetchData("../pkg/classification/db/data_categories/");
  let groupings = await fetchFile(
    "../pkg/classification/db/category_grouping.json"
  );
  return sortData(dataTypes, dataCats, groupings);
};
