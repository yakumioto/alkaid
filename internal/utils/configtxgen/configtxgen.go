/*
 * Copyright 2020. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package configtxgen

import (
	"bytes"
	"encoding/json"
	"text/template"

	cb "github.com/hyperledger/fabric-protos-go/common"

	"github.com/yakumioto/alkaid/third_party/github.com/hyperledger/fabric/common/tools/protolator"
	"github.com/yakumioto/alkaid/third_party/github.com/hyperledger/fabric/protoutil"
)

type GenesisConf struct {
	Organization     *OrganizationConf `json:"organization,omitempty"`
	ChannelName      string            `json:"channel_name,omitempty"`
	ConsortiumName   string            `json:"consortium_name,omitempty"`
	Consensus        *Consensus        `json:"consensus_type,omitempty"`
	OrdererAddresses []string          `json:"orderer_addresses,omitempty"`
}

type OrganizationConf struct {
	Name          string `json:"name,omitempty"`
	MSPIdentifier string `json:"msp_identifier,omitempty"`
	AdminCert     string `json:"admin_cert,omitempty"`
	SignRootCert  string `json:"sign_root_cert,omitempty"`
	TLSRootCert   string `json:"tls_root_cert,omitempty"`
}

type ChannelConf struct {
	Name           string              `json:"name,omitempty"`
	ConsortiumName string              `json:"consortium_name,omitempty"`
	Organizations  map[string]struct{} `json:"organizations,omitempty"`
}

type AnchorPeerConf struct {
	ChannelName      string  `json:"channel_name,omitempty"`
	OrganizationName string  `json:"organization_name,omitempty"`
	Peers            []*Peer `json:"peers,omitempty"`
}

type Consensus struct {
	Type         string       `json:"type,omitempty"`
	Consenters   []*Consenter `json:"consenters,omitempty"`
	KafkaBrokers []string     `json:"kafka_brokers,omitempty"`
}

type Consenter struct {
	Host          string `json:"host,omitempty"`
	Port          uint32 `json:"port,omitempty"`
	ClientTLSCert string `json:"client_tls_cert,omitempty"`
	ServerTLSCert string `json:"server_tls_cert,omitempty"`
}

type Peer struct {
	Host string `json:"host,omitempty"`
	Port uint32 `json:"port,omitempty"`
}

func GetGenesisBlock(conf *GenesisConf) ([]byte, error) {
	tmpl := template.New("genesis").Funcs(template.FuncMap{
		"jsonMarshal": func(v interface{}) string {
			b, err := json.Marshal(v)
			if err != nil {
				return ""
			}
			return string(b)
		},
	})

	tmpl = template.Must(tmpl.Parse(tmplGenesis))
	tmpl = template.Must(tmpl.Parse(tmplOrdererOrganization))
	tmpl = template.Must(tmpl.Parse(tmplConsensusType))
	tmpl = template.Must(tmpl.Parse(tmplOrdererAddresses))

	buf := bytes.NewBuffer(nil)
	if err := tmpl.Execute(buf, conf); err != nil {
		return nil, err
	}

	config := &cb.Config{}
	if err := protolator.DeepUnmarshalJSON(buf, config); err != nil {
		return nil, err
	}

	payloadChannelHeader := protoutil.MakeChannelHeader(cb.HeaderType_CONFIG, int32(1), conf.ChannelName, 0)
	payloadSignatureHeader := protoutil.MakeSignatureHeader(nil, protoutil.CreateNonceOrPanic())
	protoutil.SetTxID(payloadChannelHeader, payloadSignatureHeader)
	payloadHeader := protoutil.MakePayloadHeader(payloadChannelHeader, payloadSignatureHeader)
	payload := &cb.Payload{Header: payloadHeader, Data: protoutil.MarshalOrPanic(&cb.ConfigEnvelope{Config: config})}
	envelope := &cb.Envelope{Payload: protoutil.MarshalOrPanic(payload), Signature: nil}

	block := protoutil.NewBlock(0, nil)
	block.Data = &cb.BlockData{Data: [][]byte{protoutil.MarshalOrPanic(envelope)}}
	block.Header.DataHash = protoutil.BlockDataHash(block.Data)
	block.Metadata.Metadata[cb.BlockMetadataIndex_SIGNATURES] = protoutil.MarshalOrPanic(&cb.Metadata{
		Value: protoutil.MarshalOrPanic(&cb.OrdererBlockMetadata{
			LastConfig: &cb.LastConfig{Index: 0},
		}),
	})

	return protoutil.MarshalOrPanic(block), nil
}

func GetOrganizationJSON(conf *OrganizationConf) ([]byte, error) {
	tmpl := template.New("organization")
	tmpl = template.Must(tmpl.Parse(tmplPeerOrganization))

	buf := bytes.NewBuffer(nil)
	if err := tmpl.Execute(buf, conf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func GetChannelTX(conf *ChannelConf) ([]byte, error) {
	tmpl := template.New("channeltx").Funcs(template.FuncMap{
		"jsonMarshal": func(v interface{}) string {
			b, err := json.Marshal(v)
			if err != nil {
				return ""
			}
			return string(b)
		},
	})

	tmpl = template.Must(tmpl.Parse(tmplChannel))

	buf := bytes.NewBuffer(nil)
	if err := tmpl.Execute(buf, conf); err != nil {
		return nil, err
	}

	configUpdate := &cb.ConfigUpdate{}
	if err := protolator.DeepUnmarshalJSON(buf, configUpdate); err != nil {
		return nil, err
	}

	configUpdateEnv := &cb.ConfigUpdateEnvelope{
		ConfigUpdate: protoutil.MarshalOrPanic(configUpdate),
	}

	channelUpdateTx, err := protoutil.CreateSignedEnvelope(cb.HeaderType_CONFIG_UPDATE, conf.Name, nil, configUpdateEnv, 0, 0)
	if err != nil {
		return nil, err
	}

	return protoutil.MarshalOrPanic(channelUpdateTx), nil
}

func GetAnchorPeerTX(conf *AnchorPeerConf) ([]byte, error) {
	tmpl := template.New("anchorpeertx").Funcs(template.FuncMap{
		"jsonMarshal": func(v interface{}) string {
			b, err := json.Marshal(v)
			if err != nil {
				return ""
			}
			return string(b)
		},
	})
	tmpl = template.Must(tmpl.Parse(tmplAnchorPeer))

	buf := bytes.NewBuffer(nil)
	if err := tmpl.Execute(buf, conf); err != nil {
		return nil, err
	}

	configUpdate := &cb.ConfigUpdate{}
	if err := protolator.DeepUnmarshalJSON(buf, configUpdate); err != nil {
		return nil, err
	}

	configUpdateEnv := &cb.ConfigUpdateEnvelope{
		ConfigUpdate: protoutil.MarshalOrPanic(configUpdate),
	}

	channelUpdateTx, err := protoutil.CreateSignedEnvelope(cb.HeaderType_CONFIG_UPDATE, conf.ChannelName, nil, configUpdateEnv, 0, 0)
	if err != nil {
		return nil, err
	}

	return protoutil.MarshalOrPanic(channelUpdateTx), nil
}
