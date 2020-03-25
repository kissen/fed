package main

// Given maps ms, return a new map that contains the sum of all of
// those maps.
//
// Sum in this context means a map that contains all keys found in
// ms[i] for any i. If maps ms[i], ms[j], i < j, share some key k, the
// returned map will contain ms[j][k], that is it will contain the
// value from the last argument that contained a mapping for key k.
func sumMaps(ms ...map[string]interface{}) map[string]interface{} {
	sum := make(map[string]interface{})

	for _, m := range ms {
		for key, value := range m {
			sum[key] = value
		}
	}

	return sum
}
