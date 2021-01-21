package main

import (
    "fmt"
    "io/ioutil"
    "os"
    uj "github.com/nanoscopic/ujsonin/mod"
    log "github.com/sirupsen/logrus"
)

type CDevice struct {
    udid string
    width int
    height int
    clickWidth int
    clickHeight int
}

type Config struct {
    iosDeployPath      string
    mobiledevicePath   string
    iosVideoStreamPath string
    httpPort           int
    cfHost             string
    cfUsername         string
    devs               map [string] CDevice
}

func NewConfig( configPath string ) (*Config) {
    config := Config{}
    
    root := loadConfig( configPath )
    
    binPaths := root.Get("bin_paths")
    if binPaths == nil {
    }
    iosDeployNode := binPaths.Get("ios-deploy")
    if iosDeployNode == nil {
    }
    config.iosDeployPath = iosDeployNode.String()
    mobiledeviceNode := binPaths.Get("mobiledevice")
    if mobiledeviceNode == nil {
    }
    config.mobiledevicePath = mobiledeviceNode.String()
    ivsNode := binPaths.Get("ios_video_stream")
    if ivsNode == nil {
    }
    config.iosVideoStreamPath = ivsNode.String()
    
    portNode := root.Get("port")
    if root == nil {
    }
    config.httpPort = portNode.Int()
    
    cfNode := root.Get("controlfloor")
    if cfNode == nil {
    }
    cfHostNode := cfNode.Get("host")
    if cfHostNode == nil {
    }
    config.cfHost = cfHostNode.String()
    cfIdNode := cfNode.Get("username")
    if cfIdNode == nil {
    }
    config.cfUsername = cfIdNode.String()
    
    config.devs = readDevs( root )
    
    return &config
}

func readDevs( root *uj.JNode ) ( map[string]CDevice ) {
    devs := make( map[string]CDevice )
    
    devsNode := root.Get("devices")
    devsNode.ForEach( func( devNode *uj.JNode ) {
        udid := devNode.Get("udid").String()
        width := devNode.Get("mainScreenWidth").Int()
        height := devNode.Get("mainScreenHeight").Int()
        scale := devNode.Get("mainScreenScale").Int()
            
        dev := CDevice{
            udid: udid,
            width: width,
            height: height,
            clickWidth: (width/scale),
            clickHeight: (height/scale),
        }
        devs[ udid ] = dev
    } )
    return devs
}

func loadConfig( configPath string ) (*uj.JNode) {
    fh, serr := os.Stat( configPath )
    if serr != nil {
        log.WithFields( log.Fields{
            "type":        "err_read_config",
            "error":       serr,
            "config_path": configPath,
        } ).Fatal("Could not read specified config path")
    }
    configFile := configPath
    switch mode := fh.Mode(); {
        case mode.IsDir(): configFile = fmt.Sprintf("%s/config.json", configPath)
    }
    content, err := ioutil.ReadFile( configFile )
	if err != nil { log.Fatal( err ) }
	
    root, _ := uj.Parse( content )
    
    return root
}