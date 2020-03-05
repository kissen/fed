package main

func Sum(old, new map[string]interface{}) map[string]interface{} {
	sum := make(map[string]interface{})

	for key, value := range old {
		sum[key] = value
	}

	for key, value := range new {
		sum[key] = value
	}

	return sum
}
