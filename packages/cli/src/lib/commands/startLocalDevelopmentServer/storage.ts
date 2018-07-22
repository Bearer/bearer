const Router = require('koa-router')
import sqlite3 from 'sqlite3'
import * as uuidv1 from 'uuid/v1'
const router = new Router({ prefix: '/api/v1/' })
var db = new sqlite3.verbose().Database(':memory:')

db.run('CREATE TABLE records (data TEXT, refernceId TEXT);')
db.run('CREATE INDEX records_idx ON records (referenceId);')

router.delete('/api/v1/items/:reference_id', (ctx, next) => {
  const referenceId = ctx.params.reference_id
  db.run(`DELETE FROM records WHERE referenceId='${referenceId}'`)
})

router.post('/api/v1/items', (ctx, next) => {
  const referenceId = uuidv1()

  const stmnt = db.prepare('INSERT INTO records VALUES (?, ?)')
  stmnt.run(JSON.stringify(ctx.request.body), referenceId)
  stmnt.finalize()
  ctx.ok({ Item: { ...ctx.request.body, referenceId } })
})

router.put('/api/v1/items/:referenceId', (ctx, next) => {
  const referenceId = ctx.params.reference_id

  const stmnt = db.prepare(
    'INSERT INTO records VALUES (?, ?) ON CONFLICT REPLACE'
  )
  stmnt.run(JSON.stringify(ctx.request.body), referenceId)
  stmnt.finalize()
  ctx.ok({ Item: { ...ctx.request.body, referenceId } })
})

router.get('/api/v1/items/:referenceId', (ctx, next) => {
  const referenceId = ctx.params.reference_id

  db.get('SELECT data FROM records', function(err, row) {
    if (err) ctx.error(err)
    ctx.ok({ Item: { ...JSON.parse(row.data), referenceId } })
  })
})
export default router
