package function

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type Listing struct {
	Selftext  string
	Title     string
	Author    string
	Permalink string
	URL       string
	Preview   struct {
		Images []struct {
			Source struct {
				URL    string
				Width  int
				Height int
			}
		}
	}
}

type SubredditResult struct {
	Kind string
	Data struct {
		Children []struct {
			Kind string
			Data Listing
		}
	}
}

type Map struct {
	Name         string                 `json:"name"`
	DownloadCode string                 `json:"download_code"`
	Description  map[string]interface{} `json:"description"`
	Type         string                 `json:"type"`
	UserID       int64                  `json:"user_id"`
	Photos       []string               `json:"photos"`
	Tags         []string               `json:"tags"`
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)

	if len(value) == 0 {
		return fallback
	}
	return value
}

var wargroovyAPI = getenv("WARGROOVY_API", "http://localhost:4000")
var wargroovyWebAPI = getenv("WARGROOVY_WEB_API", "http://localhost:8000")
var userToken = getenv("USER_TOKEN", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBZG1pbiI6ZmFsc2UsIlVzZXJJRCI6MX0.JGdebtRFzMqs5pUDyq5P8Yzh5fIF44ttE1dEGw7IkhY")

func hasAtLeastTwoNumbers(str string) bool {
	var b = 0

	for _, value := range str {
		switch {
		case value >= '0' && value <= '9':
			b++
		}
	}

	return b > 1
}

func getMapCode(listing *Listing) (bool, string) {
	isCodeLength := regexp.MustCompile(`(\d|\w){8}`)
	code := ""
	combos := isCodeLength.FindAllString(listing.Title, -1)

	if len(combos) < 1 {
		return false, code
	}

	b := false

	for _, word := range combos {
		if hasAtLeastTwoNumbers(word) {
			b = true
			code = word
		}
	}

	return b, code
}

func SliceUniqMap(s []string) []string {
	seen := make(map[string]struct{}, len(s))
	j := 0
	for _, v := range s {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		s[j] = v
		j++
	}
	return s[:j]
}

func getMapTags(listing *Listing) []string {
	tags := make([]string, 0)
	definedTags := []string{
		"1v1",
		"2v2",
		"3v3",
		"4v4",
		"co-op",
		"coop",
		"ai",
		"funny",
		"humor",
		"challenging",
		"easy",
		"hard",
		"moba",
		"1v3",
		"1v4",
		"1v2",
		"2v1",
		"4v1",
		"3v1",
		"3v2",
		"2v3",
		"hero",
		"advance wars",
		"remake",
	}

	// check title and self text
	for _, t := range definedTags {
		if strings.Contains(listing.Title, t) {
			tags = append(tags, t)
		} else if strings.Contains(listing.Selftext, t) {
			tags = append(tags, t)
		}
	}

	return SliceUniqMap(tags)
}

func getType(listing *Listing) string {
	if strings.Contains(listing.Title, "campaign") {
		return "campaign"
	} else if strings.Contains(listing.Title, "puzzle") {
		return "puzzle"
	}

	return "skirmish"
}

func getDescriptionRaw(listing *Listing) (map[string]interface{}, error) {
	s := fmt.Sprintf("{\"markdown\":\"## %s\\n*created by %s*\\n%s\\n[reddit url](https://reddit.com%s)\"}", listing.Title, listing.Author, jsonEscape(listing.Selftext), listing.Permalink)

	// log.WithField("markdown", s).Info("sending markdown to wargroovy-web")

	var j = []byte(s)

	req, err := http.NewRequest("POST", wargroovyWebAPI+"/v1/api/description", bytes.NewBuffer(j))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	type WargroovyWebApiResponse struct {
		Success bool
		Data    map[string]interface{}
	}

	result := WargroovyWebApiResponse{}

	json.NewDecoder(resp.Body).Decode(&result)

	return result.Data, nil
}

func mapCodeIsUnique(code string) (bool, error) {
	resp, err := http.Get(wargroovyAPI + "/v1/api/map/byDownloadCode/" + code + "?jwt=" + userToken)

	if err != nil {
		return false, err
	}

	defer resp.Body.Close()

	type WargroovyMapApiResponse struct {
		Message string
		Status  bool
		Map     map[string]interface{}
	}

	result := WargroovyMapApiResponse{}

	json.NewDecoder(resp.Body).Decode(&result)

	if !result.Status {
		return true, nil
	}

	return false, nil
}

func saveMap(m *Map) error {
	j, err := json.Marshal(m)

	resp, err := http.Post(wargroovyAPI+"/v1/api/map?jwt="+userToken, "application/json", bytes.NewBuffer(j))

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	type WargroovyCreateApiResponse struct {
		Message string
		Status  bool
		Map     int
	}

	result := WargroovyCreateApiResponse{}

	json.NewDecoder(resp.Body).Decode(&result)

	if !result.Status {
		return errors.New(result.Message)
	}

	log.Info("created map")

	return nil
}

func uploadPhoto(photoURL string) (string, error) {
	parsedPhotoURL := strings.Replace(photoURL, "&amp;", "&", -1)

	log.WithField("photoURL", parsedPhotoURL).Info("grabbing photo from")

	resp, err := http.Get(parsedPhotoURL)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	v, err := ioutil.ReadAll(resp.Body)

	// resp, err := http.Post(wargroovyAPI+"/v1/api/photo", "")

	// Prepare a form that you will submit to that URL.

	var b bytes.Buffer
	var fw io.Writer

	w := multipart.NewWriter(&b)

	if fw, err = w.CreateFormFile("photos", "photo"); err != nil {
		return "", err
	}

	valueString := strings.NewReader(string(v))

	if _, err := io.Copy(fw, valueString); err != nil {
		return "", err
	}

	w.Close()

	client := &http.Client{}

	req, err := http.NewRequest("POST", wargroovyAPI+"/v1/api/photo", &b)

	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err = client.Do(req)

	if err != nil {
		return "", err
	}

	type WargroovyApiPhotoResponse struct {
		Status bool
		URLs   []string
	}

	result := WargroovyApiPhotoResponse{}

	json.NewDecoder(resp.Body).Decode(&result)

	log.WithField("body", resp.Body).Info("uploadPhoto: response body from API")

	log.WithFields(log.Fields{
		"status": result.Status,
		"urls":   result.URLs,
	}).Info("uploadPhoto: result from API")

	if result.Status {

		return result.URLs[0], nil
	}

	return "", nil
}

func jsonEscape(i string) string {
	b, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	s := string(b)
	return s[1 : len(s)-1]
}

func RedditBot(http.ResponseWriter, *http.Request) {
	subreddits := []string{
		"customgroove",
		"wargroove",
		"WargrooveCompetitive",
	}

	for _, sr := range subreddits {
		var str strings.Builder

		str.WriteString("https://reddit.com/r/")
		str.WriteString(sr)
		str.WriteString("/new.json?limit=100&show=all")

		log.WithField("url", str.String()).Info("reddit: requesting")

		client := &http.Client{}

		req, err := http.NewRequest("GET", str.String(), nil)

		if err != nil {
			log.Fatalln(err)
		}

		req.Header.Set("User-Agent", "Golang_Spider_Bot/3.0")

		resp, err := client.Do(req)

		if err != nil {
			log.Fatalln(err)
		}

		defer resp.Body.Close()

		result := SubredditResult{}

		json.NewDecoder(resp.Body).Decode(&result)

		for _, listing := range result.Data.Children {
			isMap, code := getMapCode(&listing.Data)

			if isMap {
				m := &Map{
					Name:         strings.Replace(listing.Data.Title, `\`, `\\`, -1),
					DownloadCode: code,
					Type:         getType(&listing.Data),
					UserID:       4,
					Tags:         getMapTags(&listing.Data),
				}

				isUniq, err := mapCodeIsUnique(m.DownloadCode)

				if err != nil {
					log.WithField("error", err).Error("wargroovy: could not getBy download code")
				}

				if !isUniq {
					log.Info("download code already exists in db. skipping.")
				} else {
					description, err := getDescriptionRaw(&listing.Data)

					if err == nil {
						m.Description = description
					}

					if len(listing.Data.Preview.Images) > 0 {
						if listing.Data.Preview.Images[0].Source.URL != "" {
							photoURL, _ := uploadPhoto(listing.Data.Preview.Images[0].Source.URL)

							if photoURL != "" {
								m.Photos = []string{photoURL}
							}
						}
					}

					log.WithFields(log.Fields{
						"code":        m.DownloadCode,
						"name":        m.Name,
						"description": m.Description,
					}).Info("map")

					err = saveMap(m)

					if err != nil {
						log.WithField("error", err).Error(err.Error())
					}
				}
			}
		}
	}
}
