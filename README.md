# OfflineRSS

OfflineRSS is a tool written in Go for syncing [RSS](https://en.wikipedia.org/wiki/RSS) between a hosted [feed aggregator](https://en.wikipedia.org/wiki/News_aggregator) and local RSS readers.

It generates databases for supported RSS readers that can be used offline. Any changes done to these databases will be synced back to the aggregator next time the program is run whilst also updating the local databases with everything new.

## Why?

* Allows you to use any supported reader with any supported feed aggregator
* Allows you to use your favorite reader offline. Great if you don't always have an internet connection.
* Makes all user interaction feel instant because they are synced later.

## Supported servers (feed aggregators)

* [Miniflux](https://miniflux.app) Apache-2.0
* [NewsBlur](https://newsblur.com) MIT

## Supported clients (readers)

* [FeedReader](https://jangernert.github.io/FeedReader) GPL-3.0
* [Newsboat](https://newsboat.org) MIT
* [Newsraft](https://codeberg.org/newsraft/newsraft) ISC
* [QuiteRSS](https://quiterss.org) GPL-3.0

## Get started

Install OfflineRSS (it will install to `~/go/bin/offlinerss` unless changed)

```
go install github.com/limero/offlinerss@latest
```

You can now run OfflineRSS and it will prompt you for what server and clients to use and save this information in `~/.config/offlinerss/config.json`. It will then create local databases for your chosen client(s) and symlink these to each client. If you have used any of the selected clients previously, their old databases will be renamed to `.bak` and kept.

By default, the paths for the generated databases will be in `~/.local/share/offlinerss`.

## Limitations

* Only one server can be used at a time. However, all clients can be enabled simultaneously.
* Clients can't be synced with each other, without syncing with the server first.

## License

Licensed under the MIT license. (See LICENSE)
