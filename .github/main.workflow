workflow "build + test" {
  on = "push"
  resolves = ["run test"]
}

action "install" {
  uses = "docker://node:10"
  runs = "yarn"
  args = "install --frozen-lockfile"
}

action "bootstrap" {
  uses = "docker://node:10"
  needs = "install"
  runs = "yarn"
  args = "bootstrap"
}

action "run test" {
  uses = "docker://node:10"
  needs = "bootstrap"
  runs = "yarn"
  args = "test"
}