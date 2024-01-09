package main

import (
	"fmt"
	// "time"
	// utils "github.com/0187773933/VizioController/v1/utils"
	controller "github.com/0187773933/VizioController/v1/controller"
)

func main() {

	x := controller.New( "192.168.4.194" , "Zloom5087o" )
	// config := utils.GetConfig()
	// x := controller.NewFromConfig( &config )
	// fmt.Println( x )

	// Power
	// fmt.Println( x.PowerGetState() )
	// x.PowerOn()
	// x.PowerOff()

	// Volume
	fmt.Println( x.VolumeGet() )
	// x.VolumeSet( 2 )
	// x.VolumeUp()
	// x.PowerDown()
	// x.MuteOn()
	// x.MuteOff()
	// x.MuteToggle()

	// Inputs
	// fmt.Println( x.InputGetCurrent() )
	// fmt.Println( x.InputGetAvailable() )
	// x.InputSet( "TV" )
	// x.InputSet( "hdmi1" )
	// time.Sleep( 3 * time.Second )
	// x.InputSet( "HDMI-2" )
	// x.InputCycle()

	// Audio Settings
	// fmt.Println( x.AudioGetAllSettings() )
	// fmt.Println( x.AudioGetSetting( "mute" ) )
	// x.AudioSetSetting( "mute" , "On" )
	// time.Sleep( 1 * time.Second )
	// x.AudioSetSetting( "mute" , "Off" )

	// Generic Settings
	// fmt.Println( x.SettingsGetTypes() )
	// fmt.Println( x.SettingsGetType( "network" ) )
	// fmt.Println( x.SettingsGetOptionsForType( "picture" ) )
	// fmt.Println( x.SettingsGet( "picture" , "backlight" ) )
	// x.SettingsSet( "picture" , "backlight" , 100 )

	// Apps
	// fmt.Println( x.AppGetCurrent() )
	// x.AppLaunch( "hdmi1" , 8 , "" )
	// x.AppLaunch( "1" , 3 , "" ) // netflix
}