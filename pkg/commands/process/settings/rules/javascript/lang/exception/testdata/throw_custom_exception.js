return asyncIsPermitted().then(function (result) {
  if (result === true) {
    return true
  } else {
    throw new PermissionDenied(`Error with ${current_user.email}`)
  }
})
