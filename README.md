# Website Change Notifier
A simple go app, that continously monitors a website and notifies you if it changes.

Currently, you need [Pushover](https://pushover.net/) for notifications.
It is a neat little app, which lets you send notifications to your phone via an API. There is a free trial, and a perpetual person license is super cheap at ~5USD.

## Usage
The app is configured using a rather simple file called `config.yml`. By default the app will look for it in the root directory `/`. You can customize the location of the config file via an env-var like so: `CONFIG_FILE=/my/config/path/my-config.yml`.

The structure should be pretty self-explanatory. Note that the durations have to be written in the correct format for go to parse it. Note that `h` for hour is the largest unit supported Here are some examples:


| Time        | Syntax           |
| ------------- |:-------------:|
| 10 Minutes      | 10m |
| 24 Hours      | 24h      |
| 1 Week | 168h      |

The sample config file can be found in the root directory of the project as `test-config.yml`

```yaml
site:
  name: DemoSite
  url: http://localhost:8080
  interval: 1m
  summary-interval: 24h
pushover:
  token: <your pushover token>
  user: <your pushover user>
  device: <your pushover device>
log:
  level: INFO
```

### Docker
There is also a docker image available under `jjnp/change-detector` via [Docker Hub](https://hub.docker.com/repository/docker/jjnp/change-detector).

To run simply volume mount the config file.

## Roadmap / Future features
Feature Requests and PRs are of course welcome. Since I built this project for myself, I can't say when I will expand on it. That will most likely depend on demand.
So far the roadmap includes the following features:
- [x] logging
- [x] summary notifications (so you know the app is still running)
- [ ] support different notifiers, most importantly email + webhook
- [ ] support looking for changes in multiple websites at once
- [ ] generally more config options
- [ ] more complex checking logic to make it work with PWAs/modern web framework sites
- [ ] allow monitoring uptime instead of change detection i.e. get notified if your website goes down

## Limitations
Please note that the checking is simply done by requesting the site at the URL and then calculating a hash of the response body.

This means that PWAs or modern sites powered by Angular and the like may not work as you expect. This is on the roadmap, but not a personal priority.
If this is very important to you, please let me know via feature requests.