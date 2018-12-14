package customErrors

import "errors"

// NotFoundError represents not found error with appropriate error message.
type NotFoundError struct {
	Message string
}

// Error returns error message.
func (e *NotFoundError) Error() string {
	return e.Message
}

var ErrGoodNotFound = errors.New("The specified good was not found")
var ErrIncorrectGoodID = errors.New("Incorrect goodID")
var ErrUnableCreateGood = errors.New("Unable to create a good")
var ErrGoodWithNameExists = errors.New("Good with such name already exists")
var ErrUnableToCreateFile = errors.New("Unable to create a file")
var ErrUnableUpdateGood = errors.New("Unable to update a Good")
var ErrUnableToDeleteGood = errors.New("Unable to delete a good")
var ErrWrongBodyRequest = errors.New("Wrong body request")
var ErrNoFilesUploaded = errors.New("No files were uploaded")
var ErrIncorrectArguments = errors.New("Incorrect or missed arguments")
var ErrInternalServerError = errors.New("Internal Server Error")

var badRequest = []error{
	ErrGoodWithNameExists,
	ErrWrongBodyRequest,
	ErrNoFilesUploaded,
	ErrIncorrectArguments,
}

var notFound = []error{
	ErrGoodNotFound,
}

func contains(arr []error, errToFind error) bool {
	for _, err := range arr {
		if err == errToFind {
			return true
		}
	}
	return false
}

func IsNotFoundError(err error) bool {
	if contains(notFound, err) {
		return true
	}

	_, ok := err.(*NotFoundError)
	return ok
}

func IsBadRequest(err error) bool {
	return contains(badRequest, err)
}
