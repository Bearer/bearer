// import logger from '@bearer/logger'

// export default logger('js')

// // TODO: make it use @bearer/logger
const logger = (...args: any[]) => {
  console.debug('bearer:js', ...args)
}

logger.extend = (..._args: any[]) => logger
export default logger
