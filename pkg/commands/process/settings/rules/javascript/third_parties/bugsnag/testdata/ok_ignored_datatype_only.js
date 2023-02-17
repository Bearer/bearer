Bugsnag.start({
  onError: function (e) {
    e.setUser(user.id)
    e.addMetadata('page', {
      url: currentUrl(),
    })
  }
})