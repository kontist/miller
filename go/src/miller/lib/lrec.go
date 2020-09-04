// ================================================================
// ORDERED MAP FROM STRING TO MLRVAL
//
// It is an optionally hashless implementation of insertion-ordered key-value
// pairs for Miller's fundamental record data structure. It's also an
// ordered-map data structure, suitable for Miller JSON decode/encode.
//
// ----------------------------------------------------------------
// DESIGN
//
// * It keeps a doubly-linked list of key-value pairs.
//
// * By default, no hash functions are computed when the map is written to or
//   read from.
//
// * Gets are implemented by sequential scan through the list: given a key,
//   the key-value pairs are scanned through until a match is (or is not) found.
//
// * Performance improvement of 25% percent over hash-mapping from key to entry
//   was found in the Go implementation. Test data was million-line CSV and
//   DKVP, with a dozen columns or so.
//
// Note however that an auxiliary constructor is provided which does use
// a key-to-entry hashmap in place of linear search for get/put/has/delete.
// This may be useful in certain contexts, even though it's not the default
// chosen for stream-records.
//
// ----------------------------------------------------------------
// MOTIVATION
//
// * The use case for records in Miller is that *all* fields are read from
//   strings & written to strings (split/join), while only *some* fields are
//   operated on.
//
// * Meanwhile there are few repeated accesses to a given record: the
//   access-to-construct ratio is quite low for Miller data records.  Miller
//   instantiates thousands, millions, billions of records (depending on the
//   input data) but accesses each record only once per mapping operation.
//   (This is in contrast to accumulator hashmaps which are repeatedly accessed
//   during a stats run.)
//
// * The hashed impl computes hashsums for *all* fields whether operated on or not,
//   for the benefit of the *few* fields looked up during the mapping operation.
//
// * The hashless impl only keeps string pointers.  Lookups are done at runtime
//   doing prefix search on the key names. Assuming field names are distinct,
//   this is just a few char-ptr accesses which (in experiments) turn out to
//   offer about a 10-15% performance improvement.
//
// * Added benefit: the field-rename operation (preserving field order) becomes
//   trivial.
// ================================================================

package lib

// ----------------------------------------------------------------
type Lrec struct {
	FieldCount int64
	Head       *lrecEntry
	Tail       *lrecEntry

	// Surprisingly, using this costs about 25% for cat/cut/etc tests
	// on million-line data files (CSV, DKVP) with a dozen or so columns.
	// So, the constructor allows callsites to use it, or not.
	keysToEntries map[string]*lrecEntry
}

type lrecEntry struct {
	Key   *string
	Value *Mlrval
	Prev  *lrecEntry
	Next  *lrecEntry
}

// ----------------------------------------------------------------
func NewLrec() *Lrec {
	return NewLrecNoMap()
}

// Faster on record-stream data as noted above.
func NewLrecNoMap() *Lrec {
	return &Lrec{
		FieldCount:    0,
		Head:          nil,
		Tail:          nil,
		keysToEntries: nil,
	}
}

func NewLrecMap() *Lrec {
	return &Lrec{
		FieldCount:    0,
		Head:          nil,
		Tail:          nil,
		keysToEntries: make(map[string]*lrecEntry),
	}
}

// ----------------------------------------------------------------
func newLrecEntry(key *string, value *Mlrval) *lrecEntry {
	kcopy := *key
	vcopy := *value
	return &lrecEntry{
		&kcopy,
		&vcopy,
		nil,
		nil,
	}
}