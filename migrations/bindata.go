// Package migrations Code generated by go-bindata. (@generated) DO NOT EDIT.
// sources:
// 1_create_table_consent_rule.down.sql
// 1_create_table_consent_rule.up.sql
// 2_alter_consent_record_add_version_uuid.down.sql
// 2_alter_consent_record_add_version_uuid.up.sql
// 3_rename_resource_to_data_class.down.sql
// 3_rename_resource_to_data_class.up.sql
// 4_alter_consent_record_make_valid_to_optional.down.sql
// 4_alter_consent_record_make_valid_to_optional.up.sql
// bindata.go
package migrations

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

// Name return file name
func (fi bindataFileInfo) Name() string {
	return fi.name
}

// Size return file size
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}

// Mode return file mode
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}

// ModTime return file modify time
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}

// IsDir return file whether a directory
func (fi bindataFileInfo) IsDir() bool {
	return fi.mode&os.ModeDir != 0
}

// Sys return file is sys mode
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var __1_create_table_consent_ruleDownSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x72\x09\xf2\x0f\x50\xf0\xf4\x73\x71\x8d\x50\x28\xcd\xcb\x2c\x8c\x2f\x4a\x2d\xce\x2f\x2d\x4a\x4e\xb5\xe6\x02\xcb\x84\x38\x3a\xf9\xb8\x2a\xa0\x09\xa2\x28\x4f\xce\x2f\x4a\x41\x51\x9c\x9c\x9f\x57\x9c\x9a\x57\x82\x2a\x85\xa4\xa5\x20\xb1\x24\x13\x24\x0f\x55\x87\xa2\x17\x43\x0e\x10\x00\x00\xff\xff\x55\xac\xed\x91\x9f\x00\x00\x00")

func _1_create_table_consent_ruleDownSqlBytes() ([]byte, error) {
	return bindataRead(
		__1_create_table_consent_ruleDownSql,
		"1_create_table_consent_rule.down.sql",
	)
}

func _1_create_table_consent_ruleDownSql() (*asset, error) {
	bytes, err := _1_create_table_consent_ruleDownSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "1_create_table_consent_rule.down.sql", size: 159, mode: os.FileMode(420), modTime: time.Unix(1569312499, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var __1_create_table_consent_ruleUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x84\x92\xcd\x6e\xab\x30\x10\x85\xf7\x3c\xc5\x2c\x41\x62\x75\xa5\xac\xb2\xf2\x85\x49\x8a\x4a\xec\xd6\x31\x55\xb3\x42\x14\x5c\xc5\x55\x8b\x53\x30\x55\xfb\xf6\x55\xf8\x29\xc1\xd0\x74\xed\x33\x87\x6f\xbe\x21\xe0\x48\x04\x82\x20\xff\x63\x84\x53\x66\x94\x2c\x4d\x9a\xeb\xb2\x96\xa5\x01\xd7\x01\x00\x50\x05\x3c\x10\x1e\xdc\x10\xee\xfe\x5b\xad\x3c\xb8\xe3\xd1\x8e\xf0\x03\xdc\xe2\xc1\x6f\x03\x75\xf3\xf4\x22\x73\x33\x4d\x51\x26\x80\x26\x71\xdc\x45\xf2\xa6\x36\xba\x50\x59\x79\x2d\x94\xe5\x46\x57\xcb\x01\xc7\x5b\x3b\x4e\xcf\x9a\xd0\xe8\x3e\x41\x88\x68\x88\x8f\xd0\x94\xea\x3d\xb5\xb9\x19\xb5\x57\x71\xdb\x6e\x7f\x40\xf5\x47\x20\x6f\x3d\xf4\x76\x85\xaa\xf8\xb4\xfb\xd2\x91\x7e\xa1\xf9\xb2\xc9\x99\xe8\x1c\xc6\x2b\x99\xeb\xaa\x18\x6d\x46\x54\xe0\x16\xf9\xa5\x48\x20\x89\x60\x11\x0d\x38\xee\x90\x8a\x4e\x87\x4d\x61\xdf\x81\xe3\x06\x39\xd2\x00\xf7\x33\x26\x55\x78\x5d\xc7\x47\xf6\xaa\x8a\xf4\xb9\xd2\x6f\x10\x9e\xc1\xa6\xc2\xbb\x57\xa3\x97\xde\x8e\x59\x7d\x5c\xbe\x45\x7f\x80\xeb\x27\xe9\x77\x66\xd4\xb2\xe0\x9e\x7b\x67\xa6\x2a\x59\xeb\xa6\xca\x65\xef\x68\x3a\x92\x8e\xca\x3a\xb4\x21\x9d\x9a\xaf\x93\xfc\xed\x87\x6a\x93\x1b\xc6\x31\xda\xd2\xb3\x61\x77\x56\xea\xb5\x91\x0b\x8b\xf6\xbd\x86\x08\xa3\x10\x62\x8c\x02\x21\x20\xfb\x80\x84\x7f\xae\xde\x2f\xc3\xe8\x0f\xea\xfc\xeb\xfe\x74\x0d\x6f\xed\x7c\x07\x00\x00\xff\xff\xc0\x0e\x26\xec\x8b\x03\x00\x00")

func _1_create_table_consent_ruleUpSqlBytes() ([]byte, error) {
	return bindataRead(
		__1_create_table_consent_ruleUpSql,
		"1_create_table_consent_rule.up.sql",
	)
}

func _1_create_table_consent_ruleUpSql() (*asset, error) {
	bytes, err := _1_create_table_consent_ruleUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "1_create_table_consent_rule.up.sql", size: 907, mode: os.FileMode(420), modTime: time.Unix(1575890957, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var __2_alter_consent_record_add_version_uuidDownSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x6c\x91\x41\x4f\x83\x40\x10\x85\xef\xfb\x2b\xde\x11\x12\x4e\x26\x3d\x71\x5a\x61\xaa\x1b\x61\xa9\xc3\x62\xec\x89\x90\x2e\xa6\x9b\x58\xa8\x80\xfd\xfd\x06\x6a\x4d\x25\x7b\xfe\xde\xcc\xbc\xf7\x26\xe5\x62\x07\xa5\x53\x7a\xc7\x77\xe7\xbe\xea\xa1\x3d\xf4\x83\xad\x2f\xed\x30\xba\xbe\x8b\x85\x90\x99\x21\x86\x91\x8f\x19\xe1\xd0\x77\x63\xdb\x4d\xbf\x22\x30\x69\x99\x13\x4c\xb1\x02\xf5\x74\x3a\xc7\x42\x24\x4c\xd2\x90\x7f\x34\x10\x00\xe0\x2c\x94\x36\xf4\x44\x8c\x1d\xab\x5c\xf2\x1e\x2f\xb4\x87\xac\x4c\xa1\x74\xc2\x94\x93\x36\xd1\xa2\x3c\x37\x93\x9b\xc7\x6f\x6b\x9c\xc5\x9b\xe4\xe4\x59\x72\xf0\xb0\xd9\x84\x60\xda\x12\x93\x4e\xa8\x5c\x4b\x03\x67\xc3\xeb\x8e\x4b\xf3\xe9\x6c\xfd\x31\xf4\x27\xa4\xb3\x31\x5d\x18\xe8\x2a\xcb\xee\xe9\xd4\xfb\xd8\xb1\x19\x8f\xff\xef\xdd\x38\x2a\xad\x5e\x2b\x12\x61\x2c\x84\xd2\x25\xb1\x99\x13\xad\x0b\x41\x49\x19\x25\x06\x81\xb3\x91\x27\x4a\x74\x67\x2d\xfa\x33\x12\x2d\x67\x43\x6c\xb9\xc8\xfd\x05\x2f\xbf\xf3\xd5\x7b\xe5\x3f\x01\x00\x00\xff\xff\xeb\xc2\xc4\xdc\xdb\x01\x00\x00")

func _2_alter_consent_record_add_version_uuidDownSqlBytes() ([]byte, error) {
	return bindataRead(
		__2_alter_consent_record_add_version_uuidDownSql,
		"2_alter_consent_record_add_version_uuid.down.sql",
	)
}

func _2_alter_consent_record_add_version_uuidDownSql() (*asset, error) {
	bytes, err := _2_alter_consent_record_add_version_uuidDownSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "2_alter_consent_record_add_version_uuid.down.sql", size: 475, mode: os.FileMode(420), modTime: time.Unix(1575890957, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var __2_alter_consent_record_add_version_uuidUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x8f\x31\x8b\x84\x30\x10\x85\x7b\x7f\xc5\x94\x0a\x36\x77\x60\x65\x35\x67\xe6\xee\x84\x6c\x64\x43\xb2\x6c\x27\xa2\x01\xd3\x24\x6e\x62\xfc\xfd\xcb\x2e\x6b\x63\x65\x39\x03\xef\x7b\xdf\x43\xae\x48\x82\xc2\x1f\x4e\x30\x7a\x17\x8d\x5b\xfb\x60\x46\x1f\x26\x40\xc6\xa0\xe9\xb8\xbe\x08\xd8\x4c\x88\xd6\x3b\x68\x85\xa2\x3f\x92\xc0\xe8\x17\x35\x57\xf0\x55\x67\xe7\x00\x29\xd9\x09\x6e\x28\x9b\x7f\x94\xf9\x77\x55\x15\x67\x83\x4b\x30\x9b\xf5\x29\xf6\xf3\x10\xe7\x03\x21\x6b\x24\xa1\x22\xd0\xa2\xbd\x6a\x82\x56\x30\xba\x43\x72\xf6\xf1\xe1\xf4\xbb\x76\x27\x0e\x0d\xf9\x32\xac\xf6\x75\xee\x6f\x3b\x95\x6f\xc7\x72\x9f\x5a\xd4\xcf\x00\x00\x00\xff\xff\x52\x46\xc5\x47\x1a\x01\x00\x00")

func _2_alter_consent_record_add_version_uuidUpSqlBytes() ([]byte, error) {
	return bindataRead(
		__2_alter_consent_record_add_version_uuidUpSql,
		"2_alter_consent_record_add_version_uuid.up.sql",
	)
}

func _2_alter_consent_record_add_version_uuidUpSql() (*asset, error) {
	bytes, err := _2_alter_consent_record_add_version_uuidUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "2_alter_consent_record_add_version_uuid.up.sql", size: 282, mode: os.FileMode(420), modTime: time.Unix(1575890957, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var __3_rename_resource_to_data_classDownSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x84\x91\xb1\x4e\x84\x40\x10\x86\xfb\x7d\x8a\x29\xb9\x84\xca\xe4\xaa\xab\xd6\xe5\xe7\x24\xe2\xac\x0e\x8b\xd1\x6a\x43\x96\x2d\x2e\x31\xa0\xc0\x15\xbe\xbd\x41\xf1\xc2\x49\xe1\x94\xbb\xff\x4c\xbe\x6f\xc6\x08\xb4\x03\x39\x7d\x5b\x82\x86\x38\xf6\xe7\x21\x44\x4a\x14\x7d\x57\xe8\xbb\x31\x76\x93\x1f\x62\xe8\x87\xd6\x9f\x5a\x2a\xd8\xe1\x08\x49\x97\xc0\x6f\x87\x9f\x3e\xdf\x23\x3d\x6b\x31\x77\x5a\x92\x9b\xfd\x7e\x47\x6c\x1d\x71\x5d\x96\xa9\x5a\xb2\xb9\x15\x14\x47\xa6\x7b\xbc\x26\x9b\xc1\xbb\x25\x34\x97\x20\x87\x80\x0d\xaa\x3f\x00\x94\x5c\x07\x2d\x53\x86\x12\x0e\x64\x74\x65\x74\x06\xb5\x3b\x28\xb5\x28\xd5\x5c\x3c\xd5\xa0\x82\x33\xbc\xd0\xb9\x3b\x7d\xf8\x8b\x9e\xe5\x0b\xf8\x96\x24\xbd\x96\x9a\x27\x16\x5c\x41\xdc\xec\x6e\x57\x3b\xfa\xb7\x93\x2a\x94\x30\x6e\xbb\xc5\x94\x42\xdf\x46\xca\xc5\x3e\x50\xdb\x4c\x8d\x0f\x6f\xcd\x38\x1e\x94\xca\xc4\x3e\xae\x89\xb7\x9f\x3f\x87\x5a\xbd\x7f\x05\x00\x00\xff\xff\xf9\x0e\x8e\x7a\xc1\x01\x00\x00")

func _3_rename_resource_to_data_classDownSqlBytes() ([]byte, error) {
	return bindataRead(
		__3_rename_resource_to_data_classDownSql,
		"3_rename_resource_to_data_class.down.sql",
	)
}

func _3_rename_resource_to_data_classDownSql() (*asset, error) {
	bytes, err := _3_rename_resource_to_data_classDownSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "3_rename_resource_to_data_class.down.sql", size: 449, mode: os.FileMode(420), modTime: time.Unix(1575894131, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var __3_rename_resource_to_data_classUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x74\x91\xc1\x4e\x84\x30\x10\x86\xef\x7d\x8a\x39\xb2\x09\xf1\x60\xb2\xa7\x3d\xd5\xf2\xb3\x12\x71\xaa\x43\x31\x7a\x6a\x08\xf4\xb0\x89\x01\x05\xf6\xe0\xdb\x1b\x56\x24\xec\xae\xce\xa9\xe9\xff\x75\x3a\x5f\x6b\x04\xda\x81\x9c\xbe\xcb\x41\x4d\x35\x56\xbe\x7e\xaf\x86\x81\x22\x45\xa7\xaa\xbb\x76\x08\xed\xe8\xfb\x50\x77\x7d\xe3\x0f\x0d\x65\xec\xb0\x87\xc4\x0b\xd0\x04\x7a\xd1\x62\xee\xb5\x44\xb7\xdb\xed\x86\xd8\x3a\xe2\x32\xcf\x63\x35\x23\xa9\x15\x64\x7b\xa6\x07\xbc\x45\x57\xfd\x36\x33\x34\x95\x20\x85\x80\x0d\x8a\x8b\x7b\x29\x3a\x07\x2d\x53\x82\x1c\x0e\x64\x74\x61\x74\x02\xb5\xd9\x29\x35\xbb\x94\x9c\x3d\x97\xa0\x8c\x13\xbc\xd2\xb1\x3d\x7c\xfa\x95\x97\xe5\x95\xe5\xf5\x34\xf1\xc9\x67\x6a\x96\x71\x01\x71\x93\xad\x3d\x7b\x97\xff\x8e\x50\x81\x1c\xc6\x51\x1f\x86\xee\xd8\xd7\xe1\xe6\x0f\x70\xc9\x7e\x17\x7e\xfc\xfa\x08\x94\x8a\x7d\x5c\xb2\x9d\x52\x89\xd8\xa7\xf5\xf8\x97\xd1\xcf\x67\x2d\xbb\xdf\x01\x00\x00\xff\xff\x7e\x08\x1a\x88\xc3\x01\x00\x00")

func _3_rename_resource_to_data_classUpSqlBytes() ([]byte, error) {
	return bindataRead(
		__3_rename_resource_to_data_classUpSql,
		"3_rename_resource_to_data_class.up.sql",
	)
}

func _3_rename_resource_to_data_classUpSql() (*asset, error) {
	bytes, err := _3_rename_resource_to_data_classUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "3_rename_resource_to_data_class.up.sql", size: 451, mode: os.FileMode(420), modTime: time.Unix(1575894131, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var __4_alter_consent_record_make_valid_to_optionalDownSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x74\x52\xcd\x6e\xf2\x30\x10\xbc\xfb\x29\xf6\x18\xa4\x5c\xbe\x4f\xe2\xc4\xc9\x4d\x96\xd6\xaa\x63\xd3\xc5\xa9\xca\x29\x42\x24\x15\x96\x4a\x4c\xf3\xf7\xfc\x55\x7e\x29\x21\x3d\xcf\xee\xcc\xec\xec\x84\xa4\x77\x20\x54\x88\x1f\x50\xe7\xf6\x3b\x29\xb2\x93\x2b\xd2\xa4\xc9\x8a\xd2\xba\x7c\xc3\x18\x97\x06\x09\x0c\x7f\x92\x08\x27\x97\x97\x59\x5e\x0d\x43\x40\xa8\x78\x84\x60\xf4\x0c\x48\xaa\xcb\x75\xc3\x58\x40\xc8\x0d\x2e\xaf\x7a\x0c\x00\xc0\xa6\x20\x94\xc1\x67\x24\xd8\x91\x88\x38\x1d\xe0\x15\x0f\xc0\x63\xa3\x85\x0a\x08\x23\x54\xc6\xef\x26\xaf\xc7\xca\xb6\xeb\x23\x8d\x4d\xe1\x9d\x53\xf0\xc2\xc9\xfb\xbf\x5e\xaf\x80\x70\x8b\x84\x2a\xc0\xfd\x7c\xd4\xb3\xe9\xaa\xe7\x68\x8e\x5f\x36\x4d\x3e\x0b\x77\x81\xb0\x35\xa6\xb4\x01\x15\x4b\xf9\x1b\xad\xdc\x12\x76\x3e\x96\xe7\x7b\xbd\x11\x87\x58\x89\xb7\x18\x07\x8a\x3e\xb4\xe9\xa6\x10\xb7\x3c\x96\x06\xfe\xf5\x70\x5d\xcf\x5c\x0f\xa7\x15\x59\x63\x5d\x5d\x26\x0f\x2a\x6c\x75\x4b\xb1\xd7\xf9\xfb\x53\xa0\xd5\x2c\x63\xef\x31\x33\xbf\xf3\xe0\x8f\x46\x5b\x7a\xa1\xf6\x48\xa6\xb5\x3c\xff\x22\xec\x51\x62\x60\xa0\x5d\x58\xa2\xba\xc5\xe9\x4f\xe1\xf9\x5d\x54\x93\xc0\xa8\x77\x7f\xe2\x96\x74\xb4\xdc\x98\xae\x8c\x4b\x7d\xe9\xf1\x9f\x00\x00\x00\xff\xff\xf1\x4b\x02\xf5\xac\x02\x00\x00")

func _4_alter_consent_record_make_valid_to_optionalDownSqlBytes() ([]byte, error) {
	return bindataRead(
		__4_alter_consent_record_make_valid_to_optionalDownSql,
		"4_alter_consent_record_make_valid_to_optional.down.sql",
	)
}

func _4_alter_consent_record_make_valid_to_optionalDownSql() (*asset, error) {
	bytes, err := _4_alter_consent_record_make_valid_to_optionalDownSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "4_alter_consent_record_make_valid_to_optional.down.sql", size: 684, mode: os.FileMode(420), modTime: time.Unix(1576581514, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var __4_alter_consent_record_make_valid_to_optionalUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x74\x52\xcb\x6e\xea\x30\x10\xdd\xfb\x2b\x66\x19\xa4\x6c\xee\x95\x58\xb1\x72\x93\xa1\xb5\xea\xd8\x74\x70\xaa\xb2\x8a\x10\x49\x85\xa5\x12\xd3\xbc\xbe\xbf\xca\x93\x12\xd2\xad\xcf\x9c\xc7\x1c\x4f\x48\x7a\x07\x42\x85\xf8\x01\x75\x6e\xbf\x93\x22\x3b\xb9\x22\x4d\x9a\xac\x28\xad\xcb\x37\x8c\x71\x69\x90\xc0\xf0\x27\x89\x70\x72\x79\x99\xe5\xd5\x30\x04\x84\x8a\x47\x08\x46\xcf\x80\xa4\xba\x5c\x37\x8c\x05\x84\xdc\xe0\x32\xd5\x63\x00\x00\x36\x05\xa1\x0c\x3e\x23\xc1\x8e\x44\xc4\xe9\x00\xaf\x78\x00\x1e\x1b\x2d\x54\x40\x18\xa1\x32\x7e\x37\x79\x3d\x56\xb6\xa5\x8f\x32\x36\x85\x77\x4e\xc1\x0b\x27\xef\xff\x7a\xbd\x02\xc2\x2d\x12\xaa\x00\xf7\xf3\x51\xcf\xa6\xab\x5e\xa3\x39\x7e\xd9\x34\xf9\x2c\xdc\x05\xc2\x36\x98\xd2\x06\x54\x2c\xe5\x6f\xb4\x72\x03\x36\xbd\x9f\x8f\xe5\xf9\xde\x6b\xe4\x41\xac\xc4\x5b\x8c\x03\xbd\x2f\x6c\xda\x27\xc4\x2d\x8f\xa5\x81\x7f\x3d\x5c\xd7\xb3\xc4\xc3\x5a\x45\xd6\x58\x57\x97\xc9\x83\x0b\x5b\xdd\x1a\xec\x7d\xfe\xfe\x25\xd0\x6a\xd6\xaf\xf7\xd8\x97\xdf\x65\xf0\xc7\xa0\xad\xbc\x50\x7b\x24\xd3\x46\x9e\xff\x20\xec\x51\x62\x60\xa0\x25\x2c\x49\xdd\xaa\xf4\xa7\xe2\xfc\xae\xaa\xc9\x60\xf4\xbb\x5f\x71\x4b\x3a\x5a\xbe\x96\xee\x10\x97\x6e\xa5\xc7\x7f\x02\x00\x00\xff\xff\xf8\xbf\x7f\x22\xa8\x02\x00\x00")

func _4_alter_consent_record_make_valid_to_optionalUpSqlBytes() ([]byte, error) {
	return bindataRead(
		__4_alter_consent_record_make_valid_to_optionalUpSql,
		"4_alter_consent_record_make_valid_to_optional.up.sql",
	)
}

func _4_alter_consent_record_make_valid_to_optionalUpSql() (*asset, error) {
	bytes, err := _4_alter_consent_record_make_valid_to_optionalUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "4_alter_consent_record_make_valid_to_optional.up.sql", size: 680, mode: os.FileMode(420), modTime: time.Unix(1576581499, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _bindataGo = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xd4\x9a\x5b\x6f\x23\xd7\x91\xc7\x9f\xc9\x4f\xd1\x11\x10\x83\x5c\x68\xa5\xbe\x5f\x04\x0c\xb0\x88\xed\x05\xfc\xb0\xce\x62\xed\x3c\xed\x59\x10\x7d\x39\x2d\x13\x91\x48\x99\xa4\xec\xa3\x19\xcc\x77\x5f\xfc\xaa\xaa\x75\x9b\xc9\x0c\xa5\x99\x24\xc8\x83\x44\xb2\xd9\x5d\xa7\x4e\x5d\xfe\xff\xaa\x3a\x3c\x3f\x8f\xfe\xbb\xed\xff\xda\x5e\xfa\xe8\x7a\x7d\xb9\x6b\x0f\xeb\xed\x66\x1f\x7d\xbb\x1d\x7c\x74\xe9\x37\x7e\xd7\x1e\xfc\x10\x75\x77\xd1\xe5\xf6\xdf\xbb\xf5\x66\x68\x0f\xed\x59\xb4\xf8\x8f\xfb\xaf\x96\xd1\x77\x7f\x8e\x7e\xfc\xf3\xcf\xd1\xf7\xdf\xfd\xf0\xf3\xd9\xfc\xfc\x3c\xda\x6f\x6f\x77\xbd\xdf\x5f\xf0\x3e\x59\xf5\x3b\xdf\x1e\xfc\xea\xd0\x76\x57\x7e\xd5\x6f\x37\x7b\xbf\x39\xac\x76\xb7\x57\xfe\x6c\xd8\xfe\xbe\x39\xdb\xff\x7a\xf5\xb9\xfb\x6e\x6f\xa6\xbb\xd2\x55\x7b\x75\xf0\xbb\x87\xaf\x7d\xbf\xdd\x0d\xab\x76\x18\x56\xbf\xf9\xdd\x7e\xbd\xdd\xac\x6e\x6f\xd7\xc3\x13\xc9\xc7\x3e\xf3\xb0\x4a\xb6\xda\xf9\x4d\x7b\xed\x57\x3b\xaf\x5b\x59\x1d\xb6\x2b\xf6\xbd\xea\xaf\xda\xfd\xfe\x89\xf4\xcf\xdd\xfb\x20\x35\xff\xb8\x1e\xd7\xed\x5f\xfd\xea\xb7\xf6\x6a\x3d\xf0\xe4\xf6\x06\xeb\xb7\x57\x4f\xd6\x78\xd9\x93\x0f\x2b\x4e\xde\xba\xdc\xce\x6f\x3e\x70\xf0\x7c\xbe\xbe\xbe\xd9\xee\x0e\xd1\x62\x3e\x3b\xe9\xee\x0e\x7e\x7f\x32\x9f\x9d\xf4\xdb\xeb\x9b\x9d\xdf\xef\xcf\x2f\xdf\xae\x6f\xb8\x30\x5e\x1f\x78\x59\x6f\xf5\xff\xf9\x7a\x7b\x7b\x58\x5f\xf1\x61\x2b\x0f\xdc\xb4\x87\x5f\xce\xc7\xf5\x95\xe7\x0d\x17\xf6\x87\xdd\x7a\x73\x29\xdf\x1d\xd6\xd7\xfe\x64\xbe\x9c\xcf\xc7\xdb\x4d\x3f\x69\xf3\x3f\xbe\x1d\x16\xbc\x89\xfe\xf7\xff\x58\xf6\x34\xc2\x7c\x91\x3e\xb6\x8c\x16\xd3\x55\xbf\xdb\x6d\x77\xcb\xe8\xdd\x7c\x76\xf9\x56\x3e\x45\x17\x6f\x22\xb4\x3a\xfb\xd1\xff\x8e\x10\xbf\x5b\x88\xda\x7c\xfe\xd3\xed\x38\xfa\x9d\x88\x5d\x2e\xe7\xb3\xf5\x28\x0f\xfc\xe1\x4d\xb4\x59\x5f\x21\x62\xb6\xf3\x87\xdb\xdd\x86\x8f\xa7\xd1\x78\x7d\x38\xfb\x1e\xe9\xe3\xe2\x64\xe7\xdb\x21\xfa\xe3\xaf\x17\xd1\x1f\x7f\x3b\x51\x4d\x64\xad\xe5\x7c\xf6\x7e\x3e\x9f\xfd\xd6\xee\xa2\xee\x76\x8c\x74\x1d\x5d\x64\x3e\x5b\xa9\x3a\x6f\xa2\xf5\xf6\xec\xdb\xed\xcd\xdd\xe2\x9b\xee\x76\x3c\x8d\x2e\xdf\x2e\xe7\xb3\xfe\xea\xfb\x49\xd3\xb3\x6f\xaf\xb6\x7b\xbf\x58\xce\xbf\x96\x3e\x88\x51\xf9\x7f\x43\x90\xdf\xed\x54\x6f\xbb\xd8\xdd\x8e\x67\x7f\x42\xf5\xc5\xf2\x94\x3b\xe6\xef\xe7\xf3\xc3\xdd\x8d\x8f\xda\xfd\xde\x1f\x30\xf9\x6d\x7f\x40\x8a\xec\xcf\xfc\x31\x9f\xad\x37\xe3\x36\x8a\xb6\xfb\xb3\xff\x5c\x5f\xf9\x1f\x36\xe3\xf6\xfe\x39\x73\xe1\x74\xfd\x91\x04\xf1\x61\x14\x99\x1b\xe7\xb3\xfd\xfa\xad\x7c\x5e\x6f\x0e\x65\x3e\x9f\x5d\x03\x26\xd1\xbd\xd0\xff\xda\x0e\x5e\x2e\xfe\xbc\xbe\xf6\x11\x61\x72\xc6\x3b\xd6\x39\x3f\x8f\x7e\x44\x96\x6d\x81\xc8\x12\x33\x68\x0c\x2d\xc6\xf5\x73\x25\x96\x72\xff\x62\x69\x4b\xa3\xcc\xfd\xb3\x67\xf2\xa4\x4a\xfd\x09\x8d\x1e\x4b\x45\xc5\x4f\x48\xe5\xfe\xc5\x52\x37\xf0\x54\xa8\x3c\xa8\x42\xd9\xc8\x13\xa1\x6c\xf4\x13\x42\xb9\x7f\xb1\x7c\x6c\x86\xa7\xa2\xe5\xf1\x7b\xd1\x62\x9d\x67\xd2\xd7\xe3\x9d\x18\xec\xd3\x8b\xf0\xe4\x62\xf9\x60\xd9\x0f\x56\x79\x64\xee\x1f\xf6\xdf\xad\x77\x4f\x96\xf9\xfd\x17\x7f\xf8\xc5\xef\xa2\x36\x1a\xd6\x3b\xdf\x1f\xb6\xbb\xbb\x4f\x2c\x27\xcf\x2f\x96\x51\xb7\xdd\x5e\x7d\xb8\x9b\x6f\xb6\xfb\x33\xf6\xc9\x1a\x7f\x78\x13\xc5\x93\x37\xee\xf6\x4f\x96\x5c\xef\xa3\xfd\xdd\xfe\x73\xe6\xfb\xe9\x6e\xaf\x2e\xf1\xbb\xb1\xed\xfd\xbb\xf7\x8f\xd6\xb3\xf8\x26\x65\x57\xab\x4f\xb0\xc9\x77\xdb\xdf\x37\x3f\xfd\x7a\x15\xbd\xb1\x78\x5f\x9c\xb8\x90\x8c\x2e\xd4\x9d\x0b\x71\xed\x42\x1c\x7f\xfc\x6f\x1c\x5d\xa8\x52\x17\xe2\xc6\x85\x91\xd7\xd1\x85\x82\xeb\xfc\xe5\x2e\x54\x99\x0b\x55\xe2\x42\x3d\xe8\xf5\xb4\x76\xa1\x1f\x5c\xe8\x3b\x17\xd2\xde\x85\xba\x77\x21\x1d\x5d\xc8\x5b\x17\x52\xae\x7b\xfd\xcc\x7b\xae\xe5\xde\x85\xae\x70\xc1\x97\x2e\xc4\xa9\x3e\x57\xe7\x2e\x64\xb5\x0b\x59\xeb\xc2\xd8\xb8\xd0\xd5\x2e\xa4\xad\x0b\x6d\xac\x7a\xb4\xa9\xae\x93\x8f\x0f\xf2\x44\x56\xe2\x42\x91\xb8\xd0\xf4\xf6\x87\xae\x95\xbd\x6f\xf5\x7d\x9d\xaa\xac\xba\x70\xa1\xcd\x5d\x68\x0b\x17\xd2\xd8\x85\x2e\x71\x21\xcd\x5d\x48\x32\x7d\x95\x7d\x16\x2e\xd4\x95\xae\x97\x54\x2e\xe4\x99\x0b\xb1\x77\x21\x79\x66\x9f\xd1\xee\x6d\x7b\x17\xfc\xe0\x42\x93\xe8\xda\x0f\x76\x3c\x99\xd8\xe0\x08\x17\x19\x74\x7d\x8c\x12\x26\x80\x7b\x44\x29\xf3\xd9\xec\x18\xbf\x9f\xce\x67\xb3\x93\x63\x8a\x92\x93\xd3\xf9\x6c\x49\x44\x1d\xab\x2e\x9a\xfe\x9b\x20\xeb\x63\x4d\x05\x5a\xef\xf9\xeb\xf8\x5d\x7f\x8e\x34\xee\xb1\x5e\xd0\xfa\xe2\xcd\xf3\x64\x79\x07\xf4\x5d\x44\x47\x6e\x55\xc0\xf0\x22\x4a\x8a\xe6\x54\x72\xf0\xe2\x31\x44\x2d\xf2\x34\x5e\xca\x75\x50\xe3\x42\x51\xe5\x2f\x9b\x75\x58\x24\x45\xd9\x64\x49\x9a\x37\xcd\x69\x14\x2f\xdf\xcf\x67\x2d\x9a\x7c\x23\x46\x78\x27\x3b\xbf\x88\xcc\x00\xa8\x79\x21\xff\xdf\xdf\xbb\xaf\x3d\x3d\x3a\x6b\xff\x72\xf3\xda\x9c\x25\x87\x9a\x54\x73\xb1\xf4\x2e\xb4\x9d\x0b\x59\xac\xb1\x4b\xec\x8f\x95\x0b\x59\xef\x42\x5f\x68\x9e\x92\x3b\x65\xea\x42\x55\x68\x4e\x10\xcb\x5d\xaa\x39\xcf\xfd\x79\xe3\x42\xdd\x6a\x9e\xf9\xde\x85\xa1\x74\x21\x4b\x34\xee\xbb\xcc\x85\x9c\x1c\xc9\x5d\x28\x4c\xa6\xe4\x4e\xe7\x42\x91\xe9\xba\x7c\x1e\x3b\x17\xc6\xd2\xde\x93\xd3\x8d\x0b\x7d\xe2\xc2\x10\xbb\x50\xe5\x9a\x3b\x59\xa6\x39\x57\x8e\x2e\x74\xe4\x76\xe2\x82\x8f\x5d\xc8\xd9\x6f\x6e\xf9\x6b\x7b\x2c\x33\xdd\x27\x6b\x94\xa5\x0b\x4d\x6e\x7b\x19\x34\xdf\x7d\xa7\x7b\x68\x4a\xdd\x53\xcc\x5a\x95\xbe\x62\x27\xf0\x2a\x2e\xd4\x0e\xd8\x25\xf1\x2e\x0c\xf6\xde\x7b\x17\x46\xef\x42\xd1\xb9\xd0\x0e\x7a\x0f\x38\xe4\x33\x17\x06\xf0\xce\x2b\x06\xc6\x99\x3e\xe3\x53\xdd\x0b\x7a\x73\x0d\x3b\x8e\x99\x62\x64\x9a\x2a\x4e\xb2\x37\x74\x03\xa3\xd2\xd2\x85\x3a\xd6\x57\xf0\x13\x19\x79\xa1\xf6\x6e\xb1\x6d\xe9\x42\xd7\xaa\x8e\x45\xe3\x42\xd5\x28\x66\xb2\x47\x8f\x3f\x4a\xc5\x33\xf0\x92\xfd\xf4\x95\xea\x9a\x75\x8a\xa9\xfd\xa8\x36\xc0\xb6\x1e\x2c\xf5\xea\xe3\xba\x76\xa1\xac\xf5\xb5\x1e\xf5\x7b\x91\xd9\xba\x90\x0d\x8a\xc5\x5d\xe3\x42\xd2\xe8\x7b\xd6\x10\xfd\x3a\x8d\xa3\x0a\xac\xc5\xfe\xe8\x50\xa9\x2f\xd8\x33\xcf\x0a\x1f\x94\x86\xad\xad\xfa\xb8\xcb\xd5\xef\x43\xaa\xb8\x58\x11\x8b\x89\xe1\x3a\x36\x6b\x5c\x68\x1a\xd5\x31\x21\x76\x4a\x17\xd2\xce\xae\x75\x2a\x27\x41\xe7\xc1\xf6\x9c\x6b\x3c\x24\xa5\xca\xc0\xaf\xc4\x06\x7a\xd4\x8d\x0b\x25\xdf\xe1\xdf\x41\x39\x04\x3f\x36\xb1\xc5\xae\xd7\xd8\xc2\x07\x25\x7b\x22\x4f\x12\xf5\x69\x8c\xcd\x1b\xd5\x53\x72\xa8\x52\x7f\xa5\x16\xaf\x55\xed\x42\x31\xa8\x9d\xd9\x37\x71\x8c\x6e\xc8\x65\xdf\x9d\x3d\x8b\x2d\x88\x9f\xc6\x6b\x3c\xe0\xcb\x3e\x75\xa1\xeb\xd4\x96\x6d\xe6\x42\x53\xb9\x30\x78\x8d\x21\xf1\xed\xa0\x79\x43\xbc\x13\x03\xd8\x98\xd8\x60\xff\x69\xe5\x82\xe7\x9e\x4a\x63\x7c\xc8\xd5\x6e\xde\xd6\xa8\xf0\x4b\xa5\xeb\xc0\x6b\xc5\x64\x37\xd6\x6d\x5d\x28\xe1\xd6\x51\x7d\x0e\x9f\x82\x09\xac\xcb\x77\x7c\xc6\x3f\xe4\x58\x63\x39\xd3\x12\x37\xc4\xad\x71\x99\xe4\x62\xab\xd7\x92\x4e\xfd\x43\xde\x0f\xad\xee\x97\xdc\xc7\x96\xe8\x57\x94\x1a\x47\xc4\x0d\xbe\x8e\x3b\xc5\x00\x6c\xd5\x21\xab\x54\xdc\xc2\x06\xd8\x08\xbc\xa1\x3e\x80\xf7\xd1\x41\xf2\xba\xd3\xbd\x93\xd7\xd8\xa1\xf5\x6a\x2b\x38\xbe\xcf\x34\x4e\xe0\x66\xd6\x11\x1d\x3b\xcd\x53\xf0\x03\x9f\x63\x7f\xf4\xae\x90\x5b\x7d\xc8\xd3\x7d\xac\x1c\x8e\x5f\xc1\x31\xc1\xd2\xec\x45\x3c\x2d\xa0\xfc\x75\x59\x5a\x44\x7e\x96\xa3\xb5\xc5\x3d\x9e\xa1\x45\xea\x97\xf2\xf3\xe3\xdd\xfe\x23\xd8\x79\xda\xa4\x71\x73\x13\x57\x2f\xe4\xe6\xaa\xa8\x9b\xb8\x29\xaa\xaf\xc2\xcd\x47\x4e\x51\xbe\xa4\xba\x2e\x7b\xcd\x16\xd0\x99\x8a\xb6\xce\x14\x61\x27\xa6\x26\x7b\xc9\x0a\x50\x91\x4c\x00\xdd\x92\x54\xd1\x8c\x28\x06\x79\x41\xe8\xa2\xd5\x4c\x14\xc4\xec\xec\x7d\xa3\x59\x43\xa6\x49\xb4\x37\x8a\x86\xa9\x57\x84\x68\x40\x25\x32\xb2\xd6\xac\x1b\x07\xfd\x03\x0d\xc9\x7a\x90\x32\x2d\x14\x65\x46\xcb\xc4\x1e\xa4\xea\x15\x21\x25\x8b\x0a\x95\x4f\xb6\x81\xc6\xb0\x71\xd5\x2a\x52\x82\x0a\xbe\x52\x64\x23\x63\x41\xfe\x89\x29\xd8\x27\xcc\x9a\x5a\xc6\x52\x25\xc0\x78\xdc\x8b\xdd\xd8\x3b\xba\xc2\x04\xa0\x03\xe8\x81\x9d\x60\x2d\xd8\xc9\x4f\x95\x43\xa5\x95\x00\x0c\x85\xbe\xdd\xa8\x8c\x8b\xbc\xd2\xd8\x85\xaa\x3e\xef\xb5\xca\x07\x6d\x60\x2f\x10\x23\x33\x3d\xa9\x5e\xa8\xfa\xb9\x47\x98\x2a\x56\xf4\xc9\xf2\x87\x6a\x1f\xc4\xa5\xba\x80\x2d\x61\x67\xf6\x90\xe7\x8a\x60\xc9\xa0\x15\x16\x28\x0e\x83\x27\x86\x58\xa0\xab\x74\x10\xbd\xca\x66\xff\xac\x0b\x23\x20\x07\x54\x45\x16\x15\x05\x5d\x06\x55\x46\x56\xe9\x75\x7c\xd3\x5b\x25\x54\x76\xda\xc5\x50\x59\x71\xdd\xe7\xfa\x27\xec\x91\x6a\x25\xd2\xc5\x2e\x0c\x8d\x22\x27\x2c\x08\x42\x13\x27\xc8\x22\x56\xf0\x31\xfa\xe5\x86\x78\xb0\x86\x20\x6a\xa7\x8c\xc0\x75\x2a\x16\xd8\x86\x98\x24\x0e\x40\x7a\xa9\x60\x2a\xed\x96\x40\x7b\xa9\xfa\x06\x65\x65\x6f\x9d\x19\x55\x08\x4c\xcb\xb3\xc3\xc4\xf8\xb5\xda\x1b\x5f\xd5\x86\xbe\xc4\xf2\x60\xcc\x2c\x1d\x1d\xcc\xe7\x35\xb6\xd1\x95\x98\x95\xae\x31\x57\x3f\x10\x7b\xc8\x98\x7c\xc8\xfd\xb0\x8a\xe4\x8a\x55\x29\xdc\x43\xbc\xc2\xe4\xe8\x47\x8c\xa0\x2f\x95\x2a\xb6\x66\x9f\x54\x4e\x63\xab\x8c\xce\x3a\x7c\x2e\xad\xa3\x63\xaf\xb0\x77\x5f\x5b\xec\x17\xe6\xbb\x51\xed\x31\x58\xfc\x13\xe7\xd9\xf8\x50\x3d\x3e\x66\x15\x6c\x88\x4f\xfb\x5c\x2b\xb9\xa1\x7b\x7c\xdf\x03\xab\xbc\x0c\x4e\x5e\xc1\x31\x2f\x5b\x40\x18\xe7\xa5\xc3\xe5\x67\xfc\xf3\xb2\x15\x8f\x62\xa3\x57\x59\xe9\x6b\x71\xd3\xcb\xcd\x61\x4c\x95\x57\xc5\xbf\x00\x53\x7d\x41\x47\xd9\x6b\x1e\x93\xff\xf5\x34\xa5\x79\xc4\x53\xe4\x08\x98\x29\x18\x45\x35\xd6\x2a\xb6\x09\x3e\xc7\x5a\x89\x66\x85\xe6\x9d\x2f\x15\x2f\x04\xab\x7a\x17\xca\x5c\x73\x91\xca\x96\xcf\xa9\x4d\x5d\xa4\x63\xb3\x89\x0c\xdd\x07\x3c\x33\x1a\x36\xc9\x84\xc9\x2b\x36\x82\xff\x22\xbf\x51\x6c\x03\x67\xd0\x87\x0a\x1f\xb9\x54\x91\x74\x09\x74\x8f\xe4\x29\x7b\x06\x13\xd1\x1f\xae\x4a\x2a\x9b\x60\x75\x8a\x5d\xe8\x4b\xb7\xc1\x7d\x70\x1c\x7c\x4c\xe5\xdb\xc6\x8a\x71\x74\x82\xf0\x14\x36\x03\xe7\xc0\x76\x3a\x2a\xe9\x8a\x3b\xeb\xb0\x0a\xc3\xf4\xd1\x3a\x71\xeb\xc4\x58\x2b\xb3\xce\x6a\xb4\xce\xb8\x34\x9e\xc4\xd6\x74\xc5\xe0\x78\xdc\xe8\x9e\xd3\x5a\x31\x1f\xdb\x62\xd7\xd1\xba\x01\x9e\x4b\xec\x59\xa9\x19\xac\xc3\xe7\x5e\xf8\x0d\x39\x54\xde\xa3\x55\xdb\x22\x3f\x53\x2e\xc5\x66\xd8\x14\x3e\x82\x2b\xe1\x51\x74\xa5\x4a\x87\xf7\xb1\x13\xd5\xfc\xc4\xc9\xd8\x50\x78\x06\x79\x89\x72\xaf\x74\x3d\x9d\x76\x7c\xf8\x8b\xaa\x9a\x0a\x1c\xfd\xb2\x54\x39\x8f\xfb\xe9\x84\xf1\xb5\x74\x89\xf0\x58\xa1\xef\xa5\x3e\x48\x6d\x4a\xd7\x2a\x6f\xf4\xe3\x87\x71\x27\xd3\xb5\x54\xfd\xd1\x5b\xd7\x99\xb4\x5f\x84\xaf\xaf\xad\xe0\x5f\x22\xfe\x45\xd8\xfa\xd1\xca\xfe\x25\xab\x7d\x4d\x5c\xfd\x7b\x54\xfc\x2f\x35\x84\x61\x6a\x5a\xa7\xff\x4c\x4c\xfd\xcc\x29\xe7\x97\x54\xfd\x32\x9f\x4b\xb4\xba\xc9\x0d\x09\xef\xab\xfe\x52\x51\xa8\x1a\x74\xbe\x40\x26\x77\x56\xd9\xd1\xb7\x53\xf9\x51\xe5\x53\x71\x0e\x56\x81\x93\xdd\x64\xb4\xb7\xcc\x8b\xad\x8a\xa6\xb2\x2c\x2d\x83\x9b\x52\x2b\x1f\x50\x13\x14\x07\xcd\x40\x25\x50\x04\x34\x23\xfb\xa9\xac\xc8\x70\x50\x32\xb7\x2a\xbb\xb1\x4c\x47\xef\xdc\xe6\x14\x92\xc1\xd6\xcf\x53\xe5\x82\x2c\xd9\x34\xcb\xe8\x14\x41\xd8\x07\x95\x1e\x28\x20\xfa\x25\x5a\x31\x53\x99\x25\xb9\xdd\x5b\x29\x2a\x22\x1f\x9d\x40\x13\x99\x35\x8c\xd6\xb5\xd4\x5a\x41\xb3\xcf\x09\x29\xa8\x1c\x41\x5d\x74\x63\x7d\xf4\x6c\x2a\xdd\x0b\x7a\x71\x3f\xf7\x66\x5e\x19\x20\xcd\xb4\xeb\x28\xa7\x35\x2a\x95\x03\x2a\x0b\x5a\x0e\x3a\x97\x02\x59\x4a\xab\xe6\x2b\xab\x72\x65\x76\xd8\xe8\xfd\xb0\x93\xcc\xc6\x0a\xd5\x3f\xb7\xb9\x0b\x4c\x43\x47\x94\x5a\x47\xd5\x27\xba\x1f\x2a\x53\x3a\x88\xc6\xe6\x64\xe8\x45\x97\x05\x4a\xe2\xbf\xcc\xd0\xae\xb1\xb9\x69\x6c\x95\x2a\x1d\x14\xb6\xa3\x52\xc5\x8f\xb0\xa3\x54\xb2\xd6\xb5\x50\xed\x76\x85\xa2\x29\xec\xc0\x7b\x5e\xa9\x54\x91\x43\x17\xc1\x7b\xfc\x8b\x1f\xa8\x7a\xd1\x0f\x94\x47\xff\xcc\x62\x6b\xac\x35\x46\x1a\xaf\x31\x44\x65\xcd\x35\x89\x93\x5c\xe3\x81\xf7\x4d\xab\x08\x9f\x94\xd6\x41\xd8\xbc\x92\x8e\x12\x86\xc2\x87\xc8\xa4\xc2\xee\xac\x7b\xc1\x47\x8d\x55\x08\xdc\x83\x3d\xa6\xca\x80\x6e\x4b\x98\xd9\x66\x54\x5c\xc7\x57\x45\x6c\x5d\x5c\xaf\x6c\x1c\x8f\x16\x67\xd6\x31\xd0\x55\x35\x96\x03\x54\xda\x3c\x03\xb3\xd3\xd5\x76\x93\xcf\xcd\xff\xe2\xaf\x41\x59\x93\x8a\xfe\x39\xa3\xc0\x52\xb1\xcd\xe0\xa4\x5b\x4d\x3e\xce\x28\xc7\x41\xc0\x2b\xb8\xe4\x38\xc1\xc2\x22\xc7\xfe\x30\xe3\x19\x7f\x1c\xb7\xc2\x51\xcc\xf1\x22\x2b\x7c\x2d\xce\x38\x7e\xdb\x53\x05\x9e\xbf\xf4\x1c\x07\xb6\xc8\x93\x2c\xf9\x47\xb0\xc5\xeb\x2b\xef\xca\xb8\xa2\x7f\xc4\x15\xd9\x23\xae\x90\x6a\xd7\xb8\x22\xb3\xb9\xba\x9c\xd7\x26\x5a\xc9\xf2\xb9\xad\x14\x03\x07\x3b\x63\xe8\xac\x13\xae\x6c\x62\x94\x5b\x67\x4f\x2e\xc0\x17\xe8\x24\x15\x9e\xcd\xd5\x25\x37\x0a\xc3\x72\xcb\xd3\xa4\xd3\x2a\x51\x78\xa7\xd7\x5c\xec\xbd\xe2\x25\x55\xb2\xe8\x5e\xe8\x94\xa5\xb0\x29\x46\x9c\xeb\x64\xa2\xb6\x73\x5a\xb8\xa4\x9f\x3a\xfa\xc1\x2a\xe3\x52\xaf\xcb\x59\xc1\xa8\x7b\xac\xad\x52\xcd\x0b\xdd\x0b\x3a\x4f\xd5\x27\xba\xfa\x41\xb9\x02\xde\x2c\x6c\x2a\x84\x4d\x64\xb6\x3f\xda\x6c\xb8\x50\x3d\x3b\xc3\x62\x30\x04\xcc\x03\x17\xd1\x8b\xbd\x0f\x76\x0e\xe6\xed\x5c\x1a\xbe\x02\x5b\xd8\xaf\xb7\x19\xf6\x60\xb3\x74\xb8\x96\x4a\xb7\xb7\xb3\x28\x74\x87\x6b\x5a\xe3\x08\xc1\x6d\x9b\xc5\x83\xed\x9d\x9d\xf1\x48\xe5\x3f\x68\x77\x04\x36\x0b\x47\xd8\x39\x34\x5d\xc3\xc4\x11\xf8\x93\xf8\x40\x06\x35\x00\x6b\x3e\xe1\x88\x54\xcf\x4c\x9e\x70\x84\x4d\x98\xe2\x54\xf5\xa7\x93\x90\x67\x0b\xb5\x97\x4c\x9d\x7a\x8d\x05\x78\x49\x38\xa2\xd7\xe9\x4a\x67\xd3\xaf\xce\xa6\x74\x55\xaf\x98\x8e\x6e\xdc\x0b\x47\xf0\x9e\xd7\xd1\xb8\x4d\x6a\x85\x46\xbb\x10\x39\x27\xb0\xf3\x2c\xe2\x85\xd7\xc6\xa6\x36\x22\xaf\x53\x39\x52\xff\x18\xde\xe3\xd7\xe9\xdc\xa7\x48\x2c\x27\xac\x3b\xc4\xce\x43\xa5\x3c\x2f\xbf\x0f\x20\x9e\x62\x5d\x8f\xd8\xa0\x93\xe2\x75\x34\x3d\x89\x03\xf6\x87\xdd\xee\xcf\x54\x52\x3d\xdb\xa0\x83\x20\x8e\x91\x49\xfc\x8c\x76\x06\x29\x67\x31\x89\xf9\xd0\x6c\xdb\x59\xe7\xf8\xb1\x49\x4f\x65\x5d\x1f\x1d\x08\x1d\x5f\x9f\xbd\x8a\x37\x5e\xdb\x81\x1c\x23\xf6\x28\xce\xf8\x68\xc7\x71\x8c\xf4\xaf\xc1\x17\x7f\x8f\x0e\xe3\xd8\x0d\x4f\x5c\x51\x24\xff\x4c\xae\x78\xd1\x6f\x1b\xbf\xe8\xb7\x3b\xb9\xc6\xfe\xf4\x3b\x80\x31\x7d\xe0\x0e\x99\xe2\x77\x0f\x93\x82\xa4\xd6\x29\xef\xfd\x59\xe4\x68\xe7\xda\xb9\xe6\x5c\x6e\xf5\xe0\x60\x67\x9f\x32\x75\xcf\xb4\x7e\x93\x93\x06\xeb\x55\xa6\x69\x3a\x18\xc8\xfd\xf2\xdb\x9c\xde\xce\xc3\x2b\xcd\x57\x39\x53\xf5\x7a\x7f\x6a\x78\x43\x0e\x7a\x3b\x61\xf0\xf6\x07\xd7\xa1\x53\x65\x35\x74\x91\x6b\xde\x25\xf6\x3b\x24\x6f\x67\xb1\x99\xed\xa3\xb3\x89\x77\x6a\x67\xce\xad\xe9\x5e\x4f\xe7\x95\xad\xe2\x1a\x79\x9b\x58\xfe\x53\x4f\xc3\x97\x71\x6f\x53\x96\x54\xed\x4a\xcd\xdb\xd8\x5f\xd1\xe8\xa4\x1a\x5c\xce\xed\x84\xa6\xad\xf5\x5c\x78\x9a\xb8\xc3\x97\xc8\xc8\x6b\xb5\x0f\xbc\x06\xff\xb1\x56\x51\xdb\x14\xa9\x56\x19\xa9\x71\x1b\x9c\xc3\x33\xb1\xf5\x62\x72\x8e\x6b\x3d\x81\xfc\x9e\xc0\xfa\x3a\x70\xab\xb1\xf3\x7b\xa9\x85\x07\x3b\x07\x2f\x14\x17\xa5\xf7\xc9\xf4\x6c\x15\x0e\x89\x0d\x0f\xe1\x23\x6c\xd6\x5b\x9d\x40\x8f\xd2\xda\xf4\x05\x7f\x75\xa5\x9d\x33\x67\x5a\x73\xe3\x63\x6f\xe7\xc1\x70\xcc\x18\xdb\x64\xbb\xd1\x58\xe9\x46\xb5\x43\x6b\x67\xd4\xe0\xf1\x34\x0d\x94\x93\x00\x3b\x45\xc2\xc6\x83\x9d\x18\x50\x0b\x60\x0f\x7c\x25\x93\x31\xaf\xf5\x00\xfc\x27\x9c\xef\x95\x87\xf1\xb1\x4c\xf1\x63\xe5\x2b\x6c\x21\xbd\x65\xa2\xfb\x94\x13\x8a\x46\xeb\x0e\x38\x45\x7e\xab\x91\x6a\x0d\x80\x8c\xc6\xfa\x4c\xe9\xc7\x06\x3d\x29\x93\x5e\xa0\x7e\xe8\x11\xf0\x39\x1c\x2a\xb5\x92\x57\xf9\xa2\x5f\x69\x39\x61\xbf\x8b\x41\x16\x71\x31\x7a\x3b\x2d\x8a\xed\xfc\x7f\xb4\x5e\x6d\x54\x1e\x41\x2f\x62\x03\xdb\xf3\x3d\x7f\xa9\x71\x15\xf6\xa0\x9f\xac\x4c\xf7\xce\x7e\x7f\x32\x9a\x2d\xb1\x2f\xfc\x3b\xd8\x6f\x70\x4a\x9b\x5c\x52\x87\x49\xcf\x66\x27\x10\xf8\xb9\xb6\x7e\x89\xde\xba\xb2\x93\x34\xe2\x5e\x62\x6d\xea\x07\x7b\xcd\xf9\xd4\x4e\xfb\x0a\xab\x65\x88\xcd\xd6\xf4\x97\xd3\x3e\xcb\x03\xf2\x96\xbd\x4b\x9e\xdb\x6f\x32\xe4\x1e\x3b\x99\xea\xa7\x18\xaf\x34\x2f\xbc\xfd\x96\x06\x4c\xe9\xec\xa4\xa3\xa9\xad\x07\xb3\x93\x25\xa9\xb9\xac\xf6\x6c\xfe\xc6\x34\x8f\xef\xe4\x14\xca\x4e\xe2\x5a\xfb\x0d\xc0\x73\x0e\x7d\x0d\x48\x1e\xcf\xa9\xff\x1f\x00\x00\xff\xff\xc9\x7b\xa1\xef\x00\x30\x00\x00")

func bindataGoBytes() ([]byte, error) {
	return bindataRead(
		_bindataGo,
		"bindata.go",
	)
}

func bindataGo() (*asset, error) {
	bytes, err := bindataGoBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "bindata.go", size: 28672, mode: os.FileMode(420), modTime: time.Unix(1576581519, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"1_create_table_consent_rule.down.sql":                   _1_create_table_consent_ruleDownSql,
	"1_create_table_consent_rule.up.sql":                     _1_create_table_consent_ruleUpSql,
	"2_alter_consent_record_add_version_uuid.down.sql":       _2_alter_consent_record_add_version_uuidDownSql,
	"2_alter_consent_record_add_version_uuid.up.sql":         _2_alter_consent_record_add_version_uuidUpSql,
	"3_rename_resource_to_data_class.down.sql":               _3_rename_resource_to_data_classDownSql,
	"3_rename_resource_to_data_class.up.sql":                 _3_rename_resource_to_data_classUpSql,
	"4_alter_consent_record_make_valid_to_optional.down.sql": _4_alter_consent_record_make_valid_to_optionalDownSql,
	"4_alter_consent_record_make_valid_to_optional.up.sql":   _4_alter_consent_record_make_valid_to_optionalUpSql,
	"bindata.go":                                             bindataGo,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"1_create_table_consent_rule.down.sql":                   &bintree{_1_create_table_consent_ruleDownSql, map[string]*bintree{}},
	"1_create_table_consent_rule.up.sql":                     &bintree{_1_create_table_consent_ruleUpSql, map[string]*bintree{}},
	"2_alter_consent_record_add_version_uuid.down.sql":       &bintree{_2_alter_consent_record_add_version_uuidDownSql, map[string]*bintree{}},
	"2_alter_consent_record_add_version_uuid.up.sql":         &bintree{_2_alter_consent_record_add_version_uuidUpSql, map[string]*bintree{}},
	"3_rename_resource_to_data_class.down.sql":               &bintree{_3_rename_resource_to_data_classDownSql, map[string]*bintree{}},
	"3_rename_resource_to_data_class.up.sql":                 &bintree{_3_rename_resource_to_data_classUpSql, map[string]*bintree{}},
	"4_alter_consent_record_make_valid_to_optional.down.sql": &bintree{_4_alter_consent_record_make_valid_to_optionalDownSql, map[string]*bintree{}},
	"4_alter_consent_record_make_valid_to_optional.up.sql":   &bintree{_4_alter_consent_record_make_valid_to_optionalUpSql, map[string]*bintree{}},
	"bindata.go":                                             &bintree{bindataGo, map[string]*bintree{}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
