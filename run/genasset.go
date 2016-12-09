// +build ignore

package main

import (
	"bytes"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"os"
	"time"
)

const (
	funcDecompressString = "gunzipString"

	clInput    = "i"
	clOutput   = "o"
	clPkgName  = "p"
	clFuncName = "f"
)

type args struct {
	input, output, pkgname, funcname string
}

func CompressFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)

	_, err = io.Copy(gz, file)
	clErr := gz.Close()

	if err != nil {
		return nil, err
	}
	if clErr != nil {
		return nil, clErr
	}
	return buf.Bytes(), nil
}

func writeHeader(w io.Writer, pkgname string) {
	fmt.Fprintf(w, `// automatically generated. DO NOT EDIT
// %s

package %s

import (
	"bytes"
	"compress/gzip"
	"io"
)

`, time.Now(), pkgname)
}

func writeAssetFunc(w io.Writer, funcname string, gzipped []byte) {
	fmt.Fprintf(w, "func %s() string {\n", funcname)
	fmt.Fprint(w, "\tvar compressed = []byte{")

	for i, b := range gzipped {
		if i%16 == 0 {
			fmt.Fprint(w, "\n\t\t")
		}
		fmt.Fprint(w, b, ",")
	}
	fmt.Fprintf(w, `
	}
	res, err := %s(compressed)
	if err != nil {
		panic(err)
	}
	return res
}

`, funcDecompressString)
}

func writeDecompressFunc(w io.Writer) {
	fmt.Fprintf(w, `
func %s(data []byte) (string, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()
	if err != nil {
		return "", err
	}
	if clErr != nil {
		return "", err
	}
	return buf.String(), nil
}

`, funcDecompressString)
}

func generateOutput(pathIn, pathOut, pkgname, funcname string) error {

	buf, err := CompressFile(pathIn)
	if err != nil {
		return fmt.Errorf("CompressFile: %s", err)
	}

	f, err := os.Create(pathOut)
	if err != nil {
		return err
	}
	defer f.Close()

	writeHeader(f, pkgname)
	writeAssetFunc(f, funcname, buf)
	writeDecompressFunc(f)
	return nil
}

func parseArgs() (*args, error) {
	var a args

	// command line arguments
	flag.StringVar(&a.input, clInput, "", "Input file to embed compressed.")
	flag.StringVar(&a.output, clOutput, "", "Output file. In empty will be used input file name with \".go\" suffix.")
	flag.StringVar(&a.pkgname, clPkgName, "main", "Package name of the generated code.")
	flag.StringVar(&a.funcname, clFuncName, "getAsset", "Func name of the function that returns the original asset.")

	flag.Parse()

	if a.input == "" {
		return nil, fmt.Errorf("Missing mandatory argument: %s", clInput)
	}
	if a.output == "" {
		a.output = a.input + ".go"
	}

	return &a, nil
}

func showUsage() {
	fmt.Fprintln(os.Stderr, "CMD")
	flag.PrintDefaults()
}

func run() int {
	a, err := parseArgs()
	if err == flag.ErrHelp {
		showUsage()
		return 0
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		fmt.Fprintln(os.Stderr, "Try -h for help.")
		return 2
	}
	err = generateOutput(a.input, a.output, a.pkgname, a.funcname)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}
	return 0
}

func main() {
	os.Exit(run())
}
