package tools

import (
	"strconv"
)



// return a slice wich each value unique
func RemoveDuplicateValuesString(stringSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	if len(stringSlice) > 0 {
		for _, entry := range stringSlice {
			if _, value := keys[entry]; !value {
				keys[entry] = true
				list = append(list, entry)
			}
		}
	}
	return list
}

// return a slice wich each value unique
func RemoveDuplicateValuesInt(stringSlice []int) []int {
	keys := make(map[int]bool)
	list := []int{}
	if len(stringSlice) > 0 {
		for _, entry := range stringSlice {
			if _, value := keys[entry]; !value {
				keys[entry] = true
				list = append(list, entry)
			}
		}
	}
	return list
}

// Find takes a slice and looks for an element in it. If found it will
// return it's key, otherwise it will return -1 and a bool of false.
// Find takes a slice and looks for an element in it. If found it will
// return its key, otherwise it will return -1 and a bool of false.
func FindInSliceString(slice []string, val string) (int, bool) {
	if len(slice) > 0 && len(val) > 0 {
		for i, item := range slice {
			if item == val {
				return i, true
			}
		}
	}
	return -1, false
}

// Find takes a slice and looks for an element in it. If found it will
// return it's key, otherwise it will return -1 and a bool of false.
// Find takes a slice and looks for an element in it. If found it will
// return its key, otherwise it will return -1 and a bool of false.
func FindInSliceInt(slice []int, val int) (int, bool) {
	if len(slice) > 0 {
		for i, item := range slice {
			if item == val {
				return i, true
			}
		}
	}
	return -1, false
}

// reverse slice
func ReverseSlice(a []string) []string {
	for i := len(a)/2 - 1; i >= 0; i-- {
		opp := len(a) - 1 - i
		a[i], a[opp] = a[opp], a[i]
	}

	return a
}

// fom a slice a string we get a string of int
func ConvertSliceOfStringToSliceOfInt(ss []string) []int {
	si := make([]int, len(ss))
	for i, s := range ss {
		si[i], _ = strconv.Atoi(s)
	}
	return si
}

// SliceStringCountDict returns a map with the count of each string in the slice
func SliceStringCountDict(slice []string) map[string]int {
	dict := make(map[string]int)
	for _, val := range slice {
		dict[val]++
	}
	return dict
}

func SliceIntCountDict(slice []int) map[int]int {
	//Create a   dictionary of values for each element
	var dict = make(map[int]int)
	for _, val := range slice {
		dict[val]++
	}

	return dict
}

// compairs two slice, and check if fist slice contains value in slice #2
func ContainsAny(slice1, slice2 []string) bool {
	set := make(map[string]struct{})

	// Create a set of elements from slice2 for faster lookup
	for _, val := range slice2 {
		set[val] = struct{}{}
	}

	// Check if any element from slice1 exists in the set
	for _, val := range slice1 {
		if _, exists := set[val]; exists {
			return true
		}
	}

	return false
}
