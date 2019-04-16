import * as http from 'http'
import * as fs from 'fs'
import * as path from 'path'
import * as url from 'url'

// @ts-ignore
import * as Router from 'router'
import * as finalhandler from 'finalhandler'
// import * as finalhandler = require('finalhandler')
// var http = require('http')
// var Router = require('router')

const iframePath = path.join(__dirname, '../../static/iframe.js')

export default (command: any) => {
  const router = Router()

  /**
   * API
   */
  const api = Router()
  api.use((_: any, res: http.ServerResponse, next: any) => {
    res.setHeader('Content-Type', 'application/json')
    next()
  })

  api.post(
    '/v4/backend/functions/:integration/:functionName',
    (req: http.IncomingMessage, res: http.ServerResponse) => {
      res.end('{}')
    }
  )
  api.post('/v4/functions/:int/:functionName', (req: http.IncomingMessage, res: http.ServerResponse) => {})
  api.post('/backend/:functionName', (req: http.IncomingMessage, res: http.ServerResponse) => {})
  api.post('/:functionName', (req: http.IncomingMessage, res: http.ServerResponse) => {})

  /**
   * Auth
   */
  const auth = Router()

  auth.use((_: any, res: http.ServerResponse, next: any) => {
    res.setHeader('Content-Type', 'text/html; charset=utf-8')
    next()
  })
  auth.get('/initialize', (req: http.IncomingMessage, res: http.ServerResponse) => {
    let body = '<html><head></head><body><script>'
    body += fs.readFileSync(iframePath)
    body += '</script></body></html>'
    res.end(body)
  })

  auth.get('/:integration', (req: http.IncomingMessage, res: http.ServerResponse) => {
    const query = url.parse(req.url!, true).query
    const setupId = query.setupId || 'missing'
    const authId = query.authId || 'missing'
    const integration = (req as any).params.integration
    res.end(`<html>
    <head>
      <script src="https://cdn.jsdelivr.net/npm/post-robot@8.0.28/dist/post-robot.min.js"></script>
      <script type="application/javascript">
        window.LOG_LEVEL = 'error'
        localStorage.setItem('${setupId}|${integration}', true)
        postRobot.send(window.opener, 'BEARER_AUTHORIZED', {
          integrationId: '${integration}',
          authId: '${authId}'
        })
        setTimeout(function() {
          window.close()
        }, 500)
      </script>
    </head>
    <body></body>
  </html>`)
  })

  router.use('/api', api)
  router.use('/v2/user', auth)

  return (req: http.IncomingMessage, res: http.ServerResponse) => {
    router(req, res, finalhandler(req, res))
  }
}
