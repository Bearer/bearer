workflow "Test Toolkit" {
  on       = "push"

  resolves = [
    "Test"
  ]
}

action "Install" {
  uses = "docker://node:10"
  runs = "yarn"
  args = "install --frozen-lockfile"
}

action "Test" {
  uses  = "docker://node:10"
  # postinstall is runing bootstrap
  needs = "Install"
  runs  = "yarn"

  args  = [
    "test"
  ]

  env   = {
    DEBUG = "bearer:*"
  }
}

workflow "Publish - latest" {
  on       = "push"

  resolves = [
    "Publish latest"
  ]
}

action "Install - latest" {
  uses  = "docker://node:10"
  runs  = "yarn"

  needs = [
    "tag-filter latest"
  ]

  args  = [
    "install",
    "--frozen-lockfile"
  ]
}

action "tag-filter latest" {
  uses = "actions/bin/filter@master"
  args = "tag release-latest-*"
}

action "Publish latest" {
  uses    = "docker://node:10"
  runs    = "yarn"

  needs   = [
    "Install - latest"
  ]

  args    = [
    "lerna-publish"
  ]

  secrets = [
    "NPM_TOKEN",
    "GH_TOKEN"
  ]
}

workflow "Publish - next" {
  on       = "push"

  resolves = [
    "Publish next"
  ]
}

action "Install - next" {
  uses  = "docker://node:10"
  runs  = "yarn"

  needs = [
    "tag-filter next"
  ]

  args  = [
    "install",
    "--frozen-lockfile"
  ]
}

action "tag-filter next" {
  uses = "actions/bin/filter@master"
  args = "tag release-next-*"
}

action "Publish next" {
  uses    = "docker://node:10"
  runs    = "yarn"

  needs   = [
    "Install - next"
  ]

  args    = [
    "lerna-publish",
    "--npm-tag=next"
  ]

  secrets = [
    "NPM_TOKEN",
    "GH_TOKEN"
  ]
}

workflow "Publish - canary" {
  on       = "push"

  resolves = [
    "Publish canary"
  ]
}

action "Install - canary" {
  uses  = "docker://node:10"
  runs  = "yarn"

  needs = [
    "tag-filter canary"
  ]

  args  = [
    "install",
    "--frozen-lockfile"
  ]
}

action "tag-filter canary" {
  uses = "actions/bin/filter@master"
  args = "tag release-canary-*"
}

action "Publish canary" {
  uses    = "docker://node:10"
  runs    = "yarn"

  needs   = [
    "Install - canary"
  ]

  args    = [
    "lerna",
    "publish",
    "--canary",
    "--preid",
    "canary"
  ]

  secrets = [
    "NPM_TOKEN",
    "GH_TOKEN"
  ]
}
