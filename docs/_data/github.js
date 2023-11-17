const EleventyFetch = require("@11ty/eleventy-fetch");

module.exports = async function () {
  // let repo = await EleventyFetch("https://api.github.com/repos/bearer/bearer", {
  //   duration: "1d",
  //   type: "json",
  // });
  let release = {};
  try {
    release = await EleventyFetch(
      "https://api.github.com/repos/bearer/bearer/releases/latest",
      {
        duration: "60m",
        type: "json",
      },
    );
  } catch (err) {
    console.log("Could not fetch release");
  }
  return {
    // stargazers: repo.stargazers_count,
    release: {
      name: release.tag_name || "DEV",
      url: release.html_url || "/",
    },
  };
};
