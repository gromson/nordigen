package utils

func MergeMapsOfArrays[T mergeableMapOfArraysConstraint](maps ...T) (r T) {
	r = T{}

	for _, m := range maps {
		for k, vv := range m {
			for _, v := range vv {
				r.Add(k, v)
			}
		}
	}

	return
}

type mergeableMapOfArraysConstraint interface {
	~map[string][]string
	Add(key, value string)
}
