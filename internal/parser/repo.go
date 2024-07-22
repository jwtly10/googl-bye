package parser

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/jwtly10/googl-bye/internal/common"
	"github.com/jwtly10/googl-bye/internal/models"
)

// This file handles finding repos to clone locally and parse

type RepoParser struct {
	git GitCmdLineI
	log common.Logger
}

func NewRepoParser(git GitCmdLineI, log common.Logger) *RepoParser {
	return &RepoParser{
		git: git,
		log: log,
	}
}

func (p *RepoParser) ParseRepository(repo models.RepositoryModel) ([]models.ParserLinksModel, error) {
	p.log.Infof("[%s] Parsing repo", fmt.Sprintf("%s/%s", repo.Author, repo.Name))
	tempDir, err := os.MkdirTemp("", fmt.Sprintf("%s%s%s%s%s", "repo-clone-", repo.Author, "-", repo.Name, "-"))
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tempDir)

	// Clone the repository
	err = p.git.Clone(repo.CloneUrl, tempDir)
	if err != nil {
		return nil, err
	}

	// Parse the files of cloned repository
	links, err := p.parseRepositoryFiles(repo, tempDir)
	if err != nil {
		return nil, err
	}

	return links, nil
}

const maxFileSizeMB = 10

func (p *RepoParser) parseRepositoryFiles(repo models.RepositoryModel, dest string) ([]models.ParserLinksModel, error) {
	p.log.Infof("[%s] Parsing files", fmt.Sprintf("%s/%s", repo.Author, repo.Name))
	var foundLinks []models.ParserLinksModel

	// TODO: Review the error handling
	// Currently if we find an error parsing a file, we just log it and continue
	// The application can still function as intended if a few files are unable to be processed
	// This could be because they are binary blobs, or some other minified file type, which we
	// probably dont care about as they will most likely not container shortend urls.
	err := filepath.Walk(dest, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			p.log.Errorf("Error walking dir: %v", err)
			return nil
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Check file size
		if info.Size() > int64(maxFileSizeMB*1024*1024) {
			p.log.Infof("[%s] Skipping large file: %s (%.2f MB)", fmt.Sprintf("%s/%s", repo.Author, repo.Name), path, float64(info.Size())/(1024*1024))
			return nil
		}

		// Open the file
		file, err := os.Open(path)
		if err != nil {
			p.log.Errorf("Error opening file: %v", err)
			return nil
		}
		defer file.Close()

		// TODO: Implement batch reading of binary data to reduce memory usage
		// Using bufio.Reader or io.Reader with a fixed-size buffer to prevent loading entire file to memory
		// Although current implementation works well, as we are limiting max file size and it also makes getting line numbers trivial
		scanner := bufio.NewScanner(file)
		lineNumber := 0

		errorsInFile := 0
		for scanner.Scan() {
			lineNumber++
			line := scanner.Text()

			// Check if the line contains a goo.gl link
			if strings.Contains(line, "goo.gl/") {
				relPath, _ := filepath.Rel(dest, path)
				url := extractGooGlLink(line)
				expandedUrl, err := expandGooGlLink(url)
				if err != nil {
					p.log.Errorf("Error expanding url '%s': %v", url, err)
					expandedUrl = fmt.Sprintf("ERROR: %s", err.Error())
				}
				foundLinks = append(foundLinks, models.ParserLinksModel{
					Url:         url,
					ExpandedUrl: expandedUrl,
					File:        relPath,
					LineNumber:  lineNumber,
					Path:        path,
				})
			}
		}

		if err := scanner.Err(); err != nil {
			relPath, _ := filepath.Rel(dest, path)
			// This error happens alot and is out of our control
			// (We can increase the buffer size of the scanner, but chosing against this for now)
			// TODO: Review this
			if err == bufio.ErrTooLong && strings.Contains(err.Error(), "token too long") {
				// We can skip this error, but it it happens more than 3 times we should follow up
				p.log.Debugf("Error scanning line %d in file %v: %v.", lineNumber, relPath, err)
				errorsInFile++
			} else {
				p.log.Errorf("Error scanning line %d in file %v: %v.", lineNumber, relPath, err)
				return err
			}

			if errorsInFile > 2 {
				p.log.Errorf("Error scanning line %d in file %v: %v. This is the 3rd time in this file.", lineNumber, relPath, err)
				errorsInFile = 0
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking the path %s: %v", dest, err)
	}

	return foundLinks, nil
}

func extractGooGlLink(line string) string {
	re := regexp.MustCompile(`(?i)(?:https?://)?goo\.gl/[a-zA-Z0-9_-]+`)
	match := re.FindString(line)
	return match
}

func expandGooGlLink(link string) (string, error) {
	// Check that the url starts with https
	if !strings.HasPrefix(link, "http://") && !strings.HasPrefix(link, "https://") {
		link = "https://" + link
	}

	// This just ensures we always try to call https short links
	if strings.HasPrefix(link, "http://") {
		link = strings.Replace(link, "http://", "https://", 1)
	}

	// Create a new HTTP client that doesn't automatically follow redirects
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Get(link)
	if err != nil {
		return "", fmt.Errorf("error making request to %s: %v", link, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 && resp.StatusCode < 400 {
		redirectURL := resp.Header.Get("Location")
		if redirectURL != "" {
			return redirectURL, nil
		}
		return "", fmt.Errorf("redirect URL not found for %s", link)
	}

	return "", fmt.Errorf("unexpected status code %d for %s", resp.StatusCode, link)
}
