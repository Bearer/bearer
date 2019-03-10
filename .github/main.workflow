workflow "build + test" {
  on = "push"
  resolves = ["test"]
}

action "setup" {
  uses = "docker://bearerhub/node-10.9:0.8"
  runs = "yarn"
  args = "install --frozen-lockfile"
  

}

action "bootstrap" {
  uses = "docker://bearerhub/node-10.9:0.8"
  needs = "setup"
  # runs = "yarn"
  # args = "bootstrap"
  runs = "ls"
  args = "-la packages/intents/node_modules/@bearer"
}

action "test" {
  uses = "docker://bearerhub/node-10.9:0.8"
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
