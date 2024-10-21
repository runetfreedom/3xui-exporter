# Usage

### cli 
```bash
./3xui-exporter -endpoint=http://localhost/mysecretpath -username=YOUR_USERNAME -password=YOUR_PASSWORD -listen-address=:9101
```

### docker
```bash
docker build . -t 3xui-exporter

docker run -d --rm --name 3xui-exporter -p 9101:9101 3xui-exporter 3xui-exporter -endpoint=http://localhost/mysecretpath -username=YOUR_USERNAME -password=YOUR_PASSWORD -listen-address=:9101
```

### docker-compose:

```yaml
3xui-exporter:
    build: 3xui-exporter
    container_name: 3xui-exporter
    restart: unless-stopped    
    environment:
      PANEL_ENDPOINT: http://localhost/mysecretpath
      PANEL_USERNAME: YOUR_USERNAME
      PANEL_PASSWORD: YOUR_PASSWORD
      LISTEN_ADDRESS: :9101
```

### prometheus.yml
```yaml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: '3xui'
    static_configs:
      - targets: ['3xui-exporter:9101']
```

### Example dashboard

[dashboard.json](https://github.com/runetfreedom/3xui-exporter/grafana/dashboard.json)
