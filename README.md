# Shared Registers

## Run Demo Test
We make all the build, deploy and run operations into our scripts, so you don't need to manually build anything, just follow the instruction below:

Go to the script folder and add execute permission to the script.
```shell
cd RunTestDemo
chmod +x ./deployAndTestSystem.sh
```
Then simply run the script.
```shell
./deployAndTestSystem.sh
```
The script will automatically run all the servers in the cluster and start 10 clients, perform 10 seconds Read/Write mixed operations.
After the test, the result of the test will be printed out to the console. The servers will be stopped at the end of this script. 