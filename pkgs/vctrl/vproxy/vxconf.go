package vproxy

var XrayConfStr string = `{
    "api": {
        "services": [
            "ReflectionService",
            "HandlerService",
            "LoggerService",
            "StatsService"
        ],
        "tag": "QV2RAY_API"
    },
    "dns": {
        "servers": [
            "1.1.1.1",
            "8.8.8.8",
            "8.8.4.4"
        ]
    },
    "fakedns": {
        "ipPool": "198.18.0.0/15",
        "poolSize": 65535
    },
    "inbounds": [
        {
			"port": 1080, 
			"listen": "127.0.0.1",
			"protocol": "socks",
			"settings": {
			  "udp": true
			}
		}
    ],
    "log": {
        "loglevel": "warning"
    },
    "outbounds": [
        {
            "protocol": "vmess",
            "sendThrough": "0.0.0.0",
            "settings": {
                "vnext": [
                    {
                        "address": "quranweb4.xyz",
                        "port": 443,
                        "users": [
                            {
                                "id": "8bbd91fe-a30b-4e29-bfc7-c28a44c0cb8f",
                                "security": "aes-128-gcm"
                            }
                        ]
                    }
                ]
            },
            "streamSettings": {
                "network": "ws",
                "security": "tls",
                "tlsSettings": {
                    "disableSystemRoot": false
                },
                "wsSettings": {
                    "path": "/current_time"
                },
                "xtlsSettings": {
                    "disableSystemRoot": false
                }
            },
            "tag": "PROXY"
        },
        {
            "protocol": "freedom",
            "sendThrough": "0.0.0.0",
            "settings": {
                "domainStrategy": "AsIs",
                "redirect": ":0"
            },
            "streamSettings": {
            },
            "tag": "DIRECT"
        },
        {
            "protocol": "blackhole",
            "sendThrough": "0.0.0.0",
            "settings": {
                "response": {
                    "type": "none"
                }
            },
            "streamSettings": {
            },
            "tag": "BLACKHOLE"
        }
    ],
    "policy": {
        "system": {
            "statsInboundDownlink": true,
            "statsInboundUplink": true,
            "statsOutboundDownlink": true,
            "statsOutboundUplink": true
        }
    },
    "routing": {
        "domainMatcher": "mph",
        "domainStrategy": "AsIs",
        "rules": [
            {
                "inboundTag": [
                    "QV2RAY_API_INBOUND"
                ],
                "outboundTag": "QV2RAY_API",
                "type": "field"
            },
            {
                "ip": [
                    "geoip:private"
                ],
                "outboundTag": "DIRECT",
                "type": "field"
            },
            {
                "ip": [
                    "geoip:cn"
                ],
                "outboundTag": "DIRECT",
                "type": "field"
            },
            {
                "domain": [
                    "geosite:cn"
                ],
                "outboundTag": "DIRECT",
                "type": "field"
            }
        ]
    },
    "stats": {
    }
}
`

const (
	ConfStr string = `{
		"log": {
		  "loglevel": "warning"
		},
		"dns": {
		  "servers": [
			{
			  "address": "1.1.1.1",
			  "domains": ["geosite:geolocation-!cn"]
			},
			{
			  "address": "223.5.5.5",
			  "domains": ["geosite:cn"],
			  "expectIPs": ["geoip:cn"]
			},
			{
			  "address": "114.114.114.114",
			  "domains": ["geosite:cn"]
			},
			"localhost"
		  ]
		},
		"routing": {
		  "domainStrategy": "IPIfNonMatch",
		  "rules": [
			{
			  "type": "field",
			  "domain": ["geosite:category-ads-all"],
			  "outboundTag": "block"
			},
			{
			  "type": "field",
			  "domain": ["geosite:cn"],
			  "outboundTag": "direct"
			},
			{
			  "type": "field",
			  "ip": ["geoip:cn", "geoip:private"],
			  "outboundTag": "direct"
			},
			{
			  "type": "field",
			  "domain": ["geosite:geolocation-!cn"],
			  "outboundTag": "proxy"
			},
			{
			  "type": "field",
			  "ip": ["223.5.5.5"],
			  "outboundTag": "direct"
			}
		  ]
		},
		"inbounds": [
		  {
			"tag": "socks-in",
			"protocol": "socks",
			"listen": "127.0.0.1",
			"port": 10800,
			"settings": {
			  "udp": true
			}
		  },

		  {
			"tag": "http-in",
			"protocol": "http",
			"listen": "127.0.0.1",
			"port": 10801 
		  }
		],
		"outbounds": [
			{
				"protocol": "vmess",
				"sendThrough": "0.0.0.0",
				"settings": {
					"vnext": [
						{
							"address": "quranweb4.xyz",
							"port": 443,
							"users": [
								{
									"id": "8bbd91fe-a30b-4e29-bfc7-c28a44c0cb8f",
									"security": "aes-128-gcm"
								}
							]
						}
					]
				},
				"streamSettings": {
					"network": "ws",
					"security": "tls",
					"tlsSettings": {
						"disableSystemRoot": false
					},
					"wsSettings": {
						"path": "/current_time"
					},
					"xtlsSettings": {
						"disableSystemRoot": false
					}
				},
				"tag": "PROXY"
			},
			{
				"protocol": "freedom",
				"sendThrough": "0.0.0.0",
				"settings": {
					"domainStrategy": "AsIs",
					"redirect": ":0"
				},
				"streamSettings": {
				},
				"tag": "DIRECT"
			},
			{
				"protocol": "blackhole",
				"sendThrough": "0.0.0.0",
				"settings": {
					"response": {
						"type": "none"
					}
				},
				"streamSettings": {
				},
				"tag": "BLACKHOLE"
			}
		]
	  }`
)
