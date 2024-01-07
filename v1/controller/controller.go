package controller

import (
	"fmt"
	"bytes"
	// "reflect"
	"encoding/json"
	"net/http"
	"crypto/tls"
	"io/ioutil"
	types "github.com/0187773933/VizioController/v1/types"
	utils "github.com/0187773933/VizioController/v1/utils"
)

type Controller struct {
	Config types.ConfigFile `yaml:"config"`
	HTTPClient *http.Client `yaml:"-"`
}

func New( config *types.ConfigFile ) ( result *Controller ) {
	result = &Controller{
		Config: *config ,
		HTTPClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{ InsecureSkipVerify : true } ,
			} ,
		} ,
	}
	if config.IPAddress == "" { panic( "You didn't pass the TV's IP Address in the config" ) }
	if config.DeviceID == "" { config.DeviceID = "govizio" }
	if config.DeviceName == "" { config.DeviceName = "Golang Vizio" }
	if config.AuthToken == "" {
		fmt.Println( "You didn't pass an auth token in the config , re-pairing with tv" )
		result.Pair()
		utils.SaveConfig( config.ConfigFilePath , result.Config )
		fmt.Println( "Paired Successfully , saved auth token in config file" )
	}
	return
}

func ( ctrl *Controller ) API( method string , endpoint string , request_obj interface{} , response_obj interface{} ) {
	var body_data []byte
	var err error
	if method == "PUT" && request_obj != nil {
		// fmt.Println( "We got a request object, json stringifying" )
		body_data , err = json.Marshal( request_obj )
		if err != nil { fmt.Println( err ); return }
	}
	url := fmt.Sprintf( "https://%s:7345%s" , ctrl.Config.IPAddress , endpoint )
	req , err := http.NewRequest( method , url , bytes.NewBuffer( body_data ) )
	if err != nil { fmt.Println( err ); return }
	req.Header.Set( "Content-Type" , "application/json" )
	if ctrl.Config.AuthToken != "" { req.Header.Set( "AUTH" , ctrl.Config.AuthToken ) }
	resp , err := ctrl.HTTPClient.Do( req )
	if err != nil { fmt.Println( err ); return }
	defer resp.Body.Close()
	body , err := ioutil.ReadAll( resp.Body )
	if err != nil { fmt.Println( err ); return }
	if response_obj != nil {
		json.Unmarshal( body , response_obj )
		if err != nil { fmt.Println( err ); return }
	}
	err = json.Unmarshal( body , response_obj )
	if err != nil { fmt.Println( err ); return }
	return
}

func ( ctrl *Controller ) PairingStageOne() ( result int32 ) {
	result = 0
	url_part := "/pairing/start"
	put_data := types.PairingStageOneRequest{
		URL: url_part ,
		DeviceID: ctrl.Config.DeviceID ,
		DeviceName: ctrl.Config.DeviceName ,
	}
	var response types.PairingStageOneResponse
	ctrl.API( "PUT" , url_part , &put_data , &response )
	result = response.ITEM.PAIRING_REQ_TOKEN
	return
}

func ( ctrl *Controller ) PairingStageTwo( pairing_stage_one_token int32 , code_displayed_on_tv string ) ( result string ) {
	url_part := "/pairing/pair"
	put_data := types.PairingStageTwoRequest{
		URL: url_part ,
		DeviceID: ctrl.Config.DeviceID ,
		DeviceName: ctrl.Config.DeviceName ,
		ChallengeType: 1 ,
		PairingReqToken: pairing_stage_one_token ,
		ResponseValue: code_displayed_on_tv ,
	}
	var response types.PairingStageTwoResponse
	ctrl.API( "PUT" , url_part , &put_data , &response )
	result = response.ITEM.AUTH_TOKEN
	return
}

func ( ctrl *Controller ) Pair() ( auth_token string ) {
	stage_one_token := ctrl.PairingStageOne()
	var code_displayed_on_tv string
	fmt.Println( "Please Type Code Displayed At the Top of the TV , and then Press Enter" )
	fmt.Scanln( &code_displayed_on_tv )
	auth_token = ctrl.PairingStageTwo( stage_one_token , code_displayed_on_tv )
	ctrl.Config.AuthToken = auth_token
	return
}

// Send Key Press
func ( ctrl *Controller ) KeyPress( codeset int , code int ) ( result bool ) {
	result = false
	url_part := "/key_command/"
	key_press_data := types.KeyCommand{
		URL: url_part ,
		Keylist: []types.Key{
			{
				Codeset: codeset ,
				Code: code ,
				Action: "KEYPRESS" ,
			} ,
		} ,
	}
	var response types.KeyPressResponse
	ctrl.API( "PUT" , url_part , &key_press_data , &response )
	if response.STATUS.RESULT == "SUCCESS" {
		result = true
	}
	return
}


// Power
func ( ctrl *Controller ) PowerGetState() ( result bool ) {
	var response types.PowerGetStateResponse
	ctrl.API( "GET" , "/state/device/power_mode" , nil , &response )
	switch response.ITEMS[0].VALUE {
		case 1:
			result = true
			break
		case 0:
			result = false
			break
	}
	return
}

func ( ctrl *Controller ) PowerOn() ( result bool ) {
	result = ctrl.KeyPress( 11 , 1 )
	return
}

func ( ctrl *Controller ) PowerOff() ( result bool ) {
	result = ctrl.KeyPress( 11 , 0 )
	return
}

// Volume
func ( ctrl *Controller ) VolumeGet() ( result int ) {
	var response types.VolumeGetResponse
	ctrl.API( "GET" , "/menu_native/dynamic/tv_settings/audio/volume" , nil , &response )
	result = response.ITEMS[ 0 ].VALUE
	return
}

func ( ctrl *Controller ) VolumeUp() ( result bool ) {
	result = ctrl.KeyPress( 5 , 1 )
	return
}

func ( ctrl *Controller ) VolumeDown() ( result bool ) {
	result = ctrl.KeyPress( 5 , 0 )
	return
}

func ( ctrl *Controller ) MuteOn() ( result bool ) {
	result = ctrl.KeyPress( 5 , 3 )
	return
}

func ( ctrl *Controller ) MuteOff() ( result bool ) {
	result = ctrl.KeyPress( 5 , 2 )
	return
}

func ( ctrl *Controller ) MuteToggle() ( result bool ) {
	result = ctrl.KeyPress( 5 , 4 )
	return
}

// Input
func ( ctrl *Controller ) InputGetCurrent() ( result types.Input ) {
	var response types.InputGetCurrentResponse
	ctrl.API( "GET" , "/menu_native/dynamic/tv_settings/devices/current_input" , nil , &response )
	result.Name = response.ITEMS[ 0 ].VALUE
	result.HashValue = response.ITEMS[ 0 ].HASHVAL
	return
}

func ( ctrl *Controller ) InputGetAvailable() ( result []types.Input ) {
	var response types.InputGetAvailableResponse
	ctrl.API( "GET" , "/menu_native/dynamic/tv_settings/devices/name_input" , nil , &response )
	for _ , item := range response.ITEMS {
		result = append( result , types.Input{
			Name: item.VALUE.NAME ,
			HashValue: item.HASHVAL ,
		})
	}
	return
}

func ( ctrl *Controller ) InputSet( input_name string ) ( result []types.Input ) {
	url_part := "/menu_native/dynamic/tv_settings/devices/current_input"
	current_input := ctrl.InputGetCurrent()
	var response types.InputGetAvailableResponse
	put_data := types.InputSetRequest{
		URL: url_part ,
		ItemName: "CURRENT_INPUT" ,
		VALUE: input_name ,
		HASHVAL: current_input.HashValue ,
		REQUEST: "MODIFY" ,
	}
	ctrl.API( "PUT" , url_part , &put_data , &response )
	fmt.Println( response )
	return
}

func ( ctrl *Controller ) InputCycle() ( result bool ) {
	result = ctrl.KeyPress( 7 , 1 )
	return
}

// Audio

func ( ctrl *Controller ) AudioGetSetting( setting_name string ) ( response types.AudioGetSettingResponse ) { // dynamic ??
	url_part := fmt.Sprintf( "/menu_native/dynamic/tv_settings/audio/%s" , setting_name )
	ctrl.API( "GET" , url_part , nil , &response )
	return
}

func ( ctrl *Controller ) AudioGetSettingsOption( setting_name string ) { // static ??
	var response types.AudioGetSettingsOptionResponse
	url_part := fmt.Sprintf( "/menu_native/static/tv_settings/audio/%s" , setting_name )
	ctrl.API( "GET" , url_part , nil , &response )
	// TODO , figure out why we care about this , and what to actually return
	fmt.Println( response )
}

func ( ctrl *Controller ) AudioGetAllSettings() ( response types.AudioGetAllSettingsResponse ) {
	url_part := "/menu_native/static/tv_settings/audio"
	ctrl.API( "GET" , url_part , nil , &response )
	// utils.PrettyPrint( response )
	return
}

func ( ctrl *Controller ) AudioSetSetting( setting_name string , setting_option string ) {
	current_setting := ctrl.AudioGetSetting( setting_name )
	fmt.Println( current_setting )
	url_part := fmt.Sprintf( "/menu_native/dynamic/tv_settings/audio/%s" , setting_name )
	put_data := types.InputSetRequest{
		URL: url_part ,
		ItemName: "SETTINGS" ,
		VALUE: setting_option ,
		HASHVAL: current_setting.ITEMS[ 0 ].HASHVAL ,
		REQUEST: "MODIFY" ,
	}
	var response types.AudioSetSettingResponse
	ctrl.API( "PUT" , url_part , &put_data , &response )
	fmt.Println( response )
	return
}

// Generic Settings

func ( ctrl *Controller ) SettingsGetTypes() ( response types.SettingsGetTypesResponse ) {
	url_part := "/menu_native/dynamic/tv_settings"
	ctrl.API( "GET" , url_part , nil , &response )
	utils.PrettyPrint( response )
	return
}

func ( ctrl *Controller ) SettingsGetType( settings_type string ) ( response types.SettingsGetTypeResponse ) {
	url_part := fmt.Sprintf( "/menu_native/dynamic/tv_settings/%s" , settings_type )
	ctrl.API( "GET" , url_part , nil , &response )
	utils.PrettyPrint( response )
	return
}

func ( ctrl *Controller ) SettingsGetOptionsForType( settings_type string ) ( response types.SettingsGetTypeResponse ) {
	url_part := fmt.Sprintf( "/menu_native/static/tv_settings/%s" , settings_type )
	ctrl.API( "GET" , url_part , nil , &response )
	utils.PrettyPrint( response )
	return
}

func ( ctrl *Controller ) SettingsGet( settings_type string , setting_name string ) ( response types.SettingsGetResponse ) {
	url_part := fmt.Sprintf( "/menu_native/dynamic/tv_settings/%s/%s" , settings_type , setting_name )
	ctrl.API( "GET" , url_part , nil , &response )
	utils.PrettyPrint( response )
	return
}

func ( ctrl *Controller ) SettingsSet( settings_type string , setting_name string , setting_value interface{} ) ( response types.SettingsSetResponse ) {
	current_setting := ctrl.SettingsGet( settings_type , setting_name )
	url_part := fmt.Sprintf( "/menu_native/dynamic/tv_settings/%s/%s" , settings_type , setting_name )
	ctrl.API( "GET" , url_part , nil , &response )
	switch setting_value.( type )  {
		case string:
			put_data := types.SettingsSetStringRequest{
				URL: url_part ,
				ItemName: "SETTINGS" ,
				VALUE: string( setting_value.( string ) ) ,
				HASHVAL: current_setting.ITEMS[ 0 ].HASHVAL ,
				REQUEST: "MODIFY" ,
			}
			ctrl.API( "PUT" , url_part , &put_data , &response )
			utils.PrettyPrint( response )
			break
		case int:
			put_data := types.SettingsSetIntRequest{
				URL: url_part ,
				ItemName: "SETTINGS" ,
				VALUE: int( setting_value.( int ) ) ,
				HASHVAL: current_setting.ITEMS[ 0 ].HASHVAL ,
				REQUEST: "MODIFY" ,
			}
			ctrl.API( "PUT" , url_part , &put_data , &response )
			utils.PrettyPrint( response )
			break
	}
	return
}

// App Stuff

func ( ctrl *Controller ) AppGetCurrent() ( response types.AppGetCurrentResponse ) {
	url_part := "/app/current"
	ctrl.API( "GET" , url_part , nil , &response )
	utils.PrettyPrint( response )
	return
}

// Look Here to Find APP_ID 's , Namespace Integers , and Messages
// https://github.com/vkorn/pyvizio/blob/master/pyvizio/const.py
func ( ctrl *Controller ) AppLaunch( app_id string , name_space int , message string ) ( response types.AppLaunchResponse ) {
	url_part := "/app/launch"
	put_data := types.AppLaunchRequest{
		URL: url_part ,
		VALUE: types.AppLaunchRequestValue{
			APP_ID: app_id ,
			NAME_SPACE: name_space ,
			MESSAGE: message ,
		} ,
	}
	ctrl.API( "PUT" , url_part , &put_data , &response )
	utils.PrettyPrint( response )
	return
}



/*



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

*/