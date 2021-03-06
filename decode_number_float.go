package gojay

// DecodeFloat64 reads the next JSON-encoded value from its input and stores it in the float64 pointed to by v.
//
// See the documentation for Unmarshal for details about the conversion of JSON into a Go value.
func (dec *Decoder) DecodeFloat64(v *float64) error {
	if dec.isPooled == 1 {
		panic(InvalidUsagePooledDecoderError("Invalid usage of pooled decoder"))
	}
	return dec.decodeFloat64(v)
}
func (dec *Decoder) decodeFloat64(v *float64) error {
	for ; dec.cursor < dec.length || dec.read(); dec.cursor++ {
		switch c := dec.data[dec.cursor]; c {
		case ' ', '\n', '\t', '\r', ',':
			continue
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			val, err := dec.getFloat(c)
			if err != nil {
				return err
			}
			*v = val
			return nil
		case '-':
			dec.cursor = dec.cursor + 1
			val, err := dec.getFloat(c)
			if err != nil {
				return err
			}
			*v = -val
			return nil
		case 'n':
			dec.cursor++
			err := dec.assertNull()
			if err != nil {
				return err
			}
			return nil
		default:
			dec.err = dec.makeInvalidUnmarshalErr(v)
			err := dec.skipData()
			if err != nil {
				return err
			}
			return nil
		}
	}
	return dec.raiseInvalidJSONErr(dec.cursor)
}

func (dec *Decoder) getFloat(b byte) (float64, error) {
	var end = dec.cursor
	var start = dec.cursor
	// look for following numbers
	for j := dec.cursor + 1; j < dec.length || dec.read(); j++ {
		switch dec.data[j] {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			end = j
			continue
		case '.':
			// we get part before decimal as integer
			beforeDecimal := dec.atoi64(start, end)
			// then we get part after decimal as integer
			start = j + 1
			// get number after the decimal point
			// multiple the before decimal point portion by 10 using bitwise
			for i := j + 1; i < dec.length || dec.read(); i++ {
				c := dec.data[i]
				if isDigit(c) {
					end = i
					beforeDecimal = (beforeDecimal << 3) + (beforeDecimal << 1)
					continue
				} else if c == 'e' || c == 'E' {
					afterDecimal := dec.atoi64(start, end)
					dec.cursor = i + 1
					pow := pow10uint64[end-start+2]
					floatVal := float64(beforeDecimal+afterDecimal) / float64(pow)
					exp := dec.getExponent()
					// if exponent is negative
					if exp < 0 {
						return float64(floatVal) * (1 / float64(pow10uint64[exp*-1+1])), nil
					}
					return float64(floatVal) * float64(pow10uint64[exp+1]), nil
				}
				dec.cursor = i
				break
			}
			// then we add both integers
			// then we divide the number by the power found
			afterDecimal := dec.atoi64(start, end)
			pow := pow10uint64[end-start+2]
			return float64(beforeDecimal+afterDecimal) / float64(pow), nil
		case 'e', 'E':
			dec.cursor = dec.cursor + 2
			// we get part before decimal as integer
			beforeDecimal := uint64(dec.atoi64(start, end))
			// get exponent
			exp := dec.getExponent()
			// if exponent is negative
			if exp < 0 {
				return float64(beforeDecimal) * (1 / float64(pow10uint64[exp*-1+1])), nil
			}
			return float64(beforeDecimal) * float64(pow10uint64[exp+1]), nil
		case ' ', '\n', '\t', '\r', ',', '}', ']': // does not have decimal
			dec.cursor = j
			return float64(dec.atoi64(start, end)), nil
		}
		// invalid json we expect numbers, dot (single one), comma, or spaces
		return 0, dec.raiseInvalidJSONErr(dec.cursor)
	}
	return float64(dec.atoi64(start, end)), nil
}

// DecodeFloat32 reads the next JSON-encoded value from its input and stores it in the float32 pointed to by v.
//
// See the documentation for Unmarshal for details about the conversion of JSON into a Go value.
func (dec *Decoder) DecodeFloat32(v *float32) error {
	if dec.isPooled == 1 {
		panic(InvalidUsagePooledDecoderError("Invalid usage of pooled decoder"))
	}
	return dec.decodeFloat32(v)
}
func (dec *Decoder) decodeFloat32(v *float32) error {
	for ; dec.cursor < dec.length || dec.read(); dec.cursor++ {
		switch c := dec.data[dec.cursor]; c {
		case ' ', '\n', '\t', '\r', ',':
			continue
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			val, err := dec.getFloat32(c)
			if err != nil {
				return err
			}
			*v = val
			return nil
		case '-':
			dec.cursor = dec.cursor + 1
			val, err := dec.getFloat32(c)
			if err != nil {
				return err
			}
			*v = -val
			return nil
		case 'n':
			dec.cursor++
			err := dec.assertNull()
			if err != nil {
				return err
			}
			return nil
		default:
			dec.err = dec.makeInvalidUnmarshalErr(v)
			err := dec.skipData()
			if err != nil {
				return err
			}
			return nil
		}
	}
	return dec.raiseInvalidJSONErr(dec.cursor)
}

func (dec *Decoder) getFloat32(b byte) (float32, error) {
	var end = dec.cursor
	var start = dec.cursor
	// look for following numbers
	for j := dec.cursor + 1; j < dec.length || dec.read(); j++ {
		switch dec.data[j] {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			end = j
			continue
		case '.':
			// we get part before decimal as integer
			beforeDecimal := dec.atoi32(start, end)
			// then we get part after decimal as integer
			start = j + 1
			// get number after the decimal point
			// multiple the before decimal point portion by 10 using bitwise
			for i := j + 1; i < dec.length || dec.read(); i++ {
				c := dec.data[i]
				if isDigit(c) {
					end = i
					beforeDecimal = (beforeDecimal << 3) + (beforeDecimal << 1)
					continue
				} else if c == 'e' || c == 'E' {
					afterDecimal := dec.atoi32(start, end)
					dec.cursor = i + 1
					pow := pow10uint64[end-start+2]
					floatVal := float32(beforeDecimal+afterDecimal) / float32(pow)
					exp := dec.getExponent()
					// if exponent is negative
					if exp < 0 {
						return float32(floatVal) * (1 / float32(pow10uint64[exp*-1+1])), nil
					}
					return float32(floatVal) * float32(pow10uint64[exp+1]), nil
				}
				dec.cursor = i
				break
			}
			// then we add both integers
			// then we divide the number by the power found
			afterDecimal := dec.atoi32(start, end)
			pow := pow10uint64[end-start+2]
			return float32(beforeDecimal+afterDecimal) / float32(pow), nil
		case 'e', 'E':
			dec.cursor = dec.cursor + 2
			// we get part before decimal as integer
			beforeDecimal := uint32(dec.atoi32(start, end))
			// get exponent
			exp := dec.getExponent()
			// if exponent is negative
			if exp < 0 {
				return float32(beforeDecimal) * (1 / float32(pow10uint64[exp*-1+1])), nil
			}
			return float32(beforeDecimal) * float32(pow10uint64[exp+1]), nil
		case ' ', '\n', '\t', '\r', ',', '}', ']': // does not have decimal
			dec.cursor = j
			return float32(dec.atoi32(start, end)), nil
		}
		// invalid json we expect numbers, dot (single one), comma, or spaces
		return 0, dec.raiseInvalidJSONErr(dec.cursor)
	}
	return float32(dec.atoi32(start, end)), nil
}
