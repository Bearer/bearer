import * as request from 'request'

const requestPromise = (url: string, method: string, path: string, body: any, headers = {}) =>
  new Promise<request.Response>((resolve, reject) => {
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

export default url => {
  return {
    signup: body => requestPromise(url, 'POST', 'signup', body),
    login: body => requestPromise(url, 'POST', 'login', body),
    refresh: body => requestPromise(url, 'POST', 'refresh_token', body),
    putItem: body => requestPromise(url, 'POST', 'items', body),
    assemblyScenario: (token, body) => requestPromise(url, 'POST', 'deploy', body, { Authorization: token }),
    signedUrls: (token, Keys, type) =>
      requestPromise(url, 'POST', 'signed-urls', { Keys, type }, { Authorization: token }),
    signedUrl: (token, Key, type) => requestPromise(url, 'POST', 'signed-url', { Key, type }, { Authorization: token }),
    deployScenario: (token, eventName, OrgId, scenarioId) =>
      requestPromise(url, 'POST', 'user-notifications', { eventName, OrgId, scenarioId }, { Authorization: token }),
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
