const { getRepositories } = require('./github')
const { FetchData } = require('@bearer/intents')

module.exports.action = getRepositories
module.exports.intentType = FetchData
module.exports.intentName = 'getRepositories'
