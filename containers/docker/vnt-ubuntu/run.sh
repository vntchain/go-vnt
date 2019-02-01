#!/bin/sh

port1=30311
port2=30312
port3=30313
port4=30314

rpcport1=8555
rpcport2=8556
rpcport3=8557
rpcport4=8558

nodedir1="node1"
nodedir2="node2"
nodedir3="node3"
nodedir4="node4"

account1="0x9f55d20eb4f0d9f27da0db41c0331136772a5fb0"
account2="0xd915ce08a49e70d3b0c8d6f44e10646207167223" 
account3="0x663d9b36c1aa7a1e49a8619f6976993534caea76" 
account4="0x77dcc0a74d2a37a6acffdda9ece8f4ca166a8fbd"

addr1=`/nodeaddrgen -addr  ${port1} -datadir=/nodedir/${nodedir1}/`
addr2=`/nodeaddrgen -addr  ${port2} -datadir=/nodedir/${nodedir2}/`
addr3=`/nodeaddrgen -addr  ${port3} -datadir=/nodedir/${nodedir3}/`
addr4=`/nodeaddrgen -addr  ${port4} -datadir=/nodedir/${nodedir4}/`

addr=\"${addr1}\",\"${addr2}\",\"${addr3}\",\"${addr4}\"
#echo ${addr}
genesis=`cat /genesis.json`
genesisres=${genesis/\/\*replaceaddress\*\//${addr}}
echo ${genesisres} > /genesis.json
# echo aaa${arr[1]}

echo `/gvnt init /genesis.json --datadir /nodedir/${nodedir1}`
echo `/gvnt init /genesis.json --datadir /nodedir/${nodedir2}`
echo `/gvnt init /genesis.json --datadir /nodedir/${nodedir3}`
echo `/gvnt init /genesis.json --datadir /nodedir/${nodedir4}`

echo "/gvnt --networkid 1015 --datadir /nodedir/${nodedir1} --port ${port1} --rpcport ${rpcport1} --unlock ${account1} --password /password --mine"
echo `/gvnt --networkid 1015 --datadir /nodedir/${nodedir1} --port ${port1} --rpcport ${rpcport1} --unlock ${account1} --password /password --mine`

echo "/gvnt --networkid 1015 --datadir /nodedir/${nodedir2} --port ${port2} --rpcport ${rpcport2} --unlock ${account2} --password /password --vntbootnode ${addr1} --mine"
echo `/gvnt --networkid 1015 --datadir /nodedir/${nodedir2} --port ${port2} --rpcport ${rpcport2} --unlock ${account2} --password /password --vntbootnode ${addr1} --mine`

echo "/gvnt --networkid 1015 --datadir /nodedir/${nodedir3} --port ${port3} --rpcport ${rpcport3} --unlock ${account3} --password /password --vntbootnode ${addr2} --mine"
echo `/gvnt --networkid 1015 --datadir /nodedir/${nodedir3} --port ${port3} --rpcport ${rpcport3} --unlock ${account3} --password /password --vntbootnode ${addr2} --mine`

echo "/gvnt --networkid 1015 --datadir /nodedir/${nodedir4} --port ${port4} --rpcport ${rpcport4} --unlock ${account4} --password /password --vntbootnode ${addr3} --mine"
echo `/gvnt --networkid 1015 --datadir /nodedir/${nodedir4} --port ${port4} --rpcport ${rpcport4} --unlock ${account4} --password /password --vntbootnode ${addr3} --mine`