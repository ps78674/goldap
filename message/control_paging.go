package message

import (
	"fmt"
)

//
//        ControlTypePaging ::= SEQUENCE {
//             size            INTEGER (0..maxInt),
// 							   -- requested page size from client
// 							   -- result set size estimate from server
//             cookie          OCTET STRING
//        }
//

func (control *ControlPaging) PageSize() INTEGER {
	return control.pageSize
}

func (control *ControlPaging) Cookie() OCTETSTRING {
	return control.cookie
}

func (control *ControlPaging) SetPageSize(i INTEGER) {
	control.pageSize = i
}

func (control *ControlPaging) SetCookie(s OCTETSTRING) {
	control.cookie = s
}

func (control *ControlPaging) readComponents(bytes *Bytes) (err error) {
	control.pageSize, err = readINTEGER(bytes)
	if err != nil {
		err = LdapError{fmt.Sprintf("readComponents: %s", err.Error())}
		return
	}
	control.cookie, err = readOCTETSTRING(bytes)
	if err != nil {
		err = LdapError{fmt.Sprintf("readComponents: %s", err.Error())}
		return
	}
	return
}

func (control ControlPaging) write(bytes *Bytes) (size int) {
	size += control.cookie.write(bytes)
	size += control.pageSize.write(bytes)
	size += bytes.WriteTagAndLength(classUniversal, isCompound, tagSequence, size)
	return
}

func (control ControlPaging) size() (size int) {
	size += control.cookie.size()
	size += control.pageSize.size()
	size += sizeTagAndLength(tagSequence, size)
	return
}

func ReadControlPaging(bytes *Bytes) (cp ControlPaging, err error) {
	err = bytes.ReadSubBytes(classUniversal, tagSequence, cp.readComponents)
	if err != nil {
		err = LdapError{fmt.Sprintf("ReadControlPaging: %s", err.Error())}
		return
	}
	return
}

func (cp *ControlPaging) WriteControlPaging() (bytes *Bytes, err error) {
	// if cp.cookie == nil {
	// 	cp.cookie = OCTETSTRING("")
	// }

	// Compute the needed size
	totalSize := cp.size()
	// Initialize the structure
	bytes = &Bytes{
		bytes:  make([]byte, totalSize),
		offset: totalSize,
	}

	// Go !
	size := 0
	size += cp.write(bytes)
	// Check
	if size != totalSize || bytes.offset != 0 {
		err = LdapError{fmt.Sprintf("error writing message: size is %d instead of %d, final offset is %d instead of 0", size, totalSize, bytes.offset)}
	}
	return
}
