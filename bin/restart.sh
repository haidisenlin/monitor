#!/bin/sh

tomcat_id=`ps -ef | grep tomcat | grep -v "grep" | awk '{print $2}'`
echo $tomcat_id

for id in $tomcat_id
do
    -kill -9 $id  
    echo "killed $id"
done
/opt/tomcat6/bin/startup.sh