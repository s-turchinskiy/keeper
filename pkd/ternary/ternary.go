package ternary

func Bool(condition, expressionTrue bool, expressionFalse *bool) bool {
	if condition {
		return expressionTrue
	} else {
		return *expressionFalse
	}
}
