package types

func ValidateTextAttributes(attr map[string]interface{}) bool {
	if _, ok := attr["bx"]; !ok {
		return false
	}

	if _, ok := attr["by"]; !ok {
		return false
	}

	if _, ok := attr["value"]; !ok {
		return false
	}

	if _, ok := attr["textColor"]; !ok {
		return false
	}

	if _, ok := attr["fontWidth"]; !ok {
		return false
	}

	if _, ok := attr["font"]; !ok {
		return false
	}

	if _, ok := attr["width"]; !ok {
		return false
	}

	if _, ok := attr["height"]; !ok {
		return false
	}

	if _, ok := attr["strokeWidth"]; !ok {
		return false
	}

	if _, ok := attr["strokeColor"]; !ok {
		return false
	}

	if _, ok := attr["fillColor"]; !ok {
		return false
	}

	return true
}
func ValidateCircleAttributes(attr map[string]interface{}) bool {
	if _, ok := attr["cx"]; !ok {
		return false
	}

	if _, ok := attr["cy"]; !ok {
		return false
	}

	if _, ok := attr["radius"]; !ok {
		return false
	}

	if _, ok := attr["strokeWidth"]; !ok {
		return false
	}

	if _, ok := attr["strokeColor"]; !ok {
		return false
	}

	if _, ok := attr["fillColor"]; !ok {
		return false
	}

	return true
}

func ValidateRectangleAttributes(attr map[string]interface{}) bool {
	if _, ok := attr["x"]; !ok {
		return false
	}

	if _, ok := attr["y"]; !ok {
		return false
	}

	if _, ok := attr["width"]; !ok {
		return false
	}

	if _, ok := attr["height"]; !ok {
		return false
	}

	if _, ok := attr["strokeWidth"]; !ok {
		return false
	}

	if _, ok := attr["strokeColor"]; !ok {
		return false
	}

	if _, ok := attr["fillColor"]; !ok {
		return false
	}

	return true
}
func ValidatePenAttributes(attr map[string]interface{}) bool {
	// Pen requires points, color, strokeWidth
	if _, ok := attr["points"]; !ok {
		return false
	}
	// Points must be a list/slice
	if _, ok := attr["points"].([]interface{}); !ok {
		return false
	}
	if _, ok := attr["color"]; !ok {
		if _, ok2 := attr["strokeColor"]; !ok2 {
			return false
		}
	}
	if _, ok := attr["strokeWidth"]; !ok {
		return false
	}
	return true
}

func ValidateLineAttributes(attr map[string]interface{}) bool {
	// Line/Arrow behaves like a bounding box in your frontend logic (x,y,width,height)
	if _, ok := attr["x"]; !ok {
		return false
	}
	if _, ok := attr["y"]; !ok {
		return false
	}
	if _, ok := attr["width"]; !ok {
		return false
	}
	if _, ok := attr["height"]; !ok {
		return false
	}
	if _, ok := attr["strokeWidth"]; !ok {
		return false
	}
	if _, ok := attr["strokeColor"]; !ok {
		return false
	}
	return true
}

func ValidateCreateMessage(msg map[string]interface{}) bool {
	if _, ok := msg["objectType"]; !ok {
		return false
	}

	if _, ok := msg["attributes"]; !ok {
		return false
	}

	if _, ok := msg["slideId"]; !ok {
		return false
	}

	if _, ok := msg["objectId"]; !ok {
		return false
	}

	return true
}

func ValidateDeleteMessage(msg map[string]interface{}) bool {
	if _, ok := msg["slideId"]; !ok {
		return false
	}

	if _, ok := msg["objectId"]; !ok {
		return false
	}

	if _, ok := msg["objectType"]; !ok {
		return false
	}

	return true
}

func ValidateUpdateMessage(msg map[string]interface{}) bool {
	if _, ok := msg["slideId"]; !ok {
		return false
	}

	if _, ok := msg["objectId"]; !ok {
		return false
	}

	if _, ok := msg["objectType"]; !ok {
		return false
	}

	return true
}

func ValidateCursorMoveMessage(msg map[string]interface{}) bool {
	if _, ok := msg["slideId"]; !ok {
		return false
	}

	if _, ok := msg["newCursorLocation"]; !ok {
		return false
	}

	return true

}

func ValidateSelectMessage(msg map[string]interface{}) bool {
	if _, ok := msg["slideId"]; !ok {
		return false
	}

	if _, ok := msg["objectId"]; !ok {
		return false
	}

	return true
}

func ValidateAddSlideMessage(msg map[string]interface{}) bool {
	if _, ok := msg["slideId"]; !ok {
		return false
	}

	return true
}

func ValidateRemoveSlideMessage(msg map[string]interface{}) bool {
	if _, ok := msg["slideId"]; !ok {
		return false
	}

	return true
}
