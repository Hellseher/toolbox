/*
 *
 *
 * Copyright 2012-2016 Viant.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not
 *  use this file except in compliance with the License. You may obtain a copy of
 *  the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
 *  License for the specific language governing permissions and limitations under
 *  the License.
 *
 */
package toolbox_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/viant/toolbox"
)

func TestProcessStruct(t *testing.T) {

	type User struct {
		Name        string    `column:"name"`
		DateOfBirth time.Time `column:"date" dateFormat:"2006-01-02 15:04:05.000000"`
		Id          int       `autogenrated:"true"`
		Other       string    `transient:"true"`
	}

	user := User{Id: 1, Other: "!@#", Name: "foo"}
	var userMap = make(map[string]interface{})
	toolbox.ProcessStruct(&user, func(field reflect.StructField, value interface{}) {
		userMap[field.Name] = value
	})

	assert.Equal(t, 4, len(userMap))
	assert.Equal(t, 1, userMap["Id"])
	assert.Equal(t, "!@#", userMap["Other"])
}

func TestBuildTagMapping(t *testing.T) {

	type User struct {
		Name        string    `column:"name"`
		DateOfBirth time.Time `column:"date" dateFormat:"2006-01-02 15:04:05.000000"`
		Id          int       `autogenrated:"true"`
		Other       string    `transient:"true"`
	}

	tags := []string{"column", "autogenrated"}
	result := toolbox.BuildTagMapping((*User)(nil), "column", "transient", true, true, tags)

	{
		actual := len(result)
		expected := 3
		assert.Equal(t, actual, expected, "Extract mapping count")
	}
	{
		actual, _ := result["name"]["fieldName"]
		expected := "Name"
		assert.Equal(t, actual, expected, "Extract name mapping")
	}

	{
		actual, _ := result["id"]["autogenrated"]
		expected := "true"
		assert.Equal(t, actual, expected, "Extract id flaged as autogenerated")
	}
}
