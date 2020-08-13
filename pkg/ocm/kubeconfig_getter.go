package ocm

import (
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)


// Get the kubeconfig by cluster ID
type KubecfgGetter interface {
	// Get the kubeconfig given the cluster ID
	GetKubeConfig(id string) *rest.Config

	// Init should used to keep update the cache we have that store the cluster kubeconfigs
	// including cache all the configs when startup, and add/remove the configs when cluster
	// deleted or added.
	Init() error
}

func NewConfigFileGetter(configDir string) KubecfgGetter {
	return &configFileGetter{
		configDir: configDir,
	}
}

type configFileGetter struct {
	configDir   string
	kubeconfigs map[string]*rest.Config
}

func (c *configFileGetter) GetKubeConfig(id string) *rest.Config {

	return c.kubeconfigs[id]
}

func (c *configFileGetter) Init() error {

	var configFiles []string
	err := filepath.Walk(c.configDir, func(path string, info os.FileInfo, err error) error {
		fmt.Printf("visited file or dir: %q\n", path)
		if !info.IsDir() {
			configFiles = append(configFiles, path)
		}
		return nil
	})
	if err != nil {
		return err
	}
	c.kubeconfigs = map[string]*rest.Config{}

	for _, file := range configFiles {
		fmt.Printf("load kubeconfig %s\n", file)
		// use the current context in kubeconfig
		clientConfig, err := clientcmd.BuildConfigFromFlags("", file)
		if err != nil {
			panic(err.Error())
		}
		clustername := filepath.Base(file)
		c.kubeconfigs[clustername] = clientConfig
	}

	return nil

}
