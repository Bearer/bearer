const { exec, execSync, spawn, spawnSync } = require('node:child_process');

exports.handler = async (event) => {
  exec("ls "+event["user_dir"]+"| wc -l", (err, stdout, stderr) => {
    // do something
  });

  execSync("ls "+event["user"]+"| wc -l", (err, stdout, stderr) => {
    // do something
  });

  spawn(event["query"]);

  spawnSync("grep " + event["tmp"])
};