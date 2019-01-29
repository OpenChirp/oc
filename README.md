# oc
This is a command-line utility for controlling an OpenChirp instance.

# Example Usage

## Services

### Create Service
```sh
$ oc -s https://api.openchirp.io -i myid -t mycomplextoken service create service1name "This is a simple example service"
5c4f886a3859df43aad0f1ef
```


## Generate Service Token
```sh
$ oc -s https://api.openchirp.io -i myid -t mycomplextoken service token generate 5c4f886a3859df43aad0f1ef
```

## Monitor a Service
```sh
$ oc -s https://api.openchirp.io -i myid -t mycomplextoken service monitor 5c4f886a3859df43aad0f1ef
```

# Config Files
You can save a local config file in any of the following locations:
- ~/.config/oc/occonfig.toml
- ~/.oc/occonfig.toml
- ./occonfig.toml

The content of the config file should look similar to the below:
```toml
framework-server = "https://api.openchirp.io"
mqtt-server = "tls://mqtt.openchirp.io:8883"
auth-id = "your-user-id"
auth-token = "your-user-token"
```