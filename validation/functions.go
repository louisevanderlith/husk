package validation

func happyValidation(obj interface{}, meta tagMeta) (bool, []string) {
	var issues []string

	if meta.Required && obj == nil {
		issues = append(issues, getIDMessage(meta.PropName))
	}

	isValid := len(issues) < 1

	return isValid, issues
}

func noValidation(obj interface{}, meta tagMeta) (bool, []string) {
	return false, []string{getIncorrectType(meta.PropName, obj, "Unknown")}
}

func stringValidation(obj interface{}, meta tagMeta) (bool, []string) {
	var issues []string
	val, ok := obj.(string)

	if ok {
		if meta.Required && val == "" {
			issues = append(issues, getEmptyMessage(meta.PropName))
		}

		if meta.Size > 0 && len(val) > meta.Size {
			issues = append(issues, getShortMessage(meta.PropName, meta.Size))
		}
	} else {
		issues = append(issues)
	}

	isValid := len(issues) < 1

	return isValid, issues
}

func int8Validation(obj interface{}, meta tagMeta) (bool, []string) {
	var issues []string
	val, ok := obj.(int8)

	if ok {
		if meta.Required && val < 1 {
			issues = append(issues, getIDMessage(meta.PropName))
		}
	} else {
		issues = append(issues, getIncorrectType(meta.PropName, obj, "int8"))
	}

	isValid := len(issues) < 1

	return isValid, issues
}

func int16Validation(obj interface{}, meta tagMeta) (bool, []string) {
	var issues []string
	val, ok := obj.(int16)

	if ok {
		if meta.Required && val < 1 {
			issues = append(issues, getIDMessage(meta.PropName))
		}
	} else {
		issues = append(issues, getIncorrectType(meta.PropName, obj, "int16"))
	}

	isValid := len(issues) < 1

	return isValid, issues
}

func intValidation(obj interface{}, meta tagMeta) (bool, []string) {
	var issues []string
	val, ok := obj.(int)

	if ok {
		if meta.Required && val < 1 {
			issues = append(issues, getIDMessage(meta.PropName))
		}
	} else {
		issues = append(issues, getIncorrectType(meta.PropName, obj, "int"))
	}

	isValid := len(issues) < 1

	return isValid, issues
}

func int64Validation(obj interface{}, meta tagMeta) (bool, []string) {
	var issues []string
	val, ok := obj.(int64)

	if ok {
		if meta.Required && val < 1 {
			issues = append(issues, getIDMessage(meta.PropName))
		}
	} else {
		issues = append(issues, getIncorrectType(meta.PropName, obj, "int64"))
	}

	isValid := len(issues) < 1

	return isValid, issues
}

func uint8Validation(obj interface{}, meta tagMeta) (bool, []string) {
	var issues []string
	val, ok := obj.(uint8)

	if ok {
		if meta.Required && val < 1 {
			issues = append(issues, getIDMessage(meta.PropName))
		}
	} else {
		issues = append(issues, getIncorrectType(meta.PropName, obj, "uint8"))
	}

	isValid := len(issues) < 1

	return isValid, issues
}

func uint16Validation(obj interface{}, meta tagMeta) (bool, []string) {
	var issues []string
	val, ok := obj.(uint16)

	if ok {
		if meta.Required && val < 1 {
			issues = append(issues, getIDMessage(meta.PropName))
		}
	} else {
		issues = append(issues, getIncorrectType(meta.PropName, obj, "uint16"))
	}

	isValid := len(issues) < 1

	return isValid, issues
}

func uintValidation(obj interface{}, meta tagMeta) (bool, []string) {
	var issues []string
	val, ok := obj.(uint)

	if ok {
		if meta.Required && val < 1 {
			issues = append(issues, getIDMessage(meta.PropName))
		}
	} else {
		issues = append(issues, getIncorrectType(meta.PropName, obj, "uint"))
	}

	isValid := len(issues) < 1

	return isValid, issues
}

func uint64Validation(obj interface{}, meta tagMeta) (bool, []string) {
	var issues []string
	val, ok := obj.(uint64)

	if ok {
		if meta.Required && val < 1 {
			issues = append(issues, getIDMessage(meta.PropName))
		}
	} else {
		issues = append(issues, getIncorrectType(meta.PropName, obj, "uint64"))
	}

	isValid := len(issues) < 1

	return isValid, issues
}

func floatValidation(obj interface{}, meta tagMeta) (bool, []string) {
	var issues []string
	val, ok := obj.(float32)

	if ok {
		if meta.Required && val < 1 {
			issues = append(issues, getIDMessage(meta.PropName))
		}
	} else {
		issues = append(issues, getIncorrectType(meta.PropName, obj, "float32"))
	}

	isValid := len(issues) < 1

	return isValid, issues
}

func float64Validation(obj interface{}, meta tagMeta) (bool, []string) {
	var issues []string
	val, ok := obj.(float64)

	if ok {
		if meta.Required && val < 1 {
			issues = append(issues, getIDMessage(meta.PropName))
		}
	} else {
		issues = append(issues, getIncorrectType(meta.PropName, obj, "float64"))
	}

	isValid := len(issues) < 1

	return isValid, issues
}

func structValidation(obj interface{}, meta tagMeta) (bool, []string) {
	var issues []string

	if meta.Required && obj == nil {
		issues = append(issues, getRelationMessage(meta.PropName))
	}

	isValid := len(issues) < 1

	return isValid, issues
}

func pointerValidation(obj interface{}, meta tagMeta) (bool, []string) {
	return true, nil
}

func interfaceValidation(obj interface{}, meta tagMeta) (bool, []string) {
	var issues []string

	if meta.Required && obj == nil {
		issues = append(issues, getIDMessage(meta.PropName))
	}

	return len(issues) < 1, issues
}
