package message

import (
	"fmt"
)

//
//        SimplePagedResultsControl ::= SEQUENCE {
//             size            INTEGER (0..maxInt),
// 							   -- requested page size from client
// 							   -- result set size estimate from server
//             cookie          OCTET STRING
//        }
//

func (control *SimplePagedResultsControl) PageSize() INTEGER {
	return control.pageSize
}

func (control *SimplePagedResultsControl) Cookie() OCTETSTRING {
	return control.cookie
}

func (control *SimplePagedResultsControl) readComponents(bytes *Bytes) (err error) {
	control.pageSize, err = readINTEGER(bytes)
	if err != nil {
		err = LdapError{fmt.Sprintf("SimplePagedResultsControl.readComponents: %s", err.Error())}
		return
	}
	control.cookie, err = readOCTETSTRING(bytes)
	if err != nil {
		err = LdapError{fmt.Sprintf("SimplePagedResultsControl.readComponents: %s", err.Error())}
		return
	}
	return
}

func ReadPagedResultsControl(s *OCTETSTRING) (cp SimplePagedResultsControl, err error) {
	bytes := NewBytes(0, s.Bytes())
	err = bytes.ReadSubBytes(classUniversal, tagSequence, cp.readComponents)
	if err != nil {
		err = LdapError{fmt.Sprintf("ReadPagedResultsControl: %s", err.Error())}
		return
	}
	return
}

func (control *SimplePagedResultsControl) write(bytes *Bytes) (size int) {
	size += control.cookie.write(bytes)
	size += control.pageSize.write(bytes)
	size += bytes.WriteTagAndLength(classUniversal, isCompound, tagSequence, size)
	return
}

func (control *SimplePagedResultsControl) size() (size int) {
	size += control.cookie.size()
	size += control.pageSize.size()
	size += sizeTagAndLength(tagSequence, size)
	return
}

func WritePagedResultsControl(s INTEGER, c OCTETSTRING) (v *OCTETSTRING, err error) {
	// Compute the needed size
	control := SimplePagedResultsControl{s, c}
	totalSize := control.size()
	// Initialize the structure
	bytes := &Bytes{
		bytes:  make([]byte, totalSize),
		offset: totalSize,
	}

	// Go !
	size := 0
	size += control.write(bytes)
	// Check
	if size != totalSize || bytes.offset != 0 {
		err = LdapError{fmt.Sprintf("WritePagedResultsControl: size is %d instead of %d, final offset is %d instead of 0", size, totalSize, bytes.offset)}
	} else {
		v = OCTETSTRING(string(bytes.bytes[bytes.offset:])).Pointer()
	}
	return
}
