// Code generated by "stringer -type=Severity"; DO NOT EDIT

package mylog

import "fmt"

const (
	_Severity_name_0 = "ERRORDEBUG"
	_Severity_name_1 = "FATAL"
	_Severity_name_2 = "INFO"
	_Severity_name_3 = "WARN"
)

var (
	_Severity_index_0 = [...]uint8{0, 5, 10}
	_Severity_index_1 = [...]uint8{0, 5}
	_Severity_index_2 = [...]uint8{0, 4}
	_Severity_index_3 = [...]uint8{0, 4}
)

func (i Severity) String() string {
	switch {
	case 1 <= i && i <= 2:
		i -= 1
		return _Severity_name_0[_Severity_index_0[i]:_Severity_index_0[i+1]]
	case i == 4:
		return _Severity_name_1
	case i == 8:
		return _Severity_name_2
	case i == 16:
		return _Severity_name_3
	default:
		return fmt.Sprintf("Severity(%d)", i)
	}
}
