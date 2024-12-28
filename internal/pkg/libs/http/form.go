package httpUtils

import (
	"bytes"
	"errors"
	"fmt"
	_logUtils "github.com/deeptest-com/deeptest-next/pkg/libs/log"
	"mime"
	"mime/multipart"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"
)

type TFormDataWriter struct {
	Writer  *multipart.Writer
	Payload *bytes.Buffer
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func MultipartEncoder(bodyFormData []BodyFormDataItem) (tFormWriter *TFormDataWriter, err error) {
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	tFormWriter = &TFormDataWriter{
		Writer:  writer,
		Payload: payload,
	}

	for _, item := range bodyFormData {
		formKey := item.Name
		formValue := item.Value

		if item.Type == FormDataTypeFile {
			err = tFormWriter.writeCustomFile(formKey, formValue, "", "")
			if err != nil {
				_logUtils.Infof("failed to write file: %v=@\"%v\", exit", formKey, formValue)
				return
			}

		} else if item.Type == FormDataTypeText {
			err = tFormWriter.writeCustomText(formKey, formValue, "", "")
			if err != nil {
				_logUtils.Infof("failed to write text: %v=%v, ignore", formKey, formValue)
				return
			}
		}
	}

	if err := writer.Close(); err != nil {
		_logUtils.Infof("failed to close form-data writer")
	}

	return
}

func (w *TFormDataWriter) writeCustomText(formKey, formValue, formType, formFileName string) error {
	if w.Writer == nil {
		return errors.New("form-data writer not initialized")
	}

	h := make(textproto.MIMEHeader)

	// text doesn't have Content-Type by default
	if formType != "" {
		h.Set("Content-Type", formType)
	}

	// text doesn't have filename in Content-Disposition by default
	if formFileName == "" {
		h.Set("Content-Disposition",
			fmt.Sprintf(`form-data; name="%s"`, escapeQuotes(formKey)))
	} else {
		h.Set("Content-Disposition",
			fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
				escapeQuotes(formKey), escapeQuotes(formFileName)))
	}

	part, err := w.Writer.CreatePart(h)
	if err != nil {
		return err
	}

	_, err = part.Write([]byte(formValue))
	return err
}

func (w *TFormDataWriter) writeCustomFile(formKey, formValue, formType, formFileName string) error {
	if w.Writer == nil {
		return errors.New("form-data writer not initialized")
	}

	fPath, err := filepath.Abs(formValue)
	if err != nil {
		return err
	}

	file, err := os.ReadFile(fPath)
	if err != nil {
		return err
	}

	if formType == "" {
		formType = inferFormType(formValue)
	}
	if formFileName == "" {
		formFileName = filepath.Base(formValue)
	}
	h := make(textproto.MIMEHeader)
	h.Set("Content-Type", formType)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			escapeQuotes(formKey), escapeQuotes(formFileName)))

	part, err := w.Writer.CreatePart(h)
	if err != nil {
		return err
	}

	_, err = part.Write(file)
	return err
}

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

func inferFormType(formValue string) string {
	extName := filepath.Ext(formValue)
	formType := mime.TypeByExtension(extName)
	if formType == "" {
		// file without extension name
		return "application/octet-stream"
	}
	if strings.HasPrefix(formType, "text") {
		// text/... types have the charset parameter set to "utf-8" by default.
		return strings.TrimSuffix(formType, "; charset=utf-8")
	}
	return formType
}

func MultipartContentType(w *TFormDataWriter) string {
	if w.Writer == nil {
		return ""
	}
	return w.Writer.FormDataContentType()
}

type BodyFormDataItem struct {
	Name        string       `json:"name"`
	Value       string       `json:"value"`
	Type        FormDataType `json:"type,omitempty"`
	Desc        string       `json:"desc,omitempty"`
	InterfaceId uint         `json:"interfaceId,omitempty"`
}

type FormDataType string

const (
	FormDataTypeText FormDataType = "text"
	FormDataTypeFile FormDataType = "file"
)

func (e FormDataType) String() string {
	return string(e)
}
