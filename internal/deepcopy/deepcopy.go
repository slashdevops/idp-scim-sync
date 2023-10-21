package deepcopy

// SliceOfPointers creates a new slice of pointers from the given slice.
// The new slice will have the same length and capacity as the given slice.
// and the values will be copied.
func SliceOfPointers[T any](s []*T) []*T {
	if s == nil {
		return nil
	}

	newSlice := make([]*T, 0, len(s))

	for _, v := range s {
		if v == nil {
			newSlice = append(newSlice, nil)
			continue
		}

		// Create a new struct to hold the copied data
		newStruct := new(T)
		*newStruct = *v

		// Append the new struct to the new slice
		newSlice = append(newSlice, newStruct)
	}

	return newSlice
}
