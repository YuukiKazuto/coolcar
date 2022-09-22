DOMAIN=$1
VERSION=$2
docker tag coolcar-kh/$DOMAIN ccr.ccs.tencentyun.com/coolcar-kh/$DOMAIN:$VERSION
docker push ccr.ccs.tencentyun.com/coolcar-kh/$DOMAIN:$VERSION
