package hbn

func findAvrgLatency(list []int64) int64 {
		res := int64(0)
		for _, i := range list[1:] {
				res += i
		}
		res /= int64(len(list[1:]))
		return res
}

func MinMax(values []int64) (int64, int64) {
		min_value := values[0]
		max_value := values[0]
		for _, i := range values {
				if i < min_value {
						min_value = i
				} else if i > max_value {
						max_value = i
				}
		}

		return min_value, max_value
}

func convert_nanoseconds_to_seconds(value int64) float32 {
		res := float32(value)/float32(1e9)
		return res
}
