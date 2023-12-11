package api

import (
	"regexp"
	"strings"

	vld "github.com/muonsoft/validation"
)

const base64 string = "^(?:[A-Za-z0-9+\\/]{4})*(?:[A-Za-z0-9+\\/]{2}==|[A-Za-z0-9+\\/]{3}=|[A-Za-z0-9+\\/]{4})$"

var (
	rxBase64   = regexp.MustCompile(base64)
	isNotDigit = func(c rune) bool { return c < '0' || c > '9' }
)

var (
	ErrInvalidOnlyNumbersFormat = vld.NewError(
		"invalid format: must consist only numbers",
		"This value has not the correct format.")
	ErrInvalidInsuranceChecksum = vld.NewError(
		"invalid insurance checksum",
		"This value has not the correct checksum.")
	ErrInvalidTaxpayerChecksum = vld.NewError(
		"invalid taxpayer checksum",
		"This value has not the correct checksum.")
)

func consistOnlyNumbersFormat() vld.StringFuncConstraint {
	return vld.OfStringBy(func(s string) bool {
		return strings.IndexFunc(s, isNotDigit) == -1
	}).
		WithError(ErrInvalidOnlyNumbersFormat).
		WithMessage(ErrInvalidOnlyNumbersFormat.Message())
}

func hasCorrectInsuranceChecksum() vld.StringFuncConstraint {
	return vld.OfStringBy(checksumInsurance).
		WithError(ErrInvalidInsuranceChecksum).
		WithMessage(ErrInvalidInsuranceChecksum.Message())
}

var checksumInsurance = func(s string) bool {
	var sum int
	for i := 0; i < len(s)-2; i++ {
		sum += (9 - i) * int(s[i]-'0')
	}
	mod := sum % 101
	if mod%10 != int(s[len(s)-1]-'0') {
		return false
	}
	if (mod/10)%100 != int(s[len(s)-2]-'0') {
		return false
	}
	return true
}

func hasCorrectTaxpayerChecksum() vld.StringFuncConstraint {
	return vld.OfStringBy(checksumTaxpayer).
		WithError(ErrInvalidTaxpayerChecksum).
		WithMessage(ErrInvalidTaxpayerChecksum.Message())
}

var checksumTaxpayer = func(s string) bool {
	switch len(s) {
	case 10:
		n10 := ((2*int(s[0]-'0') + 4*int(s[1]-'0') +
			10*int(s[2]-'0') + 3*int(s[3]-'0') +
			5*int(s[4]-'0') + 9*int(s[5]-'0') +
			4*int(s[6]-'0') + 6*int(s[7]-'0') +
			8*int(s[8]-'0')) % 11) % 10
		if int(s[9]-'0') != n10 {
			return false
		}
	case 12:
		n11 := ((7*int(s[0]-'0') + 2*int(s[1]-'0') +
			4*int(s[2]-'0') + 10*int(s[3]-'0') +
			3*int(s[4]-'0') + 5*int(s[5]-'0') +
			9*int(s[6]-'0') + 4*int(s[7]-'0') +
			6*int(s[8]-'0') + 8*int(s[9]-'0')) % 11) % 10
		if int(s[10]-'0') != n11 {
			return false
		}
		n12 := ((3*int(s[0]-'0') + 7*int(s[1]-'0') +
			2*int(s[2]-'0') + 4*int(s[3]-'0') +
			10*int(s[4]-'0') + 3*int(s[5]-'0') +
			5*int(s[6]-'0') + 9*int(s[7]-'0') +
			4*int(s[8]-'0') + 6*int(s[9]-'0') +
			8*n11) % 11) % 10
		if int(s[11]-'0') != n12 {
			return false
		}
	default:
		return false
	}

	return true
}

func NewBadRequestErrorFromError(err error) BadRequestError {
	if violations, ok := vld.UnwrapViolationList(err); ok {
		return BadRequestError{
			Message: violations.String(),
		}
	}
	return BadRequestError{
		Message: "bad request",
	}
}
