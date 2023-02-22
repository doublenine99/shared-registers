#!/bin/bash
for tmp in 204 200 183 185 192
do
  {
    serverName="amd${tmp}.utah.cloudlab.us"

    ssh "borisli@${serverName}" "kill \`lsof -i:50051 -t\`;\
      cd shared-registers/server; git reset --hard HEAD; git pull --no-rebase;\
      chmod +x ../shell_scripts/*;\
      source /etc/profile;\
      make clean; make build; make run;\
      cd ../client; make clean; make build;"
  }&
done
wait
echo "All servers are running now."
