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
package toolbox

import (
	"reflect"
)

//Iterator represents generic iterator.
type Iterator interface {

	//HasNext returns true if iterator has next element.
	HasNext() bool

	//Next sets item pointer with next element.
	Next(itemPointer interface{})
}



type sliceIterator struct {
	sliceValue reflect.Value
	index      int
}

func (i *sliceIterator) HasNext() bool {
	return i.index < i.sliceValue.Len()
}

func (i *sliceIterator) Next(itemPointer interface{}) {
	value := i.sliceValue.Index(i.index)
	i.index++
	itemPointerValue := reflect.ValueOf(itemPointer)
	itemPointerValue.Elem().Set(value)
}



type stringSliceIterator struct {
	sliceValue []string
	index      int
}

func (i *stringSliceIterator) HasNext() bool {
	return i.index < len(i.sliceValue)
}


func (i *stringSliceIterator) Next(itemPointer interface{}) {
	value := i.sliceValue[i.index]
	i.index++
	if stringPointer, ok  := itemPointer.(*string); ok {
		*stringPointer = value
		return
	}
	interfacePointer:= itemPointer.(*interface{})
	*interfacePointer = value
}


type interfaceSliceIterator struct {
	sliceValue []interface{}
	index      int
}

func (i *interfaceSliceIterator) HasNext() bool {
	return i.index < len(i.sliceValue)
}

func (i *interfaceSliceIterator) Next(itemPointer interface{}) {
	value := i.sliceValue[i.index]
	i.index++
	itemPointerValue := reflect.ValueOf(itemPointer)
	if value != nil {
		itemPointerValue.Elem().Set(reflect.ValueOf(value))
	} else {
		itemPointerValue.Elem().Set(reflect.Zero(reflect.TypeOf(itemPointer).Elem()))

	}
}


//NewSliceIterator creates a new slice iterator.
func NewSliceIterator(slice interface{}) Iterator {
	if aSlice, ok := slice.([]interface{});ok {
		return &interfaceSliceIterator{aSlice, 0}
	}
	if aSlice, ok := slice.([]string);ok {
		return &stringSliceIterator{aSlice, 0}
	}
	sliceValue := DiscoverValueByKind(reflect.ValueOf(slice), reflect.Slice)
	return &sliceIterator{sliceValue: sliceValue}
}
