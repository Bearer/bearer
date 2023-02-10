return asyncIsPermitted().then(function (result) {
  if (result === true) {
    return true
  } else {
    throw `${user.email}`
  }
})
