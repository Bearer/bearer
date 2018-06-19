import { GetCollection } from '@bearer/intents'

export const action = (token: string, params: any, callback: any) => {
  callback({
    collection: [
      'hello world',
      'Bonjour le monde',
      'Witaj świecie',
      'hello ao',
      'こんにちは世界'
    ]
  })
}
export const intentType = GetCollection
export const intentName = 'getHelloWorlds'
