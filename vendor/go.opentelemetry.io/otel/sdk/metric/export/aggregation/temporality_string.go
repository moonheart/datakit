// Code generated by "stringer -type=Temporality"; DO NOT EDIT.

package aggregation

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[CumulativeTemporality-1]
	_ = x[DeltaTemporality-2]
}

const _Temporality_name = "CumulativeTemporalityDeltaTemporality"

var _Temporality_index = [...]uint8{0, 21, 37}

func (i Temporality) String() string {
	i -= 1
	if i >= Temporality(len(_Temporality_index)-1) {
		return "Temporality(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _Temporality_name[_Temporality_index[i]:_Temporality_index[i+1]]
}
