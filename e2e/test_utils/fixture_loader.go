package test_utils

import (
	"database/sql"
	"github.com/go-testfixtures/testfixtures/v3"
	"io/ioutil"
	"log"
	"path"
	"regexp"
)

type FixtureLoaderBuilder interface {
	Build(db *sql.DB, fixturesDirectory string) *testfixtures.Loader
	getFilesInDirectory(fixturesDirectory string) []string
}

type fixtureLoaderBuilder struct{}

func NewFixtureLoaderBuilder() FixtureLoaderBuilder {
	return &fixtureLoaderBuilder{}
}

func (builder *fixtureLoaderBuilder) Build(db *sql.DB, fixturesDirectory string) *testfixtures.Loader {
	files := builder.getFilesInDirectory(fixturesDirectory)
	loader, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgresql"),
		testfixtures.Template(),
		testfixtures.Files(files...),
	)
	if err != nil {
		log.Fatalf("Failed in create fixtures: %s", err)
	}
	return loader
}

func (builder *fixtureLoaderBuilder) getFilesInDirectory(fixturesDirectory string) []string {
	fileNames, err := ioutil.ReadDir(fixturesDirectory)
	if err != nil {
		log.Fatalf("Error with get fixture files: %s", err)
	}
	var files []string
	fixtureExtension := ".yaml"
	for _, name := range fileNames {
		if !name.IsDir() {
			r, err := regexp.MatchString(fixtureExtension, name.Name())
			if err == nil && r {
				file := path.Join(fixturesDirectory, name.Name())
				files = append(files, file)
			}
		}
	}
	return files
}
