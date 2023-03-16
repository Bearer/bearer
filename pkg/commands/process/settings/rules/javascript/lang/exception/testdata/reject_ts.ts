function main() {
  const result = new Promise(function (resolve, reject) {
    setTimeout(() => {
      const user = {
        email: current_user.email,
      }
      reject("Error with user " + user.email)
    }, 2000)
  })

  return new Promise((resolve, reject) => {
    setTimeout(() => {
      const user = {
        email: current_user.email,
      }
      reject("Error with user " + user.email)
    }, 2000)
  })
}
