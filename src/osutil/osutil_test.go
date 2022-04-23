package osutil

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestCopy_getCommonPath(t *testing.T) {
	type fields struct {
		Source       string
		Destination  string
		Options      *CopyOptions
		sourceStat   os.FileInfo
		destStat     os.FileInfo
		currentDepth int
	}
	type args struct {
		path string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "test1",
			fields: fields{
				Source:       "/home/user/source",
				Destination:  "/home/user/destination",
				Options:      nil,
				sourceStat:   nil,
				destStat:     nil,
				currentDepth: 0,
			},
			args: args{
				path: "/home/user/source/file.txt",
			},
			want: "file.txt",
		},
		{
			name: "test2",
			fields: fields{
				Source:       "/home/user/source",
				Destination:  "/home/user/destination",
				Options:      nil,
				sourceStat:   nil,
				destStat:     nil,
				currentDepth: 0,
			},
			args: args{
				path: "/home/user/source/dir/file.txt",
			},
			want: "dir/file.txt",
		},
		{
			name: "test3",
			fields: fields{
				Source:       "/home/user/source",
				Destination:  "/home/user/destination",
				Options:      nil,
				sourceStat:   nil,
				destStat:     nil,
				currentDepth: 0,
			},
			args: args{
				path: "/home/user/source/dir/dir2/file.txt",
			},
			want: "dir/dir2/file.txt",
		},
		{
			name: "windows test",
			fields: fields{
				Source:       "C:\\home\\user\\source",
				Destination:  "C:\\home\\user\\destination",
				Options:      nil,
				sourceStat:   nil,
				destStat:     nil,
				currentDepth: 0,
			},
			args: args{
				path: "C:\\home\\user\\source\\dir\\dir2\\file.txt",
			},
			want: "dir\\dir2\\file.txt",
		},
		{
			name: "windows test2",
			fields: fields{
				Source:       "C:\\home\\user\\source",
				Destination:  "C:\\home\\user\\destination",
				Options:      nil,
				sourceStat:   nil,
				destStat:     nil,
				currentDepth: 0,
			},
			args: args{
				path: "C:\\home\\user\\source",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Copy{
				Source:      tt.fields.Source,
				Destination: tt.fields.Destination,
				Options:     tt.fields.Options,
				sourceStat:  tt.fields.sourceStat,
				destStat:    tt.fields.destStat,
			}
			if got := c.getCommonPath(tt.args.path); got != tt.want {
				t.Errorf("getCommonPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCopy_GetCalculatedDestination(t *testing.T) {
	type fields struct {
		Source       string
		Destination  string
		Options      *CopyOptions
		sourceStat   os.FileInfo
		destStat     os.FileInfo
		currentDepth int
	}
	type args struct {
		path string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "GetCalculatedDestination",
			fields: fields{
				Source:       "/home/user/source",
				Destination:  "/home/user/destination",
				Options:      nil,
				sourceStat:   nil,
				destStat:     nil,
				currentDepth: 0,
			},
			args: args{
				path: "/home/user/source/file.txt",
			},
			want: "/home/user/destination/file.txt",
		},
		{
			name: "GetCalculatedDestination 1",
			fields: fields{
				Source:       "/home/user/source",
				Destination:  "/home/user/destination",
				Options:      nil,
				sourceStat:   nil,
				destStat:     nil,
				currentDepth: 0,
			},
			args: args{
				path: "/home/user/source/",
			},
			want: "/home/user/destination/",
		},
		{
			name: "Get CalculatedDestination With Directory",
			fields: fields{
				Source:       "/home/user/source",
				Destination:  "/home/user/destination",
				Options:      nil,
				sourceStat:   nil,
				destStat:     nil,
				currentDepth: 0,
			},
			args: args{
				path: "/home/user/source/dir",
			},
			want: "/home/user/destination/dir",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Copy{
				Source:      tt.fields.Source,
				Destination: tt.fields.Destination,
				Options:     tt.fields.Options,
				sourceStat:  tt.fields.sourceStat,
				destStat:    tt.fields.destStat,
			}
			if got := c.GetCalculatedDestination(tt.args.path); got != tt.want {
				t.Errorf("GetCalculatedDestination() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getDepth(t *testing.T) {
	basePath := filepath.Join("a", "b", "c")
	type args struct {
		givenAbsolutePath string
		path              string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "TestI",
			args: args{
				givenAbsolutePath: basePath,
				path:              basePath,
			},
			want: 0,
		},
		{
			name: "TestII",
			args: args{
				givenAbsolutePath: basePath,
				path:              basePath + string(os.PathSeparator) + "abc.go",
			},
			want: 0,
		},
		{
			name: "TestIII",
			args: args{
				givenAbsolutePath: basePath,
				path:              basePath + string(os.PathSeparator) + "src" + string(os.PathSeparator) + "abc.go",
			},
			want: 1,
		},
		{
			name: "TestIV",
			args: args{
				givenAbsolutePath: basePath,
				path:              basePath + string(os.PathSeparator) + "src" + string(os.PathSeparator),
			},
			want: 1,
		},
		{
			name: "TestV",
			args: args{
				givenAbsolutePath: basePath,
				path:              basePath + string(os.PathSeparator) + "src" + string(os.PathSeparator) + "main" + string(os.PathSeparator),
			},
			want: 2,
		},
		{
			name: "TestVI",
			args: args{
				givenAbsolutePath: basePath,
				path:              basePath + string(os.PathSeparator) + "src" + string(os.PathSeparator) + "main" + string(os.PathSeparator) + "main.go",
			},
			want: 2,
		},
		{
			name: "TestVI",
			args: args{
				givenAbsolutePath: basePath,
				path:              basePath + string(os.PathSeparator) + "src" + string(os.PathSeparator) + "main" + string(os.PathSeparator) + "main.go" + string(os.PathSeparator),
			},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getDepth(tt.args.givenAbsolutePath, tt.args.path); got != tt.want {
				t.Errorf("getDepth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_copyFile(t *testing.T) {
	dir := os.TempDir()
	sourceDir := dir + "source"
	destDir := dir + string(os.PathSeparator) + "destination"

	tempSource, err := os.Create(sourceDir)
	if err != nil {
		t.Fatalf("Temp Source file didn't create. err = %v", err)
	}

	write, err := tempSource.Write([]byte("test"))
	if err != nil {
		t.Fatalf("Error occurred while writing to the file. err = %v", err)
	}
	if write != 4 {
		t.Fatalf("File byte size is wrong.")
	}

	type args struct {
		src string
		dst string
	}
	tests := []struct {
		name      string
		args      args
		checkFunc func(funcErrResult error) error
	}{
		{
			name: "Should copy file",
			args: args{
				src: sourceDir,
				dst: destDir,
			},
			checkFunc: func(funcErrResult error) error {

				if funcErrResult != nil {
					return funcErrResult
				}

				file, err2 := ioutil.ReadFile(destDir)
				if err2 != nil {
					return err2
				}
				if string(file) != "test" {
					return errors.New("destination file is wrong")
				}
				return nil
			},
		},
		{
			name: "Should give error",
			args: args{
				src: "abc.txt",
				dst: destDir,
			},
			checkFunc: func(funcErrResult error) error {
				if funcErrResult == nil {
					return errors.New("abc.txt shouldn't be exist")
				}
				return nil
			},
		},
		{
			name: "When destination directory is not exists, should give error",
			args: args{
				src: sourceDir,
				dst: dir + "newdir" + string(os.PathSeparator) + "destination",
			},
			checkFunc: func(funcErrResult error) error {
				fmt.Printf("will pass: %v", funcErrResult)
				if funcErrResult == nil {
					return errors.New("when directory is not exists should give error")
				}
				return nil
			},
		},
		{
			name: "When source is not exists, should give error",
			args: args{
				src: sourceDir + "ass",
				dst: destDir,
			},
			checkFunc: func(funcErrResult error) error {
				if funcErrResult == nil {
					return errors.New("when source is not exists, should give error")
				}
				return nil
			},
		},
		{
			name: "When source is directory, should give error",
			args: args{
				src: dir,
				dst: destDir,
			},
			checkFunc: func(funcErrResult error) error {
				if funcErrResult == nil {
					return errors.New("when source is directory, should give error")
				}
				return nil
			},
		},
		{
			name: "When destination is directory, should give error",
			args: args{
				src: sourceDir,
				dst: dir,
			},
			checkFunc: func(funcErrResult error) error {
				if funcErrResult == nil {
					return errors.New("when destination is directory, should give error")
				}
				return nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultErr := copyFile(tt.args.src, tt.args.dst)
			if err := tt.checkFunc(resultErr); err != nil {
				t.Fatalf("testname: %s, copyFile() = %v", tt.name, err)
			}
		})
	}
}

func TestCopy_loadSourceStat(t *testing.T) {

	dir := os.TempDir()
	sourceDir := dir + "testfile"
	create, err := os.Create(sourceDir)
	if err != nil {
		t.Fatalf("TestCopy_loadSourceStat is failed! err = %v", err)
	}
	write, err := create.Write([]byte("test"))
	if err != nil {
		t.Fatalf("TestCopy_loadSourceStat is failed! err = %v", err)
	}

	if write != 4 {
		t.Fatalf("TestCopy_loadSourceStat is failed! Write size is not matched with needed size.")
	}

	type fields struct {
		Source     string
		sourceStat os.FileInfo
	}
	tests := []struct {
		name      string
		fields    fields
		checkFunc func(stat os.FileInfo, responseErr error) error
	}{
		{
			name: "Should give error if stat is not loaded",
			fields: fields{
				Source:     sourceDir,
				sourceStat: nil,
			},
			checkFunc: func(sourceStat os.FileInfo, responseErr error) error {
				if responseErr != nil {
					return responseErr
				}

				if sourceStat == nil {
					return err
				}

				return nil
			},
		},
		{
			name: "Should give error if file not exists",
			fields: fields{
				Source:     "weirdpath",
				sourceStat: nil,
			},
			checkFunc: func(stat os.FileInfo, responseErr error) error {
				if responseErr == nil {
					return errors.New("if file is not exists, should give error. so there may problem")
				}
				if stat != nil {
					return errors.New("stat must be nil")
				}
				return nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Copy{
				Source:     tt.fields.Source,
				sourceStat: tt.fields.sourceStat,
			}
			respErr := c.loadSourceStat()
			if err := tt.checkFunc(tt.fields.sourceStat, respErr); err != nil {
				t.Fatalf("%s is failed. c.loadSourceStat() = %v", tt.name, err)
			}
		})
	}
}
