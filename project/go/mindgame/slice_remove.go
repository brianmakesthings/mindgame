package mindgame

import (
	"fmt"
)



func RemoveFromSlice (indices []int, slice []int) []int {
	fmt.Printf("indices %v\n", indices)
	var result []int
	if len(indices) == 1 {
		return append(slice[:indices[0]], slice[indices[0]+1:]...)
	}
	for i, val := range indices {
		if i ==0 {
			result = append(result, slice[:val]...)
		}else if i > 0 {
			result = append(result, slice[(indices[i-1]+1):val]...)	
		}
		if i == len(indices) -1 {
			result = append(result, slice[val+1:]...)
		}
	}
	fmt.Printf("result %v\n", result)
	return result
}

func RemoveFromStringSlice (indices []int, slice []string) []string {
	fmt.Printf("indices %v\n", indices)
	var result []string
	if len(indices) == 1 {
		return append(slice[:indices[0]], slice[indices[0]+1:]...)
	}
	for i, val := range indices {
		if i ==0 {
			result = append(result, slice[:val]...)
		}else if i > 0 {
			result = append(result, slice[(indices[i-1]+1):val]...)	
		}
		if i == len(indices) -1 {
			result = append(result, slice[val+1:]...)
		}
	}
	fmt.Printf("result %v\n", result)
	return result
}

func RemoveFromPlayerSlice (indices []int, slice []Player) []Player {
	fmt.Printf("indices %v\n", indices)
	var result []Player
	if len(indices) == 1 {
		return append(slice[:indices[0]], slice[indices[0]+1:]...)
	}
	for i, val := range indices {
		if i ==0 {
			result = append(result, slice[:val]...)
		}else if i > 0 {
			result = append(result, slice[(indices[i-1]+1):val]...)	
		}
		if i == len(indices) -1 {
			result = append(result, slice[val+1:]...)
		}
	}
	fmt.Printf("result %v\n", result)
	return result
}