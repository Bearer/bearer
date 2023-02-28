logger.error("something went wrong");

async function run(err) {
  const resp = await this.getResponse({
    data: {
      name: this.name,
      email: this.email
    }
  })

  if (resp.error != nil) {
    console.log(resp.error)
  }

  return data
}
