const request = require('request')

const requestPromise = (url, method, path, body, headers = {}) =>
  new Promise((resolve, reject) => {
    request(
      {
        method,
        uri: url + path,
        json: true,
        body,
        headers
      },
      (err, res) => {
        if (err) reject(err)
        else resolve(res)
      }
    )
  })

module.exports = url => {
  return {
    signup: body => requestPromise(url, 'POST', 'signup', body),
    login: body => requestPromise(url, 'POST', 'login', body),
    refresh: body => requestPromise(url, 'POST', 'refresh_token', body),
    putItem: body => requestPromise(url, 'POST', 'items', body),
    screensInvalidate: (token, body) =>
      requestPromise(url, 'POST', 'screens-invalidate', body, {
        Authorization: token
      }),
    assemblyScenario: (token, body) =>
      requestPromise(url, 'POST', 'deploy', body, { Authorization: token }),
    signedUrls: (token, Keys, type) =>
      requestPromise(
        url,
        'POST',
        'signed-urls',
        { Keys, type },
        { Authorization: token }
      ),
    signedUrl: (token, Key, type) =>
      requestPromise(
        url,
        'POST',
        'signed-url',
        { Key, type },
        { Authorization: token }
      ),
    deployScenario: (eventName, OrgId, scenarioTitle) =>
      requestPromise(
        url,
        'POST',
        '',
        {
          query: `mutation DeployScenario {
            deployScenario(state: "${eventName}", organizationIdentifier: "${OrgId}", scenarioIdentifier: "${scenarioTitle}") {
              errors {
                field
                messages
              }
              scenario {
                identifier
                state
                name
              }
            }
          }`
        },
        {}
      ),
    upload: (content, headers = {}) =>
      new Promise((resolve, reject) => {
        request(
          {
            method: 'PUT',
            url,
            body: content,
            headers
          },
          (err, data) => {
            if (err) {
              reject(err)
            } else resolve(data)
          }
        )
      })
  }
}
