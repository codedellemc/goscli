package commands

import (
	"fmt"
	"log"
	"strings"

	"github.com/emccode/goscaleio"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v1"
)

var sdcCmdV *cobra.Command

func init() {
	addCommandsSdc()
	// sdcCmd.Flags().StringVar(&sdcname, "sdcname", "", "GOSCALEIO_TEMP")
	sdcCmd.Flags().StringVar(&systemid, "systemid", "", "GOSCALEIO_SYSTEMID")
	sdcgetCmd.Flags().StringVar(&systemid, "systemid", "", "GOSCALEIO_SYSTEMID")
	sdcgetCmd.Flags().StringVar(&sdcip, "sdcip", "", "GOSCALEIO_SDCIP")
	sdcgetCmd.Flags().StringVar(&sdcguid, "sdcguid", "", "GOSCALEIO_SDCGUID")
	sdcgetCmd.Flags().StringVar(&sdcid, "sdcid", "", "GOSCALEIO_SDCID")
	sdcgetCmd.Flags().StringVar(&sdcname, "sdcname", "", "GOSCALEIO_SDCNAME")

	sdcCmdV = sdcCmd

	// initConfig(sdcCmd, "goscli", true, map[string]FlagValue{
	// 	"endpoint": {endpoint, true, false, ""},
	// 	"insecure": {insecure, false, false, ""},
	// })

	sdcCmd.Run = func(cmd *cobra.Command, args []string) {
		setGobValues(cmd, "goscli", "")
		cmd.Usage()
	}
}

func addCommandsSdc() {
	sdcCmd.AddCommand(sdcgetCmd)
	sdcCmd.AddCommand(sdclocalgetCmd)
}

var sdcCmd = &cobra.Command{
	Use:   "sdc",
	Short: "sdc",
	Long:  `sdc`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

var sdcgetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a sdc",
	Long:  `Get a sdc`,
	Run:   cmdGetSdc,
}

var sdclocalgetCmd = &cobra.Command{
	Use:   "local",
	Short: "Get local sdc",
	Long:  `Get local sdc`,
	Run:   cmdGetSdcLocal,
}

func cmdGetSdc(cmd *cobra.Command, args []string) {
	client, err := authenticate()
	if err != nil {
		log.Fatalf("error authenticating: %v", err)
	}

	initConfig(cmd, "goscli", true, map[string]FlagValue{
		"systemhref": {&systemhref, true, false, ""},
	})

	systemhref = viper.GetString("systemhref")

	system, err := client.FindSystem("", systemhref)
	if err != nil {
		log.Fatalf("err: problem getting system: %v", err)
	}

	if sdcip == "" && sdcid == "" && sdcname == "" && sdcguid == "" {
		sdcs, err := system.GetSdc()
		if err != nil {
			log.Fatalf("error getting statistics: %v", err)
		}

		yamlOutput, err := yaml.Marshal(&sdcs)
		if err != nil {
			log.Fatalf("error marshaling: %s", err)
		}
		fmt.Println(string(yamlOutput))
		return
	}

	sdc := &goscaleio.Sdc{}
	switch {
	case sdcguid != "":
		sdc, err = system.FindSdc("SdcGuid", sdcguid)
	case sdcid != "":
		sdc, err = system.FindSdc("ID", sdcid)
	case sdcname != "":
		sdc, err = system.FindSdc("Name", sdcname)
	case sdcip != "":
		sdc, err = system.FindSdc("SdcIp", sdcip)
	}

	if err != nil {
		log.Fatalf("error finding Sdc: %v", err)
	}

	if len(args) > 1 {
		log.Fatalf("Too many arguments specified")
	}

	var yamlOutput []byte
	if len(args) == 1 {
		switch args[0] {
		case "statistics":
			statistics, err := sdc.GetStatistics()
			if err != nil {
				log.Fatalf("error getting statistics: %v", err)
			}

			yamlOutput, err = yaml.Marshal(&statistics)
		case "volume":
			volumes, err := sdc.GetVolume()
			if err != nil {
				log.Fatalf("error getting statistics: %v", err)
			}

			yamlOutput, err = yaml.Marshal(&volumes)
		default:
			log.Fatalf("parameter didn't match statistics|volume")
		}
	} else {
		yamlOutput, err = yaml.Marshal(&sdc)
		if err != nil {
			log.Fatalf("error marshaling: %s", err)
		}
	}
	fmt.Println(string(yamlOutput))
	return

}

func cmdGetSdcLocal(cmd *cobra.Command, args []string) {

	client, err := authenticate()
	if err != nil {
		log.Fatalf("error authenticating: %v", err)
	}

	initConfig(cmd, "goscli", true, map[string]FlagValue{
		"systemhref": {&systemhref, true, false, ""},
	})

	systemhref = viper.GetString("systemhref")

	system, err := client.FindSystem("", systemhref)
	if err != nil {
		log.Fatalf("err: problem getting system: %v", err)
	}

	sdcguid, err := goscaleio.GetSdcLocalGUID()
	if err != nil {
		log.Fatalf("Error getting local sdc guid: %s", err)
	}

	sdc, err := system.FindSdc("SdcGuid", strings.ToUpper(sdcguid))
	if err != nil {
		log.Fatalf("Error finding Sdc %s: %s", sdcguid, err)
	}

	if len(args) > 1 {
		log.Fatalf("Too many arguments specified")
	}

	var yamlOutput []byte
	if len(args) == 1 {
		switch args[0] {
		case "statistics":
			statistics, err := sdc.GetStatistics()
			if err != nil {
				log.Fatalf("error getting statistics: %v", err)
			}

			yamlOutput, err = yaml.Marshal(&statistics)
		case "volume":
			volumes, err := sdc.GetVolume()
			if err != nil {
				log.Fatalf("error getting statistics: %v", err)
			}

			yamlOutput, err = yaml.Marshal(&volumes)
		default:
			log.Fatalf("parameter didn't match statistics|volume")
		}
	} else {
		yamlOutput, err = yaml.Marshal(&sdc)
		if err != nil {
			log.Fatalf("error marshaling: %s", err)
		}
	}
	fmt.Println(string(yamlOutput))
	return

}
