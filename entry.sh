echo off

[ -z "$TRANS_UID" ] && echo "TRANS_UID is not set" && exit 1
[ -z "$TRANS_GID" ] && echo "TRANS_GID is not set" && exit 1


echo $TRANS_UIDcode
echo $TRANS_GID

getent group $TRANS_GID

if [ $? -ne 0 ];then
   groupadd -g $TRANS_GID transmissiongroup
fi

id -u $TRANS_UID -u -n

if [ $? -ne 0 ];then
    useradd --uid $TRANS_UID --gid $TRANS_GID transmissionuser
    runuser="transmissionuser"
else
    runuser=$(id -u $TRANS_UID -u -n)
fi

mkdir -p /tmp/transmission
cp /etc/transmission/settings.json /tmp/transmission
chmod 777 -R /tmp/transmission

mkdir -p /data

su -l $runuser -c '/usr/bin/transmission-daemon -f --log-error -g /tmp/transmission'
