package funcs

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func CleanNoOp(val string) (string, error) {
	return val, nil
}

func TrimSpaces(val string) (string, error) {
	return strings.TrimSpace(val), nil
}

func CleanTimestamp(val string) (string, error) {
	t, err := time.Parse("01/02/2006 15:04:05", val)
	if err != nil {
		return "", fmt.Errorf("failed to clean timestamp: %v", err)
	}

	return t.Format(time.RFC3339), nil
}

func CleanDurationAsSeconds(duration string) (string, error) {
	if duration == "" {
		return "0", nil
	}

	var durationHandlerMap = map[string]struct {
		unitFriendlyName string
		timeMultiplier   time.Duration
	}{
		"sec": {
			unitFriendlyName: "seconds",
			timeMultiplier:   time.Second,
		},
		"min": {
			unitFriendlyName: "minutes",
			timeMultiplier:   time.Minute,
		},
		"h": {
			unitFriendlyName: "hours",
			timeMultiplier:   time.Hour,
		},
		"d": {
			unitFriendlyName: "days",
			timeMultiplier:   24 * time.Hour,
		},
	}

	result := 0 * time.Second
	segmentRegex := regexp.MustCompile("^([0-9.]+)([a-zA-Z]+)$")
	segments := strings.Split(duration, " ")
	seenUnits := map[string]bool{}
	for _, segment := range segments {
		segmentParts := segmentRegex.FindStringSubmatch(segment)
		if len(segmentParts) != 3 {
			return "", fmt.Errorf("found %d parts, but expected %d parts in segment %q in duration: %s", len(segmentParts)-1, 2, segment, duration)
		}

		unit := segmentParts[2]
		handler, found := durationHandlerMap[unit]
		if !found {
			return "", fmt.Errorf("unexpected unit %q in segment %q: %s", unit, segment, duration)
		}
		if seenUnits[unit] {
			return "", fmt.Errorf("duplicate unit %q found in segment %q: %s", handler.unitFriendlyName, segment, duration)
		}
		seenUnits[unit] = true

		valueStr := segmentParts[1]
		value, err := strconv.ParseInt(valueStr, 10, 16)
		if err != nil {
			return "", fmt.Errorf("invalid value %q in segment %q: %s", valueStr, segment, duration)
		}

		result += time.Duration(value) * handler.timeMultiplier
	}

	return strconv.FormatFloat(result.Seconds(), 'f', -1, 64), nil
}

func RemoveCommas(valueStr string) (string, error) {
	return strings.ReplaceAll(valueStr, ",", ""), nil
}

func RemoveNegativeParensFromCurrency(valueStr string) (string, error) {
	// Remove the currency symbol.
	cleanedStr, _ := strings.CutPrefix(valueStr, "$")

	// Look for parentheses in the string.
	cleanedStr, foundPrefix := strings.CutPrefix(cleanedStr, "(")
	cleanedStr, foundSuffix := strings.CutSuffix(cleanedStr, ")")

	if foundPrefix && foundSuffix {
		return "-" + cleanedStr, nil
	}
	return cleanedStr, nil
}
