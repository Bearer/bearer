async function asyncIsOk() {
  const result = await this.asyncCall({
    name: this.name,
    email: this.email
  })

  if (result.error != null) {
    throw result.error
  }
}