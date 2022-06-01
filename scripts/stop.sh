#!/bin/sh
set -e

GREP="/usr/bin/grep"

ps -ef | $GREP hstream-metrics-exporter | $GREP -v grep | awk '{print $2}' | xargs kill || echo "hstream-metrics-exporter not running"
ps -ef | $GREP http-server              | $GREP -v grep | awk '{print $2}' | xargs kill || echo "http-server not running"


docker stop hs-test-hserver || echo "hs-test-hserver not running"
docker stop hs-test-hstore  || echo "hs-test-hstore not running"
docker stop hs-test-zk      || echo "hs-test-zk not running"

docker rm hs-test-hserver   || echo "hs-test-hserver not running"

rm -rf ./nohup.out
