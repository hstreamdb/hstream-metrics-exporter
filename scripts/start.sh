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
  --name hs-test-hserver                                            \
  -v     $DATA_STORE:/data/store                                    \
    $HSTREAM_IMAGE                                                  \
      hstream-server                                                \
        --store-config /data/store/logdevice.conf --log-level debug \
        --port 6570 --server-id 0

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


node-exporter &


docker run -td --network host                                \
  --rm                                                       \
  --name hs-test-grafana                                     \
  -e GF_AUTH_ANONYMOUS_ORG_ROLE=Admin                        \
  -e GF_AUTH_ANONYMOUS_ENABLED=true                          \
  -e GF_AUTH_DISABLE_LOGIN_FORM=true                         \
  -v "$(pwd)"/configs/provisioning:/etc/grafana/provisioning \
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
