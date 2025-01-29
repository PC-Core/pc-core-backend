package errors

import "fmt"

const (
	HE_ERROR_FORMAT = "Error with the header %s with code %d"
)

const (
	HEADER_AUTHORIZATION = "Authorization"
)

// HU stands for Human Readable
const (
	HR_HEADER_MISSING_FORMAT = "%s header is missing"
)

type HeaderError struct {
	Code ErrorCode
	Kind ErrorKind
	// Header contains the name of the header caused the error
	Header string
}

func MissingHeader(header string) *HeaderError {
	return &HeaderError{
		EC_HEADER_MISSING,
		EK_HEADER,
		header,
	}
}

func (h *HeaderError) Error() string {
	return fmt.Sprintf(HE_ERROR_FORMAT, h.Header, h.Code)
}

func (h *HeaderError) GetErrorCode() ErrorCode {
	return h.Code
}

func (h *HeaderError) GetErrorKind() ErrorKind {
	return h.Kind
}

func (h *HeaderError) IntoPublic() *PublicPCCError {
	return NewPublicPCCError(
		h.Code,
		h.Kind,
		map[string]string{
			"caused_header": h.Header,
		},
		h.formatHumanReadableErrorIfPossible(fmt.Sprintf(HE_ERROR_FORMAT, h.Header, h.Code)),
	)
}

// formatHumanReadableErrorIfPossible trys to return the human readable message if possible.
// Else returns the def
func (h *HeaderError) formatHumanReadableErrorIfPossible(def string) string {
	if h.Code == EC_HEADER_MISSING {
		return fmt.Sprintf(HR_HEADER_MISSING_FORMAT, h.Header)
	} else {
		return def
	}
}
