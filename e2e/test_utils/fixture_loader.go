package test_utils

import (
	"database/sql"
	"github.com/go-testfixtures/testfixtures/v3"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"log"
	"path"
	"regexp"
	"text/template"
)

type FixtureLoaderBuilder interface {
	Build(db *sql.DB, fixturesDirectory string) *testfixtures.Loader
	getFilesInDirectory(fixturesDirectory string) []string
}

type fixtureLoaderBuilder struct {
	fixtureFunctions template.FuncMap
}

func NewFixtureLoaderBuilder() FixtureLoaderBuilder {
	fixtureFunctions := template.FuncMap{
		"GenPasswordHash": generatePasswordHash,
	}
	return &fixtureLoaderBuilder{
		fixtureFunctions: fixtureFunctions,
	}
}

func (builder *fixtureLoaderBuilder) Build(db *sql.DB, fixturesDirectory string) *testfixtures.Loader {
	files := builder.getFilesInDirectory(fixturesDirectory)
	loader, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgresql"),
		testfixtures.Template(),
		testfixtures.TemplateFuncs(builder.fixtureFunctions),
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

func generatePasswordHash(password string) string {
	passwordBytes := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.MinCost)
	if err != nil {
		log.Fatalf("Cannot generate password hash: %s", hash)
	}
	hashString := string(hash)
	return hashString
}
