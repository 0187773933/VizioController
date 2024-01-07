package main

import (
	"fmt"
	// "time"
	utils "github.com/0187773933/VizioController/v1/utils"
	controller "github.com/0187773933/VizioController/v1/controller"
)

func main() {
	config := utils.GetConfig()
	x := controller.New( &config )
	fmt.Println( x )
	// fmt.Println( x.InputGetCurrent() )
	// fmt.Println( x.InputGetAvailable() )

	// x.InputSet( "TV" )
	// time.Sleep( 3 * time.Second )
	// x.InputSet( "HDMI-2" )

	// x.AudioGetSetting( "mute" )
	// x.AudioSetSetting( "mute" , "On" )
	// time.Sleep( 1 * time.Second )
	// x.AudioSetSetting( "mute" , "Off" )

	// x.SettingsGetTypes()
	// x.SettingsGetType( "network" )
	// x.SettingsGetOptionsForType( "picture" )

	// x.SettingsGet( "picture" , "backlight" )
	// x.SettingsSet( "picture" , "backlight" , 100 )

	// x.AppGetCurrent()
	// x.AppLaunch( "hdmi1" , 8 , "" )
	x.AppLaunch( "1" , 3 , "" ) // netflix
}