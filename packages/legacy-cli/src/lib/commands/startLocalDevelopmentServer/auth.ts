const Router = require('koa-router')
import * as fs from 'fs'

const router = new Router({ prefix: '/' })

router.get('v1/user/initialize', ctx => {
  // if (!ctx.request.signedCookies || !ctx.request.signedCookies.uuid) {
  //   ctx.cookie = {
  //     uuid: uuidv1(),
  //     maxAge: 2592000000,
  //     signed: true
  //   } //expires in one month
  // }
  ctx.body = '<html><head></head><body><script>'
  ctx.body += fs.readFileSync(__dirname + '/iframe.js')
  ctx.body += '</script></body></html>'
})

router.get('v1/auth/:integration_uuid', async ctx => {
  ctx.body = `<html>
  <head>
    <script src="https://cdn.jsdelivr.net/npm/post-robot@8.0.28/dist/post-robot.min.js"></script>
    <script type="application/javascript">
      localStorage.setItem('${ctx.request.query.setupId}|${ctx.params.integration_uuid}', true)
      postRobot.send(window.opener, 'BEARER_AUTHORIZED', {
        scenarioId: '${ctx.params.integration_uuid}'
      })
      setTimeout(function() {
        window.close()
      }, 500)
    </script>
  </head>
  <body></body>
</html>`
})

export default router
