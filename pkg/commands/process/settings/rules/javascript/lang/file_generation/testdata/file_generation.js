const fs = require("fs")

function generateUsername(firstname, surname) {
  return `${firstname[0]}-${surname}`.toLowerCase()
}

const username = generateUsername(user.firstname, user.surname)
const users = [{
  email: user.email,
  first_name: user.firstname,
  username,
}]

fs.writeFile("data.csv", JSON.stringify(users), callback)
fs.writeFile("data.csv", JSON.stringify(users), "utf-8", (err) => {
  if (err) console.log(err)
  else console.log("Data saved")
})
