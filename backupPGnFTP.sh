#!/bin/bash
### Postgres Server Login Info ###
PGUSER="valued"
PGPASS="valued"
# pgsql server
PGHOST="localhost"


### FTP SERVER Login info ###
FTPUSER="WfgcC4vVmN8"
FTPPASS="britiwqeqa_0"
FTPSERVER="88.99.137.173"
PGSQL="$(which pgsql)"
PGDUMP="$(which pg_dump)"
BAK="/backup/pgsql"
GZIP="$(which gzip)"
NOW=$(date +"%d-%m-%Y")[ ! -d $BAK ] && mkdir -p $BAK || /bin/rm -f $BAK/*DBS="$($PGSQL -u $PGUSER -h $PGHOST -p$PGPASS -Bse 'show databases')"
for db in $DBS
do
  FILE=$BAK/$db.$NOW-$(date +"%T").gz
  $PGDUMP -u $PGUSER -h $PGHOST -p$PGPASS $db | $GZIP -9 > $FILE
done
lftp -u $FTPUSER,$FTPPASS -e "mkdir /pgsql/$NOW;cd /pgsql/$NOW; mput /backup/pgsql/*; quit" $FTPSERVER


ALTER USER "valued" WITH PASSWORD 'valued';
