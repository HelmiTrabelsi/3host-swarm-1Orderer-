peer=$1
if [ ${#peer} == 0 ]; then echo "Error: Enter peer number" ; exit
fi
./host$1down.sh
git pull
./host$1up.sh

