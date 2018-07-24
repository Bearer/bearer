const Koa = require('koa')
const Logger = require('koa-logger')
const cors = require('@koa/cors')
const BodyParser = require('koa-bodyparser')
const respond = require('koa-respond')

const app = new Koa()
app.use(respond())
app.use(Logger())
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

module.exports = app
