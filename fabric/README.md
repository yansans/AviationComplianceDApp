# Fabric

[Hyperledger Fabric Prerequisites](https://hyperledger-fabric.readthedocs.io/en/latest/prereqs.html)

0. Pastikan requirement sudah terinstall
1. Download script install fabric `curl -sSLO https://raw.githubusercontent.com/hyperledger/fabric/main/scripts/install-fabric.sh && chmod +x install-fabric.sh`
2. Lakukan instalasi fabric binary dan docker `./install-fabric.sh b d`
3. Masuk ke folder test-network `cd test-network`
4. Jalankan network fabric `./network down` `./network up`

export PATH=${PWD}/../bin:$PATH
