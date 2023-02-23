#!/bin/bash
for tmp in 204 200 183 185 192
do
  {
    serverName="amd${tmp}.utah.cloudlab.us"

    ssh "borisli@${serverName}" "cd shared-registers; git reset --hard HEAD; git pull --no-rebase;\
      chmod +x shell_scripts/*;\
      source /etc/profile;\
      cd client; make clean; make build;"
  }&
done
wait
echo "All nodes are ready."
