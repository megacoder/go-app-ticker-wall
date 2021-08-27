// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package polygon

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonE9bd36c2DecodeGithubComPolygonIoGoAppTickerWallPolygonClient(in *jlexer.Lexer, out *websocketTrades) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(websocketTrades, 0, 1)
			} else {
				*out = websocketTrades{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v1 websocketTrade
			(v1).UnmarshalEasyJSON(in)
			*out = append(*out, v1)
			in.WantComma()
		}
		in.Delim(']')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonE9bd36c2EncodeGithubComPolygonIoGoAppTickerWallPolygonClient(out *jwriter.Writer, in websocketTrades) {
	if in == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v2, v3 := range in {
			if v2 > 0 {
				out.RawByte(',')
			}
			(v3).MarshalEasyJSON(out)
		}
		out.RawByte(']')
	}
}

// MarshalJSON supports json.Marshaler interface
func (v websocketTrades) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonE9bd36c2EncodeGithubComPolygonIoGoAppTickerWallPolygonClient(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v websocketTrades) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonE9bd36c2EncodeGithubComPolygonIoGoAppTickerWallPolygonClient(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *websocketTrades) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonE9bd36c2DecodeGithubComPolygonIoGoAppTickerWallPolygonClient(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *websocketTrades) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonE9bd36c2DecodeGithubComPolygonIoGoAppTickerWallPolygonClient(l, v)
}
func easyjsonE9bd36c2DecodeGithubComPolygonIoGoAppTickerWallPolygonClient1(in *jlexer.Lexer, out *websocketTrade) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "ev":
			out.Event = string(in.String())
		case "i":
			out.ID = string(in.String())
		case "p":
			out.Price = float64(in.Float64())
		case "h":
			out.High = float64(in.Float64())
		case "sym":
			out.Ticker = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonE9bd36c2EncodeGithubComPolygonIoGoAppTickerWallPolygonClient1(out *jwriter.Writer, in websocketTrade) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"ev\":"
		out.RawString(prefix[1:])
		out.String(string(in.Event))
	}
	{
		const prefix string = ",\"i\":"
		out.RawString(prefix)
		out.String(string(in.ID))
	}
	{
		const prefix string = ",\"p\":"
		out.RawString(prefix)
		out.Float64(float64(in.Price))
	}
	{
		const prefix string = ",\"h\":"
		out.RawString(prefix)
		out.Float64(float64(in.High))
	}
	{
		const prefix string = ",\"sym\":"
		out.RawString(prefix)
		out.String(string(in.Ticker))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v websocketTrade) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonE9bd36c2EncodeGithubComPolygonIoGoAppTickerWallPolygonClient1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v websocketTrade) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonE9bd36c2EncodeGithubComPolygonIoGoAppTickerWallPolygonClient1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *websocketTrade) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonE9bd36c2DecodeGithubComPolygonIoGoAppTickerWallPolygonClient1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *websocketTrade) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonE9bd36c2DecodeGithubComPolygonIoGoAppTickerWallPolygonClient1(l, v)
}