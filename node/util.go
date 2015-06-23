package node

// diff returns all of the strings that exist in s1 but do not exist in s2
func diffSliceStr(s1, s2 []string) []string {
	var vs []string
	var found bool

	for _, v1 := range s1 {
		found = false
		for _, v2 := range s2 {
			if v1 == v2 {
				found = true
				break
			}
		}

		if !found {
			vs = append(vs, v1)
		}
	}

	return vs
}
