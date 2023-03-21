# cloudflare_bgp_announcement_exporter
Get Cloudflare Radar BGP announcement number of AS in prometheus
###### Made with some help of ChatGPT 4
# Usage

## Build

```bash
$ go build .
```

## Run

### Environment variables:

- ASN: Comma separated list of ASNs to monitor
- CLOUDFLARE_API_TOKEN: Cloudflare API token with Radar permissions

```bash
$ ./cloudflare_bgp_announcement_exporter
```

## Docker

```bash
$ docker run -d -p 8080:8080 --name cloudflare_bgp_announcement_exporter -e ASN=6939,3215 -e CLOUDFLARE_API_TOKEN=apikey_radar guillaumeouint2/cloudflare_bgp_announcement_exporter
```

## Metric

```
bgp_dfz_announcements
```
ASN in label.