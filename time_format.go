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

import "strings"

//DateFormatKeyword constant 'dateFormat' key
var DateFormatKeyword = "dateFormat"

//DateLayoutKeyword constant 'dateLayout' key
var DateLayoutKeyword = "dateLayout"

//DateFormatToLayout converts java date format https://docs.oracle.com/javase/6/docs/api/java/text/SimpleDateFormat.html#rfc822timezone into go date layout
func DateFormatToLayout(dateFormat string) string {

	dateFormat = strings.Replace(dateFormat, "MMMM", "January", 1)
	dateFormat = strings.Replace(dateFormat, "MMM", "Jan", 1)
	dateFormat = strings.Replace(dateFormat, "MM", "1", 1)
	dateFormat = strings.Replace(dateFormat, "M", "1", 1)

	dateFormat = strings.Replace(dateFormat, "a", "pm", 1)
	dateFormat = strings.Replace(dateFormat, "aa", "PM", 1)

	dateFormat = strings.Replace(dateFormat, "dd", "02", 1)
	dateFormat = strings.Replace(dateFormat, "d", "3", 1)

	dateFormat = strings.Replace(dateFormat, "HH", "15", 1)
	dateFormat = strings.Replace(dateFormat, "H", "15", 1)

	dateFormat = strings.Replace(dateFormat, "hh", "03", 1)
	dateFormat = strings.Replace(dateFormat, "h", "3", 1)

	dateFormat = strings.Replace(dateFormat, "mm", "04", 1)
	dateFormat = strings.Replace(dateFormat, "m", "4", 1)

	dateFormat = strings.Replace(dateFormat, "ss", "05", 1)
	dateFormat = strings.Replace(dateFormat, "s", "5", 1)

	dateFormat = strings.Replace(dateFormat, "yyyy", "2006", 1)
	dateFormat = strings.Replace(dateFormat, "yy", "06", 1)
	dateFormat = strings.Replace(dateFormat, "y", "06", 1)

	dateFormat = strings.Replace(dateFormat, "z", "MST", 1)
	dateFormat = strings.Replace(dateFormat, "zzzz", "Z0700", 1)
	dateFormat = strings.Replace(dateFormat, "zzzz", "Z07:00", 1)
	dateFormat = strings.Replace(dateFormat, "Z", "-07", 1)
	dateFormat = strings.Replace(dateFormat, "EEEE", "Monday", 1)
	dateFormat = strings.Replace(dateFormat, "E", "Mon", 1)

	dateFormat = strings.Replace(dateFormat, "SSS", "000", 1)
	return dateFormat
}

//GetTimeLayout returns time laout from passed in map, first it check if DateLayoutKeyword is defined is so it returns it, otherwise it check DateFormatKeyword and if exists converts it to  dateLayout
//If neithers keys exists it panics, please use HasTimeLayout to avoid panic
func GetTimeLayout(settings map[string]string) string {
	if value, found := settings[DateLayoutKeyword]; found {
		return value
	}
	if value, found := settings[DateFormatKeyword]; found {
		return DateFormatToLayout(value)
	}
	panic("Date format or date layout is not defined")
}

//HasTimeLayout checks if dateLayout can be taken from the passed in setting map
func HasTimeLayout(settings map[string]string) bool {
	if _, found := settings[DateLayoutKeyword]; found {
		return true
	}
	if _, found := settings[DateFormatKeyword]; found {
		return true
	}
	return false
}
