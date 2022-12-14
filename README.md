# Columbus server

The goal of the Columbus project to provide an API to discover and store every domain's every subdomain and server it nearly instant.

> :heavy_exclamation_mark: This project and the database behind it is new, it takes some time to be usable.

```bash
time curl 'https://columbus.elmasy.com/lookup/github.com'
...
real	0m0.270s
user	0m0.024s
sys     0m0.012s
```

## Usage

By default Columbus returns only the subdomain in a JSON string array:
```bash
curl 'https://columbus.elmasy.com/lookup/github.com'
```

But we think of the bash lovers, so if you don't want to mess with JSON and a newline separated list is your wish, then include the `Accept: text/plain` header.
```bash
DOMAIN="github.com"

curl -s -H "Accept: text/plain" "https://columbus.elmasy.com/lookup/$DOMAIN" | \
while read SUB
do
        if [[ "$SUB" == "" ]]
        then
                HOST="$DOMAIN"
        else
                HOST="${SUB}.${DOMAIN}"
        fi
        echo "$HOST"
done
```

**For more, see the [OpenAPI specification](https://columbus.elmasy.com/openapi.yaml)**

## Entries

Currently, entries are got from [Certificate Transparency](https://certificate.transparency.dev/).

Check the currently parsed CT logs [here](https://status.elmasy.com/status/4803b934327a1168b515).

## Build

```bash
git clone https://github.com/elmasy-com/columbus-server
make build
```

## Install

Create a new user:

```bash
adduser --system --no-create-home --disabled-login columbus
```

Copy the binary to `/usr/bin/columbus-server`.

Make it executable:
```bash
chmod +x /usr/bin/columbus-server
```

Create a directory:
```bash
mkdir /etc/columbus
```

Copy the config file to `/etc/columbus/columbus.conf`.

Set the permission to 0600.
```bash
chmod -R 0600 /etc/columbus
```

Set the owner of the config file:
```bash
chown -R columbus /etc/columbus
```

Install the service file (eg.: `/etc/systemd/system/columbus-server.service`).
```bash
cp columbus-server.service /etc/systemd/system/
```

Reload systemd:
```bash
systemctl daemon-reload
```

Start columbus:
```
systemctl start columbus-server
```

If you want to columbus start automatically:
```
systemctl enable columbus-server
```
 