package types

// ================================================================
// ABOUT DISPOSITION MATRICES/VECTORS
//
// Mlrvals can be of type MT_STRING, MT_INT, MT_FLOAT, MT_BOOLEAN, as well as
// MT_ABSENT, MT_VOID, and ERROR.  Thus when we do pairwise operations on them
// (for binary operators) or singly (for unary operators), what we do depends
// on the type.
//
// Mlrval type enums are 0-up integers precisely so that instead of if-elsing
// or switching on the types, we can instead define tables of function pointers
// and jump immediately to the right thing to do for a given type pairing.  For
// example: adding two ints, or an int and a float, or int and boolean (the
// latter being an error).
//
// The next-past-highest mlrval type enum is called MT_DIM and that is the
// dimension of the binary-operator disposition matrices and unary-operator
// disposition vectors.
//
// Note that not every operation uses disposition matrices. If something makes
// sense only for pairs of strings and nothing else, it makes sense for the
// implementing method to return an MT_STRING result if both arguments are
// MT_STRING, or MT_ERROR otherwise.
//
// Naming conventions: since these functions fit into disposition matrices, the
// names are kept quite short. Many are of the form 'plus_f_fi', 'eq_b_xs',
// etc. The conventions are:
//
// * The 'plus_', 'eq_', etc is for the name of the operator.
//
// * For binary operators, things like _f_fi indicate the type of the return
//   value (e.g. 'f') and the types of the two arguments (e.g. 'fi').
//
// * For unary operators, things like _i_i indicate the type of the return
//   value and the type of the argument.
//
// * Type names:
//   's' for string
//   'i' for int64
//   'f' for float64
//   'n' for number return types -- e.g. the auto-overflow from
//       int to float plus_n_ii returns MT_INT if integer-additio overflow
//       didn't happen, or MT_FLOAT if it did.
//   'b' for boolean
//   'x' for don't-care slots, e.g. eq_b_sx for comparing MT_STRING
//       ('s') to anything else ('x').
// ================================================================

// Function-pointer type for zary functions.
type ZaryFunc func() Mlrval

// Function-pointer type for unary-operator disposition vectors.
type UnaryFunc func(*Mlrval) Mlrval

// Helps keystroke-saving for wrapping Go math-library functions
// Examples: cos, sin, etc.
type mathLibUnaryFunc func(float64) float64
type mathLibUnaryFuncWrapper func(*Mlrval, mathLibUnaryFunc) Mlrval

// Function-pointer type for binary-operator disposition matrices.
type BinaryFunc func(*Mlrval, *Mlrval) Mlrval

// Function-pointer type for ternary functions
type TernaryFunc func(*Mlrval, *Mlrval, *Mlrval) Mlrval

// Function-pointer type for variadic functions.
type VariadicFunc func([]*Mlrval) Mlrval

// Function-pointer type for sorting. Returns < 0 if a < b, 0 if a == b, > 0 if a > b.
type ComparatorFunc func(*Mlrval, *Mlrval) int

// ================================================================
// The following are frequently used in disposition matrices for various
// operators and are defined here for re-use. The names are VERY short,
// and all the same length, so that the disposition matrices will look
// reasonable rectangular even after gofmt has been run.

// ----------------------------------------------------------------
// Return error (unary)
func _erro1(ma *Mlrval) Mlrval {
	return MlrvalFromError()
}

// Return absent (unary)
func _absn1(ma *Mlrval) Mlrval {
	return MlrvalFromAbsent()
}

// Return void (unary)
func _void1(ma *Mlrval) Mlrval {
	return MlrvalFromAbsent()
}

// Return argument (unary)
func _1u___(ma *Mlrval) Mlrval {
	return *ma
}

// ----------------------------------------------------------------
// Return error (binary)
func _erro(ma, mb *Mlrval) Mlrval {
	return MlrvalFromError()
}

// Return absent (binary)
func _absn(ma, mb *Mlrval) Mlrval {
	return MlrvalFromAbsent()
}

// Return void (binary)
func _void(ma, mb *Mlrval) Mlrval {
	return MlrvalFromVoid()
}

// Return first argument (binary)
func _1___(ma, mb *Mlrval) Mlrval {
	return *ma
}

// Return second argument (binary)
func _2___(ma, mb *Mlrval) Mlrval {
	return *mb
}

// Return first argument, as string (binary)
func _s1__(ma, mb *Mlrval) Mlrval {
	return MlrvalFromString(ma.String())
}

// Return second argument, as string (binary)
func _s2__(ma, mb *Mlrval) Mlrval {
	return MlrvalFromString(mb.String())
}

// Return integer zero (binary)
func _i0__(ma, mb *Mlrval) Mlrval {
	return MlrvalFromInt64(0)
}

// Return float zero (binary)
func _f0__(ma, mb *Mlrval) Mlrval {
	return MlrvalFromFloat64(0.0)
}
