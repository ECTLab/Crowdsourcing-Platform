package osm

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/bytedance/sonic"
	"github.com/qedus/osmpbf"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	PbfUrl        string
	PbfVersionUrl string
	Username      string
	Password      string
	Path          string
	Insecure      bool
}

var once sync.Once
var downloaded chan error = nil
var filePath string

// call as go routine (go osm.Download(osm.Config{})).
// Downloads the osm.pbf file in the background.
// call Iterate to use the downloaded data.
// if the the data is still being downloaded, Iterate will wait.
func Download(config Config) {
	once.Do(func() {
		downloaded = make(chan error)
		defer close(downloaded)
		useInsecureHttp()
		version := getLatestVersions(config, &downloaded)
		downloadOsmPbfFile(config, version, &downloaded)
		downloaded <- nil
	})
}

func useInsecureHttp() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
}

func getLatestVersions(config Config, downloaded *chan error) string {
	var version string = ""

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	client := http.DefaultClient
	req, err := http.NewRequestWithContext(ctx, "GET", config.PbfVersionUrl, nil)
	if err != nil {
		log.WithError(err).Errorf("osm: download: error creating json file download request")
		*downloaded <- err
		return version
	}

	if config.Username != "" {
		req.SetBasicAuth(config.Username, config.Password)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.WithError(err).Errorf("osm: download: error downloading version json file")
		*downloaded <- err
		return version
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.WithError(err).Errorf("osm: download: error reading version json into bytes")
		*downloaded <- err
		return version
	}
	versions := struct {
		ListOfActiveVersions []string `json:"list_of_active_versions"`
	}{}
	err = sonic.Unmarshal(body, &versions)
	if err != nil {
		log.WithError(err).Errorf("osm: download: error parsing version json")
		*downloaded <- err
		return version
	}

	lst := versions.ListOfActiveVersions
	if len(lst) > 0 {
		version = lst[len(lst)-1]
	} else {
		err := errors.New("empty osm versions list")
		log.WithError(err).Errorf("osm: download: version list is empty")
		*downloaded <- err
		return version
	}

	return version
}

func downloadOsmPbfFile(config Config, version string, downloaded *chan error) {
	PbfUrlWithVersion := strings.ReplaceAll(config.PbfUrl, "{VERSION}", version)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()
	client := http.DefaultClient
	req, err := http.NewRequestWithContext(ctx, "GET", PbfUrlWithVersion, nil)
	if err != nil {
		log.WithError(err).Errorf("osm: download: error creating pbf file download request")
		*downloaded <- err
		return
	}

	if config.Username != "" {
		req.SetBasicAuth(config.Username, config.Password)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.WithError(err).Errorf("osm: download: error downloading pbf file")
		*downloaded <- err
		return
	}

	fpath := path.Join(config.Path, "osm.pbf")

	file, err := os.Create(fpath)
	if err != nil {
		log.WithError(err).Errorf("osm: download: error creating pbf file")
		*downloaded <- err
		return
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.WithError(err).Errorf("osm: download: error writing pbf to file")
		*downloaded <- err
		return
	}

	filePath = fpath
}

// call Download before this to download data.
// otherwise Iterate will block.
func Iterate(f func(any) error) error {
	if downloaded == nil {
		return errors.New("osm.Iterate called before osm.Download")
	}

	err := <-downloaded
	if err != nil {
		return err
	}

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := osmpbf.NewDecoder(file)

	// use more memory from the start, it is faster
	decoder.SetBufferSize(512 * 1024 * 1024) // 512 Mb

	// start decoding with several goroutines, it is faster
	err = decoder.Start(runtime.GOMAXPROCS(-1))
	if err != nil {
		return err
	}

	for {
		if osmObject, err := decoder.Decode(); err == io.EOF {
			break
		} else if err != nil {
			return err
		} else {
			err := f(osmObject)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func ExampleIterFunc(v any) error {
	switch v := v.(type) {
	case *osmpbf.Node:
	case *osmpbf.Way:
	case *osmpbf.Relation:
	default:
		return fmt.Errorf("unknown type %T", v)
	}
	return nil
}
