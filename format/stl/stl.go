package stl

// https://www.fabbers.com/tech/STL_Format
// https://paulbourke.net/dataformats/stl/

import (
	"embed"

	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	// "github.com/wader/fq/pkg/scalar"
)

////go:embed stl.jq
////go:embed stl.md
var stlFS embed.FS

func init() {
	interp.RegisterFormat(
		format.STL,
		&decode.Format{
			Description: "Stereolithography",
			DecodeFn:    decodeSTL,
			// Functions:   []string{"torepr"},
		})
	interp.RegisterFS(stlFS)
}

const (
	headerLength = 80
)

func decodeVector(d *decode.D) {
	d.FieldF32("x")
	d.FieldF32("y")
	d.FieldF32("z")
}

func decodeSTLModel(d *decode.D) {
	// TODO support color of Materialise Magics
	d.FieldUTF8NullFixedLen("header", headerLength)
	// d.FieldUTF8("header", headerLength)

	// TODO parse header components? not guaranteed to be null-terminated
	// d.SeekAbs(0)
	// d.FramedFn(headerLength * 8, func(d *decode.D) {
	// 	d.FieldArray("headerValues", func(d *decode.D) {
	// 		for !d.End() {
	// 			// d.FieldUTF8Null("header")
	// 			header := d.UTF8Null();
	// 			if header != "" {
	// 				// TODO add/track position of header values
	// 				d.FieldValueStr("header", header)
	// 				// d.FieldScalarStr("header", header)
	// 			}
	// 		}
	// 	})
	// })

	triangleCount := d.FieldU32("numberOfFacets")
	d.FieldStructNArray("facets", "facet", int64(triangleCount), func(d *decode.D) {
		d.FieldStruct("normal", decodeVector)
		d.FieldStructNArray("vertices", "vertex", 3, decodeVector)
		// TODO support color of VisCAM and SolidView
		// TODO support color of Materialise Magics
		attributeByteCount := d.FieldU16("attributeByteCount")
		if attributeByteCount > 0 {
			d.FieldRawLen("attribute", int64(attributeByteCount) * 8)
		}
	})

}

func decodeSTL(d *decode.D) any {
	d.Endian = decode.LittleEndian

	decodeSTLModel(d)

	return nil
}
