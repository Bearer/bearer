import { Server } from 'http'
import { Socket } from 'net'
import * as express from 'express'
import getPort from 'get-port'

export type TDestroyableServer = Server & { destroy: (cb: () => void) => void }

// We want to stop the server as fast possible and for that we need to have no connection opened
// here we manually register incoming connections and close them when destroy is called
export async function startServer(port: number, app: express.Express) {
  return new Promise<TDestroyableServer>(async (resolve, reject) => {
    const connections: Record<string, Socket> = {}
    try {
      const listenablePort = await getPort({ port })
      if (listenablePort !== port) {
        reject(new UnavailablePort(port))
      } else {
        const server = app.listen(port, () => {
          resolve(server as TDestroyableServer)
        })

        server.on('connection', (conn: Socket) => {
          const key = [conn.remoteAddress, conn.remotePort].join(':')
          connections[key] = conn
          conn.on('close', () => {
            delete connections[key]
          })
        })
        ;(server as TDestroyableServer).destroy = cb => {
          server.close(cb)
          // force closing all connections
          for (const key in connections) {
            connections[key].destroy()
          }
        }
      }
    } catch (e) {
      reject(e)
    }
  })
}

export class UnavailablePort extends Error {
  constructor(port: number) {
    super(`Can not start server on port ${port}`)
  }
}
