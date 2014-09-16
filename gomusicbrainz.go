/*
 * Copyright (c) 2014 Michael Wendland
 *
 * Permission is hereby granted, free of charge, to any person obtaining a
 * copy of this software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation
 * the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the
 * Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
 * FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
 * IN THE SOFTWARE.
 *
 * 	Authors:
 * 		Michael Wendland <michael@michiwend.com>
 */

/*
Package gomusicbrainz implements a MusicBrainz WS2 client library.
*/
package gomusicbrainz

import (
	"encoding/xml"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

// NewWS2Client returns a new instance of WS2Client with default values.
func NewWS2Client() *WS2Client {
	c := WS2Client{}

	c.WS2RootURL, _ = url.Parse("https://musicbrainz.org/ws/2")

	// Provide meaningful User-Agent informations.
	c.SetClientInfo(
		"GoMusicBrainz - a Golang WS2 client",
		"0.0.1-beta",
		"michael@michiwend.com",
	)

	return &c
}

// WS2Client defines a Go client for the MusicBrainz Web Service 2.
type WS2Client struct {
	WS2RootURL *url.URL // The API root URL

	userAgentHeader string
}

func (c *WS2Client) getReqeust(data interface{}, params url.Values, endpoint string) error {

	client := &http.Client{}

	req, err := http.NewRequest("GET", c.WS2RootURL.String()+endpoint+"?"+params.Encode(), nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("User-Agent", c.userAgentHeader)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	decoder := xml.NewDecoder(resp.Body)

	if err = decoder.Decode(data); err != nil {
		return err
	}
	return nil
}

// intParamToString returns an empty string for -1.
func intParamToString(i int) string {
	if i == -1 {
		return ""
	} else {
		return strconv.Itoa(i)
	}
}

// SetClientInfo sets the HTTP user-agent header of the WS2Client. Please
// provide meaningful information about your application as described at:
// https://musicbrainz.org/doc/XML_Web_Service/Rate_Limiting#Provide_meaningful_User-Agent_strings
func (c *WS2Client) SetClientInfo(application string, version string, contact string) {
	c.userAgentHeader = application + "/" + version + " ( " + contact + " ) "
}

// SearchArtist queries MusicBrainz' Search Server for Artists.
// searchTerm follows the Apache Lucene syntax. If no fields were specified the
// Search Server searches for searchTerm in any of the fields artist, sortname
// and alias. For a list of all valid search fields visit
// https://musicbrainz.org/doc/Development/XML_Web_Service/Version_2/Search#Artist
// limit defines how many entries will be returned by the server (allowed
// range 1-100, defaults to 25). offset can be used for result pagination. -1
// can be set for both limit and offset to use the default values.
func (c *WS2Client) SearchArtist(searchTerm string, limit int, offset int) (*ArtistResponse, error) {

	result := artistResult{}
	endpoint := "/artist"
	params := url.Values{
		"query":  {searchTerm},
		"limit":  {intParamToString(limit)},
		"offset": {intParamToString(offset)},
	}

	if err := c.getReqeust(&result, params, endpoint); err != nil {
		return nil, err
	}

	return &result.Resonse, nil
}

// SearchRelease queries MusicBrainz' Search Server for Releases.
// searchTerm follows the Apache Lucene syntax. If no fields were specified the
// Search Server searches the release field only. For a list of all valid
// search fields visit
// https://musicbrainz.org/doc/Development/XML_Web_Service/Version_2/Search#Release
// limit defines how many entries will be returned by the server (allowed
// range 1-100, defaults to 25). offset can be used for result pagination. -1
// can be set for both limit and offset to use the default values.
func (c *WS2Client) SearchRelease(searchTerm string, limit int, offset int) (*ReleaseResponse, error) {

	result := releaseResult{}
	endpoint := "/release"
	params := url.Values{
		"query":  {searchTerm},
		"limit":  {intParamToString(limit)},
		"offset": {intParamToString(offset)},
	}

	if err := c.getReqeust(&result, params, endpoint); err != nil {
		return nil, err
	}

	return &result.Response, nil
}

// SearchReleaseGroup queries MusicBrainz' Search Server for ReleaseGroups.
// searchTerm follows the Apache Lucene syntax. If no fields were specified the
// Search Server searches the releasegroup field only. For a list of all valid
// search fields visit
// https://musicbrainz.org/doc/Development/XML_Web_Service/Version_2/Search#Release_Group
// limit defines how many entries will be returned by the server (allowed
// range 1-100, defaults to 25). offset can be used for result pagination. -1
// can be set for both limit and offset to use the default values.
func (c *WS2Client) SearchReleaseGroup(searchTerm string, limit int, offset int) (*ReleaseGroupResponse, error) {

	result := releaseGroupResult{}
	endpoint := "/release-group"
	params := url.Values{
		"query":  {searchTerm},
		"limit":  {intParamToString(limit)},
		"offset": {intParamToString(offset)},
	}

	if err := c.getReqeust(&result, params, endpoint); err != nil {
		return nil, err
	}

	return &result.Response, nil
}

// SearchTag queries MusicBrainz' Search Server for Tags.
// searchTerm follows the Apache Lucene syntax. The Tag index contains only the
// tag field. For more information visit
// https://musicbrainz.org/doc/Development/XML_Web_Service/Version_2/Search#Tag
// limit defines how many entries will be returned by the server (allowed
// range 1-100, defaults to 25). offset can be used for result pagination. -1
// can be set for both limit and offset to use the default values.
func (c *WS2Client) SearchTag(searchTerm string, limit int, offset int) (*TagResponse, error) {

	result := tagResult{}
	endpoint := "/tag"
	params := url.Values{
		"query":  {searchTerm},
		"limit":  {intParamToString(limit)},
		"offset": {intParamToString(offset)},
	}

	if err := c.getReqeust(&result, params, endpoint); err != nil {
		return nil, err
	}

	return &result.Response, nil
}
