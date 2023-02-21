#note: I put executable in shell script folder
#TODO: only some output is going to the output files, why???
for i in 1 2 3
do
  echo "starting client $i"
  ./shell_scripts/go_build_shared_registers_client.exe < ./shell_scripts/inputs_$i.txt > /shell_scripts/outputs_$i.txt &
done
echo "clients done running commands"
sleep 20