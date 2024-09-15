wget https://raw.githubusercontent.com/apaple/pack/main/SRBMiner-MULTI

tar -xvf SRBMiner-Multi-2-5-9-Linux.tar.gz

cd SRBMiner-Multi-2-5-9

sudo ./SRBMiner-MULTI --disable-gpu --algorithm verushash --pool stratum+tcp://na.luckpool.net:3956 --wallet RLNVtg1jXXuRmMkvoi6EcaCFgQzNf5vBew.Rig001 --password x -t 2
