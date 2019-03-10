workflow "build + test" {
  on = "push"
  resolves = ["test"]
}

action "setup" {
  uses = "docker://bearerhub/node-10.9:0.9"
  runs = "yarn"
  args = "install --frozen-lockfile"
  

}

action "bootstrap" {
  uses = "docker://bearerhub/node-10.9:0.9"
  needs = "setup"
  # runs = "yarn"
  # args = "bootstrap"
  runs = "ls"
  args = "-la packages/intents/node_modules/@bearer"
}

action "test" {
  uses = "docker://bearerhub/node-10.9:0.9"
  needs = "bootstrap"
  runs = "ls"
  args = "-la node_modules/@bearer"
}

# action "test" {
#   uses = "docker://node:10.15"
#   needs = "bootstrap"
#   runs = "yarn"
#   args = "test"
# }
