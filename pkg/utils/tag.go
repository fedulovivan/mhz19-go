package utils

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

type TagName string

var seqs = make(map[string]*atomic.Int32, 10)
var seqslock sync.Mutex

type Tag interface {
	With(string, ...any) Tag
	WithTid(string) Tag
	F(format string, a ...any) string
}

type tag struct {
	strings.Builder
}

const (
	SPACE = ' '
	HASH  = '#'
)

func NewTag(first TagName) Tag {
	res := &tag{}
	res.WriteString(string(first))
	res.WriteRune(SPACE)
	return res
}

// format existing
func (t *tag) F(format string, a ...any) string {
	return t.String() + fmt.Sprintf(format, a...)
}

// extend "tags chain" with arbitrary formatted string
func (t *tag) With(format string, a ...any) Tag {
	res := &tag{}
	res.WriteString(t.String())
	fmt.Fprintf(&res.Builder, format, a...)
	res.WriteRune(SPACE)
	return res
}

// extend "tags chain" with tid
func (t *tag) WithTid(ns string) Tag {
	seqslock.Lock()
	defer seqslock.Unlock()
	if _, exist := seqs[ns]; !exist {
		seqs[ns] = &atomic.Int32{}
	}
	res := &tag{}
	res.WriteString(t.String())
	res.WriteString(ns)
	res.WriteRune(HASH)
	res.WriteString(strconv.FormatInt(
		int64(seqs[ns].Add(1)),
		10,
	))
	res.WriteRune(SPACE)
	return res
}

// type tag struct {
// 	tags []string
// }

// func NewTag(first TagName) Tag {
// 	return &tag{
// 		tags: []string{string(first)},
// 	}
// }

// func (t *tag) With(format string, a ...any) Tag {
// 	// return NewTag("")
// 	tagscopy := make([]string, len(t.tags)+1)
// 	copy(tagscopy, t.tags)
// 	tagscopy[len(t.tags)] = fmt.Sprintf(format, a...)
// 	return &tag{
// 		tags: tagscopy,
// 	}
// }

// func (t *tag) WithTid(ns string) Tag {
// 	// return NewTag("")
// 	secsMu.Lock()
// 	defer secsMu.Unlock()
// 	if _, exist := nsSequences[ns]; !exist {
// 		nsSequences[ns] = utils.NewSeq(0)
// 	}
// 	return t.With("%s#%v", ns, nsSequences[ns].Inc())
// }

// func (t *tag) F(format string, a ...any) string {
// 	// return ""
// 	return strings.Join(
// 		append(t.tags, fmt.Sprintf(format, a...)),
// 		" ",
// 	)
// }
