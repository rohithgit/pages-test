package persist

import (
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/models"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/speca/mdb"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/global"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/constants"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/utils"
	"time"
	"fmt"
)

type IUrlDB interface {
	GetAllUrls() ([]string, error)
	GetAllUrlsForApp( appId string) ([]models.SystemUrl, error)
	GetUrlCountForApp( appId string) (int, error)
	HasUrl(url string, appId string) bool
	InsertUrl(results *models.SystemUrl) (string, error)
	RemoveUrl(id string) error
	GetUrl( appId, url string) (*models.SystemUrl, error)
}

type UrlDB struct {
	url     mdb.IDb
}

func NewUrlDB() (IUrlDB, error) {
	var err error
	resDB := new(UrlDB)
	if resDB.url, err = initUrlDB(); err != nil {
		return nil, err
	}
	return resDB, nil
}

func initUrlDB() (mdb.IDb, error) {
	var (
		iDb mdb.IDb
		err error
	)
	iDb, err = mdb.NewMdb(global.Session, constants.DB_NAME, constants.SPECTRE_TEST_URLS)
	if err != nil {
		utils.SpectreLog.Errorf("InitDB Error: " + err.Error())
		return nil, err
	}
	return iDb, nil
}

func (ref *UrlDB) HasUrl(name string, appId string) bool {
	qp := map[string]interface{}{"url": name, "appId":appId}
	return ref.url.HasItem(qp)
}

func (ref *UrlDB) InsertUrl(url *models.SystemUrl) (string, error) {
	// if duplicate app then error
	if ref != nil && ref.HasUrl(url.Url, url.AppId) {
		utils.SpectreLog.Debug("Url %s exists in DB", url.Url)
		return "", models.NewError(409, "duplicate")
	}

	// if not, then create a new App and return the App Id
	id := ref.url.NewId()
	url.CreatedOn = time.Now()
	err := ref.url.UpsertId(id, url)
	if err != nil {
		utils.SpectreLog.Errorf("Insert app errored: %v", err)
		return "", err
	}

	return id, nil
}

func (ref *UrlDB) GetAllUrls() ([]string, error) {
	var result []models.SystemUrl
	// get all apps for this tenantid that are not deleted
	qp := map[string]interface{}{
	}
	err := ref.url.FindAll(qp, &result)

	if err != nil {
		return nil, err
	}
	set := make(map[string]string, len(result))
	var sysUrls []string
	sysUrls = make([]string, 0)
	for _, sysUrl := range result {
		if( set[sysUrl.Url] == ""){
			sysUrls = append(sysUrls, sysUrl.Url)
		}
		set[sysUrl.Url] = sysUrl.Url
	}
	return sysUrls, nil
}

func (ref *UrlDB) GetAllUrlsForApp( appId string) ([]models.SystemUrl, error) {
	var result []models.SystemUrl
	// get all apps for this tenantid that are not deleted
	qp := map[string]interface{}{
		"appId":  appId,
	}
	err := ref.url.FindAll(qp, &result)

	if err != nil {
		return nil, err
	}
	/*set := make(map[string]string, len(result))
	var sysUrls []string
	sysUrls = make([]string, 0)
	for _, sysUrl := range result {
		if( set[sysUrl.Url] == ""){
			sysUrls = append(sysUrls, sysUrl.Url)
		}
		set[sysUrl.Url] = sysUrl.Url
	}*/
	return result, nil
}

func (ref *UrlDB) GetUrl( url, appId string) (*models.SystemUrl, error) {
	var result *models.SystemUrl
	// get all apps for this tenantid that are not deleted
	qp := map[string]interface{}{
		"appId":  appId,
		"url": url,
	}
	err := ref.url.FindOne(qp, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (ref *UrlDB) GetUrlCountForApp( appId string) (int, error) {
	var result []*models.SystemUrl
	// get all apps for this tenantid that are not deleted
	qp := map[string]interface{}{
		"appId":  appId,
	}
	err := ref.url.FindAll(qp, &result)

	if err != nil {
		return 0, err
	}
	count := len(result)
	return count, nil
}

func (ref *UrlDB) RemoveUrl(url string) error {
	qp := map[string]interface{}{
		"url":  url,
	}

	err := ref.url.Remove(qp)
	if err != nil {
		utils.SpectreLog.Errorf("RemoveTest errored: %#v", err)
		return fmt.Errorf("RemoveTest errored: %#v", err)
	}

	return nil
}
