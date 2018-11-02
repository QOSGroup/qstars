package config

import (
"bytes"
"path/filepath"
"text/template"

cmn "github.com/tendermint/tendermint/libs/common"
)

var (
	defaultConfigDir     = "config"
	defaultDataDir       = "data"
	defaultConfigFileName  = "config.toml"
	defaultGenesisJSONName = "genesis.json"
	defaultConfigFilePath  = filepath.Join(defaultConfigDir, defaultConfigFileName)
)

var configTemplate *template.Template

func init() {
	var err error
	if configTemplate, err = template.New("configFileTemplate").Parse(defaultConfigTemplate); err != nil {
		panic(err)
	}
}

// EnsureRoot creates the root, config, and data directories if they don't exist,
// and panics if it fails.
func EnsureRoot(rootDir string) {
	if err := cmn.EnsureDir(rootDir, 0700); err != nil {
		cmn.PanicSanity(err.Error())
	}
	if err := cmn.EnsureDir(filepath.Join(rootDir, defaultConfigDir), 0700); err != nil {
		cmn.PanicSanity(err.Error())
	}
	if err := cmn.EnsureDir(filepath.Join(rootDir, defaultDataDir), 0700); err != nil {
		cmn.PanicSanity(err.Error())
	}

	configFilePath := filepath.Join(rootDir, defaultConfigFilePath)

	// Write default config file if missing.
	if !cmn.FileExists(configFilePath) {
		writeDefaultConfigFile(configFilePath)
	}
}

// XXX: this func should probably be called by cmd/tendermint/commands/init.go
// alongside the writing of the genesis.json and priv_validator.json
func writeDefaultConfigFile(configFilePath string) {
	WriteConfigFile(configFilePath, DefaultConfig())
}

// WriteConfigFile renders config using the template and writes it to configFilePath.
func WriteConfigFile(configFilePath string, config *CLIConfig) {
	var buffer bytes.Buffer

	if err := configTemplate.Execute(&buffer, config); err != nil {
		panic(err)
	}

	cmn.MustWriteFile(configFilePath, buffer.Bytes(), 0644)
}

// Note: any changes to the comments/variables/mapstructure
// must be reflected in the appropriate struct in config/config.go
const defaultConfigTemplate = `# This is a TOML config file.
# For more information, see https://github.com/toml-lang/toml

# Path to the JSON file containing the initial validator set and other meta data
qos_chain_id = "{{ .QOSChainID }}"

qsc_chain_id = "{{ .QSCChainID }}"

qos_node_uri = "{{ .QOSNodeURI }}"

qstars_node_uri = "{{ .QSTARSNodeURI }}"

direct_to_qos = "{{ .DirectTOQOS }}"
`



