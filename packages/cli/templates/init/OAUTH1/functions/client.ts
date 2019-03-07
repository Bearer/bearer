import { OAuth } from 'oauth'

const config: IConfig = require('../auth.config.json')

interface IConfig {
  requestTokenURL: string
  accessTokenURL: string
  signatureMethod: string
  callbackURL: string
}

export default function({
  consumerKey,
  consumerSecret,
  userAgent = 'Bearer'
}: {
  consumerKey: string
  consumerSecret: string
  userAgent?: string
}) {
  return new OAuth(
    config.requestTokenURL,
    config.accessTokenURL,
    consumerKey,
    consumerSecret,
    '1.0A',
    config.callbackURL,
    config.signatureMethod,
    null,
    {
      Accept: 'application/json',
      'User-Agent': userAgent
    }
  )
}
