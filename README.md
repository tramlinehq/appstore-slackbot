<img height="240" width="240" alt="logo" src="https://github.com/tramlinehq/appstore-slackbot/assets/50663/e7f61d61-20d2-4f49-8582-d60dd2775712" />

# App Store Slackbot

[![License](https://img.shields.io/badge/license-MIT-green.svg?style=flat)](https://github.com/tramlinehq/appstore-slackbot/blob/master/LICENSE)
[![Discord](https://img.shields.io/discord/974284993641725962?label=discord%20chat)](https://discord.gg/u7VwyvBV2Z)

Talk to the App Store directly from your Slack workspace.

Head over to 🌏 [appstoreslackbot.com](https://appstoreslackbot.com) and setup your account to get started.

## Development 👩‍💻

* Install `go`
* Setup local certs.

```
bin/ssl
```

* Create config.
```
cp .env.sample .env
```

Change the values of the config to be meaningful.

* Run the service.
```
go run .
```


## Thanks 🥰

### Uses
- [gin](https://github.com/gin-gonic/gin "gin")
- [applelink](https://github.com/tramlinehq/applelink "applelink") 
- [picniccss](https://picnicss.com/ "picnic-css") 

### Infrastructure 
Deployed on [render.com](https://render.com) 
