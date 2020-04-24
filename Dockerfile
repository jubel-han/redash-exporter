FROM golang:1.13

ENV REDASH_API_BASE_URL ""
ENV REDASH_API_KEY ""
ENV REDASH_PROBE_QUERY_ID ""
ENV REDASH_PROBE_ALERT_ID ""
ENV REDASH_PROBE_INTERVAL=1h

COPY . /redash-exporter

RUN cd /redash-exporter \
 && go build ./... \
 && cp redash-exporter /usr/bin/redash-exporter \
 && chmod 755 /usr/bin/redash-exporter \
 && cd / \
 && rm -rf /redash-exporter

ENTRYPOINT ["/usr/bin/redash-exporter"]
