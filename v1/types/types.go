package types

type ConfigFile struct {
	ConfigFilePath string `yaml:"-"`
	IPAddress string `yaml:"ip_address"`
	AuthToken string `yaml:"auth_token"`
	DeviceID string `yaml:"device_id"`
	DeviceName string `yaml:"device_name"`
}

type APIResponse struct {
	Body []byte
	StatusCode int
	Err error
	Success bool
}

type PairingStageOneRequest struct {
	URL string `json:"_url"`
	DeviceID string `json:"DEVICE_ID"`
	DeviceName string `json:"DEVICE_NAME"`
}

type PairingStageOneResponse struct {
	STATUS struct {
		RESULT string `json:"RESULT"`
		DETAIL string `json:"DETAIL"`
	} `json:"STATUS"`
	ITEM struct {
		CHALLENGE_TYPE int `json:"CHALLENGE_TYPE"`
		PAIRING_REQ_TOKEN int32 `json:"PAIRING_REQ_TOKEN"`
	} `json:"ITEM"`
}

type PairingStageTwoRequest struct {
	URL string `json:"_url"`
	DeviceID string `json:"DEVICE_ID"`
	DeviceName string `json:"DEVICE_NAME"`
	ChallengeType int `json:"CHALLENGE_TYPE"`
	PairingReqToken int32 `json:"PAIRING_REQ_TOKEN"`
	ResponseValue string `json:"RESPONSE_VALUE"`
}

type PairingStageTwoResponse struct {
	STATUS struct {
		RESULT string `json:"RESULT"`
		DETAIL string `json:"DETAIL"`
	} `json:"STATUS"`
	ITEM struct {
		AUTH_TOKEN string `json:"AUTH_TOKEN"`
	} `json:"ITEM"`
}

type PowerGetStateResponse struct {
	//HASHLIST []int `json:"HASHLIST"`
	ITEMS []struct {
		CNAME string `json:"CNAME"`
		ENABLED string `json:"ENABLED"`
		HASHVAL int64 `json:"HASHVAL"`
		NAME string `json:"NAME"`
		TYPE string `json:"TYPE"`
		VALUE int `json:"VALUE"`
	} `json:"ITEMS"`
	// PARAMETERS struct {
	// 	FLAT string `json:"FLAT"`
	// 	HASHONLY string `json:"HASHONLY"`
	// 	HELPTEXT string `json:"HELPTEXT"`
	// } `json:PARAMETERS`
	STATUS struct {
		DETAIL string `json:"DETAIL"`
		RESULT string `json:"RESULT"`
	} `json:STATUS`
	URI string `json:"URI"`
}


type Key struct {
	Codeset int `json:"CODESET"`
	Code int `json:"CODE"`
	Action string `json:"ACTION"`
}

type KeyCommand struct {
	URL string `json:"_url"`
	Keylist []Key `json:"KEYLIST"`
}

type KeyPressResponse struct {
	STATUS struct {
		RESULT string `json:"RESULT"`
		DETAIL string `json:"DETAIL"`
	} `json:"STATUS"`
	URI string `json:"URI"`
}


type VolumeGetResponse struct {
	HASHLIST []int64 `json:"HASHLIST"`
	ITEMS []struct {
		CNAME string `json:"CNAME"`
		ENABLED string `json:"ENABLED"`
		HASHVAL int64 `json:"HASHVAL"`
		NAME string `json:"NAME"`
		TYPE string `json:"TYPE"`
		VALUE int `json:"VALUE"`
	} `json:"ITEMS"`
	PARAMETERS struct {
		FLAT string `json:"FLAT"`
		HASHONLY string `json:"HASHONLY"`
		HELPTEXT string `json:"HELPTEXT"`
	} `json:PARAMETERS`
	STATUS struct {
		DETAIL string `json:"DETAIL"`
		RESULT string `json:"RESULT"`
	} `json:STATUS`
	URI string `json:"URI"`
}

type Input struct {
	Name string `json:"NAME"`
	HashValue int64 `json:"HASHVAL"`
}

type InputGetCurrentResponse struct {
	HASHLIST []int64 `json:"HASHLIST"`
	ITEMS []struct {
		CNAME string `json:"CNAME"`
		ENABLED string `json:"ENABLED"`
		HASHVAL int64 `json:"HASHVAL"`
		NAME string `json:"NAME"`
		TYPE string `json:"TYPE"`
		VALUE string `json:"VALUE"`
	} `json:"ITEMS"`
	PARAMETERS struct {
		FLAT string `json:"FLAT"`
		HASHONLY string `json:"HASHONLY"`
		HELPTEXT string `json:"HELPTEXT"`
	} `json:PARAMETERS`
	STATUS struct {
		DETAIL string `json:"DETAIL"`
		RESULT string `json:"RESULT"`
	} `json:STATUS`
	URI string `json:"URI"`
}

type InputGetAvailableResponse struct {
	HASHLIST []int64 `json:"HASHLIST"`
	GROUP string `json:"GROUP"`
	NAME string `json:"NAME"`
	ITEMS []struct {
		CNAME string `json:"CNAME"`
		ENABLED string `json:"ENABLED"`
		HASHVAL int64 `json:"HASHVAL"`
		NAME string `json:"NAME"`
		TYPE string `json:"TYPE"`
		VALUE struct {
			NAME string `json:"NAME"`
			METADATA string `json:"METADATA"`
		} `json:VALUE`
	} `json:"ITEMS"`
	PARAMETERS struct {
		FLAT string `json:"FLAT"`
		HASHONLY string `json:"HASHONLY"`
		HELPTEXT string `json:"HELPTEXT"`
	} `json:PARAMETERS`
	STATUS struct {
		DETAIL string `json:"DETAIL"`
		RESULT string `json:"RESULT"`
	} `json:STATUS`
	URI string `json:"URI"`
}

type InputSetRequest struct {
	URL string `json:"_url"`
	ItemName string `json:"item_name"`
	VALUE string `json:"VALUE"`
	HASHVAL int64 `json:"HASHVAL"`
	REQUEST string `json:"REQUEST"`
}

type InputSetRequestString struct {
	URL string `json:"_url"`
	ItemName string `json:"item_name"`
	VALUE string `json:"VALUE"`
	HASHVAL int64 `json:"HASHVAL"`
	REQUEST string `json:"REQUEST"`
}

type InputSetRequestInt struct {
	URL string `json:"_url"`
	ItemName string `json:"item_name"`
	VALUE int `json:"VALUE"`
	HASHVAL int64 `json:"HASHVAL"`
	REQUEST string `json:"REQUEST"`
}

type AudioSettingItem struct {
	INDEX int
	HASHVAL int64
	NAME string
	VALUE interface{}
	CNAME string
	TYPE string
}
type AudioGetSettingResponse struct {
	STATUS struct {
		RESULT string
		DETAIL string
	}
	HASHLIST []int64
	GROUP string
	NAME string
	PARAMETERS struct {
		FLAT string
		HELPTEXT string
		HASHONLY string
	}
	ITEMS []AudioSettingItem
	URI string
	CNAME string
	TYPE string
}

type AudioSettingsOption struct {
	TYPE string
	CNAME string
	GROUP string
	NAME string
	ELEMENTS []string
	INCREMENT int
	INCMARKER string
	DECMARKER string
	MAXIMUM int
	CENTER int
}
type AudioGetAllSettingsResponse struct {
	STATUS struct {
		RESULT string
		DETAIL string
	}
	HASHVAL int64
	GROUP string
	NAME string
	PARAMETERS struct {
		FLAT string
		HELPTEXT string
		HASHONLY string
	}
	ITEMS []AudioSettingsOption
	URI string
	CNAME string
	TYPE string
}


type AudioGetSettingsOptionResponse struct {
	STATUS struct {
		RESULT string
		DETAIL string
	}
	PARAMETERS struct {
		FLAT string
		HELPTEXT string
		HASHONLY string
	}
	ITEMS []AudioSettingsOption
	HASHVAL int64
	URI string
}

type AudioSettingResultItem struct {
	HASHVAL int64
	NAME string
}

type AudioSetSettingResponse struct {
	STATUS struct {
		RESULT string
		DETAIL string
	}
	PARAMETERS struct {
		HASHVAL int64
		REQUEST string
		VALUE string
	}
	ITEMS []AudioSettingResultItem
	HASHLIST []int64
	URI string
}

type AudioSetSettingResponseString struct {
	STATUS struct {
		RESULT string
		DETAIL string
	}
	PARAMETERS struct {
		HASHVAL int64
		REQUEST string
		VALUE string
	}
	ITEMS []AudioSettingResultItem
	HASHLIST []int64
	URI string
}

type AudioSetSettingResponseInt struct {
	STATUS struct {
		RESULT string
		DETAIL string
	}
	PARAMETERS struct {
		HASHVAL int64
		REQUEST string
		VALUE int
	}
	ITEMS []AudioSettingResultItem
	HASHLIST []int64
	URI string
}


type GetSettingsTypeResultItem struct {
	HASHVAL int64
	CNAME string
	TYPE string
	NAME string
}
type SettingsGetTypesResponse struct {
	STATUS struct {
		RESULT string
		DETAIL string
	}
	HASHLIST []int64
	NAME string
	PARAMETERS struct {
		FLAT string
		HELPTEXT string
		HASHONLY string
	}
	ITEMS []GetSettingsTypeResultItem
	URI string
	CNAME string
	TYPE string
}


type SettingsGetTypeResponseItem struct {
	INDEX int
	HASHVAL int64
	NAME string
	ENABLED string
	READONLY string
	CNAME string
	HIDDEN string
	TYPE string
	VALUE []interface{}
}
type SettingsGetTypeResponse struct {
	STATUS struct {
		RESULT string
		DETAIL string
	}
	HASHLIST []int64
	GROUP string
	NAME string
	PARAMETERS struct {
		FLAT string
		HELPTEXT string
		HASHONLY string
	}
	ITEMS []SettingsGetTypeResponseItem
	URI string
	CNAME string
	TYPE string
}

type SettingsGetItem struct {
	HASHVAL int64
	CNAME string
	TYPE string
	NAME string
	VALUE int
}
type SettingsGetResponse struct {
	STATUS struct {
		RESULT string
		DETAIL string
	}
	PARAMETERS struct {
		FLAT string
		HELPTEXT string
		HASHONLY string
	}
	ITEMS []SettingsGetItem
	HASHLIST []int64
	URI string
}

type SettingsSetIntRequest struct {
	URL string `json:"_url"`
	ItemName string `json:"item_name"`
	VALUE int
	HASHVAL int64
	REQUEST string
}
type SettingsSetStringRequest struct {
	URL string `json:"_url"`
	ItemName string `json:"item_name"`
	VALUE string
	HASHVAL int64
	REQUEST string
}


type SettingsSetResponseItem struct {
	HASHVAL int64
	NAME string
}
type SettingsSetResponse struct {
	STATUS struct {
		RESULT string
		DETAIL string
	}
	PARAMETERS struct {
		HASHVAL int64
		REQUEST string
		VALUE interface{}
	}
	ITEMS []SettingsSetResponseItem
	HASHLIST []int64
	URI string
}

type AppGetCurrentResponse struct {
	STATUS struct {
		RESULT string
		DETAIL string
	}
	ITEM struct {
		TYPE string
		VALUE struct {
			MESSAGE string
			NAME_SPACE int
			APP_ID string
		}
	}
	URI string
}

type AppLaunchRequestValue struct {
	APP_ID string `json:"APP_ID"`
	NAME_SPACE int `json:"NAME_SPACE"`
	MESSAGE string `json:"MESSAGE"`
}
type AppLaunchRequest struct {
	URL string `json:"_url"`
	VALUE AppLaunchRequestValue `json:"VALUE"`
}

type AppLaunchResponse struct {
	STATUS struct {
		RESULT string
		DETAIL string
	}
}
