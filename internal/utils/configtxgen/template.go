/*
 * Copyright 2020. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package configtxgen

var (
	tmplGenesis = `{
  "channel_group": {
    "groups": {
      "Application": {
        "groups": {
          {{template "orderer_organization" .Organization}}
        },
        "mod_policy": "Admins",
        "policies": {
          "Admins": {
            "mod_policy": "Admins",
            "policy": {
              "type": 3,
              "value": {
                "rule": "MAJORITY",
                "sub_policy": "Admins"
              }
            },
            "version": "0"
          },
          "Readers": {
            "mod_policy": "Admins",
            "policy": {
              "type": 3,
              "value": {
                "rule": "ANY",
                "sub_policy": "Readers"
              }
            },
            "version": "0"
          },
          "Writers": {
            "mod_policy": "Admins",
            "policy": {
              "type": 3,
              "value": {
                "rule": "ANY",
                "sub_policy": "Writers"
              }
            },
            "version": "0"
          }
        },
        "values": {
          "Capabilities": {
            "mod_policy": "Admins",
            "value": {
              "capabilities": {
                "V1_4_2": {}
              }
            },
            "version": "0"
          }
        },
        "version": "0"
      },
      "Consortiums": {
        "groups": {
          "{{.ConsortiumName}}": {
            "groups": {},
            "mod_policy": "/Channel/Orderer/Admins",
            "policies": {},
            "values": {
              "ChannelCreationPolicy": {
                "mod_policy": "/Channel/Orderer/Admins",
                "value": {
                  "type": 3,
                  "value": {
                    "rule": "ANY",
                    "sub_policy": "Admins"
                  }
                },
                "version": "0"
              }
            },
            "version": "0"
          }
        },
        "mod_policy": "/Channel/Orderer/Admins",
        "policies": {
          "Admins": {
            "mod_policy": "/Channel/Orderer/Admins",
            "policy": {
              "type": 1,
              "value": {
                "identities": [],
                "rule": {
                  "n_out_of": {
                    "n": 0,
                    "rules": []
                  }
                },
                "version": 0
              }
            },
            "version": "0"
          }
        },
        "values": {},
        "version": "0"
      },
      "Orderer": {
        "groups": {
          {{template "orderer_organization" .Organization}}
        },
        "mod_policy": "Admins",
        "policies": {
          "Admins": {
            "mod_policy": "Admins",
            "policy": {
              "type": 3,
              "value": {
                "rule": "MAJORITY",
                "sub_policy": "Admins"
              }
            },
            "version": "0"
          },
          "BlockValidation": {
            "mod_policy": "Admins",
            "policy": {
              "type": 3,
              "value": {
                "rule": "ANY",
                "sub_policy": "Writers"
              }
            },
            "version": "0"
          },
          "Readers": {
            "mod_policy": "Admins",
            "policy": {
              "type": 3,
              "value": {
                "rule": "ANY",
                "sub_policy": "Readers"
              }
            },
            "version": "0"
          },
          "Writers": {
            "mod_policy": "Admins",
            "policy": {
              "type": 3,
              "value": {
                "rule": "ANY",
                "sub_policy": "Writers"
              }
            },
            "version": "0"
          }
        },
        "values": {
          "BatchSize": {
            "mod_policy": "Admins",
            "value": {
              "absolute_max_bytes": 103809024,
              "max_message_count": 10,
              "preferred_max_bytes": 524288
            },
            "version": "0"
          },
          "BatchTimeout": {
            "mod_policy": "Admins",
            "value": {
              "timeout": "2s"
            },
            "version": "0"
          },
          "Capabilities": {
            "mod_policy": "Admins",
            "value": {
              "capabilities": {
                "V1_4_2": {}
              }
            },
            "version": "0"
          },
          "ChannelRestrictions": {
            "mod_policy": "Admins",
            "value": null,
            "version": "0"
          },
          {{template "consensus_type" .Consensus}}
        },
        "version": "0"
      }
    },
    "mod_policy": "Admins",
    "policies": {
      "Admins": {
        "mod_policy": "Admins",
        "policy": {
          "type": 3,
          "value": {
            "rule": "MAJORITY",
            "sub_policy": "Admins"
          }
        },
        "version": "0"
      },
      "Readers": {
        "mod_policy": "Admins",
        "policy": {
          "type": 3,
          "value": {
            "rule": "ANY",
            "sub_policy": "Readers"
          }
        },
        "version": "0"
      },
      "Writers": {
        "mod_policy": "Admins",
        "policy": {
          "type": 3,
          "value": {
            "rule": "ANY",
            "sub_policy": "Writers"
          }
        },
        "version": "0"
      }
    },
    "values": {
      "BlockDataHashingStructure": {
        "mod_policy": "Admins",
        "value": {
          "width": 4294967295
        },
        "version": "0"
      },
      "Capabilities": {
        "mod_policy": "Admins",
        "value": {
          "capabilities": {
            "V1_4_2": {}
          }
        },
        "version": "0"
      },
      "HashingAlgorithm": {
        "mod_policy": "Admins",
        "value": {
          "name": "SHA256"
        },
        "version": "0"
      },
      {{template "orderer_addresses" .OrdererAddresses}}
    },
    "version": "0"
  },
  "sequence": "0"
}`
	tmplPeerOrganization = `{
    "groups":{},
    "mod_policy":"Admins",
    "policies":{
        "Admins":{
            "mod_policy":"Admins",
            "policy":{
                "type":1,
                "value":{
                    "identities":[
                        {
                            "principal":{
                                "msp_identifier":"{{.MSPIdentifier}}",
                                "role":"ADMIN"
                            },
                            "principal_classification":"ROLE"
                        }
                    ],
                    "rule":{
                        "n_out_of":{
                            "n":1,
                            "rules":[
                                {
                                    "signed_by":0
                                }
                            ]
                        }
                    },
                    "version":0
                }
            },
            "version":"0"
        },
        "Readers":{
            "mod_policy":"Admins",
            "policy":{
                "type":1,
                "value":{
                    "identities":[
                        {
                            "principal":{
                                "msp_identifier":"{{.MSPIdentifier}}",
                                "role":"ADMIN"
                            },
                            "principal_classification":"ROLE"
                        },
                        {
                            "principal":{
                                "msp_identifier":"{{.MSPIdentifier}}",
                                "role":"PEER"
                            },
                            "principal_classification":"ROLE"
                        },
                        {
                            "principal":{
                                "msp_identifier":"{{.MSPIdentifier}}",
                                "role":"CLIENT"
                            },
                            "principal_classification":"ROLE"
                        }
                    ],
                    "rule":{
                        "n_out_of":{
                            "n":1,
                            "rules":[
                                {
                                    "signed_by":0
                                },
                                {
                                    "signed_by":1
                                },
                                {
                                    "signed_by":2
                                }
                            ]
                        }
                    },
                    "version":0
                }
            },
            "version":"0"
        },
        "Writers":{
            "mod_policy":"Admins",
            "policy":{
                "type":1,
                "value":{
                    "identities":[
                        {
                            "principal":{
                                "msp_identifier":"{{.MSPIdentifier}}",
                                "role":"ADMIN"
                            },
                            "principal_classification":"ROLE"
                        },
                        {
                            "principal":{
                                "msp_identifier":"{{.MSPIdentifier}}",
                                "role":"CLIENT"
                            },
                            "principal_classification":"ROLE"
                        }
                    ],
                    "rule":{
                        "n_out_of":{
                            "n":1,
                            "rules":[
                                {
                                    "signed_by":0
                                },
                                {
                                    "signed_by":1
                                }
                            ]
                        }
                    },
                    "version":0
                }
            },
            "version":"0"
        }
    },
    "values":{
        "MSP":{
            "mod_policy":"Admins",
            "value":{
                "config":{
                    "admins":[
                        "{{.AdminCert}}"
                    ],
                    "crypto_config":{
                        "identity_identifier_hash_function":"SHA256",
                        "signature_hash_family":"SHA2"
                    },
                    "fabric_node_ous":{
                        "admin_ou_identifier":null,
                        "client_ou_identifier":{
                            "certificate":"{{.SignRootCert}}",
                            "organizational_unit_identifier":"client"
                        },
                        "enable":true,
                        "orderer_ou_identifier":null,
                        "peer_ou_identifier":{
                            "certificate":"{{.SignRootCert}}",
                            "organizational_unit_identifier":"peer"
                        }
                    },
                    "intermediate_certs":[

                    ],
                    "name":"{{.MSPIdentifier}}",
                    "organizational_unit_identifiers":[

                    ],
                    "revocation_list":[

                    ],
                    "root_certs":[
                        "{{.SignRootCert}}"
                    ],
                    "signing_identity":null,
                    "tls_intermediate_certs":[

                    ],
                    "tls_root_certs":[
                        "{{.TLSRootCert}}"
                    ]
                },
                "type":0
            },
            "version":"0"
        }
    },
    "version":"0"
}`
	tmplAnchorPeer = `{
    "channel_id":"{{.ChannelName}}",
    "isolated_data":{

    },
    "read_set":{
        "groups":{
            "Application":{
                "groups":{
                    "{{.OrganizationName}}":{
                        "groups":{

                        },
                        "mod_policy":"",
                        "policies":{
                            "Admins":{
                                "mod_policy":"",
                                "policy":null,
                                "version":"0"
                            },
                            "Readers":{
                                "mod_policy":"",
                                "policy":null,
                                "version":"0"
                            },
                            "Writers":{
                                "mod_policy":"",
                                "policy":null,
                                "version":"0"
                            }
                        },
                        "values":{
                            "MSP":{
                                "mod_policy":"",
                                "value":null,
                                "version":"0"
                            }
                        },
                        "version":"0"
                    }
                },
                "mod_policy":"",
                "policies":{

                },
                "values":{

                },
                "version":"1"
            }
        },
        "mod_policy":"",
        "policies":{

        },
        "values":{

        },
        "version":"0"
    },
    "write_set":{
        "groups":{
            "Application":{
                "groups":{
                    "{{.OrganizationName}}":{
                        "groups":{

                        },
                        "mod_policy":"Admins",
                        "policies":{
                            "Admins":{
                                "mod_policy":"",
                                "policy":null,
                                "version":"0"
                            },
                            "Readers":{
                                "mod_policy":"",
                                "policy":null,
                                "version":"0"
                            },
                            "Writers":{
                                "mod_policy":"",
                                "policy":null,
                                "version":"0"
                            }
                        },
                        "values":{
                            "AnchorPeers":{
                                "mod_policy":"Admins",
                                "value":{
                                    "anchor_peers": {{.Peers | jsonMarshal}}
                                },
                                "version":"0"
                            },
                            "MSP":{
                                "mod_policy":"",
                                "value":null,
                                "version":"0"
                            }
                        },
                        "version":"1"
                    }
                },
                "mod_policy":"",
                "policies":{

                },
                "values":{

                },
                "version":"1"
            }
        },
        "mod_policy":"",
        "policies":{

        },
        "values":{

        },
        "version":"0"
    }
}`
	tmplChannel = `{
    "channel_id":"{{.Name}}",
    "isolated_data":{

    },
    "read_set":{
        "groups":{
            "Application":{
                "groups":{{.Organizations | jsonMarshal}},
                "mod_policy":"",
                "policies":{

                },
                "values":{

                },
                "version":"0"
            }
        },
        "mod_policy":"",
        "policies":{

        },
        "values":{
            "Consortium":{
                "mod_policy":"",
                "value":null,
                "version":"0"
            }
        },
        "version":"0"
    },
    "write_set":{
        "groups":{
            "Application":{
                "groups":{{.Organizations | jsonMarshal}},
                "mod_policy":"Admins",
                "policies":{
                    "Admins":{
                        "mod_policy":"Admins",
                        "policy":{
                            "type":3,
                            "value":{
                                "rule":"MAJORITY",
                                "sub_policy":"Admins"
                            }
                        },
                        "version":"0"
                    },
                    "Readers":{
                        "mod_policy":"Admins",
                        "policy":{
                            "type":3,
                            "value":{
                                "rule":"ANY",
                                "sub_policy":"Readers"
                            }
                        },
                        "version":"0"
                    },
                    "Writers":{
                        "mod_policy":"Admins",
                        "policy":{
                            "type":3,
                            "value":{
                                "rule":"ANY",
                                "sub_policy":"Writers"
                            }
                        },
                        "version":"0"
                    }
                },
                "values":{
                    "Capabilities":{
                        "mod_policy":"Admins",
                        "value":{
                            "capabilities":{
                                "V1_4_2":{

                                }
                            }
                        },
                        "version":"0"
                    }
                },
                "version":"1"
            }
        },
        "mod_policy":"",
        "policies":{

        },
        "values":{
            "Consortium":{
                "mod_policy":"",
                "value":{
                    "name":"{{.ConsortiumName}}"
                },
                "version":"0"
            }
        },
        "version":"0"
    }
}`
	tmplOrdererOrganization = `{{define "orderer_organization"}}"{{.Name}}":{
    "groups":{},
    "mod_policy":"Admins",
    "policies":{
        "Admins":{
            "mod_policy":"Admins",
            "policy":{
                "type":1,
                "value":{
                    "identities":[
                        {
                            "principal":{
                                "msp_identifier":"{{.MSPIdentifier}}",
                                "role":"ADMIN"
                            },
                            "principal_classification":"ROLE"
                        }
                    ],
                    "rule":{
                        "n_out_of":{
                            "n":1,
                            "rules":[
                                {
                                    "signed_by":0
                                }
                            ]
                        }
                    },
                    "version":0
                }
            },
            "version":"0"
        },
        "Readers":{
            "mod_policy":"Admins",
            "policy":{
                "type":1,
                "value":{
                    "identities":[
                        {
                            "principal":{
                                "msp_identifier":"{{.MSPIdentifier}}",
                                "role":"MEMBER"
                            },
                            "principal_classification":"ROLE"
                        }
                    ],
                    "rule":{
                        "n_out_of":{
                            "n":1,
                            "rules":[
                                {
                                    "signed_by":0
                                }
                            ]
                        }
                    },
                    "version":0
                }
            },
            "version":"0"
        },
        "Writers":{
            "mod_policy":"Admins",
            "policy":{
                "type":1,
                "value":{
                    "identities":[
                        {
                            "principal":{
                                "msp_identifier":"{{.MSPIdentifier}}",
                                "role":"MEMBER"
                            },
                            "principal_classification":"ROLE"
                        }
                    ],
                    "rule":{
                        "n_out_of":{
                            "n":1,
                            "rules":[
                                {
                                    "signed_by":0
                                }
                            ]
                        }
                    },
                    "version":0
                }
            },
            "version":"0"
        }
    },
    "values":{
        "MSP":{
            "mod_policy":"Admins",
            "value":{
                "config":{
                    "admins":[
                        "{{.AdminCert}}"
                    ],
                    "crypto_config":{
                        "identity_identifier_hash_function":"SHA256",
                        "signature_hash_family":"SHA2"
                    },
                    "fabric_node_ous":null,
                    "intermediate_certs":[

                    ],
                    "name":"{{.MSPIdentifier}}",
                    "organizational_unit_identifiers":[

                    ],
                    "revocation_list":[

                    ],
                    "root_certs":[
                        "{{.SignRootCert}}"
                    ],
                    "signing_identity":null,
                    "tls_intermediate_certs":[

                    ],
                    "tls_root_certs":[
                        "{{.TLSRootCert}}"
                    ]
                },
                "type":0
            },
            "version":"0"
        }
    },
    "version":"0"
}{{end}}`
	tmplConsensusType = `{{define "consensus_type"}}"ConsensusType":{
    "mod_policy":"Admins",
    "value":{
        "metadata":{{if eq .Type "etcdraft"}}{
            "consenters":{{.Consenters | jsonMarshal}},
            "options":{
                "election_tick":10,
                "heartbeat_tick":1,
                "max_inflight_blocks":5,
                "snapshot_interval_size":20971520,
                "tick_interval":"500ms"
            }
        }{{else}}null{{end}},
        "state":"STATE_NORMAL",
        "type":"{{.Type}}"
    },
    "version":"0"
}{{if eq .Type "kafka"}},
"KafkaBrokers": {
    "mod_policy": "Admins",
    "value": {
        "brokers": {{.KafkaBrokers | jsonMarshal}}
    },
    "version": "0"
}{{end}}{{end}}`
	tmplOrdererAddresses = `{{define "orderer_addresses"}}"OrdererAddresses":{
    "mod_policy":"/Channel/Orderer/Admins",
    "value":{
        "addresses":{{. | jsonMarshal}}
    },
    "version":"0"
}{{end}}`
)
