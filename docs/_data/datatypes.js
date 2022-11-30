const { readdir, readFile } = require("node:fs/promises");
const path = require("path");

async function fetchData(dir) {
  try {
    const files = await readdir(dir);
    let result = await Promise.all(
      files.map(async (file) => {
        let data = await readFile(path.join(dir, file), { encoding: "utf8" });
        // console.log(data);
        return JSON.parse(data);
      })
    );
    return result;
  } catch (err) {
    throw err;
  }
}

function sortData(types, cats) {
  let workingData = {};
  let output = [];
  cats.forEach((cat) => {
    workingData[cat.uuid] = { name: cat.name, types: [] };
  });

  types.forEach((item) => {
    workingData[item.category_uuid].types.push(item);
  });

  for (const key in workingData) {
    output.push({ ...workingData[key] });
  }
  return output;
}
// example();
module.exports = async function () {
  let dataTypes = await fetchData("../pkg/classification/db/data_types/");
  let dataCats = await fetchData("../pkg/classification/db/data_categories/");
  return sortData(dataTypes, dataCats);
};
