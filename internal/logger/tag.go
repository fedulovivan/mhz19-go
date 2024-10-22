package logger

import (
	"fmt"
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
	tags []string
}

func NewTag(first TagName) Tag {
	return &tag{
		tags: []string{string(first)},
	}
}

func (t *tag) With(format string, a ...any) Tag {
	// return NewTag("")
	tagscopy := make([]string, len(t.tags)+1)
	copy(tagscopy, t.tags)
	tagscopy[len(t.tags)] = fmt.Sprintf(format, a...)
	return &tag{
		tags: tagscopy,
	}
}

func (t *tag) WithTid(ns string) Tag {
	// return NewTag("")
	secsMu.Lock()
	defer secsMu.Unlock()
	if _, exist := nsSequences[ns]; !exist {
		nsSequences[ns] = utils.NewSeq(0)
	}
	return t.With("%s#%v", ns, nsSequences[ns].Inc())
}

func (t *tag) F(format string, a ...any) string {
	// return ""
	return strings.Join(
		append(t.tags, fmt.Sprintf(format, a...)),
		" ",
	)
}

// type tag struct {
// 	tags *strings.Builder
// }

// const (
// 	SPACE = byte(' ')
// 	HASH  = byte('#')
// )

// func NewTag(first TagName) Tag {
// 	b := &strings.Builder{}
// 	b.WriteString(string(first))
// 	b.WriteByte(SPACE)
// 	return &tag{
// 		tags: b,
// 	}
// }

// // format
// func (t *tag) F(format string, a ...any) string {
// 	return t.tags.String() + fmt.Sprintf(format, a...)
// }

// // extend with arbitrary string
// func (t *tag) With(format string, a ...any) Tag {
// 	bcopy := &strings.Builder{}
// 	bcopy.WriteString(t.tags.String())
// 	fmt.Fprintf(bcopy, format, a...)
// 	bcopy.WriteByte(SPACE)
// 	return &tag{
// 		tags: bcopy,
// 	}
// }

// // extend with tid
// func (t *tag) WithTid(ns string) Tag {
// 	secsMu.Lock()
// 	defer secsMu.Unlock()
// 	if _, exist := nsSequences[ns]; !exist {
// 		nsSequences[ns] = utils.NewSeq(0)
// 	}
// 	bcopy := &strings.Builder{}
// 	bcopy.WriteString(t.tags.String())
// 	bcopy.WriteString(ns)
// 	bcopy.WriteByte(HASH)
// 	bcopy.WriteString(strconv.FormatInt(
// 		int64(nsSequences[ns].Inc()),
// 		10,
// 	))
// 	bcopy.WriteByte(SPACE)
// 	return &tag{
// 		tags: bcopy,
// 	}
// }
