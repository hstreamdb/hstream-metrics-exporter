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
