package utils

import (
	fmt "fmt"
	os "os"
	ioutil "io/ioutil"
	"encoding/json"
	filepath "path/filepath"
	yaml "gopkg.in/yaml.v2"
	types "github.com/0187773933/VizioController/v1/types"
)

func PrettyPrint( data interface{} ) {
	indented , _ := json.MarshalIndent( data , "" , "    " )
	fmt.Println( string( indented ) )
}

func ParseConfig( file_path string ) ( result types.ConfigFile ) {
	config_file , _ := ioutil.ReadFile( file_path )
	error := yaml.Unmarshal( config_file , &result )
	if error != nil { panic( error ) }
	return
}

func GetConfig() ( result types.ConfigFile ) {
	config_file_path , _ := filepath.Abs( "./config.yaml" )
	if len( os.Args ) > 1 { config_file_path , _ = filepath.Abs( os.Args[ 1 ] ) }
	result = ParseConfig( config_file_path )
	fmt.Printf( "Loaded Config File From : %s\n" , config_file_path )
	result.ConfigFilePath = config_file_path
	return
}

func SaveConfig( file_path string , config types.ConfigFile ) {
	yaml_data , _ := yaml.Marshal( config )
	ioutil.WriteFile( file_path , yaml_data , 0644 )
}