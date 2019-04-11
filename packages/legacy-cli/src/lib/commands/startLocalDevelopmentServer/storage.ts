const Router = require('koa-router')
import * as knex from 'knex'
import * as uuidv1 from 'uuid/v1'

const filename = process.env.BEARER_LOCAL_DATABASE || ':memory:'
import debug from '../../logger'
const logger = debug.extend('start')

const db = knex({
  dialect: 'sqlite3',
  connection: {
    filename
  },
  useNullAsDefault: true,
  debug: /bearer:/.test(process.env.DEBUG || '')
})

export async function getRows(referenceId) {
  const rows = await db
    .table('records')
    .select('data')
    .where({ referenceId })
    .limit(1)

  const { data } = rows[0]

  return data
}

export default () => {
  const router = new Router({ prefix: '/api/' })

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

  const getItem = async ctx => {
    const referenceId = ctx.params.referenceId

    try {
      const rows = await db
        .table('records')
        .select('data')
        .where({ referenceId })
        .limit(1)

      const { data } = rows[0]
      logger('%j', data)
      ctx.ok({ Item: { ...JSON.parse(data), referenceId } })
    } catch (e) {
      ctx.notFound(e)
    }
  }

  const deleteItem = async (ctx, _next) => {
    const referenceId = ctx.params.referenceId
    return db
      .table('records')
      .where({ referenceId })
      .delete()
  }

  const putItem = async (ctx, _next) => {
    const referenceId = ctx.params.referenceId

    logger('%j', ctx.request)
    logger('BODY is: %j', ctx.request.body)
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
  }

  const postItem = async (ctx, _next) => {
    const referenceId = uuidv1()

    logger('%j', ctx.request.body)
    try {
      await db.table('records').insert({ data: JSON.stringify(ctx.request.body), referenceId })
      ctx.ok({ Item: { ...ctx.request.body, referenceId } })
    } catch (e) {
      ctx.badRequest(e)
    }
  }
  router.get('v1/items/:referenceId', getItem)
  router.get('v2/items/:referenceId', getItem)
  router.get('v4/items/:referenceId', getItem)

  router.delete('v1/items/:referenceId', deleteItem)
  router.delete('v2/items/:referenceId', deleteItem)
  router.delete('v4/items/:referenceId', deleteItem)

  router.put('v1/items/:referenceId', putItem)
  router.put('v2/items/:referenceId', putItem)
  router.put('v4/items/:referenceId', putItem)

  router.post('v1/items', postItem)
  router.post('v2/items', postItem)
  router.post('v4/items', postItem)

  return router
}
