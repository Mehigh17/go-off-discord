<div align="center">
    <h1>Go Off Discord (GOD)</h1>
    <div>
        <a href="#usage">Usage</a> •
        <a href="#notice">Notice</a>
    </div>
</div>

![License](https://img.shields.io/github/license/Mehigh17/go-off-discord)
![Issues](https://img.shields.io/github/issues/Mehigh17/SharpCatch)

GOD allows you to remove all your messages on an entire discord server or just a single channel with a simple command.

# Usage

## Creating an account configuration

Create a json file in which you will put your authentication token and your user id, just like this:

```json
{
    "authToken": "",
    "userId": ""
}
```

## Removing messages from a channel

`go-off-discord del -a accountInfo.json channel --id channel_id_here`

## Removing messages from a server

`go-off-discord del -a accountInfo.json server --id server_id_here`

# Notice

You will not be able to remove the messages that aren't readable to you, for example if you sent messages to a channel then you've been removed access to type there, you will no longer be able to remove them until you're granted read access again.