package shared

import "time"

const (
	// CHARS for setting short random string
	CHARS = "abcdefghijklmnopqrstuvwxyz0123456789"
	// ErrorDataNotFound error message when data doesn't exist
	ErrorDataNotFound = "Data %s tidak ditemukan"
	// ErrorDataAlreadyExist error message when data already exist
	ErrorDataAlreadyExist = "Data %s sudah terdaftar"
	// ErrorParameterInvalid error message for parameter is invalid
	ErrorParameterInvalid = "Parameter %s tidak valid"
	// ErrorParameterRequired error message for parameter is missing
	ErrorParameterRequired = "Parameter %s dibutuhkan"
	// ErrorParameterDuplicate error message for parameter is duplicate
	ErrorParameterDuplicate = "Parameter %s duplikat"
	// ErrorHeader error message for header is missing
	ErrorHeader = "Header %s dibutuhkan"
	// ErrorParameterLength error message for parameter length is invalid
	ErrorParameterLength = "Panjang parameter %s melebihi batas %d"
	// ErrorUnauthorized error message for unauthorized user
	ErrorUnauthorized = "Anda tidak memiliki hak akses"
	// SuccessMessage message for success process
	SuccessMessage = "Berhasil memproses data %s"
	// SuccessSaveEdit message for success process
	SuccessSaveEdit = "Berhasil menyimpan/mengubah data %s"
	// SuccessGetList message for get list
	SuccessGetList = "Berhasil mendapatkan daftar data %s"
	// SuccessGetDetail message for get detail
	SuccessGetDetail = "Berhasil mendapatkan data %s"
	// StatusSuccess message for success status
	StatusSuccess = "OK"
	// ErrorBadRequest message for bad request
	ErrorBadRequest = "bad request"
	// ErrorUnknown message for unknown error
	ErrorUnknown = "Kesalahan tidak diketahui"
	// ErrorParseData message for failed parse data
	ErrorParseData = "Sistem tidak dapat memproses data %s"
	// ErrorDataNotActived message for data not active
	ErrorDataNotActived = "Data %s tidak aktif"
	// ErrorDataIsActived message for data not active
	ErrorDataIsActived = "Data %s dalam keadaan aktif"
	// ErrorMessage message for failed process
	ErrorMessage = "Gagal memproses data %s"
	// DateFormat date formatting
	DateFormat = "2006-01-02T15:04:05Z"
)

var (
	// IndonesianMonth mapping system month to indonesian month
	IndonesianMonth = map[string]string{
		time.January.String():   "Januari",
		time.February.String():  "Februari",
		time.March.String():     "Maret",
		time.April.String():     "April",
		time.May.String():       "Mei",
		time.June.String():      "Juni",
		time.July.String():      "Juli",
		time.August.String():    "Agustus",
		time.September.String(): "September",
		time.October.String():   "Oktober",
		time.November.String():  "November",
		time.December.String():  "Desember",
	}
)
