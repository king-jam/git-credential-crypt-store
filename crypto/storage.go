package crypto

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
)

const (
	// storageDelimiter is used to space the component of the storageLayout
	// value. DO NOT CHANGE THIS A LOT OR WE WILL HAVE TO HANDLE MIGRATIONS.
	storageDelimiter = "_"
	// versionIndex stores the index in the storageLayout where version is
	nonceIndex = 0
	// nonceIndex stores the index in the storageLayout where the nonce is
	valueIndex = 1
)

var (
	// base64EncodingType stores the type of encoding used for encode/decode operations
	base64EncodingType = base64.RawStdEncoding
	// ErrStorageLayoutEncodingInput returns an error when encoding has failed for
	// invalid input errors
	ErrStorageLayoutEncodingInput = errors.New("encoding failure: invalid inputs")
	// ErrStorageLayoutDecodingInput returns an error when decoding has failed for
	// invalid input errors
	ErrStorageLayoutDecodingInput = errors.New("decoding failure: invalid inputs")
)

// storageLayout handles how the hash will be stored in the database
// This helper type does Encoding/Decoding operations and accessor
// methods to simplify working with the encrypted values.
// the format is:
// <nonce - base64><delimiter><ciphertext - base64>
// the delimiter cannot be any character within the base64 encoding range
type storageLayout struct {
	nonce []byte
	value []byte
}

// Encode will take the values within a storageLayout and output a ready to store
// byte slice. It will validate that values have been filled.
func (sl storageLayout) Encode() ([]byte, error) {
	if len(sl.nonce) == 0 {
		return []byte{}, ErrStorageLayoutEncodingInput
	}
	if len(sl.value) == 0 {
		return []byte{}, ErrStorageLayoutEncodingInput
	}
	delimiter := []byte(storageDelimiter)
	buffer := new(bytes.Buffer)
	// base64 encode it so we don't overlap our delimiters
	nonceEncoder := base64.NewEncoder(base64EncodingType, buffer)
	// append on the nonce value - <nonce>
	_, err := nonceEncoder.Write(sl.nonce) // write through the encoder to the buffer
	if err != nil {
		return []byte{}, ErrStorageLayoutEncodingInput
	}
	nonceEncoder.Close() // close it so everything flushes
	// append on the third delimiter - <nonce><delimiter>
	_, err = buffer.Write(delimiter)
	if err != nil {
		return []byte{}, ErrStorageLayoutEncodingInput
	}
	// base64 encode it so we don't overlap our delimiters
	valueEncoder := base64.NewEncoder(base64EncodingType, buffer)
	// append on the ciphertext value - <nonce><delimiter><value>
	_, err = valueEncoder.Write(sl.value) // write through the encoder to the buffer
	if err != nil {
		return []byte{}, ErrStorageLayoutEncodingInput
	}
	valueEncoder.Close() // close it so everything flushes
	return buffer.Bytes(), nil
}

// Decode will take a byte slice in the storageLayout format and marshal it
// into a storageLayout struct to be accessed with the accessor methods
func (sl *storageLayout) Decode(b []byte) error {
	// ensure we didn't get an empty array
	if len(b) == 0 {
		return ErrStorageLayoutDecodingInput
	}
	// split it based on the delimiter, if we change this for a system in
	// production and update....we are going to have a bad time.
	parts := bytes.Split(b, []byte(storageDelimiter))
	// check that we ONLY get 2 parts, more or less is an error
	if len(parts) != 2 {
		return ErrStorageLayoutDecodingInput
	}
	// make a temp slice to read things into
	nonceBuf := make([]byte, base64EncodingType.DecodedLen(len(parts[nonceIndex])))
	// decode the base64 slice into the nonce array - don't care about number of bytes
	_, err := base64EncodingType.Decode(nonceBuf, parts[nonceIndex])
	if err != nil {
		return ErrStorageLayoutDecoding(err.Error())
	}
	sl.nonce = nonceBuf
	// make a temp slice to read things into
	valueBuf := make([]byte, base64EncodingType.DecodedLen(len(parts[valueIndex])))
	_, err = base64EncodingType.Decode(valueBuf, parts[valueIndex])
	if err != nil {
		return ErrStorageLayoutDecoding(err.Error())
	}
	sl.value = valueBuf
	return nil
}

// Value will return the current stored Value
func (sl *storageLayout) Value() []byte {
	return sl.value
}

// Nonce will return the current stored Nonce
func (sl *storageLayout) Nonce() []byte {
	return sl.nonce
}

// ErrStorageLayoutDecoding is returned when a general decoding error occurs
type ErrStorageLayoutDecoding string

// Error returns the formatted config error
func (sld ErrStorageLayoutDecoding) Error() string {
	return fmt.Sprintf("decoding failure: %s", string(sld))
}
