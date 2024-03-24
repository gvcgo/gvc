package asciinema

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/cirocosta/asciinema-edit/cast"
	"github.com/cirocosta/asciinema-edit/commands/transformer"
	"github.com/cirocosta/asciinema-edit/editor"
	"github.com/pkg/errors"
)

type quantizeTransformation struct {
	ranges []editor.QuantizeRange
}

func (t *quantizeTransformation) Transform(c *cast.Cast) (err error) {
	err = editor.Quantize(c, t.ranges)
	return
}

// ParseQuantizeRange takes an input string that represents
// a quantization range and converts it into a QuantizeRange
// instance.
//
// It allows both bounded and unbounded ranges.
//
// For instance:
// - bounded: 1,2
// - unbounded: 1
//
// Fails if the input can't be converted to a QuantizeRange.
func ParseQuantizeRange(input string) (res editor.QuantizeRange, err error) {
	cols := strings.Split(input, ",")

	if len(cols) > 2 {
		err = errors.Errorf(
			"invalid range format: must be `value[,value]`")
		return
	}

	if len(cols) == 2 {
		res.To, err = strconv.ParseFloat(cols[1], 64)
		if err != nil {
			err = errors.Errorf(
				"malformed range: second element is not a float '%s'", cols[1])
			return
		}
	}

	res.From, err = strconv.ParseFloat(cols[0], 64)
	if err != nil {
		err = errors.Errorf(
			"malformed range: first element is not a float '%s'", cols[0])
		return
	}

	if res.To == 0 {
		res.To = math.MaxFloat64
	}

	if res.From < 0 {
		err = errors.Errorf(
			"constraint not verified: from >= 0")
		return
	}

	if res.To <= res.From {
		err = errors.Errorf(
			"constraint not verified: from < to")
		return
	}

	return
}

func parseQuantizeRanges(inputs []string) (ranges []editor.QuantizeRange, err error) {
	ranges = make([]editor.QuantizeRange, 0)

	var (
		qRange editor.QuantizeRange
		input  string
	)

	for _, input = range inputs {
		qRange, err = ParseQuantizeRange(input)
		if err != nil {
			err = errors.Wrapf(err, "failed to parse range %s",
				input)
			return
		}

		ranges = append(ranges, qRange)
	}

	return
}

/*
The command acts on the delays between the frames, reducing such
timings to the lowest value defined in a given range that they
lie in.
*/

func (a *Asciinema) Quantize(fPath, outFilePath string, ranges []string) (err error) {
	if len(ranges) == 0 {
		return fmt.Errorf("a range must be specified")
	}

	transformation := &quantizeTransformation{}
	transformation.ranges, err = parseQuantizeRanges(ranges)
	if err != nil {
		return err
	}

	var t *transformer.Transformer
	t, err = transformer.New(transformation, fPath, outFilePath)
	if err != nil {
		return err
	}
	defer t.Close()

	return t.Transform()
}
