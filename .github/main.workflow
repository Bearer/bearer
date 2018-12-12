workflow "Build, test and publish" {
  on = "push"
  resolves = [
    "docker://node:10-2",
  ]
}

action "Run test" {
  uses = "docker://node:10"
  args = "test"
  runs = "yarn"
}

action "docker://node:10" {
  uses = "docker://node:10"
  runs = "yarn"
  args = "install"
  secrets = ["NPM_TOKEN"]
}

action "docker://node:10-1" {
  uses = "docker://node:10"
  needs = ["docker://node:10"]
  runs = "yarn"
  args = "test"
}

action "Filters for GitHub Actions" {
  uses = "actions/bin/filter@95c1a3b"
  needs = ["docker://node:10-1"]
  args = "branch release"
}

action "docker://node:10-2" {
  uses = "docker://node:10"
  needs = ["Filters for GitHub Actions"]
  runs = "yarn"
  args = "fail"
}
