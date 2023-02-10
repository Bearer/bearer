return asyncIsPermitted().then(function (result) {
  if (result === true) {
    return true
  } else {
    Promise.reject(new PermissionDenied("fail" + user.email))
  }
})
