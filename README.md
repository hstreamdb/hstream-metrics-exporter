# HStream Metrics Exporter

## Start Grafana

The command

```shell
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
```

will start a Grafana instance which will listen at `localhost:3000`.

The options

```shell
  -e GF_AUTH_ANONYMOUS_ORG_ROLE=Admin                                 \
  -e GF_AUTH_ANONYMOUS_ENABLED=true                                   \
  -e GF_AUTH_DISABLE_LOGIN_FORM=true                                  \
```

are used for skip log in screen and log in as anonymous user with admin permission so that data sources can be added to
Grafana.

The option

```shell
  -e GF_DEFAULT_APP_MODE=development                                  \
```

is used for enable unsigned Grafana plugins.

## Provisioning

The option

```shell
  -v "$PWD"/configs/provisioning:/etc/grafana/provisioning            \
```

will prepare **data sources** and **dashboards** when starting Grafana.

The configuration

```shell
  -v "$PWD"/configs/defaults.ini:/usr/share/grafana/conf/defaults.ini \
```

is the same with default configuration in the image, excepted the minimum refresh interval are set lower than 5s.

## Use HStream data source plugin

### Build

```shell
go install github.com/magefile/mage@latest # if mage is not installed
cd ./grafana-plugins/hstreamdb-grafana-plugin && yarn install && make # in the Makefile, `GOBIN` is set to `~/go/bin` by default
```

### Add a panel

Configurations can be done by either using the provisioning utils above or setting manually.

Manually settings:

1. Add the data source provide by HStream plugin: first click the bottom left settings button, then choose the data
   sources option (which requires log in as admin and `GF_DEFAULT_APP_MODE=development`)
2. In any dashboard, add a panel, choose the `Table` format at top right.
3. Test query, if data is got correctly, near the query builder click `Transform` and search `Filter by name`,
   ignore `time` and `values` fields.
