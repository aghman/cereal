package config

/*
 * Example yaml config
 *
  location:
	  name: near_pond
		serialport: /dev/usb01
		output_file: /var/log/cereal.out
    tags:
      - airquality
      - yard
      - outdoors
*/

type Location struct {
	Name       string   `yaml:"name"`
	SerialPort string   `yaml:"serialport"`
	OutputFile string   `yaml:"outputfile"`
	Tags       []string `yaml:"tags"`
}

type CerealConfig struct {
	Location Location
}

func NewCerealConfig(input map[string]interface{}) *CerealConfig {
	var newConfig CerealConfig
	newConfig.Location.Name = input["name"].(string)
	newConfig.Location.SerialPort = input["serialport"].(string)
	newConfig.Location.OutputFile = input["outputfile"].(string)

	return &newConfig
}
