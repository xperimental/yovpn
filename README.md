# yovpn

**Note:** I have not worked on this project for _a long time_. I'm not actually using this as a VPN server implementation anymore (I've switched to Wireguard for my VPN needs), but reimplementing this project with Wireguard might be something I'd do in the future at some point. For now, I recommend not to use this software when you actually need a VPN, but only for educational or entertainment purposes.

yovpn is a server that can be used to create VPN endpoints on public cloud infrastructure. Currently only DigitalOcean is supported as a provider.

The aim of VPNs created using this tool is mainly to avoid IP geolocation and not provide secure internet access. As such, the security of the VPN tunnel is not very high and no public-key cryptography is used to keep setup simple.

This tool will only provide the VPN server. You still need an [OpenVPN](https://openvpn.net/) client for your platform to connect to the provisioned endpoints.

*Note:* This tool is still in the early stages. While the `yovpn-cli` command line version is already quite useable, the HTTP server still has many unimplemented features (most notably authentication and secure transfer of VPN credentials) and should be considered unsafe.

## Usage

First you will need credentials for [DigitalOcean](https://www.digitalocean.com/):

1. Create an account on digitalocean
2. Go to the [API settings](https://cloud.digitalocean.com/settings/applications)
3. Create a "personal access token" with write access (the name does not matter)
4. Take note of the token (a long hexadecimal number)

Now you need to decide if you want to use the command-line or start a HTTP server

### Command-line

Assuming a correctly setup Go environment you can get the sources and build the command-line client with the following command:

```bash
go get -u github.com/xperimental/yovpn/cmd/yovpn-cli
```

You should now have a `yovpn-cli` binary in your path. You can get usage information using:

```bash
yovpn-cli -help
```

To get available regions for new VPN endpoints run:

```bash
yovpn-cli -token $token
```

You can now provision a new VPN endpoint in one of the available regions. The tool will create a new virtual server, run configuration commands on it and write the client configuration to a file on your computer:

```bash
yovpn-cli -token $token -region $region -output yovpn.ovpn
# For example
yovpn-cli -token c0ff33 -region nyc1 -output yovpn.ovpn
```

When the command completes successfully (takes approximately 2-5min) the directory you ran the command in should contain a new file called `yovpn.ovpn` (or what you selected with the `-output` parameter). That file should contain the client configuration needed for your VPN.

### HTTP server

Assuming a correctly setup Go environment you can get the sources and build the HTTP server with the following command:

```bash
go get -u github.com/xperimental/yovpn/cmd/yovpn-server
```

You should now have a `yovpn-server` binary in your path. The HTTP server can be started using the following command:

```bash
yovpn-server -port $port -token $token
# For example
yovpn-server -port 8080 -token c0ff33
```

*Note:* The HTTP server is currently not suitable for use as very important features are missing.

## Heroku

If you are really daring you can already deploy the HTTP server to [Heroku](https://www.heroku.com).
You will still need to create the DigitalOcean token beforehand.

[![Deploy](https://www.herokucdn.com/deploy/button.png)](https://heroku.com/deploy?template=https://github.com/xperimental/yovpn/tree/heroku)

## HTTP Endpoints

The HTTP server currently has the following endpoints:

| Method | Path                       | Description                                 |
|--------|----------------------------|---------------------------------------------|
| GET    | `/regions`                 | List all available regions.                 |
| GET    | `/cleanup`                 | Remove all known endpoints.                  |
| GET    | `/endpoint/:id`            | Return information about endpoint with `id` |
| PUT    | `/endpoint?region=:region` | Create a new endpoint in selected region.   |
| DELETE | `/endpoint/:id`            | Remove endpoint with `id`                   |
