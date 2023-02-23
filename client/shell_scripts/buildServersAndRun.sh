#!/bin/bash
for tmp in 204 200 183 185 192
do
  {
    serverName="amd${tmp}.utah.cloudlab.us"

    ssh "borisli@${serverName}" "cd shared-registers/server;\
      chmod +x ../shell_scripts/*;\
      source /etc/profile;\
      make clean; make build; make run"
  }&
done
wait
echo "All servers are running now."