# Gorge

![Release](https://github.com/whitewater-guide/gorge/workflows/Release/badge.svg?branch=master&event=push)

Gorge is a service which harvests hydrological data (river's discharge and water level) on schedule.
Harvested data is stored in database and can be queried later.

## Usage

Gorge is distributed as docker image with two binary files:

- `gorge-server` (_entrypoint_) - web server with REST API
- `gorge-cli` - command-line client for this server. Since image is distroless, use `docker exec gorge gorge-cli` to call it

### Launching

`gorge-server` accepts configuration via cli arguments (use `gorge-server --help`). You can pass them via docker-compose command field, like this:

```yaml
command:
  [
    "--pg-db",
    "gorge",
    "--debug",
    "--log-format",
    "plain",
    "--db-chunk-size",
    "1000",
  ]
```

Here is the list of available flags:

```
--cache string             Either 'inmemory' or 'redis' (default "redis")
--db string                Either 'inmemory' or 'postgres' (default "postgres")
--db-chunk-size int        Measurements will be saved to db in chunks of this size. When set to 0, they will be saved in one chunk, which can cause errors
--debug                    Enables debug mode, sets log level to debug
--endpoint string          Endpoint path (default "/")
--http-timeout int         Request timeout in seconds (default 60)
--http-user-agent string   User agent for requests sent from scripts. Leave empty to use fake browser agent (default "whitewater.guide robot")
--log-format string        Set this to 'json' to output log in json (default "json")
--log-level string         Log level. Leave empty to discard logs (default "warn")
--pg-db string             Postgres database (default "postgres")
--pg-host string           Postgres host (default "db")
--pg-password string       Postgres password
--pg-user string           Postgres user (default "postgres")
--port string              Port (default "7080")
--redis-host string        Redis host (default "redis")
--redis-port string        Redis port (default "6379")
```

Postgres and redis can also be configured using folowing environment variables:

- POSTGRES_HOST
- POSTGRES_DB
- POSTGRES_USER
- POSTGRES_PASSWORD
- REDIS_HOST
- REDIS_PORT

Environment variables have lower priority than cli flags.

Gorge uses database to store harvested measurements and scheduled jobs. It comes with postgres and sqlite drivers. Postgres with timescaledb extension is recommended for production. Gorge will initialize all the required tables. Check out sql migration file if you're curious about db schema.

Gorge uses cache to store safe-to-lose data: latest measurement each gauge and harvest statuses. It comes with redis (recommended) and embedded redis drivers.

Gorge server is supposed to be running in private network. It doesn't support HTTPS. If you want to expose it to public, use reverse proxy.

### Working with API

Below is the list of emdpoints exposed by gorge server. You can use `request.http` files in project root and script directories to play with running server.

- `GET /version`

  Returns running server version:
  ```json
  {
    "version": "1.0.0"
  }
  ```

- `GET /scripts`

  Returns array of available scripts with their harvest modes:
  ```json
  [
    {
      "name": "sepa",
      "mode": "oneByOne"
    },
    {
      "name": "switzerland",
      "mode": "allAtOnce"
    }
  ]
  ```

- `POST /upstream/{script}/gauges`

  Lists gauges available for harvest in an upsteam source.

  URL parameters:
    
  - `script` - script name for usptream source
  
  POST body contains JSON that contains script-specific parameters. For example, it can contain authentication credentials for protected sources. Another example is `all_at_once` test script, which accepts `gauges` JSON parameter to specify number of gauges to return.

  Returns JSON array of gauges. For example:

  ```json
  [
    {
      "script": "tirol", // script name
      "code": "201012", // gauge code in upstream source
      "name": "Lech / Steeg", // gauge name
      "url": "https://apps.tirol.gv.at/hydro/#/Wasserstand/?station=201012", // upstream gauge webpage for humans 
      "levelUnit": "cm", // units of water level measurement, if gauge provides water level 
      "flowUnit": "cm", // units of water discharge measurement, if gauge provides discharge
      "location": { // gauge location in EPSG4326 coordinate system, if provided
        "latitude": 47.24192,
        "longitude": 10.2935,
        "altitude": 1109
      }
    },
  ]
  ```

- `POST /upstream/{script}/measurements?codes=[codes]&since=[since]`

  Harvests measurements directly from upstream source without saving them.

  URL parameters:
    
  - `script` - script name for usptream source
  - `codes` - comma-separated list of gauge codes to return. This paramter is required for one-by-one scripts. For all-at-once scripts it's optional, and without it all gauges will be returned.
  - `since` - optional unix timstamp indicating start of the period you want to get measurements from. This is passed directly to upstream, if it support such parameter (very few actually do)
  
  POST body contains JSON that contains script-specific parameters. For example, it can contain authentication credentials for protected sources. Another example is `all_at_once` test script, which accepts `min`, `max` and `value` JSON parameters to control produced values.

  Returns JSON array of measurements. For example:

  ```json
  [
    {
      "script": "tirol", // script name
      "code": "201178", // gauge code
      "timestamp": "2020-02-25T17:15:00Z", // timestamp in RFC3339
      "level": 212.3, // water level value, if provided, otherwise null
      "flow": null // water discharge value, if provided, otherwise null
    }
  ]
  ```

- `GET /jobs`

  Returns array of running jobs:

  ```json
  [
    {
      "id": "3382456e-4242-11e8-aa0e-134a9bf0be3b", // unique job id
      "script": "norway", // job script
      "gauges": { // array of gauges that this job harvests
        "100.1": null,
        "103.1": { // it's possible to set script-specific parameter for each individual gauge
          "version": 2
        }
      },
      "cron": "38 * * * *", // job's cron schedule, for all-at-once jobs
      "options": { // script-specific parameters
        "csv": true
      },
      "status": { // information about running job
        "success": true, // whether latest execution was successful
        "timestamp": "2020-02-25T17:44:00Z", // latest execution timestamp
        "count": 10, // number of measurements harvested during latest execution
        "next": "2020-02-25T17:46:00Z", // next execution timestamp
        "error": "somethin went wrong" // latest execution error, omitted when success = true
      }
    }
  ]
  ```

- `GET /jobs/{jobId}`

  URL parameters:
    
  - `jobId` - harvest job id

  Returns the job description. It's same as item in `/jobs` array, but without `status`

  ```json
  {
    "id": "3382456e-4242-11e8-aa0e-134a9bf0be3b",
    "script": "norway",
    "gauges": {
      "100.1": null,
      "103.1": {
        "version": 2
      }
    },
    "cron": "38 * * * *",
    "options": null
  }
  ```

- `GET /jobs/{jobId}/gauges`

  URL parameters:
    
  - `jobId` - harvest job 
  
  Returns map object with gauge statuses, where keys are gauge codes and values are statuses:

  ```json
  [
    {
    "010802": {
      "success": false,
      "timestamp": "2020-02-24T18:00:00Z",
      "count": 0,
      "next": "2020-02-25T18:00:00Z"
    }
  ]
  ```

- `POST /jobs`

  Adds new job.

  POST body must contain JSON job description. For example:

  ```json
  {
    "id": "78a9e166-2a73-4be2-a3fb-71d254eb7868", // unique id, must be set by client
    "script": "one_by_one", // script for this job
    "gauges": { // list of gauges
      "g000": null, // set to null if gauge has no script-specific options
      "g001": { "version": 2 } // or pass script-specific options
    },
    "options": { // optional, common script-specific options
      "auth": "some_token"
    },
    "cron": "10 * * * *" // cron schedule required for all-at-once scripts
  }
  ```

  Returns same object in case of success, error object otherwise

- `DELETE /jobs/{jobId}`
  
  URL parameters:
    
  - `jobId` - harvest job id
  
  Stop the job and deletes it from schedule

- `GET /measurements/{script}/{code}?from=[from]&to=[to]`

  URL parameters:
    
  - `script` - script name
  - `code` - optional, gauge code
  - `from` - optional unix timstamp indicating start of the period you want to get measurements from. Default to 30 days from now.
  - `to` - optional unix timstamp indicating end of the period you want to get measurements from. Defaults to now.

  Returns array of measurements that were harvested and stored in gorge database for given script (and gauge). Resulting JSON is same as in `/upstream/{script}/measurements`

- `GET /measurements/{script}/{code}/latest`

  URL parameters:

  - `script` - script name, required
  - `code` - gauge code, optional

  Returns array of measurements for given script or gauge. For each gauge, only latest measurement will be returned. Resulting JSON is same as in `/upstream/{script}/measurements`

- `GET /measurements/latest?scripts=[scripts]`

  URL parameters:

  - `scripts` - comma-separated list of script names, required

  Same as `GET /measurements/{script}/{code}/latest` but allows to return latest measurements from multiple scripts at once.

## Development

Preferred way of development is to develop inside docker container. I do this in [VS Code](https://code.visualstudio.com/docs/remote/containers). There's a compose file for this purpose.

There's a [modd](https://github.com/cortesi/modd) tool installed in dev image, which enables live reloading and tests. Start it using `make run`.

If you want to develop on host machine, you'll need following tools installed on it (they're installed in docker image, see Dockerfile for more info):

- [libproj](https://proj.org/) shared library, to convert coordinate systems

Some tests require postgres. You cannot run them inside docker container (unless you want to mess with docker-inside-docker). They're excluded from main test set, I run them using `make test-nodocker` from host machine or CI environment.

### Writing scripts

Here are some recommendations for writing scripts for new sources

- Write tests, but when testing, **do not use** calls to real URLs, because unit tests can flood upstream with requests
- Round locations to 5 digits precision [link](https://en.wikipedia.org/wiki/Decimal_degrees), round levels and flows to what seems reasonable
- When converting coordinates, use `core.ToEPSG4326` utility function. It uses [PROJ](https://proj.org/) internally
- Use `core.Client` http client, which sets timeout, user-agent and has various helpers
- Do not bother with sorting results - this is done by script consumers
- Do not filter by `codes` and `since` inside worker. They are meant to be passed to upstream. Empty `codes` for all-at-once script must return all available measurements.
- Return null value (`nulltype.NullFloat64{}`) for level/flow when it's not provided
- Pay extra attention to time zones!
- Pass variables like access keys via script options
- Provide sample http requests (see `requests.http` files)
- Be forgiving when handling errors: only exit harvest function on real stoppers. If a single JSON object/CSV line causes error - log it then process next entry.

## TODO

- Build this using github actions _without_ docker. Problem: ubuntu 18.04 has old version of libproj-dev
- Virtual gauges
  - Statuses
  - What happens when one component is broken?
- Authorization
- Pushing
- Subscriptions
- Advanced scheduling, new harvest mode: batched
- Scripts as Go plugins
- Send logs to sentry
- Per-script binaries for third-party consumption
