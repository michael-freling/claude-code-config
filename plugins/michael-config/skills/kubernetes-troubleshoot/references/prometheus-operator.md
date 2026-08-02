# Prometheus Operator

- Creates StatefulSets for Prometheus/Alertmanager CRs
- Uses label selectors on informers
- Check `.status.selector` on Prometheus CR for expected labels
