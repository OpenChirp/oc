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