package dqueue

import "time"

// maxNextInterval shows our queue can hold tasks up to 24 hours by default,
// since executing golang application longer than it can lead to memleaks
// and funny things done by GC.
const maxNextInterval = 24 * time.Hour
