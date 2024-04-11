package tmdb

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

const base_url string = "http://api.themoviedb.org/3"

type TMDb struct {
	api_key string
	config  *tmdbConfig
}

func Init(api_key string) *TMDb {
	return &TMDb{api_key: api_key}
}

type filtered_output struct {
	Title        string `json:"title"`
	Artwork      string `json:"artwork"`
	Release_date string `json:"year"`
}

type tmdbResponse struct {
	Page          int
	Results       []tmdbResult
	Total_pages   int
	Total_results int
}

type tmdbResult struct {
	Adult          bool
	Name           string
	Backdrop_path  string
	Id             int
	Original_name  string
	Original_title string
	First_air_date string
	Release_date   string
	Poster_path    string
	Title          string
	Media_type     string
	Profile_path   string
}

type tmdbConfig struct {
	Images imageConfig
}

type imageConfig struct {
	Base_url        string
	Secure_base_url string

	Backdrop_sizes []string
	Logo_sizes     []string
	Poster_sizes   []string
	Profile_sizes  []string
	Still_sizes    []string
}

type movieMetadata struct {
	Id            int
	Media_type    string
	Backdrop_path string
	Poster_path   string
	Credits       tmdbCredits
	Config        *tmdbConfig
	Imdb_id       string
	Overview      string
	Title         string
	Release_date  string
}

type tmdbCredits struct {
	Id   int
	Cast []tmdbCast
	Crew []tmdbCrew
}

type tmdbCast struct {
	Character    string
	Name         string
	Profile_path string
}

type tmdbCrew struct {
	Department   string
	Name         string
	Job          string
	Profile_path string
}

func (tmdb *TMDb) MovieData(media_name string) (tmdbResponse, error) {
	results, err := tmdb.searchMovie(media_name)
	if err != nil {
		return results, err
	}

	return results, nil
}

func (tmdb *TMDb) searchTmdbMulti(media_name string) (tmdbResponse, error) {
	var resp tmdbResponse
	res, err := http.Get(base_url + "/search/multi?api_key=" + tmdb.api_key + "&query=" + url.QueryEscape(media_name))
	if err != nil {
		return resp, err
	}
	if res.StatusCode != 200 {
		return resp, error_status(res.StatusCode)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return tmdbResponse{}, err
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return tmdbResponse{}, err
	}
	return resp, nil
}

func (tmdb *TMDb) searchMovie(media_name string) (tmdbResponse, error) {
	var resp tmdbResponse
	res, err := http.Get(base_url + "/search/movie?api_key=" + tmdb.api_key + "&query=" + url.QueryEscape(media_name))
	if err != nil {
		return resp, err
	}
	if res.StatusCode != 200 {
		return resp, error_status(res.StatusCode)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return tmdbResponse{}, err
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return tmdbResponse{}, err
	}
	return resp, nil
}

func (tmdb *TMDb) searchTmdbTv(media_name string) (tmdbResponse, error) {
	var resp tmdbResponse
	res, err := http.Get(base_url + "/search/tv?api_key=" + tmdb.api_key + "&query=" + url.QueryEscape(media_name))
	if err != nil {
		return resp, err
	}
	if res.StatusCode != 200 {
		return resp, error_status(res.StatusCode)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return tmdbResponse{}, err
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return tmdbResponse{}, err
	}
	return resp, nil
}

func (tmdb *TMDb) getConfig() (*tmdbConfig, error) {
	if tmdb.config == nil || tmdb.config.Images.Base_url == "" {
		var conf = &tmdbConfig{}
		res, err := http.Get(base_url + "/configuration?api_key=" + tmdb.api_key)
		if err != nil {
			return conf, err
		}
		if res.StatusCode != 200 {
			return conf, error_status(res.StatusCode)
		}
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return &tmdbConfig{}, err
		}
		if err := json.Unmarshal(body, &conf); err != nil {
			return &tmdbConfig{}, err
		}
		tmdb.config = conf
	}
	return tmdb.config, nil
}

func (tmdb *TMDb) getMovieDetails(MediaId string) (movieMetadata, error) {
	var met movieMetadata
	res, err := http.Get(base_url + "/movie/" + MediaId + "?api_key=" + tmdb.api_key)
	if err != nil {
		return met, err
	}
	if res.StatusCode != 200 {
		return met, error_status(res.StatusCode)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return movieMetadata{}, err
	}
	if err := json.Unmarshal(body, &met); err != nil {
		return movieMetadata{}, err
	}
	return met, nil
}

func (tmdb *TMDb) getMovieCredits(MediaId string) (tmdbCredits, error) {
	var cred tmdbCredits
	res, err := http.Get(base_url + "/movie/" + MediaId + "/credits?api_key=" + tmdb.api_key)
	if err != nil {
		return cred, err
	}
	if res.StatusCode != 200 {
		return cred, error_status(res.StatusCode)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return tmdbCredits{}, err
	}
	if err := json.Unmarshal(body, &cred); err != nil {
		return tmdbCredits{}, err
	}
	return cred, nil
}

func (tmdb *TMDb) getTmdbTvDetails(MediaId string) (movieMetadata, error) {
	var met movieMetadata
	res, err := http.Get(base_url + "/tv/" + MediaId + "?api_key=" + tmdb.api_key)
	if err != nil {
		return met, err
	}
	if res.StatusCode != 200 {
		return met, error_status(res.StatusCode)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return movieMetadata{}, err
	}
	if err := json.Unmarshal(body, &met); err != nil {
		return movieMetadata{}, err
	}
	return met, nil
}

func (tmdb *TMDb) getTmdbTvCredits(MediaId string) (tmdbCredits, error) {
	var cred tmdbCredits
	res, err := http.Get(base_url + "/tv/" + MediaId + "/credits?api_key=" + tmdb.api_key)
	if err != nil {
		return cred, err
	}
	if res.StatusCode != 200 {
		return cred, error_status(res.StatusCode)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return tmdbCredits{}, err
	}
	if err := json.Unmarshal(body, &cred); err != nil {
		return tmdbCredits{}, err
	}
	return cred, nil
}

func (tmdb *TMDb) ToJSON(data string) (string, error) {
	var f filtered_output
	var det movieMetadata

	if err := json.Unmarshal([]byte(data), &det); err != nil {
		return "", err
	}

	f.Title = det.Title
	f.Release_date = det.Release_date
	if len(det.Release_date) > 4 {
		f.Release_date = det.Release_date[0:4]
	}
	size := det.poster_size("w154")
	f.Artwork = det.Config.Images.Base_url + size + det.Poster_path

	metadata, err := json.Marshal(f)
	if err != nil {
		return "", err
	}
	return string(metadata), nil
}

func (md *movieMetadata) poster_size(size string) string {
	if len(md.Config.Images.Poster_sizes) == 0 {
		return "original"
	}
	for i := range md.Config.Images.Poster_sizes {
		if md.Config.Images.Poster_sizes[i] == size {
			return size
		}
	}
	return md.Config.Images.Poster_sizes[0]
}

func error_status(status int) error {
	return errors.New(fmt.Sprintf("Status Code %d received from TMDb", status))
}

func main() {
	tmdb := TMDb{api_key: "71e68428e0a8d7f642158c4cc4c74f4c"}
	resp, err := tmdb.searchMovie("The Matrix")
	if err != nil {
		fmt.Println(err)
	}
	for i := range resp.Results {
		fmt.Println(resp.Results[i].Title)
		fmt.Println(resp.Results[i].Id)
		fmt.Println("https://vidsrc.xyz/embed/movie?tmdb=" + strconv.Itoa(resp.Results[i].Id))
	}
}
