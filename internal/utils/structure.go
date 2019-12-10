package utils

type SystemBlock struct {
	Data struct {
		Data []struct {
			Payload struct {
				Data struct {
					Config     *SystemConfig          `json:"config"`
					LastUpdate map[string]interface{} `json:"last_update"`
				} `json:"data"`
				Header map[string]interface{} `json:"header"`
			} `json:"payload"`
			Signature string `json:"signature"`
		} `json:"data"`
	} `json:"data"`
	Header   map[string]interface{} `json:"header"`
	Metadata map[string]interface{} `json:"metadata"`
}

type Block struct {
	Data struct {
		Data []struct {
			Payload struct {
				Data struct {
					Config     *Config                `json:"config"`
					LastUpdate map[string]interface{} `json:"last_update"`
				} `json:"data"`
				Header map[string]interface{} `json:"header"`
			} `json:"payload"`
			Signature string `json:"signature"`
		} `json:"data"`
	} `json:"data"`
	Header   map[string]interface{} `json:"header"`
	Metadata map[string]interface{} `json:"metadata"`
}

type Config struct {
	ChannelGroup struct {
		Groups struct {
			Application struct {
				Groups    map[string]interface{} `json:"groups"`
				ModPolicy string                 `json:"mod_policy"`
				Policies  map[string]interface{} `json:"policies"`
				Values    map[string]interface{} `json:"values"`
				Version   string                 `json:"version"`
			} `json:"Application"`
			Orderer struct {
				Groups    map[string]interface{} `json:"groups"`
				ModPolicy string                 `json:"mod_policy"`
				Policies  map[string]interface{} `json:"policies"`
				Values    map[string]interface{} `json:"values"`
				Version   string                 `json:"version"`
			} `json:"Orderer"`
		} `json:"groups"`
		ModPolicy string                 `json:"mod_policy"`
		Policies  map[string]interface{} `json:"policies"`
		Values    map[string]interface{} `json:"values"`
		Version   string                 `json:"version"`
	} `json:"channel_group"`
	Sequence string `json:"sequence"`
}

type SystemConfig struct {
	ChannelGroup struct {
		Groups struct {
			Consortiums struct {
				Groups struct {
					SampleConsortium struct {
						Groups    map[string]interface{} `json:"groups"`
						ModPolicy string                 `json:"mod_policy"`
						Policies  map[string]interface{} `json:"policies"`
						Values    map[string]interface{} `json:"values"`
						Version   string                 `json:"version"`
					} `json:"SampleConsortium"`
				} `json:"groups"`
			} `json:"Consortiums"`
			Orderer struct {
				Groups    map[string]interface{} `json:"groups"`
				ModPolicy string                 `json:"mod_policy"`
				Policies  map[string]interface{} `json:"policies"`
				Values    map[string]interface{} `json:"values"`
				Version   string                 `json:"version"`
			} `json:"Orderer"`
		} `json:"groups"`
		ModPolicy string                 `json:"mod_policy"`
		Policies  map[string]interface{} `json:"policies"`
		Values    map[string]interface{} `json:"values"`
		Version   string                 `json:"version"`
	} `json:"channel_group"`
	Sequence string `json:"sequence"`
}
