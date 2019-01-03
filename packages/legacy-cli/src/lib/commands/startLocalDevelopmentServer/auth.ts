import * as fs from 'fs'
import * as path from 'path'
import * as Router from 'koa-router'

const router = new Router({ prefix: '/' })
const auth = '123'
const iframePath = path.join(__dirname, '../../../../static/iframe.js')

router.get('v1/user/initialize', ctx => {
  ctx.body = '<html><head></head><body><script>'
  ctx.body += fs.readFileSync(iframePath)
  ctx.body += '</script></body></html>'
})

router.get('v2/user/initialize', ctx => {
  ctx.body = '<html><head></head><body><script>'
  ctx.body += fs.readFileSync(iframePath)
  ctx.body += '</script></body></html>'
})

const authScenarioV1 = async ctx => {
  ctx.body = `<html>
  <head>
    <script src="https://cdn.jsdelivr.net/npm/post-robot@8.0.28/dist/post-robot.min.js"></script>
    <script type="application/javascript">
      window.LOG_LEVEL = 'error'
      localStorage.setItem('${ctx.request.query.setupId}|${ctx.params.integration_uuid}', true)
      postRobot.send(window.opener, 'BEARER_AUTHORIZED', {
        scenarioId: '${ctx.params.integration_uuid}',
        authId: '${auth}'
      })
      setTimeout(function() {
        window.close()
      }, 500)
    </script>
  </head>
  <body></body>
</html>`
}

const authScenarioV2 = async ctx => {
  ctx.body = `<html>
  <head>
    <script src="https://cdn.jsdelivr.net/npm/post-robot@8.0.28/dist/post-robot.min.js"></script>
    <script type="application/javascript">
      window.LOG_LEVEL = 'error'
      localStorage.setItem('${ctx.request.query.setupId}|${ctx.params.integration_uuid}', true)
      postRobot.send(window.opener, 'BEARER_AUTHORIZED', {
        scenarioId: '${ctx.params.integration_uuid}',
        authId: '${ctx.request.query.authId || auth}'
      })
      setTimeout(function() {
        window.close()
      }, 500)
    </script>
  </head>
  <body></body>
</html>`
}
router.get('v1/auth/:integration_uuid', authScenarioV1)
router.get('v2/auth/:integration_uuid', authScenarioV2)

export default router
