package types

import (
	"miller/lib"
)

// ================================================================
func MlrvalTypeof(ma *Mlrval) Mlrval {
	return MlrvalFromString(ma.GetTypeName())
}

// ----------------------------------------------------------------
func string_to_int(ma *Mlrval) Mlrval {
	i, ok := lib.TryInt64FromString(ma.printrep)
	if ok {
		return MlrvalFromInt64(i)
	} else {
		return MlrvalFromError()
	}
}

func float_to_int(ma *Mlrval) Mlrval {
	return MlrvalFromInt64(int64(ma.floatval))
}

func bool_to_int(ma *Mlrval) Mlrval {
	if ma.boolval == true {
		return MlrvalFromInt64(1)
	} else {
		return MlrvalFromInt64(0)
	}
}

var to_int_dispositions = [MT_DIM]UnaryFunc{
	/*ERROR  */ _erro1,
	/*ABSENT */ _absn1,
	/*VOID   */ _void1,
	/*STRING */ string_to_int,
	/*INT    */ _1u___,
	/*FLOAT  */ float_to_int,
	/*BOOL   */ bool_to_int,
	/*ARRAY  */ _erro1,
	/*MAP    */ _erro1,
}

func MlrvalToInt(ma *Mlrval) Mlrval {
	return to_int_dispositions[ma.mvtype](ma)
}

// ----------------------------------------------------------------
func string_to_float(ma *Mlrval) Mlrval {
	f, ok := lib.TryFloat64FromString(ma.printrep)
	if ok {
		return MlrvalFromFloat64(f)
	} else {
		return MlrvalFromError()
	}
}

func int_to_float(ma *Mlrval) Mlrval {
	return MlrvalFromFloat64(float64(ma.intval))
}

func bool_to_float(ma *Mlrval) Mlrval {
	if ma.boolval == true {
		return MlrvalFromFloat64(1.0)
	} else {
		return MlrvalFromFloat64(0.0)
	}
}

var to_float_dispositions = [MT_DIM]UnaryFunc{
	/*ERROR  */ _erro1,
	/*ABSENT */ _absn1,
	/*VOID   */ _void1,
	/*STRING */ string_to_float,
	/*INT    */ int_to_float,
	/*FLOAT  */ _1u___,
	/*BOOL   */ bool_to_float,
	/*ARRAY  */ _erro1,
	/*MAP    */ _erro1,
}

func MlrvalToFloat(ma *Mlrval) Mlrval {
	return to_float_dispositions[ma.mvtype](ma)
}

// ----------------------------------------------------------------
func string_to_boolean(ma *Mlrval) Mlrval {
	b, ok := lib.TryBoolFromBoolString(ma.printrep)
	if ok {
		return MlrvalFromBool(b)
	} else {
		return MlrvalFromError()
	}
}

func int_to_bool(ma *Mlrval) Mlrval {
	return MlrvalFromBool(ma.intval != 0)
}

func float_to_bool(ma *Mlrval) Mlrval {
	return MlrvalFromBool(ma.floatval != 0.0)
}

var to_boolean_dispositions = [MT_DIM]UnaryFunc{
	/*ERROR  */ _erro1,
	/*ABSENT */ _absn1,
	/*VOID   */ _void1,
	/*STRING */ string_to_boolean,
	/*INT    */ int_to_bool,
	/*FLOAT  */ float_to_bool,
	/*BOOL   */ _1u___,
	/*ARRAY  */ _erro1,
	/*MAP    */ _erro1,
}

func MlrvalToBoolean(ma *Mlrval) Mlrval {
	return to_boolean_dispositions[ma.mvtype](ma)
}