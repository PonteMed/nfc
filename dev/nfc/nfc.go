// This package wraps the libnfc to provide an API for Go. Most documentation
// was taken unchanged from the documentation inside the libnfc. Some functions
// and names have been altered to fit the conventions and idioms used in Go.
package nfc

import "fmt"

// Maximum length for an NFC connection string
const BUFSIZE_CONNSTRING = 1024

// Properties for (*Device).SetPropertyInt() and (*Device).SetPropertyBool().
const (
	// Default command processing timeout
	// Property value's (duration) unit is ms and 0 means no timeout (infinite).
	// Default value is set by driver layer
	TIMEOUT_COMMAND = iota

	// Timeout between ATR_REQ and ATR_RES
	// When the device is in initiator mode, a target is considered as mute
	// if no valid ATR_RES is received within this timeout value.
	// Default value for this property is 103 ms on PN53x based devices.
	TIMEOUT_ATR

	// Timeout value to give up reception from the target in case of no answer.
	// Default value for this property is 52 ms).
	TIMEOUT_COM

	// Let the PN53X chip handle the CRC bytes. This means that the chip
	// appends the CRC bytes to the frames that are transmitted. It will
	// parse the last bytes from received frames as incoming CRC bytes. They
	// will be verified against the used modulation and protocol. If a frame
	// is expected with incorrect CRC bytes this option should be disabled.
	// Example frames where this is useful are the ATQA and UID+BCC that are
	// transmitted without CRC bytes during the anti-collision phase of the
	// ISO14443-A protocol.
	HANDLE_CRC

	// Parity bits in the network layer of ISO14443-A are by default
	// generated and validated in the PN53X chip. This is a very convenient
	// feature. On certain times though it is useful to get full control of
	// the transmitted data. The proprietary MIFARE Classic protocol uses
	// for example custom (encrypted) parity bits. For interoperability it
	// is required to be completely compatible, including the arbitrary
	// parity bits. When this option is disabled, the functions to
	// communicating bits should be used.
	HANDLE_PARITY

	// This option can be used to enable or disable the electronic field of
	// the NFC device.
	ACTIVATE_FIELD

	// The internal CRYPTO1 co-processor can be used to transmit messages
	// encrypted. This option is automatically activated after a successful
	// MIFARE Classic authentication.
	ACTIVATE_CRYPTO1

	// The default configuration defines that the PN53X chip will try
	// indefinitely to invite a tag in the field to respond. This could be
	// desired when it is certain a tag will enter the field. On the other
	// hand, when this is uncertain, it will block the application. This
	// option could best be compared to the (NON)BLOCKING option used by
	// (socket)network programming.
	INFINITE_SELECT

	// If this option is enabled, frames that carry less than 4 bits are
	// allowed. According to the standards these frames should normally be
	// handles as invalid frames.
	ACCEPT_INVALID_FRAMES

	// If the NFC device should only listen to frames, it could be useful to
	// let it gather multiple frames in a sequence. They will be stored in
	// the internal FIFO of the PN53X chip. This could be retrieved by using
	// the receive data functions. Note that if the chip runs out of bytes
	// (FIFO = 64 bytes long), it will overwrite the first received frames,
	// so quick retrieving of the received data is desirable.
	ACCEPT_MULTIPLE_FRAMES

	// This option can be used to enable or disable the auto-switching mode
	// to ISO14443-4 is device is compliant.
	// In initiator mode, it means that NFC chip will send RATS
	// automatically when select and it will automatically poll for
	// ISO14443-4 card when ISO14443A is requested.
	// In target mode, with a NFC chip compliant (ie. PN532), the chip will
	// emulate a 14443-4 PICC using hardware capability.
	AUTO_ISO14443_4

	// Use automatic frames encapsulation and chaining.
	EASY_FRAMING

	// Force the chip to switch in ISO14443-A
	FORCE_ISO14443_A

	// Force the chip to switch in ISO14443-B
	FORCE_ISO14443_B

	// Force the chip to run at 106 kbps
	FORCE_SPEED_106
)

// NFC modulation types
const (
	ISO14443A = iota + 1
	JEWEL
	ISO14443B
	ISO14443BI   // pre-ISO14443B aka ISO/IEC 14443 B' or Type B'
	ISO14443B2SR // ISO14443-2B ST SRx
	ISO14443B2CT // ISO14443-2B ASK CTx
	FELICA
	DEP
)

// NFC baud rates. UNDEFINED is also a valid baud rate, albeit defined
// further below.
const (
	NBR_106 = iota + 1
	NBR_212
	NBR_424
	NBR_847
)

// NFC modes. An NFC device can either be a target or an initiator.
const (
	TARGET = iota
	INITIATOR
)

// NFC modulation structure. Use the supplied constants.
type Modulation struct {
	Type     int
	BaudRate int
}

// An error as reported by various methods of Device. If device returns an error
// that is not castable to Error, something outside on the Go side went wrong.
type Error int

// Returns the same strings as nfc_errstr except if the error is not among the
// known errors. Instead of reporting an "Unknown error", Error() will return
// something like "Error -123".
func (e Error) Error() string {
	if errorMessages[int(e)] == "" {
		return fmt.Sprintf("Error %d", int(e))
	}

	return errorMessages[int(e)]
}

// Error codes. Casted to Error, these yield all possible errors this package
// provides. Use nfc.Error(errorcode).Error() to get a descriptive string for an
// error code.
const (
	SUCCESS      = 0   // Success (no error)
	EIO          = -1  // Input / output error, device may not be usable anymore without re-open it
	EINVARG      = -2  // Invalid argument(s)
	EDEVNOTSUPP  = -3  // Operation not supported by device
	ENOTSUCHDEV  = -4  // No such device
	EOVFLOW      = -5  // Buffer overflow
	ETIMEOUT     = -6  // Operation timed out
	EOPABORTED   = -7  // Operation aborted (by user)
	ENOTIMPL     = -8  // Not (yet) implemented
	ETGRELEASED  = -10 // Target released
	ERFTRANS     = -20 // Error while RF transmission
	EMFCAUTHFAIL = -30 // MIFARE Classic: authentication failed
	ESOFT        = -80 // Software error (allocation, file/pipe creation, etc.)
	ECHIP        = -90 // Device's internal chip error
)

// replicate error messages here because the libnfc is incapable of giving
// direct access to the error strings. Stupidly, only the error string for the
// error code of an nfc_device can be read out.
var errorMessages = map[int]string{
	SUCCESS:      "Success",
	EIO:          "Input / Output Error",
	EINVARG:      "Invalid argument(s)",
	EDEVNOTSUPP:  "Not Supported by Device",
	ENOTSUCHDEV:  "No Such Device",
	EOVFLOW:      "Buffer Overflow",
	ETIMEOUT:     "Timeout",
	EOPABORTED:   "Operation Aborted",
	ENOTIMPL:     "Not (yet) Implemented",
	ETGRELEASED:  "Target Released",
	EMFCAUTHFAIL: "Mifare Authentication Failed",
	ERFTRANS:     "RF Transmission Error",
	ECHIP:        "Device's Internal Chip Error",
}

// the global library context
var theContext *context = &context{}