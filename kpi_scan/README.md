# KPI Scan docker image

This docker image is ubuntu with a script to download the latest Bearer CLI
and run it for a given REPOSITORY_URL and API_KEY.

## Building

The image must be built and deployed manually. For MacOS:

```sh
$ docker buildx build --platform=linux/amd64 -t bearersh/kpi-scan .
$ docker push bearersh/kpi-scan:latest
```
