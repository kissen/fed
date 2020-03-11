package fedutil

func SumMaps(ms ...map[string]interface{}) map[string]interface{} {
	sum := make(map[string]interface{})

	for _, m := range ms {
		for key, value := range m {
			sum[key] = value
		}
	}

	return sum
}
