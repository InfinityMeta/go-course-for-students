package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type Options struct {
	From      string
	To        string
	Offset    int64
	Limit     int64
	BlockSize int64
	Conv      string
	// todo: add required flags
}

type ConvFile struct {
	Options
	buffer []byte
}

type ReaderWriter interface {
	Read() (int, error)
	Write() (int, error)
}

type Converter interface {
	Convert() (int, error)
}

func (convFile *ConvFile) Read() (int, error) {
	var (
		input io.Reader
		err   error
	)

	if fileFrom := convFile.From; fileFrom != "" {
		f, err := os.Open(fileFrom)
		if err != nil {
			err = errors.New("can not open file for reading")
			return 0, err
		}

		defer f.Close()

		fStat, err := f.Stat()
		if err != nil {
			err = errors.New("can not check file size")
			return 0, err
		}

		if convFile.Offset >= fStat.Size() {
			err = errors.New("offset is bigger than file size")
			return 0, err
		}

		input = f

	} else {
		input = os.Stdin
	}

	if convFile.Limit != -1 {
		input = io.LimitReader(input, convFile.Limit+convFile.Offset)
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(input)

	if err != nil {
		return 0, err
	}

	convFile.buffer = buf.Bytes()[convFile.Offset:]

	return len(convFile.buffer), nil

}

func (ConvFile *ConvFile) Write() (int, error) {
	var (
		output io.Writer
		err    error
	)
	if fileTo := ConvFile.To; fileTo != "" {

		if _, err := os.Stat(fileTo); err == nil {
			err = errors.New("copying to existing files if prohibited")
			return 0, err
		}

		f, err := os.Create(fileTo)
		if err != nil {
			err = errors.New("can not create file")
			return 0, err
		}
		defer f.Close()

		output = f

	} else {
		output = os.Stdout
	}

	_, err = output.Write(ConvFile.buffer)

	if err != nil {
		return 0, err
	}

	return len(ConvFile.buffer), nil

}

func (convFile *ConvFile) Convert() (int, error) {
	var err error

	if convFile.Conv == "" {
		return 0, nil
	}

	conversionsSlice := strings.Split(convFile.Conv, ",")
	for _, conv := range conversionsSlice {
		if conv != "upper_case" && conv != "lower_case" && conv != "trim_spaces" {
			err = errors.New("invalid conversions")
			return 0, err
		}
	}

	if strings.Contains(convFile.Conv, "upper_case") && strings.Contains(convFile.Conv, "lower_case") {
		err = errors.New("upper case and lower case cannot be applied both")
		return 0, err
	}

	strBytes := string(convFile.buffer)

	if strings.Contains(convFile.Conv, "upper_case") {
		strBytes = strings.ToUpper(strBytes)
	}

	if strings.Contains(convFile.Conv, "lower_case") {
		strBytes = strings.ToLower(strBytes)
	}

	if strings.Contains(convFile.Conv, "trim_spaces") {
		strBytes = strings.TrimSpace(strBytes)
	}

	convFile.buffer = []byte(strBytes)

	return len(convFile.buffer), nil

}

func ParseFlags() (*Options, error) {
	var opts Options

	flag.StringVar(&opts.From, "from", "", "file to read. by default - stdin")
	flag.StringVar(&opts.To, "to", "", "file to write. by default - stdout")
	flag.Int64Var(&opts.Offset, "offset", 0, "number of bytes in input to skip. by default - 0")
	flag.Int64Var(&opts.Limit, "limit", -1, "maximum of bytes to read. by default - size of file")
	flag.Int64Var(&opts.BlockSize, "block_size", -1, "size of block for process. by default - all available memory")
	flag.StringVar(&opts.Conv, "conv", "", "convertions applying to text. by default - none")

	// todo: parse and validate all flags

	flag.Parse()

	return &opts, nil
}

func main() {
	opts, err := ParseFlags()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "can not parse flags:", err)
		os.Exit(1)
	}

	var fileConverter = &ConvFile{*opts, []byte{}}

	_, err = fileConverter.Read()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "can not read from file:", err)
		os.Exit(1)
	}

	_, err = fileConverter.Convert()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "can not convert file:", err)
		os.Exit(1)
	}

	_, err = fileConverter.Write()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "can not write to file:", err)
		os.Exit(1)
	}

	// todo: implement the functional requirements described in read.me
}
