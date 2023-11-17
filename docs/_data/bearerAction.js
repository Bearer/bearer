const EleventyFetch = require("@11ty/eleventy-fetch")
const yaml = require("js-yaml")
module.exports = async function () {
  let action
  try {
    const response = await EleventyFetch(
      "https://raw.githubusercontent.com/Bearer/bearer-action/main/action.yml",
      {
        duration: "60m",
        type: "text",
      },
    )
    action = yaml.load(response)
  } catch (err) {
    console.log("Could not fetch release")
  }
  return {
    ...action,
  }
}
