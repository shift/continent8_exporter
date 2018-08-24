# Continent 8 Prometheus Exporter

This exporter uses the Continent 8 support portal API to expose bandwidth and environmental metrics.

```bash
export C8_USERNAME=company.username
export C8_TOKEN=SupportToken
./continent8_exporter
```

Your token can be [generated here](https://support.continent8.com/api/token.php).

## Example Output

```
# HELP bandwidth C8 Bandwidth
# TYPE bandwidth counter
bandwidth{datacentre="IDCA",network="COMPANY_A35",rack="sid-12669",time="1535119222",type="in"} 6.34262072016e+12
bandwidth{datacentre="IDCA",network="COMPANY_A35",rack="sid-12669",time="1535119222",type="out"} 2.9760324199788e+13
bandwidth{datacentre="IDCA",network="total",rack="sid-12669",time="1535119222",type="in"} 6.34262072016e+12
bandwidth{datacentre="IDCA",network="total",rack="sid-12669",time="1535119222",type="out"} 2.9760324199788e+13
bandwidth{datacentre="MLT2",network="COMPANY_A7",rack="sid-19125",time="1535119222",type="in"} 1.14829481699e+11
bandwidth{datacentre="MLT2",network="COMPANY_A7",rack="sid-19125",time="1535119222",type="out"} 5.0825158356e+10
bandwidth{datacentre="MLT2",network="total",rack="sid-19125",time="1535119222",type="in"} 1.14829481699e+11
bandwidth{datacentre="MLT2",network="total",rack="sid-19125",time="1535119222",type="out"} 5.0825158356e+10
bandwidth{datacentre="TWN2",network="COMPANY_A2",rack="sid-15332",time="1535119222",type="in"} 4.108357128463e+12
bandwidth{datacentre="TWN2",network="COMPANY_A2",rack="sid-15332",time="1535119222",type="out"} 3.5743970609905e+13
bandwidth{datacentre="TWN2",network="total",rack="sid-15332",time="1535119222",type="in"} 4.108357128463e+12
bandwidth{datacentre="TWN2",network="total",rack="sid-15332",time="1535119222",type="out"} 3.5743970609905e+13
# HELP continent8_exporter_build_info A metric with a constant '1' value labeled by version, revision, branch, and goversion from which continent8_exporter was built.
# TYPE continent8_exporter_build_info gauge
continent8_exporter_build_info{branch="",goversion="go1.10.3",revision="",version=""} 1
# HELP environment C8 Environment
# TYPE environment gauge
environment{datacentre="IDCA",rack="S-1306",type="humidity"} 0
environment{datacentre="IDCA",rack="S-1306",type="power"} 2.497
environment{datacentre="IDCA",rack="S-1306",type="temperature"} 20.33
environment{datacentre="MLT2",rack="S-1207",type="humidity"} 0
environment{datacentre="MLT2",rack="S-1207",type="power"} 0.372
environment{datacentre="MLT2",rack="S-1207",type="temperature"} 0
environment{datacentre="TWN2",rack="S-1116",type="humidity"} 0
environment{datacentre="TWN2",rack="S-1116",type="power"} 0.401
environment{datacentre="TWN2",rack="S-1116",type="temperature"} 0
```

Sponsored by [Booming-Games](https://booming-games.com).
