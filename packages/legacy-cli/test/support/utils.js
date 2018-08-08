const { spawn, exec } = require('child_process')
const unzip = require('unzip-stream')

const { Transform, Writable } = require('stream')
const fs = require('fs')
const fx = require('mkdir-recursive')
const path = require('path')

class StreamCache extends Writable {
  constructor() {
    super([arguments])
    this.data = ''
  }

  write(chunk, encoding, callback) {
    this.data = this.data + chunk.toString()
    if (typeof callback === 'function') {
      callback()
    }
  }
}

const extractContent = (readStream, fileName, stream) =>
  readStream
    .pipe(unzip.Parse())
    .on('error', err => {
      // We don't really care in which state we will end up
      console.log(err.toString())
    })
    .pipe(
      Transform({
        objectMode: true,
        transform(entry, e, cb) {
          const filePath = entry.path
          if (filePath === fileName) {
            // TODO what will happen if we don't find file?
            entry.pipe(stream).on('finish', cb)
          } else {
            entry.autodrain()
            cb()
          }
        }
      })
    )

const makeTmpDir = () => {
  const randomDirName = Math.random()
    .toString(36)
    .substring(2, 15)

  const dir = `../tmp/bearer-cli-${randomDirName}`
  fx.mkdirSync(dir)
  return path.resolve(dir)
}

const rmTmpDir = dir => {
  exec(`rm -rf ${dir}`, (err, stdout, stderr) => {
    if (err) console.error(err)
    if (stdout !== '') console.log(stdout)
    if (stderr !== '') console.log(stderr)
  })
}

module.exports = {
  StreamCache,
  extractContent,
  makeTmpDir,
  rmTmpDir
}
