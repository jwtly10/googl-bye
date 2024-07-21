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

type FoundLinks struct {
	Url         string
	ExpandedUrl string
	File        string
	LineNumber  int
	Path        string
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

	// Parse the cloned repository
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

	err := filepath.Walk(dest, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Check file size
		if info.Size() > int64(maxFileSizeMB*1024*1024) {
			p.log.Infof("Skipping large file: %s (%.2f MB)", path, float64(info.Size())/(1024*1024))
			return nil
		}

		// Open the file
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		// TODO: Batch files ...
		// This does work great for now however as it makes getting line numbers really easy
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
					expandedUrl = "error"
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
			// We can just log the error on a line and continue to the next line, but only if there less than 3 errors in a row
			relPath, _ := filepath.Rel(dest, path)
			p.log.Warnf("Error scanning line %d in file %v: %v. Continuing.", lineNumber, relPath, err)
			errorsInFile++
			if errorsInFile > 3 {
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
