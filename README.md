# Job Router

[Project's page](https://musicmanager.gitpages.windmaker.net/Job-Router)

[Actual Repo](https://git.windmaker.net/musicmanager/Job-Router)

 [![pipeline status](https://git.windmaker.net/musicmanager/Job-Router/badges/master/pipeline.svg)](https://git.windmaker.net/musicmanager/Job-Router/-/commits/master) [![coverage report](https://git.windmaker.net/musicmanager/Job-Router/badges/master/coverage.svg)](https://git.windmaker.net/musicmanager/Job-Router/-/commits/master) [![Quality Gate Status](https://sonarqube.windmaker.net/api/project_badges/measure?project=music-manager-job-router&metric=alert_status)](https://sonarqube.windmaker.net/dashboard?id=music-manager-job-router)

Service that routes jobs to Wrappers and Job Manager. When job finishes status is sended to **Status Manager**, if job finishes successfully it is also sended to **Storage Manager**.

See [Job Routing Docs](https://musicmanager.gitpages.windmaker.net/Music-Manager-Docs/job-routing/) for more info.

## Service Config

This service requires the following config:

### sever
Contains Rabbitmq server access credentials.

### wrappers
Contains Rabbitmq queue configuration for each wrapper that will consume jobs.

### jobmanager
Contains Rabbitmq queue configuration for jobs queue where JobManager sends jobs to be routed by JobRouter

### wrapperoutput
Contains Rabbitmq queue configuration for jobs queue where wrappers send jobs to be routed or finished by JobRouter


### status
Contains StatusManager service name

### storage
Contains StorageManager service name

## Config example
This service will look for its config in **/etc/music-manager/config.toml**, parent folder can be changed setting the environment variable **MUSIC_MANAGER_SERVICE_CONFIG_FILE_LOCATION**.
```toml
[server]

host = "localhost"
port = 5672
user = "guest"
password = "pass"

[wrappers]

  [wrappers.firstwrapper]
  name = "firstwrapper"
  
  [wrappers.secondwrapper]
  name = "secondwrapper"

[wrapperoutput]
name = "wrapperoutput"

[jobmanager]
name = "jobmanager"
durable = true

[status]
name = "status"

[storage]
name = "storage"

```
