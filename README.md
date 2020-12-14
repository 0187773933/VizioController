# Vizio Controller

```
import (
	"context"
	"time"
	"fmt"
	"github.com/go-redis/redis/v8"
	viziocontroller "github.com/0187773933/VizioController"
)

func get_redis_connection( address string , db int , password string ) ( redis_connection *redis.Client ) {
	redis_connection = redis.NewClient(&redis.Options{
		Addr: address ,
		Password: password ,
		DB: db ,
	})
	return
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
	power_state := viziocontroller.GetPowerState( ip_address , auth_token )
	fmt.Println( power_state )
	viziocontroller.PowerOff( ip_address , auth_token )
	time.Sleep( 2 * time.Second )
	viziocontroller.PowerOn( ip_address , auth_token )

	// Volume
	volume := viziocontroller.GetVolume( ip_address , auth_token )
	fmt.Println( volume )
	viziocontroller.VolumeDown( ip_address , auth_token )
	time.Sleep( 1 * time.Second )
	viziocontroller.VolumeUp( ip_address , auth_token )
	time.Sleep( 1 * time.Second )
	viziocontroller.MuteOff( ip_address , auth_token )
	time.Sleep( 1 * time.Second )
	viziocontroller.MuteOn( ip_address , auth_token )
	time.Sleep( 1 * time.Second )
	viziocontroller.MuteToggle( ip_address , auth_token )

	// Inputs
	current_input := viziocontroller.GetCurrentInput( ip_address , auth_token )
	fmt.Println( current_input )
	time.Sleep( 1 * time.Second )
	available_inputs := viziocontroller.GetAvailableInputs( ip_address , auth_token )
	fmt.Println( available_inputs )
	viziocontroller.SetInput( ip_address , auth_token , "HDMI-2" )
	time.Sleep( 2 * time.Second )
	viziocontroller.SetInput( ip_address , auth_token , "HDMI-1" )
	time.Sleep( 2 * time.Second )
	viziocontroller.CycleInput( ip_address , auth_token )

	// Audio Settings
	audio_settings_tv_speakers := viziocontroller.GetAudioSetting( ip_address , auth_token , "tv_speakers" )
	fmt.Println( audio_settings_tv_speakers )
	viziocontroller.GetAllAudioSettingsOptions( ip_address , auth_token )
	tv_speakers := viziocontroller.GetAudioSettingsOption( ip_address , auth_token , "tv_speakers" )
	fmt.Println( tv_speakers )
	viziocontroller.SetAudioSetting( ip_address , auth_token , "mute" , "Off" )

	// Other Settings
	viziocontroller.GetSettingsTypes( ip_address , auth_token )
	viziocontroller.GetAllSettingsForType( ip_address , auth_token , "network" )
	viziocontroller.GetAllSettingsOptionsForType( ip_address , auth_token , "picture" )
	viziocontroller.GetSetting( ip_address , auth_token , "picture" , "backlight" )
	viziocontroller.SetSettingsOption( ip_address , auth_token , "picture" , "backlight" , 100 )

	// Smart Apps
	viziocontroller.GetCurrentApp( ip_address , auth_token )
	// Look Here to Find APP_ID 's , Namespace Integers , and Messages
	// https://github.com/vkorn/pyvizio/blob/master/pyvizio/const.py
	// viziocontroller.LaunchApp( ip_address , auth_token , "75" , 4 , "https://cd-dmgz.bamgrid.com/bbd/vizio_tv/index.html" ) // Disney+
	viziocontroller.LaunchApp( ip_address , auth_token , "hdmi1" , 8 , "None" )

}
```