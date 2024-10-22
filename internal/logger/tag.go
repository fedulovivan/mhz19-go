package logger

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/fedulovivan/mhz19-go/pkg/utils"
)

var nsSequences = make(map[string]utils.Seq)
var secsMu = new(sync.Mutex)

type Tag interface {
	With(string, ...any) Tag
	WithTid(string) Tag
	F(format string, a ...any) string
}

type tag struct {
	tags strings.Builder
}

const (
	SPACE = ' '
	HASH  = '#'
)

func NewTag(first TagName) Tag {
	res := &tag{}
	res.tags.WriteString(string(first))
	res.tags.WriteRune(SPACE)
	return res
}

// format existing
func (t *tag) F(format string, a ...any) string {
	return t.tags.String() + fmt.Sprintf(format, a...)
}

// extend "tags chain" with arbitrary formatted string
func (t *tag) With(format string, a ...any) Tag {
	res := &tag{}
	res.tags.WriteString(t.tags.String())
	fmt.Fprintf(&res.tags, format, a...)
	res.tags.WriteRune(SPACE)
	return res
}

// extend "tags chain" with tid
func (t *tag) WithTid(ns string) Tag {
	secsMu.Lock()
	defer secsMu.Unlock()
	if _, exist := nsSequences[ns]; !exist {
		nsSequences[ns] = utils.NewSeq(0)
	}
	res := &tag{}
	res.tags.WriteString(t.tags.String())
	res.tags.WriteString(ns)
	res.tags.WriteRune(HASH)
	res.tags.WriteString(strconv.FormatInt(
		int64(nsSequences[ns].Inc()),
		10,
	))
	res.tags.WriteRune(SPACE)
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
