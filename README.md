# prom-source-http

An HTTP server which serves content of the file for further use by Prometheus. 
For that use GET `/_metrics`

It can also fetch Prometheus text-formatted (*.prom) endpoint and dump the result as JSON.
For that use GET `/_metrics.json?url=<url to *.prom endpoin>`

Use `prom-source-http --help` for details.