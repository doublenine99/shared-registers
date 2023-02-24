#!/bin/bash
for tmp in 204 200 183 185 192
do
  {
    serverName="amd${tmp}.utah.cloudlab.us"

    ssh -t borisli@${serverName} > /dev/null 2>&1 << EOF
      kill $(lsof -i:50051 -t)
      cd shared-registers/server
      git reset --hard HEAD
      git pull --no-rebase
      source /etc/profile
      make clean
      make build
      nohup make run
EOF
    echo "Sever ${serverName} starts at $(date)"

#    ssh "borisli@${serverName}" "kill \`lsof -i:50051 -t\`;\
#      cd shared-registers/server; git reset --hard HEAD; git pull --no-rebase;\
#      chmod +x ../shell_scripts/*;\
#      source /etc/profile;\
#      make clean; make build; nohup make run > /dev/null 2>&1 &;\
#      echo "Starting Server ${serverName} at $(date)";"
  }&
done
wait

sleep 1 # wait for servers start
echo "All servers are running now. "

clientName="amd183.utah.cloudlab.us"
ssh "borisli@${clientName}" "cd shared-registers/client/protocol;\
  source /etc/profile;\
  go test -run ^$ -bench=BenchmarkWriteThenReadForDemo"
echo "all test are done"
