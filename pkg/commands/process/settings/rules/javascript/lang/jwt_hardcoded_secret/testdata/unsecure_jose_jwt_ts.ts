import * as jose from 'jose'

const secret = "my-hardcoded-secret"

const jwt = await new jose.SignJWT({ 'urn:example:claim': true })
  .setIssuedAt()
  .setExpirationTime('2h')
  .sign(secret)

const jwt2 = await new jose.SignJWT().sign(config.secret)

console.log(jwt)