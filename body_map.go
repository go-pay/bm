package bm

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"net/url"
	"sort"
	"strings"
)

type BodyMap map[string]any

type xmlMapMarshal struct {
	XMLName xml.Name
	Value   any `xml:",cdata"`
}

type xmlMapUnmarshal struct {
	XMLName xml.Name
	Value   string `xml:",cdata"`
}

type File struct {
	Name    string `json:"name"`
	Content []byte `json:"content"`
}

// 设置参数
func (bm BodyMap) Set(key string, value any) BodyMap {
	bm[key] = value
	return bm
}

func (bm BodyMap) SetBodyMap(key string, value func(bm BodyMap)) BodyMap {
	_bm := make(BodyMap)
	value(_bm)
	bm[key] = _bm
	return bm
}

// 设置 FormFile
func (bm BodyMap) SetFormFile(key string, file *File) BodyMap {
	bm[key] = file
	return bm
}

// 获取参数转换string
func (bm BodyMap) GetString(key string) string {
	if bm == nil {
		return ""
	}
	value, ok := bm[key]
	if !ok {
		return ""
	}
	v, ok := value.(string)
	if !ok {
		return convertToString(value)
	}
	return v
}

// Deprecated: Use GetAny instead.
// 获取原始参数
func (bm BodyMap) GetInterface(key string) any {
	if bm == nil {
		return nil
	}
	return bm[key]
}

// 获取原始参数
func (bm BodyMap) GetAny(key string) any {
	if bm == nil {
		return nil
	}
	return bm[key]
}

// 解析到结构体指针
func (bm BodyMap) Decode(key string, ptr any) error {
	if bm == nil {
		return errors.New("BodyMap is nil")
	}
	if err := json.Unmarshal([]byte(bm.GetString(key)), ptr); err != nil {
		return err
	}
	return nil
}

// 删除参数
func (bm BodyMap) Remove(key string) {
	delete(bm, key)
}

// 置空BodyMap
func (bm BodyMap) Reset() {
	for k := range bm {
		delete(bm, k)
	}
}

func (bm BodyMap) JsonBody() (jb string) {
	bs, err := json.Marshal(bm)
	if err != nil {
		return ""
	}
	jb = string(bs)
	return jb
}

// unmarshal JSON bytes in BodyMap
func (bm BodyMap) Unmarshal(jsonBs []byte) error {
	if err := json.Unmarshal(jsonBs, &bm); err != nil {
		return err
	}
	return nil
}

// unmarshal JSON string in BodyMap

func (bm BodyMap) UnmarshalString(jsonStr string) error {
	if err := json.Unmarshal([]byte(jsonStr), &bm); err != nil {
		return err
	}
	return nil
}

func (bm BodyMap) MarshalXML(e *xml.Encoder, start xml.StartElement) (err error) {
	if len(bm) == 0 {
		return nil
	}
	start.Name = xml.Name{Space: "", Local: "xml"}
	if err = e.EncodeToken(start); err != nil {
		return
	}
	for k := range bm {
		if v := bm.GetString(k); v != "" {
			_ = e.Encode(xmlMapMarshal{XMLName: xml.Name{Local: k}, Value: v})
		}
	}
	return e.EncodeToken(start.End())
}

func (bm *BodyMap) UnmarshalXML(d *xml.Decoder, _ xml.StartElement) (err error) {
	for {
		var e xmlMapUnmarshal
		err = d.Decode(&e)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		bm.Set(e.XMLName.Local, e.Value)
	}
}

// ("bar=baz&foo=quux") sorted by key.
func (bm BodyMap) EncodeURLParams() string {
	if bm == nil {
		return ""
	}
	var (
		buf  strings.Builder
		keys []string
	)
	for k := range bm {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if v := bm.GetString(k); v != "" {
			buf.WriteString(url.QueryEscape(k))
			buf.WriteByte('=')
			buf.WriteString(url.QueryEscape(v))
			buf.WriteByte('&')
		}
	}
	if buf.Len() <= 0 {
		return ""
	}
	return buf.String()[:buf.Len()-1]
}

func (bm BodyMap) CheckEmptyError(keys ...string) error {
	var emptyKeys []string
	for _, k := range keys {
		if v := bm.GetString(k); v == "" {
			emptyKeys = append(emptyKeys, k)
		}
	}
	if len(emptyKeys) > 0 {
		return errors.New(strings.Join(emptyKeys, ", ") + " : cannot be empty")
	}
	return nil
}

func (bm BodyMap) Range(f func(k string, v any) bool) {
	for k, v := range bm {
		if !f(k, v) {
			break
		}
	}
}

func convertToString(v any) (str string) {
	if v == nil {
		return ""
	}
	var (
		bs  []byte
		err error
	)
	if bs, err = json.Marshal(v); err != nil {
		return ""
	}
	str = string(bs)
	return
}
