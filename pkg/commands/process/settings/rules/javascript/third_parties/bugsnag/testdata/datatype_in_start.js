Bugsnag.start({
  onError: function (e) {
    e.setUser(user.id, user.email, user.name)
    e.addMetadata('user location', {
      country: user.home_country,
    })
  },
  onSession: function (session) {
    session.setUser(user.email)
  }
})