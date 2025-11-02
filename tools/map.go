package tools

// remove duplicate in a map[string]string
func MapRemoveDuplicates(input map[string]string) map[string]string {
	output := make(map[string]string)
	//keep track of keys that have already been added
	keys := make(map[string]bool) 
	//iterate through original map
    for key, value := range input {
        //check if key already exists in keys
        if _, exists := keys[key]; !exists {
            keys[key] = true
            output[key] = value
        }
    }
	return output
}

// get the key with the maximum value in a map[string]int 
func GetKeyWithMaxValue(data map[string]int) string {
	var maxKey string
	maxValue := 0

	for key, value := range data {
		if value > maxValue {
			maxKey = key
			maxValue = value
		}
	}

	return maxKey
}

// isMaxValueForMultipleKeys checks if the maximum value in the given map is set for multiple keys.
// It returns true if more than one key has the maximum value, otherwise, it returns false.
func IsMaxValueForMultipleKeys(m map[string]int) bool {
    var maxValue int
    var countMaxValue int

    // Iterate over the map to find the maximum value
    for _, v := range m {
        if v > maxValue {
            maxValue = v
            countMaxValue = 1
        } else if v == maxValue {
            countMaxValue++
        }
    }

    // Check if more than one key has the maximum value
    return countMaxValue > 1
}