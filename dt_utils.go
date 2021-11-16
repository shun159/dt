// SPDX-License-Idenfier: BSD-2-Clause
// Author: Eishun Kondoh <dreamdiagnosis@gmail.com>

package dt

const DEFAULT_OLD_VCLOCK = 86400
const DEFAULT_YOUNG_VCLOCK = 20
const DEFAULT_BIG_VCLOCK = 50
const DEFAULT_SMALL_VCLOCK = 50

func GetOldVclock(prop map[string]int) int {
	return GetParams(prop, "old_vclock", DEFAULT_OLD_VCLOCK)
}

func GetYoungVclock(prop map[string]int) int {
	return GetParams(prop, "young_vclock", DEFAULT_YOUNG_VCLOCK)
}

func GetBigVclock(prop map[string]int) int {
	return GetParams(prop, "big_vclock", DEFAULT_BIG_VCLOCK)
}

func GetSmallVclock(prop map[string]int) int {
	return GetParams(prop, "small_vclock", DEFAULT_SMALL_VCLOCK)
}

func GetParams(d map[string]int, key string, default_value int) int {
	if val, ok := d[key]; ok {
		return val
	} else {
		return default_value
	}
}
