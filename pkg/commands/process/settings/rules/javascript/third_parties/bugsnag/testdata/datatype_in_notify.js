riskyCode(() => {
  try {
    // risky business
  } catch (e) {
    Bugsnag.notify(user.ip_address + " : " + e)
  }
})
