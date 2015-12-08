package config

// CloudConfig contains the configuration that is sent to the server to start the configuration process.
const CloudConfig = `
#cloud-config
packages:
    - openvpn
write_files:
  - path: /etc/openvpn/yovpn-server.conf
    content: |
        dev tun
        secret secret.key
        user nobody
        group nogroup
        ifconfig 10.42.86.1 10.42.86.2
        cipher AES-128-CBC
        comp-lzo
        persist-key
        persist-tun
  - path: /root/enable-nat.sh
    content: |
        #!/bin/bash
        NAT_INTERN=tap0
        NAT_EXTERN=eth0
        echo 1 > /proc/sys/net/ipv4/ip_forward
        iptables -t nat -A POSTROUTING -o $NAT_EXTERN -j MASQUERADE
        iptables -A FORWARD -i $NAT_INTERN -o $NAT_EXTERN -m state --state RELATED,ESTABLISHED -j ACCEPT
        iptables -A FORWARD -i $NAT_EXTERN -o $NAT_INTERN -j ACCEPT
runcmd:
    - echo "net.ipv4.ip_forward=1" >> /etc/sysctl.conf
    - sysctl -w net.ipv4.ip_forward=1
    - sh /root/enable-nat.sh
    - openvpn --genkey --secret /etc/openvpn/secret.key
    - service openvpn restart
    - touch /root/yovpn.ready
`
