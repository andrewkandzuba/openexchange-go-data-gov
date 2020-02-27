package model

import (
	"bufio"
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"strings"
	"testing"
)

func Test_ParseArticleFromJson_Success(t *testing.T) {
	handle, err := os.Open("testdata/article.json")
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	text := ""
	scanner := bufio.NewScanner(handle)
	for scanner.Scan() {
		text += scanner.Text()
	}

	var a Article
	r := bytes.NewReader([]byte(text))
	err = json.NewDecoder(r).Decode(&a)

	assert.Nil(t, err)
	assert.Equal(t, "news", a.Type)
	assert.Equal(t, 48, a.Nid)
	assert.Equal(t, "Former Secretary of Commerce Pritzker’s Official Portrait Unveiled", a.Label)
	assert.Equal(t, 1508360246, a.Created)
	assert.Equal(t, 1513969347, a.Update)
	assert.Equal(t, "https://www.commerce.gov/news/press-releases/2017/10/former-secretary-commerce-pritzkers-official-portrait-unveiled", a.Href)
	assert.True(t, strings.Contains(a.Body, "Commerce Wilbur Ross attended the unveiling of former"))
	assert.Equal(t, "FOR IMMEDIATE RELEASE", a.Status)
	assert.Equal(t, "86e95078-92cf-4ff7-a55d-cb2fbffd8b85", a.UUID)
	assert.NotNil(t, a.AdminOfficials)
	assert.Equal(t, 2, len(a.AdminOfficials))
	assert.Equal(t, "9", a.AdminOfficials[0].Id)
	assert.Equal(t, "Wilbur Ross", a.AdminOfficials[0].Label)
	assert.Equal(t, "https://www.commerce.gov/about/leadership/wilbur-ross", a.AdminOfficials[0].Href)
	assert.Equal(t, "10", a.AdminOfficials[1].Id)
	assert.Equal(t, "James Bond", a.AdminOfficials[1].Label)
	assert.Equal(t, "https://www.commerce.gov/about/leadership/james-bond", a.AdminOfficials[1].Href)
}