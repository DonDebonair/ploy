package utils

func Find[E any](s []E, f func(E) bool) *E {
	for _, v := range s {
		if f(v) {
			return &v
		}
	}
	return nil
}
