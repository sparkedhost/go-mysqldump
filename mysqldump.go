package mysqldump

import (
	"database/sql"
	"errors"
	"io"
	"os"
	"path"
)

// Register a new dumper.
func Register(db *sql.DB, dir, filename string) (*Data, error) {
	if !isDir(dir) {
		return nil, errors.New("Invalid directory")
	}

	p := path.Join(dir, filename+".sql")

	// Check dump directory
	if e, _ := exists(p); e {
		return nil, errors.New("Dump '" + filename + ".sql" + "' already exists.")
	}

	// Create .sql file
	f, err := os.Create(p)

	if err != nil {
		return nil, err
	}

	return &Data{
		Out:        f,
		Connection: db,
	}, nil
}

// Dump Creates a MYSQL dump from the connection to the stream.
func Dump(db *sql.DB, out io.Writer) error {
	return (&Data{
		Connection: db,
		Out:        out,
	}).Dump()
}

// Close the dumper.
// Will also close the database the dumper is connected to as well as the out stream if it has a Close method.
//
// Not required.
func (data *Data) Close() error {
	defer func() {
		data.Connection = nil
		data.Out = nil
	}()
	if out, ok := data.Out.(io.Closer); ok {
		out.Close()
	}
	return data.Connection.Close()
}

func exists(p string) (bool, os.FileInfo) {
	f, err := os.Open(p)
	if err != nil {
		return false, nil
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return false, nil
	}
	return true, fi
}

func isDir(p string) bool {
	if e, fi := exists(p); e {
		return fi.Mode().IsDir()
	}
	return false
}
