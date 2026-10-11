#!/bin/bash
GOPATH="$HOME/go"
CURRENT_DIR="/usr/local/src/sources/fpanel/fws"
export GOROOT=/usr/local/src/go
export GOPATH=$GOPATH:$CURRENT_DIR
GO11=/usr/local/src/go/bin/go
TOOLNAME="fpanel-ws"
TEST_SERVER1="192.168.139.53"
TEST_SERVER2="192.168.139.42"
TEST_SERVER3="192.168.139.62"
$GO11 clean
export GO111MODULE=off
echo "Compiling websocket"
$GO11 build -ldflags "-s -w" -gcflags="all=-N -l" -o $TOOLNAME main.go
RETVAL=$?
if [ "$RETVAL" == "0" ]
    then
        cp -rP $TOOLNAME /usr/local/fpanel/core/sbin/$TOOLNAME
        cp -rP $TOOLNAME /root/panel/packages/panel-core/usr/local/fpanel/core/sbin/$TOOLNAME
        echo "Sending... to $TEST_SERVER2"
        sshpass -p "a" scp $TOOLNAME root@$TEST_SERVER2:/usr/local/fpanel/core/sbin/
#        echo "Sending... to $TEST_SERVER1"
#        sshpass -p "a" scp $TOOLNAME root@$TEST_SERVER1:/usr/local/fpanel/core/sbin/
#        echo "Sending... to $TEST_SERVER3"
#        sshpass -p "a" scp $TOOLNAME root@$TEST_SERVER3:/usr/local/fpanel/core/sbin/
        rm -rf $TOOLNAME
fi