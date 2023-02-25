#!/bin/bash


# Path to the private key file
keyfile=./privateKeyForCluster
chmod 400 $keyfile
# Add the key to the SSH agent
eval `ssh-agent -s` > /dev/null 2>&1
ssh-add $keyfile > /dev/null 2>&1

echo "Starting the servers."

for tmp in 204 200 183 185 192
do
  {
    serverName="amd${tmp}.utah.cloudlab.us"
    ssh -o StrictHostKeyChecking=no -t borisli@${serverName} > /dev/null 2>&1 << EOF
      kill $(lsof -i:50051 -t)
      sleep 1
      cd shared-registers/server
      git reset --hard HEAD
      git pull --no-rebase
      source /etc/profile
      make clean
      make build
      nohup make run
EOF
  }&
done


sleep 5 # wait for servers start
echo "All servers are running now. At $(date)"
echo "Running test now, it takes about 10 sec. "
clientName="amd183.utah.cloudlab.us"
ssh -o StrictHostKeyChecking=no "borisli@${clientName}" "cd shared-registers/client/protocol;\
  source /etc/profile;\
  go test -run ^$ -bench=BenchmarkWriteThenReadForDemo"
echo "All tests are done"

for tmp in 204 200 183 185 192
do
  {
    serverName="amd${tmp}.utah.cloudlab.us"

    ssh -o StrictHostKeyChecking=no "borisli@${serverName}" "kill \`lsof -i:50051 -t\`;"
  }&
done
wait
echo "All servers are stopped now."