const { readdir, readFile } = require("node:fs/promises");
const path = require("path");
const PDID = "e1d3135b-3c0f-4b55-abce-19f27a26cbb3";

async function fetchFile(location) {
  return readFile(location, { encoding: "utf8" }).then((file) =>
    JSON.parse(file)
  );
}

async function fetchData(dir) {
  const files = await readdir(dir);
  let result = await Promise.all(
    files.map(async (file) => {
      return await fetchFile(path.join(dir, file));
    })
  );
  return result;
}

function sortData(typesFile, catsFile, groupsFile) {
  let output = {};
  let counts = {
    types: typesFile.length,
  };

  // setup groups
  // makes output[key] per group where key is UUID of group(PD, pii, etc)
  for (const key in groupsFile.groups) {
    output[key] = {
      uuid: key,
      categories: {},
      ...groupsFile.groups[key],
    };
  }

  // add categories to each group
  for (const key in groupsFile.category_mapping) {
    for (const groupUUID of groupsFile.category_mapping[key].group_uuids) {
      output[groupUUID].categories[key] = {
        types: [],
        uuid: key,
        ...groupsFile.category_mapping[key],
      };
      // if group has personal data as a parent, also add category to parent
      // note: update logic when needed in future where multiple parents exist
      if (output[groupUUID].parent_uuids.includes(PDID)) {
        output[PDID].categories[key] = {
          types: [],
          uuid: key,
          ...groupsFile.category_mapping[key],
        };
      }
    }
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
  let dataTypes = await fetchData("../internal/classification/db/data_types/");
  let dataCats = await fetchData("../internal/classification/db/data_categories/");
  let groupings = await fetchFile(
    "../internal/classification/db/category_grouping.json"
  );
  return sortData(dataTypes, dataCats, groupings);
};
