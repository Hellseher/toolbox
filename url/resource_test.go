package url_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/viant/toolbox/url"
	"github.com/viant/toolbox"
	"os"
	"fmt"
	"path"
)

func TestNewResource(t *testing.T) {

	var resource = url.NewResource("https://raw.githubusercontent.com/viant/toolbox/master/LICENSE.txt")
	assert.EqualValues(t, resource.ParsedURL.String(), "https://raw.githubusercontent.com/viant/toolbox/master/LICENSE.txt")
	data, err := resource.Download()
	assert.Nil(t, err)
	assert.NotNil(t, data)

}



func TestResource_YamlDecode(t *testing.T) {
	var filename = path.Join(os.Getenv("TMPDIR"), "resource.yaml")
	toolbox.RemoveFileIfExist(filename)
	defer toolbox.RemoveFileIfExist(filename)
	var aMap = map[string]interface{}{
		"a": 1,
		"b": "123",
		"c": []int{1, 3, 6},
	}
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0644)
	fmt.Printf("%v\n", filename)
	if assert.Nil(t, err) {
		err =toolbox.NewYamlEncoderFactory().Create(file).Encode(aMap)
		assert.Nil(t, err)
	}

	var resource = url.NewResource(filename)
	assert.EqualValues(t, resource.ParsedURL.String(), toolbox.FileSchema+filename)

	var resourceData = make(map[string]interface{})
	err = resource.YamlDecode(&resourceData)
	assert.Nil(t, err)

	assert.EqualValues(t, resourceData["a"], 1)
	assert.EqualValues(t, resourceData["b"], "123")

}



func TestResource_JsonDecode(t *testing.T) {
	var filename = path.Join(os.Getenv("TMPDIR"), "resource.json")
	toolbox.RemoveFileIfExist(filename)
	defer toolbox.RemoveFileIfExist(filename)
	var aMap = map[string]interface{}{
		"a": 1,
		"b": "123",
	}
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0644)
	fmt.Printf("%v\n", filename)
	if assert.Nil(t, err) {
		err =toolbox.NewJSONEncoderFactory().Create(file).Encode(aMap)
		assert.Nil(t, err)
	}

	var resource = url.NewResource(filename)
	assert.EqualValues(t, resource.ParsedURL.String(), toolbox.FileSchema+filename)

	var resourceData = make(map[string]interface{})
	err = resource.JsonDecode(&resourceData)
	assert.Nil(t, err)

	assert.EqualValues(t, resourceData["a"], 1)
	assert.EqualValues(t, resourceData["b"], "123")

}

func TestResource_LoadCredential(t *testing.T) {

	{
		var filename = path.Join(os.Getenv("TMPDIR"), "resource_secret.json")
		toolbox.RemoveFileIfExist(filename)
		defer toolbox.RemoveFileIfExist(filename)
		var aMap = map[string]interface{}{
			"username": "uzytkownik",
			"password": "haslo",
		}
		file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0644)
		fmt.Printf("%v\n", filename)
		if assert.Nil(t, err) {
			err = toolbox.NewJSONEncoderFactory().Create(file).Encode(aMap)
			assert.Nil(t, err)
		}
		var resource= url.NewResource("https://raw.githubusercontent.com/viant/toolbox/master/LICENSE.txt", filename)

		username, password, err := resource.LoadCredential(false)
		assert.Nil(t, err)
		assert.Equal(t, username, "uzytkownik")
		assert.Equal(t, password, "haslo")
	}

	{//error case

		var resource= url.NewResource("https://raw.githubusercontent.com/viant/toolbox/master/LICENSE.txt")
		_, _, err := resource.LoadCredential(true)
		assert.NotNil(t, err)

	}

	{//error case

		var resource= url.NewResource("https://raw.githubusercontent.com/viant/toolbox/master/LICENSE.txt", "bogus 343")
		_, _, err := resource.LoadCredential(true)
		assert.NotNil(t, err)

	}


}