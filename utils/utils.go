package utils

func Ternary[E any](condition bool, trueValue E, falseValue E) E {
	if condition {
		return trueValue
	}

	return falseValue
}
