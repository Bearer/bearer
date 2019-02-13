const proxy = require('http-proxy-middleware')
module.exports = function expressMiddleware(router) {
  router.use(
    '/build',
    proxy({
      target: 'http://localhost:3333',
      changeOrigin: true
    })
  )

  router.use(
    '/~dev-server',
    proxy({
      target: 'http://localhost:3333',
      changeOrigin: true,
      ws: true
    })
  )
}
