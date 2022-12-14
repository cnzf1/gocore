/*
 * @Author: cnzf1
 * @Date: 2022-09-13 20:32:00
 * @LastEditors: cnzf1
 * @LastEditTime: 2022-09-14 10:33:27
 * @Description:
 */
package host

import (
	"strconv"
	"strings"

	"github.com/cnzf1/gocore/filex"
)

func GetBtime() int64 {
	file := "/proc/stat"
	lines, err := filex.ReadLinesOffsetN(file, 0, -1)
	if err != nil {
		return 0
	}

	for _, line := range lines {
		if strings.HasPrefix(line, "btime") {
			f := strings.Fields(line)
			if len(f) != 2 {
				return 0
			}

			b, err := strconv.ParseInt(f[1], 10, 64)
			if err != nil {
				return 0
			}

			t := int64(b)
			return t
		}
	}

	return 0
}
