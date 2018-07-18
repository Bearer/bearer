const semver = require('semver')
const term = require('terminal-kit').terminal
const updateNotifier = require('update-notifier')

const pkg = require('../../package.json')
const version = pkg.engines.node

// Check if a new version is available
const notifier = updateNotifier({
  pkg,
  updateCheckInterval: 1000 * 60 * 60 * 24 // everyday for the moment
})
notifier.notify()

if (!semver.satisfies(process.version, version)) {
  term.white('Bearer: ')
  term.red(
    `Required node version ${version} not satisfied with current version ${
      process.version
    }.`
  )
  term('\n')
  process.exit(1)
}
