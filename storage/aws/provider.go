package aws

import (
	"github.com/viant/toolbox/storage"
	"github.com/viant/toolbox/url"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const ProviderScheme = "s3"

func init() {
	storage.NewStorageProvider().Registry[ProviderScheme] = serviceProvider
}

func serviceProvider(credentialFile string) (storage.Service, error) {
	if !strings.HasPrefix(credentialFile, "/") {
		dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		credentialFile = path.Join(dir, credentialFile)
	}
	s3config := &Config{}
	resource := url.NewResource(credentialFile)
	err := resource.Decode(credentialFile)
	if err != nil {
		return nil, err
	}
	return NewService(s3config), nil
}
