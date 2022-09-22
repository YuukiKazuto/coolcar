DOMAIN=$1
cd ../server
docker build -t coolcar-kh/$DOMAIN -f ../deployment/$DOMAIN/Dockerfile .