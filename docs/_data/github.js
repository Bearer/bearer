const EleventyFetch = require("@11ty/eleventy-fetch");

module.exports = async function () {
  // let repo = await EleventyFetch("https://api.github.com/repos/bearer/curio", {
  //   duration: "1d",
  //   type: "json",
  // });

  let release = await EleventyFetch(
    "https://api.github.com/repos/bearer/curio/releases/latest",
    {
      duration: "60m",
      type: "json",
    }
  );
  return {
    // stargazers: repo.stargazers_count,
    release: {
      name: release.tag_name,
      url: release.html_url,
    },
  };
};
