package main

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	const testDirectory = "testdata"
	const testFileExtension = ".txt"
	const postfix = "_new"

	testTable := []struct {
		name           string
		sourceFileName string
		offset         int64
		limit          int64
		error          bool
	}{
		{
			name:           "out_offset0_limit0",
			sourceFileName: "input.txt",
			offset:         0,
			limit:          0,
			error:          false,
		},
		{
			name:           "out_offset0_limit10",
			sourceFileName: "input.txt",
			offset:         0,
			limit:          10,
			error:          false,
		},
		{
			name:           "out_offset0_limit1000",
			sourceFileName: "input.txt",
			offset:         0,
			limit:          1000,
			error:          false,
		},
		{
			name:           "out_offset0_limit10000",
			sourceFileName: "input.txt",
			offset:         0,
			limit:          10000,
			error:          false,
		},
		{
			name:           "out_offset100_limit1000",
			sourceFileName: "input.txt",
			offset:         100,
			limit:          1000,
			error:          false,
		},
		{
			name:           "out_offset6000_limit1000",
			sourceFileName: "input.txt",
			offset:         6000,
			limit:          1000,
			error:          false,
		},
		{
			name:           "the_offset_is_larger_than_the_file_size",
			sourceFileName: "input.txt",
			offset:         10000,
			limit:          1000,
			error:          true,
		},
		{
			name:           "file with unknown size",
			sourceFileName: "random",
			offset:         0,
			limit:          0,
			error:          true,
		},
	}
	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			_, testFileName, _, _ := runtime.Caller(0)
			testDirectory := filepath.Join(filepath.Dir(testFileName), testDirectory)
			goldenFilePath := filepath.Join(testDirectory, filepath.FromSlash(tc.name)+testFileExtension)
			sourceFilePath := filepath.Join(testDirectory, tc.sourceFileName)
			resultFilePath := filepath.Join(testDirectory, filepath.FromSlash(tc.name+postfix)+testFileExtension)
			defer os.Remove(resultFilePath)

			err := Copy(sourceFilePath, resultFilePath, tc.offset, tc.limit)
			if !tc.error {
				require.NoError(t, err, "failed copying")
				goldenFileBytes, err := os.ReadFile(goldenFilePath)
				require.NoError(t, err, "failed reading golden file")

				resultFileBytes, err := os.ReadFile(resultFilePath)
				require.NoError(t, err, "failed reading result file")

				require.Equal(t, goldenFileBytes, resultFileBytes)
				return
			}
			require.Error(t, err, "an error is expected ...")
		})
	}
}
