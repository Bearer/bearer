workflow "Build, test and publish" {
  on = "push"
  resolves = ["Publish"]
}

action "GitHub Action for npm" {
  uses = "docker://node:10"
  secrets = ["NPM_TOKEN"]
  args = "install"
  runs = "yarn"
}

action "Lerna bootstrap" {
  uses = "docker://node:10"
  needs = ["GitHub Action for npm"]
  secrets = ["NPM_TOKEN"]
  args = "lerna bootstrap"
  runs = "yarn"
}

action "Run test" {
  uses = "docker://node:10"
  needs = ["Lerna bootstrap"]
  args = "test"
  runs = "yarn"
}

action "Release only" {
  uses = "actions/bin/filter@95c1a3b"
  args = "branch release"
  needs = ["Run test"]
}

action "Publish" {
  uses = "actions/npm@6309cd9"
  needs = ["Release only"]
  runs = "sdqsdqsd"
  args = "--yes "
}
