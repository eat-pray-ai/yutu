```shell
Arguments: Specify one of the following:
  * console - runs an agent in console mode.
  * web - starts web server with additional sub-servers specified by sublaunchers
Details:
  console
  -streaming_mode string
        defines streaming mode (none|sse) (default "sse")

  web
  -idle-timeout duration
        Server idle timeout (i.e. '10s', '2m' - see time.ParseDuration for details) - for waiting for the next request (only when keep-alive is enabled) (default 1m0s)
  -port int
        Localhost port for the server (default 8080)
  -read-timeout duration
        Server read timeout (i.e. '10s', '2m' - see time.ParseDuration for details) - for reading the whole request including body (default 15s)
  -shutdown-timeout duration
        Server shutdown timeout (i.e. '10s', '2m' - see time.ParseDuration for details) - for waiting for active requests to finish during shutdown (default 15s)
  -write-timeout duration
        Server write timeout (i.e. '10s', '2m' - see time.ParseDuration for details) - for writing the response after reading the headers & body (default 15s)
  You may specify sublaunchers:
    * api - starts ADK REST API server, accepting origins specified by webui_address (CORS)
    * a2a - starts A2A server which handles jsonrpc requests on /a2a/invoke path
    * webui - starts ADK Web UI server which provides UI for interacting with ADK REST API
  Sublaunchers syntax:
    api
    -sse-write-timeout duration
        SSE server write timeout (i.e. '10s', '2m' - see time.ParseDuration for details) - for writing the SSE response after reading the headers & body (default 2m0s)
    -webui_address string
        ADK WebUI address as seen from the user browser. It's used to allow CORS requests. Please specify only hostname and (optionally) port. (default "localhost:8080")

    a2a
    -a2a_agent_url string
        A2A host URL as advertised in the public agent card. It is used by A2A clients as a connection endpoint. (default "http://localhost:8080")

    webui
    -api_server_address string
        ADK REST API server address as seen from the user browser. Please specify the whole URL, i.e. 'http://localhost:8080/api'. (default "http://localhost:8080/api")
```
