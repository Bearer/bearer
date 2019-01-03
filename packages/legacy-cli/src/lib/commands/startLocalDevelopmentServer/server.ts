import * as Koa from 'koa'
import * as cors from '@koa/cors'
import * as BodyParser from 'koa-bodyparser'
import * as respond from 'koa-respond'

const app = new Koa()
app.use(respond())
app.use(
  cors({
    credentials: true
  })
)
app.use(
  BodyParser({
    enableTypes: ['json'],
    jsonLimit: '5mb',
    strict: true,
    onerror: function(err, ctx) {
      ctx.throw('body parse error', 422)
    }
  })
)

export default app
