import Koa from 'koa'
import cors from '@koa/cors'
import BodyParser from 'koa-bodyparser'
import respond from 'koa-respond'

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
