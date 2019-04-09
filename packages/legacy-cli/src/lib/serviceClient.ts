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
    login: body => requestPromise(url, 'POST', 'login', body)
  }
}
