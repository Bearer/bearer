const Router = require('koa-router')
import * as knex from 'knex'
import * as uuidv1 from 'uuid/v1'

export default () => {
  const router = new Router({ prefix: '/api/v1/' })

  const filename = process.env.BEARER_LOCAL_DATABASE || ':memory:'
  const debug = process.env.BEARER_DEBUG === '*'

  const db = knex({
    dialect: 'sqlite3',
    connection: {
      filename
    },
    useNullAsDefault: true,
    debug
  })

  db.schema
    .hasTable('records')
    .then(async exists => {
      if (!exists) {
        await db.schema.createTable('records', table => {
          table.index(['referenceId'], 'ref_id_idx')
          table.string('referenceId')
          table.text('data')
        })
      }
    })
    .catch(() =>
      db.schema.createTable('records', table => {
        table.index(['referenceId'], 'ref_id_idx')
        table.string('referenceId')
        table.text('data')
      })
    )

  router.get('items/:referenceId', async ctx => {
    const referenceId = ctx.params.referenceId

    try {
      const rows = await db
        .table('records')
        .select('data')
        .where({ referenceId })
        .limit(1)

      const { data } = rows[0]
      console.log(data)
      ctx.ok({ Item: { ...JSON.parse(data), referenceId } })
    } catch (e) {
      ctx.notFound(e)
    }
  })
  router.delete('items/:referenceId', async (ctx, next) => {
    const referenceId = ctx.params.referenceId
    return await db
      .table('records')
      .where({ referenceId })
      .delete()
  })

  router.put('items/:referenceId', async (ctx, next) => {
    const referenceId = ctx.params.referenceId

    console.log(ctx.request)
    console.log('BODY is: ', ctx.request.body)
    await db
      .table('records')
      .where({ referenceId })
      .delete()
      .then(async () => {
        try {
          await db.table('records').insert({ data: JSON.stringify(ctx.request.body), referenceId })
          ctx.ok({ Item: { ...ctx.request.body, referenceId } })
        } catch (e) {
          ctx.badRequest(e)
        }
      })
  })

  router.post('items', async (ctx, next) => {
    const referenceId = uuidv1()

    console.log(ctx.request.body)
    try {
      await db.table('records').insert({ data: JSON.stringify(ctx.request.body), referenceId })
      ctx.ok({ Item: { ...ctx.request.body, referenceId } })
    } catch (e) {
      ctx.badRequest(e)
    }
  })
  return router
}
