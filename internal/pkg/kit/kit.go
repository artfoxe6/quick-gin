package kit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"math/rand"
	"regexp"
	"strings"
	"time"
	"unicode"
)

func String2Json(str string) map[string]any {
	var res map[string]any
	_ = json.Unmarshal([]byte(str), &res)
	return res
}

func String2JsonArray(str string) []map[string]any {
	var res []map[string]any
	_ = json.Unmarshal([]byte(str), &res)
	return res
}

func Json2String(object any) string {
	res, _ := json.Marshal(object)
	return string(res)
}

func MustSplit(str, delimiter string) []string {
	if str == "" {
		return nil
	}
	return strings.Split(str, delimiter)
}

func MustJoin(arr []string, delimiter string) string {
	if len(arr) == 0 {
		return ""
	}
	return strings.Join(arr, delimiter)
}

func Slug(input string) string {
	re := regexp.MustCompile(`\s+`)
	input = re.ReplaceAllString(input, " ")

	var output string

	for _, char := range input {
		if unicode.IsLetter(char) || unicode.IsDigit(char) {
			output += string(char)
		} else if unicode.IsSpace(char) {
			output += "-"
		}
	}

	return strings.ToLower(output)
}

func InArray(str string, arr []string) bool {
	for _, item := range arr {
		if item == str {
			return true
		}
	}
	return false
}

func GenCode(length int) string {

	randGen := rand.New(rand.NewSource(time.Now().UnixNano()))
	var code = ""
	for i := 0; i < length; i++ {
		code += fmt.Sprintf("%d", randGen.Intn(10))
	}
	return code
}

func Compress(content []byte, typ string) ([]byte, error) {
	var img image.Image
	var err error
	old := bytes.NewReader(content)
	buf := new(bytes.Buffer)

	switch typ {
	case "image/png":
		img, err = png.Decode(old)
		if err != nil {
			return nil, err
		}
	case "image/jpeg":
		img, err = jpeg.Decode(old)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported image type: %s", typ)
	}

	if err = jpeg.Encode(buf, img, &jpeg.Options{Quality: 50}); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
