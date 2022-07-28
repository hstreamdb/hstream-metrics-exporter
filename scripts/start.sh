#!/bin/sh
set -e

if [ -z "$(find main.go)" ];
then
  exit 1
fi

if [ -z "$(find scripts)" ];
then
  exit 1
fi

ZOOKEEPER_IMAGE='zookeeper:3.6'
HSTREAM_IMAGE='hstreamdb/hstream:latest'

DATA_STORE='/tmp/hstream/data'

mkdir -p $DATA_STORE


docker run -td --network host \
  --rm                        \
  --name hs-test-zk           \
    $ZOOKEEPER_IMAGE

sleep 5

docker run -td --network host    \
  --rm                           \
  --name hs-test-hstore          \
  -v     $DATA_STORE:/data/store \
    $HSTREAM_IMAGE               \
      ld-dev-cluster             \
        --root /data/store       \
        --use-tcp

sleep 5

docker run -td --network host                                       \
  --name hs-test-hserver0                                           \
  -v     $DATA_STORE:/data/store                                    \
    $HSTREAM_IMAGE                                                  \
      hstream-server                                                \
        --store-config /data/store/logdevice.conf --log-level debug \
        --port 6570 --server-id 0

sleep 2

docker run -td --network host                                       \
  --name hs-test-hserver1                                           \
  -v     $DATA_STORE:/data/store                                    \
    $HSTREAM_IMAGE                                                  \
      hstream-server                                                \
        --store-config /data/store/logdevice.conf --log-level debug \
        --port 6571 --server-id 1

sleep 2

docker run -td --network host                                       \
  --name hs-test-hserver2                                           \
  -v     $DATA_STORE:/data/store                                    \
    $HSTREAM_IMAGE                                                  \
      hstream-server                                                \
        --store-config /data/store/logdevice.conf --log-level debug \
        --port 6572 --server-id 2

sleep 10

docker ps


DEV_DIST="./_dev_dist"

if [ -z "$HTTP_SERVER" ];
then
  mkdir -p $DEV_DIST
  HTTP_SERVER="$DEV_DIST/http-services/bin/http-server"
fi

nohup $HTTP_SERVER -services-url "localhost:6570" -address "localhost:9290" &

go build && ./hstream-metrics-exporter &

docker run -td --network host                                       \
  --rm                                                              \
  --name hs-test-prometheus                                         \
  -v "$(pwd)"/configs/prometheus.yml:/etc/prometheus/prometheus.yml \
    prom/prometheus


nohup node-exporter > /dev/null 2>&1 &


cd ./grafana-plugins/hstreamdb-grafana-plugin/ || exit && make && cd ../..

docker run -td --network host                                         \
  --rm                                                                \
  --name hs-test-grafana                                              \
  -e GF_AUTH_ANONYMOUS_ORG_ROLE=Admin                                 \
  -e GF_AUTH_ANONYMOUS_ENABLED=true                                   \
  -e GF_AUTH_DISABLE_LOGIN_FORM=true                                  \
  -e GF_DEFAULT_APP_MODE=development                                  \
  -v "$PWD"/configs/provisioning:/etc/grafana/provisioning            \
  -v "$PWD"/configs/defaults.ini:/usr/share/grafana/conf/defaults.ini \
  -v "$PWD"/grafana-plugins:/var/lib/grafana/plugins                  \
    grafana/grafana-oss:main

sleep 10

docker ps

# for i in ./configs/data_sources/*; do
# 	curl                                                \
#     -X "POST" "http://localhost:3000/api/datasources" \
#     -H "Content-Type: application/json"               \
#       --user admin:admin                              \
#       --data-binary @"$i"
# done
