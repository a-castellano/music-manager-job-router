# Job Router

[Project's page](https://musicmanager.gitpages.windmaker.net/Job-Router)

[Actual Repo](https://git.windmaker.net/musicmanager/Job-Router)

 [![pipeline status](https://git.windmaker.net/musicmanager/Job-Router/badges/master/pipeline.svg)](https://git.windmaker.net/musicmanager/Job-Router/-/commits/master) [![coverage report](https://git.windmaker.net/musicmanager/Job-Router/badges/master/coverage.svg)](https://git.windmaker.net/musicmanager/Job-Router/-/commits/master) [![Quality Gate Status](https://sonarqube.windmaker.net/api/project_badges/measure?project=music-manager-job-router&metric=alert_status)](https://sonarqube.windmaker.net/dashboard?id=music-manager-job-router)

Service that routes jobs to Wrappers and Job Manager. When job finishes status is sended to **Status Manager**, if job finishes successfully it is also sended to **Storage Manager**.

See [Job Routing Docs](https://musicmanager.gitpages.windmaker.net/Music-Manager-Docs/job-routing/) for more info.

## Config example

```toml
[server]

host = "localhost"
port = 5672
user = "guest"
password = "pass"

[wrappers]

  [wrappers.firstwrapper]
  name = "firstwrapper"
  durable = true
  delete_when_unused = false
  exclusive = false
  no_wait = false
  auto_ack = false
  
  [wrappers.secondwrapper]
  name = "secondwrapper"
  durable = true
  delete_when_unused = false
  exclusive = false
  no_wait = false
  auto_ack = false

[jobmanager]
name = "jobmanager"
durable = true
delete_when_unused = false
exclusive = false
no_wait = false
auto_ack = false

[status]
name = "status"

[storage]
name = "storage"
```
