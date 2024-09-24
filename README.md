[![Sensu Bonsai Asset](https://img.shields.io/badge/Bonsai-Download%20Me-brightgreen.svg?colorB=89C967&logo=sensu)](https://bonsai.sensu.io/assets/DoctorOgg/sensu-fastapi-appsumo-healthcheck)
![goreleaser](https://github.com/DoctorOgg/sensu-fastapi-appsumo-healthcheck/workflows/goreleaser/badge.svg)

# sensu-fastapi-appsumo-healthcheck

## Table of Contents

- [Overview](#overview)
- [Files](#files)
- [Usage examples](#usage-examples)
- [Configuration](#configuration)
  - [Asset registration](#asset-registration)
  - [Check definition](#check-definition)
- [Installation from source](#installation-from-source)

## Overview

This is a simple health check for FastAPI AppSumo applications. It makes a GET request to the specified URL and expects a JSON response. If the response contains all keys with value "working", the check passes.

```json
{
  "status":"ok",
  "backends":
    {
      "Cache backend: default":"working",
      "DatabaseBackend":"working",
      "DefaultFileStorageHealthCheck":"working",
      "MigrationsHealthCheck":"working",
      "ProductsHealthCheckBackend":"working"
    }
}
```

## Files

- sensu-fastapi-appsumo-healthcheck

## Usage examples

```bash
sensu-fastapi-appsumo-healthcheck -u https://example.com/api/v1/health

```

Help:

```bash
sensu-fastapi-appsumo-healthcheck -h

Check FastAPI AppSumo health status

Usage:
  sensu-fastapi-appsumo-healthcheck [flags]
  sensu-fastapi-appsumo-healthcheck [command]

Available Commands:
  help        Help about any command
  version     Print the version number of this plugin

Flags:
  -d, --debug                  Enable debug mode
  -h, --help                   help for sensu-fastapi-appsumo-healthcheck
  -i, --insecure-skip-verify   Skip TLS certificate verification (not recommended!)
  -T, --timeout int            Request timeout in seconds (default 15)
  -u, --url string             URL to test (default "http://localhost:80/")
```

## Configuration

### Asset registration

[Sensu Assets][10] are the best way to make use of this plugin. If you're not using an asset, please consider doing so! If you're using sensuctl 5.13 with Sensu Backend 5.13 or later, you can use the following command to add the asset:

```bash
sensuctl asset add DoctorOgg/sensu-fastapi-appsumo-healthcheck
```

If you're using an earlier version of sensuctl, you can find the asset on the [Bonsai Asset Index][https://bonsai.sensu.io/assets/DoctorOgg/sensu-fastapi-appsumo-healthcheck].

### Check definition

```yml
---
api_version: core/v2
type: CheckConfig
metadata:
  name: check-api
  labels:
    sensu.io/workflow: ci_action
spec:
  runtime_assets:
    - sensu-django-healthcheck
  command: sensu-django-healthcheck -u https://example.com/api/v1/health
  subscriptions:
    - ecs-worker
  interval: 120
  round_robin: true
  proxy_entity_name: round_robin
  handlers:
    - notify_all
```
