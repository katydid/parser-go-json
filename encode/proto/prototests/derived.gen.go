// Code generated by goderive DO NOT EDIT.

package prototests

import (
	"bytes"
)

func deriveEqualSimple(this, that *Simple) bool {
	return (this == nil && that == nil) ||
		this != nil && that != nil &&
			((this.Field1 == nil && that.Field1 == nil) || (this.Field1 != nil && that.Field1 != nil && *(this.Field1) == *(that.Field1))) &&
			((this.Field2 == nil && that.Field2 == nil) || (this.Field2 != nil && that.Field2 != nil && *(this.Field2) == *(that.Field2))) &&
			((this.Field3 == nil && that.Field3 == nil) || (this.Field3 != nil && that.Field3 != nil && *(this.Field3) == *(that.Field3))) &&
			((this.Field4 == nil && that.Field4 == nil) || (this.Field4 != nil && that.Field4 != nil && *(this.Field4) == *(that.Field4))) &&
			((this.Field5 == nil && that.Field5 == nil) || (this.Field5 != nil && that.Field5 != nil && *(this.Field5) == *(that.Field5))) &&
			((this.Field6 == nil && that.Field6 == nil) || (this.Field6 != nil && that.Field6 != nil && *(this.Field6) == *(that.Field6))) &&
			((this.Field7 == nil && that.Field7 == nil) || (this.Field7 != nil && that.Field7 != nil && *(this.Field7) == *(that.Field7))) &&
			((this.Field8 == nil && that.Field8 == nil) || (this.Field8 != nil && that.Field8 != nil && *(this.Field8) == *(that.Field8))) &&
			((this.Field9 == nil && that.Field9 == nil) || (this.Field9 != nil && that.Field9 != nil && *(this.Field9) == *(that.Field9))) &&
			((this.Field10 == nil && that.Field10 == nil) || (this.Field10 != nil && that.Field10 != nil && *(this.Field10) == *(that.Field10))) &&
			((this.Field11 == nil && that.Field11 == nil) || (this.Field11 != nil && that.Field11 != nil && *(this.Field11) == *(that.Field11))) &&
			((this.Field12 == nil && that.Field12 == nil) || (this.Field12 != nil && that.Field12 != nil && *(this.Field12) == *(that.Field12))) &&
			((this.Field13 == nil && that.Field13 == nil) || (this.Field13 != nil && that.Field13 != nil && *(this.Field13) == *(that.Field13))) &&
			((this.Field14 == nil && that.Field14 == nil) || (this.Field14 != nil && that.Field14 != nil && *(this.Field14) == *(that.Field14))) &&
			bytes.Equal(this.Field15, that.Field15) &&
			deriveEqualSliceOffloat64(this.Fields1, that.Fields1) &&
			deriveEqualSliceOffloat32(this.Fields2, that.Fields2) &&
			deriveEqualSliceOfint32(this.Fields3, that.Fields3) &&
			deriveEqualSliceOfint64(this.Fields4, that.Fields4) &&
			deriveEqualSliceOfuint32(this.Fields5, that.Fields5) &&
			deriveEqualSliceOfuint64(this.Fields6, that.Fields6) &&
			deriveEqualSliceOfint32(this.Fields7, that.Fields7) &&
			deriveEqualSliceOfint64(this.Fields8, that.Fields8) &&
			deriveEqualSliceOfuint32(this.Fields9, that.Fields9) &&
			deriveEqualSliceOfint32(this.Fields10, that.Fields10) &&
			deriveEqualSliceOfuint64(this.Fields11, that.Fields11) &&
			deriveEqualSliceOfint64(this.Fields12, that.Fields12) &&
			deriveEqualSliceOfbool(this.Fields13, that.Fields13) &&
			deriveEqualSliceOfstring(this.Fields14, that.Fields14) &&
			deriveEqualSliceOfSliceOfbyte(this.Fields15, that.Fields15) &&
			bytes.Equal(this.XXX_unrecognized, that.XXX_unrecognized)
}

func deriveEqualNested(this, that *Nested) bool {
	return (this == nil && that == nil) ||
		this != nil && that != nil &&
			this.One.Equal(that.One) &&
			deriveEqualSliceOfPtrToSimple(this.Many, that.Many) &&
			bytes.Equal(this.XXX_unrecognized, that.XXX_unrecognized)
}

func deriveEqualSliceOffloat64(this, that []float64) bool {
	if this == nil || that == nil {
		return this == nil && that == nil
	}
	if len(this) != len(that) {
		return false
	}
	for i := 0; i < len(this); i++ {
		if !(this[i] == that[i]) {
			return false
		}
	}
	return true
}

func deriveEqualSliceOffloat32(this, that []float32) bool {
	if this == nil || that == nil {
		return this == nil && that == nil
	}
	if len(this) != len(that) {
		return false
	}
	for i := 0; i < len(this); i++ {
		if !(this[i] == that[i]) {
			return false
		}
	}
	return true
}

func deriveEqualSliceOfint32(this, that []int32) bool {
	if this == nil || that == nil {
		return this == nil && that == nil
	}
	if len(this) != len(that) {
		return false
	}
	for i := 0; i < len(this); i++ {
		if !(this[i] == that[i]) {
			return false
		}
	}
	return true
}

func deriveEqualSliceOfint64(this, that []int64) bool {
	if this == nil || that == nil {
		return this == nil && that == nil
	}
	if len(this) != len(that) {
		return false
	}
	for i := 0; i < len(this); i++ {
		if !(this[i] == that[i]) {
			return false
		}
	}
	return true
}

func deriveEqualSliceOfuint32(this, that []uint32) bool {
	if this == nil || that == nil {
		return this == nil && that == nil
	}
	if len(this) != len(that) {
		return false
	}
	for i := 0; i < len(this); i++ {
		if !(this[i] == that[i]) {
			return false
		}
	}
	return true
}

func deriveEqualSliceOfuint64(this, that []uint64) bool {
	if this == nil || that == nil {
		return this == nil && that == nil
	}
	if len(this) != len(that) {
		return false
	}
	for i := 0; i < len(this); i++ {
		if !(this[i] == that[i]) {
			return false
		}
	}
	return true
}

func deriveEqualSliceOfbool(this, that []bool) bool {
	if this == nil || that == nil {
		return this == nil && that == nil
	}
	if len(this) != len(that) {
		return false
	}
	for i := 0; i < len(this); i++ {
		if !(this[i] == that[i]) {
			return false
		}
	}
	return true
}

func deriveEqualSliceOfstring(this, that []string) bool {
	if this == nil || that == nil {
		return this == nil && that == nil
	}
	if len(this) != len(that) {
		return false
	}
	for i := 0; i < len(this); i++ {
		if !(this[i] == that[i]) {
			return false
		}
	}
	return true
}

func deriveEqualSliceOfSliceOfbyte(this, that [][]byte) bool {
	if this == nil || that == nil {
		return this == nil && that == nil
	}
	if len(this) != len(that) {
		return false
	}
	for i := 0; i < len(this); i++ {
		if !(bytes.Equal(this[i], that[i])) {
			return false
		}
	}
	return true
}

func deriveEqualSliceOfPtrToSimple(this, that []*Simple) bool {
	if this == nil || that == nil {
		return this == nil && that == nil
	}
	if len(this) != len(that) {
		return false
	}
	for i := 0; i < len(this); i++ {
		if !(this[i].Equal(that[i])) {
			return false
		}
	}
	return true
}
