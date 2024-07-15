const EleventyFetch = require("@11ty/eleventy-fetch")

module.exports = async function () {
  // let repo = await EleventyFetch("https://api.github.com/repos/bearer/bearer", {
  //   duration: "1d",
  //   type: "json",
  // });
  let release = {}
  let rulesRelease = {}
  try {
    release = await EleventyFetch(
      "https://api.github.com/repos/bearer/bearer/releases/latest",
      {
        duration: "60m",
        type: "json",
      },
    )
  } catch (err) {
    console.log("Could not fetch release")
    if (process.env.ELEVENTY_PRODUCTION) {
      throw err
    }
  }
  try {
    rulesRelease = await EleventyFetch(
      "https://api.github.com/repos/bearer/bearer-rules/releases/latest",
      {
        duration: "60m",
        type: "json",
      },
    )
  } catch (err) {
    console.log("Could not fetch rulesRelease")
    if (process.env.ELEVENTY_PRODUCTION) {
      throw err
    }
  }
  return {
    // stargazers: repo.stargazers_count,
    release: {
      name: release.tag_name || "DEV",
      url: release.html_url || "/",
    },
    rules: {
      name: rulesRelease.tag_name || "DEV_RULES",
    },
  }
}
