const { exec, execSync, spawn, spawnSync } = require('node:child_process');

exports.handler = async (event) => {
  exec("ls -lh /usr", (err, stdout, stderr) => {
    // do something
  });

  execSync("ls -lh /usr", (err, stdout, stderr) => {
    // do something
  });

  spawn("ls -lh /usr");
  spawnSync("ls -lh /usr")
};