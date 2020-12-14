# Vizio Controller

import (
	"time"
	"fmt"
	"github.com/go-redis/redis/v8"
	viziocontroller "github.com/0187773933/VizioController"
)

```

func get_redis_connection( address string , db int , password string ) ( redis_connection *redis.Client ) {
	redis_connection = redis.NewClient(&redis.Options{
		Addr: address ,
		Password: password ,
		DB: db ,
	})
	return
}

func main() {
	// RegenerateAuthToken()

	// Get IP , Auth Token , and Mac Address From Redis
	var ctx = context.Background()
	redis_connection := get_redis_connection( "localhost:6379", 3 , "" )
	ip_address , err := redis_connection.Get( ctx , "STATE.VIZIO_TV.IP_ADDRESS" ).Result()
	if err != nil { panic( err ) }
	auth_token , err := redis_connection.Get( ctx , "STATE.VIZIO_TV.AUTH_TOKEN" ).Result()
	if err != nil { panic( err ) }
	// mac_address , err = redis_connection.Get( ctx , "CONFIG.VIZIO_TV.MAC_ADDRESS" ).Result()
	// if err != nil { panic( err ) }

	// Power
	power_state := GetPowerState( ip_address , auth_token )
	fmt.Println( power_state )
	PowerOff( ip_address , auth_token )
	time.Sleep( 2 * time.Second )
	PowerOn( ip_address , auth_token )

	// Volume
	volume := GetVolume( ip_address , auth_token )
	fmt.Println( volume )
	VolumeDown( ip_address , auth_token )
	time.Sleep( 1 * time.Second )
	VolumeUp( ip_address , auth_token )
	time.Sleep( 1 * time.Second )
	MuteOff( ip_address , auth_token )
	time.Sleep( 1 * time.Second )
	MuteOn( ip_address , auth_token )
	time.Sleep( 1 * time.Second )
	MuteToggle( ip_address , auth_token )

	// Inputs
	current_input := GetCurrentInput( ip_address , auth_token )
	fmt.Println( current_input )
	time.Sleep( 1 * time.Second )
	available_inputs := GetAvailableInputs( ip_address , auth_token )
	fmt.Println( available_inputs )
	SetInput( ip_address , auth_token , "HDMI-2" )
	time.Sleep( 2 * time.Second )
	SetInput( ip_address , auth_token , "HDMI-1" )
	time.Sleep( 2 * time.Second )
	CycleInput( ip_address , auth_token )

	// Audio Settings
	audio_settings_tv_speakers := GetAudioSetting( ip_address , auth_token , "tv_speakers" )
	fmt.Println( audio_settings_tv_speakers )
	GetAllAudioSettingsOptions( ip_address , auth_token )
	tv_speakers := GetAudioSettingsOption( ip_address , auth_token , "tv_speakers" )
	fmt.Println( tv_speakers )
	SetAudioSetting( ip_address , auth_token , "mute" , "Off" )

	// Other Settings
	GetSettingsTypes( ip_address , auth_token )
	GetAllSettingsForType( ip_address , auth_token , "network" )
	GetAllSettingsOptionsForType( ip_address , auth_token , "picture" )
	GetSetting( ip_address , auth_token , "picture" , "backlight" )
	SetSettingsOption( ip_address , auth_token , "picture" , "backlight" , 100 )

	// Smart Apps
	GetCurrentApp( ip_address , auth_token )
	// Look Here to Find APP_ID 's , Namespace Integers , and Messages
	// https://github.com/vkorn/pyvizio/blob/master/pyvizio/const.py
	// LaunchApp( ip_address , auth_token , "75" , 4 , "https://cd-dmgz.bamgrid.com/bbd/vizio_tv/index.html" )
	LaunchApp( ip_address , auth_token , "hdmi1" , 8 , "None" ) // Disney+

}
```