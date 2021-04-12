package log

import (
	"bytes"

	"github.com/gmcoringa/tswitch/pkg/util"
	"github.com/sirupsen/logrus"
)

type Format struct{}

func (f *Format) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer

	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	if util.IsNotBlank(entry.Message) {
		b.WriteString(entry.Message)
		b.WriteByte('\n')
	}

	return b.Bytes(), nil
}
