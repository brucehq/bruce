package loader

import (
	"bruce/exe"
	"bytes"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"path"
)

type PageLink struct {
	Target string
	Text   string
}

func CopyFile(src, dest string, perm os.FileMode, overwrite bool) error {
	// if filemode is 0, set it to 0644
	if perm == 0 {
		perm = 0644
	}
	sd, _, err := GetRemoteData(src)
	if err != nil {
		log.Error().Err(err).Msg("cannot open source file")
		return err
	}
	// create a io.reader from sd
	source := bytes.NewReader(sd)
	// check if the destination does not start with . or /, then it's a remote path
	if dest[0] != '.' && dest[0] != '/' {
		// now if it starts with s3 we upload it to s3
		if dest[0:5] == "s3://" {
			return uploadToS3(dest, source)
		}
		// if it starts with http we upload it to http
		if dest[0:4] == "http" {
			return errors.New("http upload not supported")
		}
	}
	if exe.FileExists(dest) {
		if overwrite {
			log.Err(exe.DeleteFile(dest))
		} else {
			log.Error().Msgf("file %s already exists", dest)
			return fmt.Errorf("file %s already exists", dest)
		}
	} else {
		// check if the directories exist to render the file
		if !exe.FileExists(path.Dir(dest)) {
			log.Err(os.MkdirAll(path.Dir(dest), perm))
		}
	}

	destination, err := os.OpenFile(dest, os.O_RDWR|os.O_CREATE, perm)
	if err != nil {
		log.Error().Err(err).Msgf("could not open file for writing copy: %s", dest)
		return err
	}
	printSrc := src
	if len(src) > 32 {
		printSrc = "..." + src[len(src)-32:]
	}
	printDest := dest
	if len(dest) > 32 {
		printDest = "..." + dest[len(dest)-32:]
	}
	log.Info().Msgf("copying %s ==> %s", printSrc, printDest)

	sln, err := io.Copy(destination, source)
	if err != nil {
		log.Error().Err(err).Msg("could not copy file")
		log.Err(destination.Close())
		return err
	}
	log.Debug().Msgf("copied %d bytes", sln)
	log.Err(destination.Close())
	return nil
}

func uploadToS3(dest string, source *bytes.Reader) error {
	// dest should be in format s3://bucket/key, validate this format then split the dest into bucket and key
	// then upload the source to the bucket with the key
	if dest[0:5] != "s3://" {
		return errors.New("invalid s3 destination, must use format s3://<bucket>/<key>")
	}
	// read the source to bytes:
	data, err := io.ReadAll(source)
	if err != nil {
		log.Error().Err(err).Msg("could not read source data")
		return err
	}
	return WriteToS3(dest, data)
}

func RecursiveCopy(src string, baseDir, dest string, overwrite bool, ignores []string, isFlatCopy bool, maxDepth, maxConcurrent int) error {
	if src[0:4] == "http" {
		// This is a remote http copy
		return recursiveHttpCopy(src, baseDir, dest, overwrite, ignores, isFlatCopy, maxDepth, maxConcurrent)
	}
	if src[0:5] == "s3://" {
		// This is a remote s3 copy
		return recursiveS3Copy(src, baseDir, dest, overwrite, ignores, isFlatCopy, maxDepth, maxConcurrent)
	}
	return recursiveNotSupported(src, baseDir, dest, overwrite, ignores, isFlatCopy, maxDepth)
}

func recursiveNotSupported(_ string, _, _ string, _ bool, _ []string, _ bool, _ int) error {
	log.Error().Msg("recursive copy not supported for this source")
	return fmt.Errorf("recursive copy not supported for this source")
}
