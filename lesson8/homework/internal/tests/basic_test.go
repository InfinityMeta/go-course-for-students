package tests

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"homework8/internal/app"
)

func TestCreateAd(t *testing.T) {
	client := getTestClient()

	_, _ = client.createUser("Bob", "bob@box.com")

	response, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)
	assert.Zero(t, response.Data.ID)
	assert.Equal(t, response.Data.Title, "hello")
	assert.Equal(t, response.Data.Text, "world")
	assert.Equal(t, response.Data.AuthorID, int64(0))
	assert.False(t, response.Data.Published)
}

func TestChangeAdStatus(t *testing.T) {
	client := getTestClient()

	_, _ = client.createUser("Bob", "bob@box.com")

	response, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)

	response, err = client.changeAdStatus(0, response.Data.ID, true)
	assert.NoError(t, err)
	assert.True(t, response.Data.Published)

	response, err = client.changeAdStatus(0, response.Data.ID, false)
	assert.NoError(t, err)
	assert.False(t, response.Data.Published)

	response, err = client.changeAdStatus(0, response.Data.ID, false)
	assert.NoError(t, err)
	assert.False(t, response.Data.Published)
}

func TestUpdateAd(t *testing.T) {
	client := getTestClient()

	_, _ = client.createUser("Bob", "bob@box.com")

	response, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)

	response, err = client.updateAd(0, response.Data.ID, "привет", "мир")
	assert.NoError(t, err)
	assert.Equal(t, response.Data.Title, "привет")
	assert.Equal(t, response.Data.Text, "мир")
}

func TestGetAdByID(t *testing.T) {
	client := getTestClient()

	_, _ = client.createUser("Bob", "bob@box.com")

	response, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)

	gotAd, err := client.getAdByID(response.Data.ID)
	assert.NoError(t, err)
	assert.Equal(t, gotAd.Data.ID, response.Data.ID)
	assert.Equal(t, gotAd.Data.Title, response.Data.Title)
	assert.Equal(t, gotAd.Data.Text, response.Data.Text)
	assert.Equal(t, gotAd.Data.AuthorID, gotAd.Data.AuthorID)
	assert.False(t, gotAd.Data.Published)
}

func TestCreateUser(t *testing.T) {
	client := getTestClient()

	response, err := client.createUser("Oleg", "boss@gmail.com")
	assert.NoError(t, err)
	assert.Zero(t, response.Data.ID)
	assert.Equal(t, response.Data.Nickname, "Oleg")
	assert.Equal(t, response.Data.Email, "boss@gmail.com")
}

func TestUpdateUser(t *testing.T) {
	client := getTestClient()

	response, err := client.createUser("Oleg", "boss@gmail.com")
	assert.NoError(t, err)

	response, err = client.updateUser(response.Data.ID, "Vladimir", "newboss@gmail.com")
	assert.NoError(t, err)
	assert.Equal(t, response.Data.Nickname, "Vladimir")
	assert.Equal(t, response.Data.Email, "newboss@gmail.com")
}

func TestSearchAdByName(t *testing.T) {
	client := getTestClient()

	_, _ = client.createUser("Bob", "bob@box.com")

	createdAd, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)

	response, err := client.searchAdByName(createdAd.Data.Title)
	assert.NoError(t, err)
	assert.Equal(t, response.Data.AuthorID, createdAd.Data.AuthorID)
	assert.Equal(t, response.Data.Title, createdAd.Data.Title)
	assert.Equal(t, response.Data.Text, createdAd.Data.Text)

	createdAd, err = client.createAd(0, "red sedan mercedes", "buy red sedan mercedes good condition expensive")
	assert.NoError(t, err)

	response, err = client.searchAdByName("sedan")
	assert.NoError(t, err)
	assert.Equal(t, response.Data.AuthorID, createdAd.Data.AuthorID)
	assert.Equal(t, response.Data.Title, createdAd.Data.Title)
	assert.Equal(t, response.Data.Text, createdAd.Data.Text)
}

func TestFilterAdsOneUser(t *testing.T) {
	client := getTestClient()

	_, _ = client.createUser("Bob", "bob@box.com")
	_, _ = client.createUser("Dob", "dob@box.com")

	ad1, err := client.createAd(0, "hello1", "world1")
	assert.NoError(t, err)

	_, err = client.changeAdStatus(0, ad1.Data.ID, true)
	assert.NoError(t, err)

	ad2, err := client.createAd(0, "hello2", "world2")
	assert.NoError(t, err)

	_, err = client.changeAdStatus(0, ad2.Data.ID, true)
	assert.NoError(t, err)

	ad3, err := client.createAd(0, "hello3", "world3")
	assert.NoError(t, err)

	_, err = client.changeAdStatus(0, ad3.Data.ID, true)
	assert.NoError(t, err)

	ad4, err := client.createAd(1, "hello4", "world4")
	assert.NoError(t, err)

	_, err = client.changeAdStatus(1, ad4.Data.ID, true)
	assert.NoError(t, err)

	ads, err := client.filterAds(app.WithAuthorID(0))
	assert.NoError(t, err)
	assert.Len(t, ads.Data, 3)
}

func TestFilterAdsAfter(t *testing.T) {
	client := getTestClient()

	_, _ = client.createUser("Bob", "bob@box.com")

	ad1, err := client.createAd(0, "hello1", "world1")
	assert.NoError(t, err)

	_, err = client.changeAdStatus(0, ad1.Data.ID, true)
	assert.NoError(t, err)

	ad2, err := client.createAd(0, "hello2", "world2")
	assert.NoError(t, err)

	_, err = client.changeAdStatus(0, ad2.Data.ID, true)
	assert.NoError(t, err)

	timePoint := time.Now().UTC()

	ad3, err := client.createAd(0, "hello3", "world3")
	assert.NoError(t, err)

	_, err = client.changeAdStatus(0, ad3.Data.ID, true)
	assert.NoError(t, err)

	ad4, err := client.createAd(0, "hello4", "world4")
	assert.NoError(t, err)

	_, err = client.changeAdStatus(0, ad4.Data.ID, true)
	assert.NoError(t, err)

	ad5, err := client.createAd(0, "hello5", "world5")
	assert.NoError(t, err)

	_, err = client.changeAdStatus(0, ad5.Data.ID, true)
	assert.NoError(t, err)

	ads, err := client.filterAds(app.WithPublishedAfter(timePoint))
	assert.NoError(t, err)
	assert.Len(t, ads.Data, 3)
}

func TestFilterAdsBefore(t *testing.T) {
	client := getTestClient()

	_, _ = client.createUser("Bob", "bob@box.com")

	ad1, err := client.createAd(0, "hello1", "world1")
	assert.NoError(t, err)

	_, err = client.changeAdStatus(0, ad1.Data.ID, true)
	assert.NoError(t, err)

	ad2, err := client.createAd(0, "hello2", "world2")
	assert.NoError(t, err)

	_, err = client.changeAdStatus(0, ad2.Data.ID, true)
	assert.NoError(t, err)

	timePoint := time.Now().UTC()

	ad3, err := client.createAd(0, "hello3", "world3")
	assert.NoError(t, err)

	_, err = client.changeAdStatus(0, ad3.Data.ID, true)
	assert.NoError(t, err)

	ad4, err := client.createAd(0, "hello4", "world4")
	assert.NoError(t, err)

	_, err = client.changeAdStatus(0, ad4.Data.ID, true)
	assert.NoError(t, err)

	ad5, err := client.createAd(0, "hello5", "world5")
	assert.NoError(t, err)

	_, err = client.changeAdStatus(0, ad5.Data.ID, true)
	assert.NoError(t, err)

	ads, err := client.filterAds(app.WithPublishedBefore(timePoint))
	assert.NoError(t, err)
	assert.Len(t, ads.Data, 2)
}

func TestFilterAdsMultiple(t *testing.T) {
	client := getTestClient()

	_, _ = client.createUser("Bob", "bob@box.com")
	_, _ = client.createUser("Dob", "dob@box.com")

	ad1, err := client.createAd(0, "hello1", "world1")
	assert.NoError(t, err)

	_, err = client.changeAdStatus(0, ad1.Data.ID, true)
	assert.NoError(t, err)

	timeAfter := time.Now().UTC()

	ad3, err := client.createAd(0, "hello3", "world3")
	assert.NoError(t, err)

	_, err = client.changeAdStatus(0, ad3.Data.ID, true)
	assert.NoError(t, err)

	ad4, err := client.createAd(1, "hello4", "world4")
	assert.NoError(t, err)

	_, err = client.changeAdStatus(1, ad4.Data.ID, true)
	assert.NoError(t, err)

	timeBefore := time.Now().UTC()

	ad5, err := client.createAd(1, "hello5", "world5")
	assert.NoError(t, err)

	_, err = client.changeAdStatus(1, ad5.Data.ID, true)
	assert.NoError(t, err)

	ads, err := client.filterAds(app.WithAuthorID(0), app.WithPublishedAfter(timeAfter), app.WithPublishedBefore(timeBefore))
	assert.NoError(t, err)
	assert.Len(t, ads.Data, 1)
}

func TestUserExist(t *testing.T) {

	client := getTestClient()

	_, _ = client.createUser("Bob", "bob@box.com")
	_, err := client.createAd(0, "hello1", "world1")
	assert.NoError(t, err)

	_, err = client.createAd(1, "hello1", "world1")
	assert.ErrorIs(t, err, ErrNotFound)

}
