package share

var ClientTemplate string = `
dev tun
user nobody
group nogroup
ifconfig 10.42.86.2 10.42.86.1
cipher AES-128-CBC
comp-lzo
persist-key
persist-tun
route 0.0.0.0 0.0.0.0 vpn_gateway
`
