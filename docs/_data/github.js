const EleventyFetch = require("@11ty/eleventy-fetch");

module.exports = async function () {
  let repo = await EleventyFetch("https://api.github.com/repos/bearer/curio", {
    duration: "1d",
    type: "json",
  });

  let release = await EleventyFetch(
    "https://api.github.com/repos/bearer/curio/releases/latest",
    {
      duration: "1d",
      type: "json",
    }
  );
  console.log(repo.stargazers, release.tag_name);
  return {
    stargazers: repo.stargazers_count,
    version: release.tag_name,
  };
};
