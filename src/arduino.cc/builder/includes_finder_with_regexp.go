/*
 * This file is part of Arduino Builder.
 *
 * Arduino Builder is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 2 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to the Free Software
 * Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA  02110-1301  USA
 *
 * As a special exception, you may use this file as part of a free software
 * library without restriction.  Specifically, if other files instantiate
 * templates or use macros or inline functions from this file, or you compile
 * this file and link it with other files to produce an executable, this
 * file does not by itself cause the resulting executable to be covered by
 * the GNU General Public License.  This exception does not however
 * invalidate any other reasons why the executable file might be covered by
 * the GNU General Public License.
 *
 * Copyright 2015 Arduino LLC (http://www.arduino.cc/)
 */

package builder

import (
	"arduino.cc/builder/constants"
	"arduino.cc/builder/utils"
	"regexp"
	"strings"
)

var INCLUDE_REGEXP = regexp.MustCompile("(?ms)^\\s*#[ \t]*include\\s*[<\"](\\S+)[\">]")

type IncludesFinderWithRegExp struct {
	ContextField string
}

func (s *IncludesFinderWithRegExp) Run(context map[string]interface{}) error {
	source := context[s.ContextField].(string)

	matches := INCLUDE_REGEXP.FindAllStringSubmatch(source, -1)
	includes := []string{}
	for _, match := range matches {
		includes = append(includes, strings.TrimSpace(match[1]))
	}

	if len(includes) == 0 {
		include := findIncludesForOldCompilers(source)
		if include != "" {
			includes = append(includes, include)
		}
	}

	context[constants.CTX_INCLUDES_JUST_FOUND] = includes

	if !utils.MapHas(context, constants.CTX_INCLUDES) {
		context[constants.CTX_INCLUDES] = includes
		return nil
	}

	context[constants.CTX_INCLUDES] = utils.AddStringsToStringsSet(context[constants.CTX_INCLUDES].([]string), includes)

	return nil
}

func findIncludesForOldCompilers(source string) string {
	firstLine := strings.Split(source, "\n")[0]
	splittedLine := strings.Split(firstLine, ":")
	for i, _ := range splittedLine {
		if strings.Contains(splittedLine[i], "fatal error") {
			return strings.TrimSpace(splittedLine[i+1])
		}
	}
	return ""
}
