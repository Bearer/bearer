workflow "Build, test and publish" {
  on = "push"
  resolves = ["Publish"]
}

action "GitHub Action for npm" {
  uses = "actions/npm@6309cd9"
  secrets = ["NPM_TOKEN"]
  args = "install"
}

action "Lerna bootstrap" {
  uses = "actions/npm@6309cd9"
  needs = ["GitHub Action for npm"]
  secrets = ["NPM_TOKEN"]
  runs = "lerna bootstrap"
}

action "Run test" {
  uses = "actions/npm@6309cd9"
  needs = ["Lerna bootstrap"]
  args = "test"
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
