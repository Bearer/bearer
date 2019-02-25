workflow "build + test" {
  on = "push"
  resolves = ["test"]
}

action "setup" {
  uses = "docker://node:10"
  runs = "yarn"
  args = "install --frozen-lockfile"
}

action "test" {
  uses = "docker://node:10"
  needs = "setup"
  args = "yarn"
  runs = "test"
}