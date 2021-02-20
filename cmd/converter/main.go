package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"kubesphere.io/monitoring-dashboard/tools/converter"
)

// a converter container
type ConverterContainer struct {
	// target paths that the converter takes jobs to Converter
	Inputs []string
	// inner useful json filepaths to parse
	JsonFilePaths []string
	// default: "json";  a suffix string for json path
	Suffix string
	// output path for target manifests
	Output string
}

var inputPath string
var outputPath string
var isClusterCrd bool
var namespace string
var name string

// init
func init() {
	flag.StringVar(&inputPath, "inputPath", "./manifests/inputs", "a input path for the converter to look for jobs")
	flag.StringVar(&outputPath, "outputPath", "./manifests/outputs", "a output path for the converter to store manifests")
	flag.BoolVar(&isClusterCrd, "isClusterCrd", false, "a flag that defines whether build the cluster dashboard resource or not")
	flag.StringVar(&namespace, "namespace", "default", "namespace of the dashboard resource")
	flag.StringVar(&name, "name", "", "name of the dashboard resource")
}

// main function
func main() {
	// parse the params
	flag.Parse()

	// init a Converter Container
	c := NewConverterContainer(inputPath)
	// fills with a logger
	logger, err := createLogger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not create logger: %s", err)
		os.Exit(1)
	}

	// finds json files from the given input path
	for _, jsonPath := range c.Inputs {
		c.getJsonFiles(jsonPath)
	}

	// exits if it could not get a json file
	if len(c.JsonFilePaths) == 0 {
		fmt.Fprintf(os.Stderr, "Could not get a json file: %s\n", err)
		os.Exit(1)
	}

	// sets a gr for each json file
	// once compeleted, each manifest will fill in the target path
	var wg sync.WaitGroup

	for _, fi := range c.JsonFilePaths {
		wg.Add(1)
		go func(inputFile string, logger *zap.Logger) {
			c.toKubesphereDashboard(inputFile, logger, isClusterCrd, namespace, name)
			wg.Done()
		}(fi, logger)
	}

	wg.Wait()
	logger.Info("Finished processing")

}

// new a Converter Container struct pointer by a list of inputs
func NewConverterContainer(inputs ...string) *ConverterContainer {
	if len(inputs) <= 0 {
		inputs = append(inputs, ".")
	}

	return &ConverterContainer{
		Inputs:        inputs,
		Suffix:        "json",
		JsonFilePaths: make([]string, 0),
		Output:        outputPath,
	}
}

func createLogger() (*zap.Logger, error) {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	cfg := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:       false,
		DisableStacktrace: true,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "console",
		EncoderConfig:    encoderCfg,
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	return cfg.Build()
}

// find out all json paths under the given path
func (c *ConverterContainer) getJsonFiles(dirPath string) error {

	UpperSuffix := strings.ToUpper(c.Suffix)

	// termination conditions
	_, err := ioutil.ReadFile(dirPath)
	if err == nil {
		// needs to confirm whether was a json file
		if isJsonFile(name, c.Suffix, UpperSuffix) {
			c.JsonFilePaths = append(c.JsonFilePaths, dirPath)
			return nil
		}
	}

	dir, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return errors.New("not a dir")
	}

	pthSep := string(os.PathSeparator)
	// recursive algorithm
	for _, f := range dir {
		name := f.Name()
		fString := strings.Join([]string{dirPath, name}, pthSep)
		if f.IsDir() {
			c.getJsonFiles(fString)
		} else {
			ok := isJsonFile(name, c.Suffix, UpperSuffix)
			if ok {
				c.JsonFilePaths = append(c.JsonFilePaths, fString)
			}
		}

	}

	return nil

}

// ConverterContainers a json file to a k8s manifest
func (c *ConverterContainer) toKubesphereDashboard(inputFile string, logger *zap.Logger, isClusterCrd bool, ns string, name string) {
	input, err := os.Open(inputFile)
	if err != nil {
		logger.Fatal("Could not open input file", zap.Error(err))
	}

	_, fileName := filepath.Split(inputFile)
	prevFileName := strings.Split(fileName, ".")[0]

	if ns == "" {
		ns = "default"
	}

	// inner name
	if name == "" {
		name = strings.Replace(prevFileName, "_", "-", -1)
	}

	outputFile := filepath.Join(c.Output, prevFileName+".yaml")
	output, err := os.OpenFile(outputFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755)
	if err != nil {
		logger.Fatal("Could not open output file", zap.Error(err))
	}

	conv := converter.NewConverter(logger)

	if err := conv.ConvertKubsphereDashboard(input, output, isClusterCrd, ns, name); err != nil {
		logger.Fatal("Could not convert dashboard", zap.Error(err))
	}

	logger.Info("Successfully convert a input json file to a manifest", zap.Any("srcPath", inputFile), zap.Any("targetPath", outputFile))

}

// confirms it was a json file
func isJsonFile(name string, suffix string, upSuffix string) bool {
	return strings.HasSuffix(name, suffix) || strings.HasSuffix(name, upSuffix)
}
