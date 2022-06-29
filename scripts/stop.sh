#!/bin/sh
set -e

GREP="/usr/bin/grep"

ps -ef | $GREP node-exporter            | $GREP -v grep | awk '{print $2}' | xargs kill || echo "node-exporter not running"
ps -ef | $GREP hstream-metrics-exporter | $GREP -v grep | awk '{print $2}' | xargs kill || echo "hstream-metrics-exporter not running"
ps -ef | $GREP http-server              | $GREP -v grep | awk '{print $2}' | xargs kill || echo "http-server not running"


docker stop hs-test-grafana    || echo "hs-test-grafana not running"
docker stop hs-test-prometheus || echo "hs-test-prometheus not running"

docker stop hs-test-hserver0 || echo "hs-test-hserver0 not running"
docker stop hs-test-hserver1 || echo "hs-test-hserver1 not running"
docker stop hs-test-hserver2 || echo "hs-test-hserver2 not running"

docker stop hs-test-hstore  || echo "hs-test-hstore not running"
docker stop hs-test-zk      || echo "hs-test-zk not running"

docker rm hs-test-hserver0 || echo "hs-test-hserver0 not running"
docker rm hs-test-hserver1 || echo "hs-test-hserver1 not running"
docker rm hs-test-hserver2 || echo "hs-test-hserver2 not running"

rm -rf ./nohup.out
