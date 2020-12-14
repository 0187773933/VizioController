package viziocontroller

import (
	"fmt"
	"time"
	"context"
	"bytes"
	"reflect"
	"encoding/json"
	"net/http"
	"crypto/tls"
	"io/ioutil"
)



// Auth Functions

func pairing_stage_one( ip_address string ) ( result int32 ) {
	result = 0
	put_data , _ := json.Marshal(map[string]string{
		"_url": "/pairing/start" ,
		"DEVICE_ID": "pyvizio" ,
		"DEVICE_NAME": "Python Vizio" ,
	})
	url := fmt.Sprintf( "https://%s:7345/pairing/start" , ip_address )
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true ,
			} ,
		} ,
	}
	request , request_error := http.NewRequest( "PUT" , url , bytes.NewBuffer( put_data ) )
	request.Header.Set( "Content-Type" , "application/json" )
	if request_error != nil { fmt.Println( request_error ); return result }
	response , response_error := client.Do( request )
	if response_error != nil { fmt.Println( response_error ); return result }
	defer response.Body.Close()
	body , body_error := ioutil.ReadAll( response.Body )
	if body_error != nil { fmt.Println( body_error ); return result }
	var result_json struct {
		STATUS struct {
			RESULT string `json:"RESULT"`
			DETAIL string `json:"DETAIL"`
		} `json:"STATUS"`
		ITEM struct {
			CHALLENGE_TYPE int `json:"CHALLENGE_TYPE"`
			PAIRING_REQ_TOKEN int32 `json:"PAIRING_REQ_TOKEN"`
		} `json:"ITEM"`
	}
	fmt.Println( string( body[:] ) )
	json_decode_error := json.Unmarshal( body , &result_json )
	if json_decode_error != nil { fmt.Println( json_decode_error ); return result }
	fmt.Println( result_json )
	fmt.Println( result_json.STATUS )
	result = result_json.ITEM.PAIRING_REQ_TOKEN
	return result
}

type PairingData struct {
	_url string
	DEVICE_ID string
	DEVICE_NAME string
	CHALLENGE_TYPE int
	PAIRING_REQ_TOKEN int32
	RESPONSE_VALUE string
}
func paring_stage_two( ip_address string , pairing_request_token int32 , code_displayed_on_tv string ) ( result string ) {
	fmt.Println("Starting Stage 2")
	result = "failed"
	put_data , _ := json.Marshal(PairingData{
		_url: "/pairing/pair" ,
		DEVICE_ID: "pyvizio" ,
		DEVICE_NAME: "Python Vizio" ,
		CHALLENGE_TYPE: 1 ,
		PAIRING_REQ_TOKEN: pairing_request_token ,
		RESPONSE_VALUE: code_displayed_on_tv ,
	})
	url := fmt.Sprintf( "https://%s:7345/pairing/pair" , ip_address )
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true ,
			} ,
		} ,
	}
	request , request_error := http.NewRequest( "PUT" , url , bytes.NewBuffer( put_data ) )
	request.Header.Set( "Content-Type" , "application/json" )
	if request_error != nil { fmt.Println( request_error ); return result }
	response , response_error := client.Do( request )
	if response_error != nil { fmt.Println( response_error ); return result }
	defer response.Body.Close()
	body , body_error := ioutil.ReadAll( response.Body )
	if body_error != nil { fmt.Println( body_error ); return result }
	var result_json struct {
		STATUS struct {
			RESULT string `json:"RESULT"`
			DETAIL string `json:"DETAIL"`
		} `json:"STATUS"`
		ITEM struct {
			AUTH_TOKEN string `json:"AUTH_TOKEN"`
		} `json:"ITEM"`
	}
	fmt.Println( string( body[:] ) )
	json_decode_error := json.Unmarshal( body , &result_json )
	if json_decode_error != nil { fmt.Println( json_decode_error ); return result }
	fmt.Println( result_json )
	fmt.Println( result_json.ITEM.AUTH_TOKEN )
	result = result_json.ITEM.AUTH_TOKEN
	return result
}

func RegenerateAuthToken() {
	var ctx = context.Background()
	redis_connection := get_redis_connection( "localhost:6379", 3 , "" )
	ip_address , err := redis_connection.Get( ctx , "STATE.VIZIO_TV.IP_ADDRESS" ).Result()
	if err != nil { panic( err ) }
	fmt.Println( "STATE.VIZIO_TV.IP_ADDRESS" , ip_address )
	pairing_request_token := pairing_stage_one( ip_address )
	fmt.Println( "Enter Code Displayed on TV")
	var code_displayed_on_tv string
	fmt.Scanln( &code_displayed_on_tv )
	auth_token := paring_stage_two( ip_address , pairing_request_token , code_displayed_on_tv )
	err = redis_connection.Set( ctx , "STATE.VIZIO_TV.AUTH_TOKEN", auth_token , 0 ).Err()
	if err != nil { fmt.Println( err ) }
}

// Control Functions

type power_state_result struct {
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
func GetPowerState( ip_address string , auth_token string ) ( result int ) {
	result = -1
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true ,
			} ,
		} ,
	}
	url := fmt.Sprintf( "https://%s:7345/state/device/power_mode" , ip_address )
	request , request_error := http.NewRequest( "GET" , url , nil )
	request.Header.Set( "Content-Type" , "application/json" )
	request.Header.Set( "AUTH" , auth_token )
	if request_error != nil { fmt.Println( request_error ); return result }
	response , response_error := client.Do( request )
	if response_error != nil { fmt.Println( response_error ); return result }
	defer response.Body.Close()
	body , body_error := ioutil.ReadAll( response.Body )
	if body_error != nil { fmt.Println( body_error ); return result }
	//fmt.Println( string( body[:] ) )
	var result_json power_state_result
	json_decode_error := json.Unmarshal( body , &result_json )
	if json_decode_error != nil { fmt.Println( json_decode_error ); return result }
	result = result_json.ITEMS[0].VALUE
	return result
}


type key struct {
	CODESET int
	CODE int
	ACTION string
}
type key_command struct {
	URL string `json:"_url"`
	KEYLIST []key `json:KEYLIST"`
}
func PowerOn( ip_address string , auth_token string ) ( result string ) {
	result = "failed"
	power_on_key := []key{
		key{
			CODESET: 11 ,
			CODE: 1 ,
			ACTION: "KEYPRESS" ,
		} ,
	}
	put_data , _ := json.Marshal( key_command{
		URL: "/key_command/" ,
		KEYLIST: power_on_key ,
	})
	fmt.Println( string( put_data[:] ) )
	url := fmt.Sprintf( "https://%s:7345/key_command/" , ip_address )
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true ,
			} ,
		} ,
	}
	request , request_error := http.NewRequest( "PUT" , url , bytes.NewBuffer( put_data ) )
	request.Header.Set( "Content-Type" , "application/json" )
	request.Header.Set( "AUTH" , auth_token )
	if request_error != nil { fmt.Println( request_error ); return result }
	response , response_error := client.Do( request )
	if response_error != nil { fmt.Println( response_error ); return result }
	defer response.Body.Close()
	body , body_error := ioutil.ReadAll( response.Body )
	if body_error != nil { fmt.Println( body_error ); return result }
	var result_json struct {
		STATUS struct {
			RESULT string `json:"RESULT"`
			DETAIL string `json:"DETAIL"`
		} `json:"STATUS"`
		URI string `json:"URI"`
	}
	fmt.Println( string( body[:] ) )
	json_decode_error := json.Unmarshal( body , &result_json )
	if json_decode_error != nil { fmt.Println( json_decode_error ); return result }
	result = result_json.STATUS.RESULT
	return result
}

func PowerOff( ip_address string , auth_token string ) ( result string ) {
	result = "failed"
	power_on_key := []key{
		key{
			CODESET: 11 ,
			CODE: 0 ,
			ACTION: "KEYPRESS" ,
		} ,
	}
	put_data , _ := json.Marshal( key_command{
		URL: "/key_command/" ,
		KEYLIST: power_on_key ,
	})
	fmt.Println( string( put_data[:] ) )
	url := fmt.Sprintf( "https://%s:7345/key_command/" , ip_address )
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true ,
			} ,
		} ,
	}
	request , request_error := http.NewRequest( "PUT" , url , bytes.NewBuffer( put_data ) )
	request.Header.Set( "Content-Type" , "application/json" )
	request.Header.Set( "AUTH" , auth_token )
	if request_error != nil { fmt.Println( request_error ); return result }
	response , response_error := client.Do( request )
	if response_error != nil { fmt.Println( response_error ); return result }
	defer response.Body.Close()
	body , body_error := ioutil.ReadAll( response.Body )
	if body_error != nil { fmt.Println( body_error ); return result }
	var result_json struct {
		STATUS struct {
			RESULT string `json:"RESULT"`
			DETAIL string `json:"DETAIL"`
		} `json:"STATUS"`
		URI string `json:"URI"`
	}
	fmt.Println( string( body[:] ) )
	json_decode_error := json.Unmarshal( body , &result_json )
	if json_decode_error != nil { fmt.Println( json_decode_error ); return result }
	result = result_json.STATUS.RESULT
	return result
}

type VolumeSetting struct {
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
func GetVolume( ip_address string , auth_token string ) ( result int ) {
	result = -1
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true ,
			} ,
		} ,
	}
	url := fmt.Sprintf( "https://%s:7345/menu_native/dynamic/tv_settings/audio/volume" , ip_address )
	request , request_error := http.NewRequest( "GET" , url , nil )
	request.Header.Set( "Content-Type" , "application/json" )
	request.Header.Set( "AUTH" , auth_token )
	if request_error != nil { fmt.Println( request_error ); return result }
	response , response_error := client.Do( request )
	if response_error != nil { fmt.Println( response_error ); return result }
	defer response.Body.Close()
	body , body_error := ioutil.ReadAll( response.Body )
	if body_error != nil { fmt.Println( body_error ); return result }
	fmt.Println( string( body[:] ) )
	var result_json VolumeSetting
	json_decode_error := json.Unmarshal( body , &result_json )
	if json_decode_error != nil { fmt.Println( json_decode_error ); return result }
	result = result_json.ITEMS[0].VALUE
	return result
}

func VolumeDown( ip_address string , auth_token string ) ( result string ) {
	result = "failed"
	power_on_key := []key{
		key{
			CODESET: 5 ,
			CODE: 0 ,
			ACTION: "KEYPRESS" ,
		} ,
	}
	put_data , _ := json.Marshal( key_command{
		URL: "/key_command/" ,
		KEYLIST: power_on_key ,
	})
	fmt.Println( string( put_data[:] ) )
	url := fmt.Sprintf( "https://%s:7345/key_command/" , ip_address )
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true ,
			} ,
		} ,
	}
	request , request_error := http.NewRequest( "PUT" , url , bytes.NewBuffer( put_data ) )
	request.Header.Set( "Content-Type" , "application/json" )
	request.Header.Set( "AUTH" , auth_token )
	if request_error != nil { fmt.Println( request_error ); return result }
	response , response_error := client.Do( request )
	if response_error != nil { fmt.Println( response_error ); return result }
	defer response.Body.Close()
	body , body_error := ioutil.ReadAll( response.Body )
	if body_error != nil { fmt.Println( body_error ); return result }
	var result_json struct {
		STATUS struct {
			RESULT string `json:"RESULT"`
			DETAIL string `json:"DETAIL"`
		} `json:"STATUS"`
		URI string `json:"URI"`
	}
	fmt.Println( string( body[:] ) )
	json_decode_error := json.Unmarshal( body , &result_json )
	if json_decode_error != nil { fmt.Println( json_decode_error ); return result }
	result = result_json.STATUS.RESULT
	return result
}

func VolumeUp( ip_address string , auth_token string ) ( result string ) {
	result = "failed"
	power_on_key := []key{
		key{
			CODESET: 5 ,
			CODE: 1 ,
			ACTION: "KEYPRESS" ,
		} ,
	}
	put_data , _ := json.Marshal( key_command{
		URL: "/key_command/" ,
		KEYLIST: power_on_key ,
	})
	fmt.Println( string( put_data[:] ) )
	url := fmt.Sprintf( "https://%s:7345/key_command/" , ip_address )
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true ,
			} ,
		} ,
	}
	request , request_error := http.NewRequest( "PUT" , url , bytes.NewBuffer( put_data ) )
	request.Header.Set( "Content-Type" , "application/json" )
	request.Header.Set( "AUTH" , auth_token )
	if request_error != nil { fmt.Println( request_error ); return result }
	response , response_error := client.Do( request )
	if response_error != nil { fmt.Println( response_error ); return result }
	defer response.Body.Close()
	body , body_error := ioutil.ReadAll( response.Body )
	if body_error != nil { fmt.Println( body_error ); return result }
	var result_json struct {
		STATUS struct {
			RESULT string `json:"RESULT"`
			DETAIL string `json:"DETAIL"`
		} `json:"STATUS"`
		URI string `json:"URI"`
	}
	fmt.Println( string( body[:] ) )
	json_decode_error := json.Unmarshal( body , &result_json )
	if json_decode_error != nil { fmt.Println( json_decode_error ); return result }
	result = result_json.STATUS.RESULT
	return result
}

func MuteOff( ip_address string , auth_token string ) ( result string ) {
	result = "failed"
	power_on_key := []key{
		key{
			CODESET: 5 ,
			CODE: 2 ,
			ACTION: "KEYPRESS" ,
		} ,
	}
	put_data , _ := json.Marshal( key_command{
		URL: "/key_command/" ,
		KEYLIST: power_on_key ,
	})
	fmt.Println( string( put_data[:] ) )
	url := fmt.Sprintf( "https://%s:7345/key_command/" , ip_address )
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true ,
			} ,
		} ,
	}
	request , request_error := http.NewRequest( "PUT" , url , bytes.NewBuffer( put_data ) )
	request.Header.Set( "Content-Type" , "application/json" )
	request.Header.Set( "AUTH" , auth_token )
	if request_error != nil { fmt.Println( request_error ); return result }
	response , response_error := client.Do( request )
	if response_error != nil { fmt.Println( response_error ); return result }
	defer response.Body.Close()
	body , body_error := ioutil.ReadAll( response.Body )
	if body_error != nil { fmt.Println( body_error ); return result }
	var result_json struct {
		STATUS struct {
			RESULT string `json:"RESULT"`
			DETAIL string `json:"DETAIL"`
		} `json:"STATUS"`
		URI string `json:"URI"`
	}
	fmt.Println( string( body[:] ) )
	json_decode_error := json.Unmarshal( body , &result_json )
	if json_decode_error != nil { fmt.Println( json_decode_error ); return result }
	result = result_json.STATUS.RESULT
	return result
}

func MuteOn( ip_address string , auth_token string ) ( result string ) {
	result = "failed"
	power_on_key := []key{
		key{
			CODESET: 5 ,
			CODE: 3 ,
			ACTION: "KEYPRESS" ,
		} ,
	}
	put_data , _ := json.Marshal( key_command{
		URL: "/key_command/" ,
		KEYLIST: power_on_key ,
	})
	fmt.Println( string( put_data[:] ) )
	url := fmt.Sprintf( "https://%s:7345/key_command/" , ip_address )
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true ,
			} ,
		} ,
	}
	request , request_error := http.NewRequest( "PUT" , url , bytes.NewBuffer( put_data ) )
	request.Header.Set( "Content-Type" , "application/json" )
	request.Header.Set( "AUTH" , auth_token )
	if request_error != nil { fmt.Println( request_error ); return result }
	response , response_error := client.Do( request )
	if response_error != nil { fmt.Println( response_error ); return result }
	defer response.Body.Close()
	body , body_error := ioutil.ReadAll( response.Body )
	if body_error != nil { fmt.Println( body_error ); return result }
	var result_json struct {
		STATUS struct {
			RESULT string `json:"RESULT"`
			DETAIL string `json:"DETAIL"`
		} `json:"STATUS"`
		URI string `json:"URI"`
	}
	fmt.Println( string( body[:] ) )
	json_decode_error := json.Unmarshal( body , &result_json )
	if json_decode_error != nil { fmt.Println( json_decode_error ); return result }
	result = result_json.STATUS.RESULT
	return result
}

func MuteToggle( ip_address string , auth_token string ) ( result string ) {
	result = "failed"
	power_on_key := []key{
		key{
			CODESET: 5 ,
			CODE: 4 ,
			ACTION: "KEYPRESS" ,
		} ,
	}
	put_data , _ := json.Marshal( key_command{
		URL: "/key_command/" ,
		KEYLIST: power_on_key ,
	})
	fmt.Println( string( put_data[:] ) )
	url := fmt.Sprintf( "https://%s:7345/key_command/" , ip_address )
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true ,
			} ,
		} ,
	}
	request , request_error := http.NewRequest( "PUT" , url , bytes.NewBuffer( put_data ) )
	request.Header.Set( "Content-Type" , "application/json" )
	request.Header.Set( "AUTH" , auth_token )
	if request_error != nil { fmt.Println( request_error ); return result }
	response , response_error := client.Do( request )
	if response_error != nil { fmt.Println( response_error ); return result }
	defer response.Body.Close()
	body , body_error := ioutil.ReadAll( response.Body )
	if body_error != nil { fmt.Println( body_error ); return result }
	var result_json struct {
		STATUS struct {
			RESULT string `json:"RESULT"`
			DETAIL string `json:"DETAIL"`
		} `json:"STATUS"`
		URI string `json:"URI"`
	}
	fmt.Println( string( body[:] ) )
	json_decode_error := json.Unmarshal( body , &result_json )
	if json_decode_error != nil { fmt.Println( json_decode_error ); return result }
	result = result_json.STATUS.RESULT
	return result
}

type InputType struct {
	name string
	hash_value int64
}
type CurrentInputSetting struct {
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
func GetCurrentInput( ip_address string , auth_token string ) ( result InputType ) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true ,
			} ,
		} ,
	}
	url := fmt.Sprintf( "https://%s:7345/menu_native/dynamic/tv_settings/devices/current_input" , ip_address )
	request , request_error := http.NewRequest( "GET" , url , nil )
	request.Header.Set( "Content-Type" , "application/json" )
	request.Header.Set( "AUTH" , auth_token )
	if request_error != nil { fmt.Println( request_error ); return result }
	response , response_error := client.Do( request )
	if response_error != nil { fmt.Println( response_error ); return result }
	defer response.Body.Close()
	body , body_error := ioutil.ReadAll( response.Body )
	if body_error != nil { fmt.Println( body_error ); return result }
	fmt.Println( string( body[:] ) )
	var result_json CurrentInputSetting
	json_decode_error := json.Unmarshal( body , &result_json )
	if json_decode_error != nil { fmt.Println( json_decode_error ); return result }
	result.name = result_json.ITEMS[0].VALUE
	result.hash_value = result_json.ITEMS[0].HASHVAL
	return result
}

type InputList []InputType
type AvailableInputSettings struct {
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
func GetAvailableInputs( ip_address string , auth_token string ) ( result InputList ) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true ,
			} ,
		} ,
	}
	url := fmt.Sprintf( "https://%s:7345/menu_native/dynamic/tv_settings/devices/name_input" , ip_address )
	request , request_error := http.NewRequest( "GET" , url , nil )
	request.Header.Set( "Content-Type" , "application/json" )
	request.Header.Set( "AUTH" , auth_token )
	if request_error != nil { fmt.Println( request_error ); return result }
	response , response_error := client.Do( request )
	if response_error != nil { fmt.Println( response_error ); return result }
	defer response.Body.Close()
	body , body_error := ioutil.ReadAll( response.Body )
	if body_error != nil { fmt.Println( body_error ); return result }
	fmt.Println( string( body[:] ) )
	var result_json AvailableInputSettings
	json_decode_error := json.Unmarshal( body , &result_json )
	if json_decode_error != nil { fmt.Println( json_decode_error ); return result }
	for _ , item := range result_json.ITEMS {
		result = append( result , InputType{
			name: item.VALUE.NAME ,
			hash_value: item.HASHVAL ,
		})
	}
	return result
}

type InputSet struct {
	_url string
	item_name string
	VALUE string
	HASHVAL int64
	REQUEST string
}
func SetInput( ip_address string , auth_token string , input_name string ) ( result string ) {
	current_input := GetCurrentInput( ip_address , auth_token )
	result = "failed"
	put_data , _ := json.Marshal(InputSet{
		_url: "/menu_native/dynamic/tv_settings/devices/current_input" ,
		item_name: "CURRENT_INPUT" ,
		VALUE: input_name ,
		HASHVAL: current_input.hash_value ,
		REQUEST: "MODIFY" ,
	})
	url := fmt.Sprintf( "https://%s:7345/menu_native/dynamic/tv_settings/devices/current_input" , ip_address )
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true ,
			} ,
		} ,
	}
	request , request_error := http.NewRequest( "PUT" , url , bytes.NewBuffer( put_data ) )
	request.Header.Set( "Content-Type" , "application/json" )
	request.Header.Set( "AUTH" , auth_token )
	if request_error != nil { fmt.Println( request_error ); return result }
	response , response_error := client.Do( request )
	if response_error != nil { fmt.Println( response_error ); return result }
	defer response.Body.Close()
	body , body_error := ioutil.ReadAll( response.Body )
	if body_error != nil { fmt.Println( body_error ); return result }
	var result_json struct {
		STATUS struct {
			RESULT string `json:"RESULT"`
			DETAIL string `json:"DETAIL"`
		} `json:"STATUS"`
		ITEM struct {
			AUTH_TOKEN string `json:"AUTH_TOKEN"`
		} `json:"ITEM"`
	}
	fmt.Println( string( body[:] ) )
	json_decode_error := json.Unmarshal( body , &result_json )
	if json_decode_error != nil { fmt.Println( json_decode_error ); return result }
	fmt.Println( result_json )
	//result = result_json.ITEM.AUTH_TOKEN
	return result
}

// You have to Call This Twice, Because the First Time it just Brings Up the Input Menu
func CycleInput( ip_address string , auth_token string ) ( result string ) {
	result = "failed"
	key_press := []key{
		key{
			CODESET: 7 ,
			CODE: 1 ,
			ACTION: "KEYPRESS" ,
		} ,
	}
	put_data , _ := json.Marshal( key_command{
		URL: "/key_command/" ,
		KEYLIST: key_press ,
	})
	fmt.Println( string( put_data[:] ) )
	url := fmt.Sprintf( "https://%s:7345/key_command/" , ip_address )
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true ,
			} ,
		} ,
	}
	request , request_error := http.NewRequest( "PUT" , url , bytes.NewBuffer( put_data ) )
	request.Header.Set( "Content-Type" , "application/json" )
	request.Header.Set( "AUTH" , auth_token )
	if request_error != nil { fmt.Println( request_error ); return result }
	response , response_error := client.Do( request )
	if response_error != nil { fmt.Println( response_error ); return result }
	defer response.Body.Close()
	body , body_error := ioutil.ReadAll( response.Body )
	if body_error != nil { fmt.Println( body_error ); return result }
	var result_json struct {
		STATUS struct {
			RESULT string `json:"RESULT"`
			DETAIL string `json:"DETAIL"`
		} `json:"STATUS"`
		URI string `json:"URI"`
	}
	fmt.Println( string( body[:] ) )
	json_decode_error := json.Unmarshal( body , &result_json )
	if json_decode_error != nil { fmt.Println( json_decode_error ); return result }
	result = result_json.STATUS.RESULT
	return result
}

type AudioSettingItem struct {
	INDEX int
	HASHVAL int64
	NAME string
	VALUE interface{}
	CNAME string
	TYPE string
}
type AudioSetting struct {
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
//type CurrentAudioSettings interface {}
func GetAudioSetting( ip_address string , auth_token string , audio_setting_name string) ( result AudioSetting ) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true ,
			} ,
		} ,
	}
	url := fmt.Sprintf( "https://%s:7345/menu_native/dynamic/tv_settings/audio/%s" , ip_address , audio_setting_name )
	request , request_error := http.NewRequest( "GET" , url , nil )
	request.Header.Set( "Content-Type" , "application/json" )
	request.Header.Set( "AUTH" , auth_token )
	if request_error != nil { fmt.Println( request_error ); return result }
	response , response_error := client.Do( request )
	if response_error != nil { fmt.Println( response_error ); return result }
	defer response.Body.Close()
	body , body_error := ioutil.ReadAll( response.Body )
	if body_error != nil { fmt.Println( body_error ); return result }
	fmt.Println( string( body[:] ) )
	var result_json AudioSetting
	json_decode_error := json.Unmarshal( body , &result_json )
	if json_decode_error != nil { fmt.Println( json_decode_error ); return result }
	fmt.Println( result_json )
	result = result_json
	return result
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
type AudioSettingsOptions struct {
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
func GetAllAudioSettingsOptions( ip_address string , auth_token string ) ( result AudioSettingsOptions ) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true ,
			} ,
		} ,
	}
	url := fmt.Sprintf( "https://%s:7345/menu_native/static/tv_settings/audio" , ip_address )
	request , request_error := http.NewRequest( "GET" , url , nil )
	request.Header.Set( "Content-Type" , "application/json" )
	request.Header.Set( "AUTH" , auth_token )
	if request_error != nil { fmt.Println( request_error ); return result }
	response , response_error := client.Do( request )
	if response_error != nil { fmt.Println( response_error ); return result }
	defer response.Body.Close()
	body , body_error := ioutil.ReadAll( response.Body )
	if body_error != nil { fmt.Println( body_error ); return result }
	fmt.Println( string( body[:] ) )
	var result_json AudioSettingsOptions
	json_decode_error := json.Unmarshal( body , &result_json )
	if json_decode_error != nil { fmt.Println( json_decode_error ); return result }
	fmt.Println( result_json )
	result = result_json
	return result
}

type SingleAudioSettingOption struct {
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
func GetAudioSettingsOption( ip_address string , auth_token string , audio_setting_name string ) ( result SingleAudioSettingOption ) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true ,
			} ,
		} ,
	}
	url := fmt.Sprintf( "https://%s:7345/menu_native/static/tv_settings/audio/%s" , ip_address, audio_setting_name )
	request , request_error := http.NewRequest( "GET" , url , nil )
	request.Header.Set( "Content-Type" , "application/json" )
	request.Header.Set( "AUTH" , auth_token )
	if request_error != nil { fmt.Println( request_error ); return result }
	response , response_error := client.Do( request )
	if response_error != nil { fmt.Println( response_error ); return result }
	defer response.Body.Close()
	body , body_error := ioutil.ReadAll( response.Body )
	if body_error != nil { fmt.Println( body_error ); return result }
	fmt.Println( string( body[:] ) )
	var result_json SingleAudioSettingOption
	json_decode_error := json.Unmarshal( body , &result_json )
	if json_decode_error != nil { fmt.Println( json_decode_error ); return result }
	//fmt.Println( result_json )
	result = result_json
	return result
}

type AudioSettingResultItem struct {
	HASHVAL int64
	NAME string
}
type SetAudioSettingResult struct {
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
func SetAudioSetting( ip_address string , auth_token string , audio_setting_name string , audio_setting_option string ) ( result SetAudioSettingResult ) {
	// For this one, you first need the hash value of the audio setting we are modifying
	audio_setting := GetAudioSetting( ip_address , auth_token , audio_setting_name )
	put_data , _ := json.Marshal(InputSet{
		_url: fmt.Sprintf( "/menu_native/dynamic/tv_settings/audio/%s" , audio_setting_name ) ,
		item_name: "SETTINGS" ,
		VALUE: audio_setting_option ,
		HASHVAL: audio_setting.ITEMS[0].HASHVAL ,
		REQUEST: "MODIFY" ,
	})
	url := fmt.Sprintf( "https://%s:7345/menu_native/dynamic/tv_settings/audio/%s" , ip_address , audio_setting_name )
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true ,
			} ,
		} ,
	}
	request , request_error := http.NewRequest( "PUT" , url , bytes.NewBuffer( put_data ) )
	request.Header.Set( "Content-Type" , "application/json" )
	request.Header.Set( "AUTH" , auth_token )
	if request_error != nil { fmt.Println( request_error ); return result }
	response , response_error := client.Do( request )
	if response_error != nil { fmt.Println( response_error ); return result }
	defer response.Body.Close()
	body , body_error := ioutil.ReadAll( response.Body )
	if body_error != nil { fmt.Println( body_error ); return result }
	fmt.Println( string( body[:] ) )
	var result_json SetAudioSettingResult
	json_decode_error := json.Unmarshal( body , &result_json )
	if json_decode_error != nil { fmt.Println( json_decode_error ); return result }
	//fmt.Println( result_json )
	result = result_json
	return result
}

type GetSettingsTypeResultItem struct {
	HASHVAL int64
	CNAME string
	TYPE string
	NAME string
}
type GetSettingsTypesResult struct {
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
func GetSettingsTypes( ip_address string , auth_token string ) ( result GetSettingsTypesResult ) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true ,
			} ,
		} ,
	}
	url := fmt.Sprintf( "https://%s:7345/menu_native/dynamic/tv_settings" , ip_address )
	request , request_error := http.NewRequest( "GET" , url , nil )
	request.Header.Set( "Content-Type" , "application/json" )
	request.Header.Set( "AUTH" , auth_token )
	if request_error != nil { fmt.Println( request_error ); return result }
	response , response_error := client.Do( request )
	if response_error != nil { fmt.Println( response_error ); return result }
	defer response.Body.Close()
	body , body_error := ioutil.ReadAll( response.Body )
	if body_error != nil { fmt.Println( body_error ); return result }
	fmt.Println( string( body[:] ) )
	var result_json GetSettingsTypesResult
	json_decode_error := json.Unmarshal( body , &result_json )
	if json_decode_error != nil { fmt.Println( json_decode_error ); return result }
	fmt.Println( result_json )
	result = result_json
	return result
}

// type GetSettingsForTypeResultItemValue interface{}
// type GetSettingsForTypeResultItem struct {
// 	INDEX int
// 	HASHVAL int64
// 	NAME string
// 	ENABLED string
// 	READONLY string
// 	CNAME string
// 	HIDDEN string
// 	TYPE string
// 	VALUE []GetSettingsForTypeResultItemValue
// }
// type GetSettingsForTypeResult struct {
// 	STATUS struct {
// 		RESULT string
// 		DETAIL string
// 	}
// 	HASHLIST []int64
// 	GROUP string
// 	NAME string
// 	PARAMETERS struct {
// 		FLAT string
// 		HELPTEXT string
// 		HASHONLY string
// 	}
// 	ITEMS []GetSettingsForTypeResultItem
// 	URI string
// 	CNAME string
// 	TYPE string
// }
type GetSettingsForTypeResult interface{}
func GetAllSettingsForType( ip_address string , auth_token string , setting_type string ) ( result GetSettingsForTypeResult ) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true ,
			} ,
		} ,
	}
	url := fmt.Sprintf( "https://%s:7345/menu_native/dynamic/tv_settings/%s" , ip_address , setting_type )
	request , request_error := http.NewRequest( "GET" , url , nil )
	request.Header.Set( "Content-Type" , "application/json" )
	request.Header.Set( "AUTH" , auth_token )
	if request_error != nil { fmt.Println( request_error ); return result }
	response , response_error := client.Do( request )
	if response_error != nil { fmt.Println( response_error ); return result }
	defer response.Body.Close()
	body , body_error := ioutil.ReadAll( response.Body )
	if body_error != nil { fmt.Println( body_error ); return result }
	fmt.Println( string( body[:] ) )
	var result_json GetSettingsForTypeResult
	json_decode_error := json.Unmarshal( body , &result_json )
	if json_decode_error != nil { fmt.Println( json_decode_error ); return result }
	fmt.Println( result_json )
	result = result_json
	return result
}

type GetAllSettingsForTypeResultItem interface {}
type GetAllSettingsForTypeResult struct {
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
	ITEMS []GetAllSettingsForTypeResultItem
	URI string
	CNAME string
	TYPE string
}
func GetAllSettingsOptionsForType( ip_address string , auth_token string , setting_type string ) ( result GetAllSettingsForTypeResult ) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true ,
			} ,
		} ,
	}
	url := fmt.Sprintf( "https://%s:7345/menu_native/static/tv_settings/%s" , ip_address , setting_type )
	request , request_error := http.NewRequest( "GET" , url , nil )
	request.Header.Set( "Content-Type" , "application/json" )
	request.Header.Set( "AUTH" , auth_token )
	if request_error != nil { fmt.Println( request_error ); return result }
	response , response_error := client.Do( request )
	if response_error != nil { fmt.Println( response_error ); return result }
	defer response.Body.Close()
	body , body_error := ioutil.ReadAll( response.Body )
	if body_error != nil { fmt.Println( body_error ); return result }
	fmt.Println( string( body[:] ) )
	var result_json GetAllSettingsForTypeResult
	json_decode_error := json.Unmarshal( body , &result_json )
	if json_decode_error != nil { fmt.Println( json_decode_error ); return result }
	fmt.Println( result_json )
	result = result_json
	return result
}

type GetSettingResultItem struct {
	HASHVAL int64
	CNAME string
	TYPE string
	NAME string
	VALUE int
}
type GetSettingResult struct {
	STATUS struct {
		RESULT string
		DETAIL string
	}
	PARAMETERS struct {
		FLAT string
		HELPTEXT string
		HASHONLY string
	}
	ITEMS []GetSettingResultItem
	HASHLIST []int64
	URI string
}
func GetSetting( ip_address string , auth_token string , setting_type string , setting_name string ) ( result GetSettingResult ) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true ,
			} ,
		} ,
	}
	url := fmt.Sprintf( "https://%s:7345/menu_native/dynamic/tv_settings/%s/%s" , ip_address , setting_type , setting_name )
	request , request_error := http.NewRequest( "GET" , url , nil )
	request.Header.Set( "Content-Type" , "application/json" )
	request.Header.Set( "AUTH" , auth_token )
	if request_error != nil { fmt.Println( request_error ); return result }
	response , response_error := client.Do( request )
	if response_error != nil { fmt.Println( response_error ); return result }
	defer response.Body.Close()
	body , body_error := ioutil.ReadAll( response.Body )
	if body_error != nil { fmt.Println( body_error ); return result }
	fmt.Println( string( body[:] ) )
	var result_json GetSettingResult
	json_decode_error := json.Unmarshal( body , &result_json )
	if json_decode_error != nil { fmt.Println( json_decode_error ); return result }
	fmt.Println( result_json )
	result = result_json
	return result
}

type SetSettingsOptionInt struct {
	_url string
	item_name string
	VALUE int
	HASHVAL int64
	REQUEST string
}
type SetSettingsOptionString struct {
	_url string
	item_name string
	VALUE string
	HASHVAL int64
	REQUEST string
}
func resolve_put_data( ip_address string , auth_token string , setting_type string , setting_name string , setting_value interface{} ) ( put_data []uint8 ) {
	current_setting := GetSetting( ip_address , auth_token , setting_type , setting_name )
	if reflect.TypeOf( setting_value ).String() == "string" {
		put_data , _ := json.Marshal(SetSettingsOptionString{
			_url: fmt.Sprintf( "/menu_native/dynamic/tv_settings/%s/%s" , setting_type , setting_name ) ,
			item_name: "SETTINGS" ,
			VALUE: string( setting_value.(string) ) ,
			HASHVAL: current_setting.ITEMS[0].HASHVAL ,
			REQUEST: "MODIFY" ,
		})
		return put_data
	} else if reflect.TypeOf( setting_value ).String() == "int" {
		put_data , _ := json.Marshal(SetSettingsOptionInt{
			_url: fmt.Sprintf( "/menu_native/dynamic/tv_settings/%s/%s" , setting_type , setting_name ) ,
			item_name: "SETTINGS" ,
			VALUE: int( setting_value.(int) ) ,
			HASHVAL: current_setting.ITEMS[0].HASHVAL ,
			REQUEST: "MODIFY" ,
		})
		return put_data
	}
	return
}
type SettingResultItem struct {
	HASHVAL int64
	NAME string
}
type SetSettingsForTypeResult struct {
	STATUS struct {
		RESULT string
		DETAIL string
	}
	PARAMETERS struct {
		HASHVAL int64
		REQUEST string
		VALUE interface{}
	}
	ITEMS []SettingResultItem
	HASHLIST []int64
	URI string
}
func SetSettingsOption( ip_address string , auth_token string  , setting_type string , setting_name string , setting_value interface{} ) ( result SetSettingsForTypeResult ) {
	put_data := resolve_put_data( ip_address , auth_token , setting_type , setting_name , setting_value )
	url := fmt.Sprintf( "https://%s:7345/menu_native/dynamic/tv_settings/%s/%s" , ip_address , setting_type , setting_name )
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true ,
			} ,
		} ,
	}
	request , request_error := http.NewRequest( "PUT" , url , bytes.NewBuffer( put_data ) )
	request.Header.Set( "Content-Type" , "application/json" )
	request.Header.Set( "AUTH" , auth_token )
	if request_error != nil { fmt.Println( request_error ); return result }
	response , response_error := client.Do( request )
	if response_error != nil { fmt.Println( response_error ); panic( response_error ) }
	defer response.Body.Close()
	body , body_error := ioutil.ReadAll( response.Body )
	if body_error != nil { fmt.Println( body_error ); return result }
	fmt.Println( string( body[:] ) )
	var result_json SetSettingsForTypeResult
	json_decode_error := json.Unmarshal( body , &result_json )
	if json_decode_error != nil { fmt.Println( json_decode_error ); return result }
	fmt.Println( result_json )
	result = result_json
	return result
}

type GetCurrentAppResult struct {
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
func GetCurrentApp( ip_address string , auth_token string ) ( result GetCurrentAppResult ) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true ,
			} ,
		} ,
	}
	url := fmt.Sprintf( "https://%s:7345/app/current" , ip_address )
	request , request_error := http.NewRequest( "GET" , url , nil )
	request.Header.Set( "Content-Type" , "application/json" )
	request.Header.Set( "AUTH" , auth_token )
	if request_error != nil { fmt.Println( request_error ); return result }
	response , response_error := client.Do( request )
	if response_error != nil { fmt.Println( response_error ); return result }
	defer response.Body.Close()
	body , body_error := ioutil.ReadAll( response.Body )
	if body_error != nil { fmt.Println( body_error ); return result }
	fmt.Println( string( body[:] ) )
	var result_json GetCurrentAppResult
	json_decode_error := json.Unmarshal( body , &result_json )
	if json_decode_error != nil { fmt.Println( json_decode_error ); return result }
	fmt.Println( result_json )
	result = result_json
	return result
}

// Look Here to Find APP_ID 's , Namespace Integers , and Messages
// https://github.com/vkorn/pyvizio/blob/master/pyvizio/const.py
type LauchAppDataValue struct {
	APP_ID string
	NAME_SPACE int
	MESSAGE string
}
type LaunchAppData struct {
	_url string
	VALUE LauchAppDataValue
}
type LaunchAppResult struct {
	STATUS struct {
		RESULT string
		DETAIL string
	}
}
func LaunchApp( ip_address string , auth_token string  , app_id string , name_space int , message string ) ( result LaunchAppResult ) {
	put_data , _ := json.Marshal(LaunchAppData{
		_url: "/app/launch" ,
		VALUE: LauchAppDataValue{
			APP_ID: app_id ,
			NAME_SPACE: name_space ,
			MESSAGE: message ,
		} ,
	})
	url := fmt.Sprintf( "https://%s:7345/app/launch" , ip_address )
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true ,
			} ,
		} ,
	}
	request , request_error := http.NewRequest( "PUT" , url , bytes.NewBuffer( put_data ) )
	request.Header.Set( "Content-Type" , "application/json" )
	request.Header.Set( "AUTH" , auth_token )
	if request_error != nil { fmt.Println( request_error ); return result }
	response , response_error := client.Do( request )
	if response_error != nil { fmt.Println( response_error ); panic( response_error ) }
	defer response.Body.Close()
	body , body_error := ioutil.ReadAll( response.Body )
	if body_error != nil { fmt.Println( body_error ); return result }
	fmt.Println( string( body[:] ) )
	var result_json LaunchAppResult
	json_decode_error := json.Unmarshal( body , &result_json )
	if json_decode_error != nil { fmt.Println( json_decode_error ); return result }
	fmt.Println( result_json )
	result = result_json
	return result
}