# ovpn-po-notify
This sends a notification over Pushover whenever a user connects to their OpenVPN.
This assumes you have device names and have set up a configuration file thusly at `./config.json`:

```
{
    "user":"YOUR_USER_TOKEN",
    "app":"YOUR_APP_SPECIFIC_TOKEN",
    "location":"THE_LOCATION_OF_YOUR_LOG_APPEND_FILE"
}
```
# Build Instructions
Use `go build`.
