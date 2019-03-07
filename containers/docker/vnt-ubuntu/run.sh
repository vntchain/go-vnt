#!/bin/sh
killall gvnt
rm -rf /nodedir
peernumber=$1
echo $peernumber
if [[ $peernumber -lt 4 ]]  
then  
     peernumber=4
fi   
baseport=30311
baserpcport=8555
basenodedir="/nodedir/node"

echo "create keystore and p2paddress"
password="12345678"
echo ${password} > password
declare -a keystore
declare -a p2paddress
for ((i=0; i<${peernumber}; i ++))  
do 
    #cd ${basenodedir}${i} && rm -rf *
    #cd /
    echo "/gvnt account new --password ${password} --datadir ${basenodedir}${i}"
    newaddress=`/gvnt account new --password /password  --datadir ${basenodedir}${i}`
    newaddress="0x"${newaddress:10:40}
    echo ${newaddress} 
    let port=baseport+i
    p2paddr=`/nodeaddrgen -addr  ${port} --datadir ${basenodedir}${i}`
    echo ${p2paddr}
    keystore[i]=${newaddress} 
    p2paddress[i]=${p2paddr}
done

cd /

echo "replace alloc,witnesses,p2paddress in genesis.json"

wlist=""
alist=""
plist=""
for data in ${keystore[@]}
do
   alist=${alist},\"${data}\"\:\{"\"balance\":\"0x200000000000000000000000000000000000000000000000000000000000000\""}
   wlist=${wlist},\"${data}\"
done
alist=${alist:1:${#alist}}
wlist=${wlist:1:${#wlist}}
echo $wlist
echo $alist

for data in ${p2paddress[@]}
do
    echo ${data}
    plist=${plist},\"${data}\"
done
plist=${plist:1:${#plist}}
echo ${plist}

genesis=`cat /genesis.json.templet`
genesis=${genesis/\/\*replaceaddress\*\//${plist}}
genesis=${genesis/\/\*replacewitnesses\*\//${wlist}}
genesis=${genesis/\/\*replacealloc\*\//${alist}}
genesis=${genesis/\/\*replacenumber\*\//${peernumber}}
echo ${genesis} > /genesis.json

#init gnvt
echo "init gvnt"

for ((i=0; i<$peernumber; i ++))
do
    echo `/gvnt init /genesis.json --datadir  ${basenodedir}${i}`
done

# start peer
echo "start peer"
for ((i=0; i<$peernumber; i ++))  
do  
    let port=baseport+i
    let rpcport=baserpcport+i
    if [ $i -eq 0 ]  
    then  
      nohup /gvnt --verbosity 5 --networkid 1015 --datadir  ${basenodedir}${i} --port ${port} --rpcport ${rpcport} --unlock ${keystore[${i}]} --password /password --mine > /tmp/node${i}.log 2>&1 &  
    else  
      let j=i-1
      nohup /gvnt --verbosity 5 --networkid 1015 --datadir  ${basenodedir}${i} --port ${port} --rpcport ${rpcport} --unlock ${keystore[${i}]} --password /password --vntbootnode ${p2paddress[${j}]} --mine > /tmp/node${i}.log 2>&1 &        
    fi   
done

tail -f /tmp/node0.log