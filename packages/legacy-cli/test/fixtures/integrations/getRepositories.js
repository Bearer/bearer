const { getRepositories } = require('./github')
const { FetchData } = require('@bearer/functions')

module.exports.action = getRepositories
module.exports.functionType = FetchData
module.exports.functionName = 'getRepositories'
