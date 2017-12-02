// Package persist handles provider data persistence to an etcd data store.
package persist

import (
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/models"
	// "bitbucket-eng-sjc1.cisco.com/bitbucket/specnl/spectre-base-microservice/log"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/speca/mdb"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/constants"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/global"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/utils"
	"time"
	"errors"
	"fmt"
)

type ITestResultDB interface {
	FindTests(query interface{}) ([]*models.ResultsStorage, error)
	GetUrl(id string) (*models.ResultsStorage, error)
	GetTest(id string) (*models.ResultsStorage, error)
	GetAllTests(url string) ([]*models.ResultsStorage, error)
	GetAllTestResults(tenantId string) ([]*models.ResultsStorage, error)
	HasTest(test *models.ResultsStorage) bool
	HasTestId(id string) bool
	InsertTest(results *models.ResultsStorage) (string, error)
	RemoveTest(id string) error
	UpdateTest(id string, results *models.ResultsStorage) error
}

type TestResultDB struct {
	testResults     mdb.IDb
}

func NewTestResultDB() (ITestResultDB, error) {
	var err error
	resDB := new(TestResultDB)
	if resDB.testResults, err = initDB(); err != nil {
		return nil, err
	}
	return resDB, nil
}

func initDB() (mdb.IDb, error) {
	var (
		iDb mdb.IDb
		err error
	)
	iDb, err = mdb.NewMdb(global.Session, constants.DB_NAME, constants.TEST_COLLECTION)
	if err != nil {
		utils.SpectreLog.Errorf("InitDB Error: " + err.Error())
		return nil, err
	}
	return iDb, nil
}

func (ref *TestResultDB) HasTest(result *models.ResultsStorage) bool {
	qp := map[string]interface{}{"testUrl": result.LookupUrl}
	return ref.testResults.HasItem(qp)
}

func (ref *TestResultDB) HasTestId(id string) bool {
	return ref.testResults.HasItemId(id)
}

func (ref *TestResultDB) GetTest(name string) (*models.ResultsStorage, error) {
	var result = new(models.ResultsStorage)
	qp := map[string]interface{}{"lookupUrl": name}
	err := ref.testResults.FindOne(qp, result)
	if err != nil {
		return nil, err
	}

	// don't return deleted element
	if result.DeletedOn != ref.testResults.ZeroTime() {
		return nil, errors.New("not found")
	}

	return result, nil
}

func (ref *TestResultDB) GetUrl(name string) (*models.ResultsStorage, error) {
	var result = new(models.ResultsStorage)
	qp := map[string]interface{}{"url": name}
	err := ref.testResults.FindOne(qp, result)
	if err != nil {
		return nil, err
	}

	// don't return deleted element
	if result.DeletedOn != ref.testResults.ZeroTime() {
		return nil, errors.New("not found")
	}

	return result, nil
}

func (ref *TestResultDB) GetAllTests(tenantId string) ([]*models.ResultsStorage, error) {
	var result []*models.ResultsStorage
	// get all apps for this tenantid that are not deleted
	qp := map[string]interface{}{
		"tenantId":  tenantId,
		"deletedOn": ref.testResults.ZeroTime(),
	}
	err := ref.testResults.FindAll(qp, &result)
	//ref.testResults.
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (ref *TestResultDB) GetAllTestResults(tenantId string) ([]*models.ResultsStorage, error) {
	var result []*models.ResultsStorage
	// get all apps for this tenantid that are not deleted
	qp := map[string]interface{}{}
	err := ref.testResults.FindAll(qp, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (ref *TestResultDB) FindTests(query interface{}) ([]*models.ResultsStorage, error) {
	var result []*models.ResultsStorage
	err := ref.testResults.FindAll(query, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (ref *TestResultDB) InsertTest(test *models.ResultsStorage) (string, error) {
	// if duplicate app then error
	if ref.HasTest(test) {
		return "", models.NewError(409, "duplicate")
	}

	// if not, then create a new App and return the App Id
	id := ref.testResults.NewId()
	test.TestID = id
	test.CreatedOn = time.Now()
	test.UpdatedOn = time.Now()
	test.DeletedOn = ref.testResults.ZeroTime()
	err := ref.testResults.UpsertId(id, test)
	if err != nil {
		utils.SpectreLog.Errorf("Insert app errored: %v", err)
		return "", err
	}

	return id, nil
}

func (ref *TestResultDB) UpdateTest(id string, test *models.ResultsStorage) error {
	if id != test.TestID {
		return errors.New("UpdateAppId Error: Invalid App Id - " + id)
	}

	// verify app
	if err := verifyTestResults(ref.testResults, id); err != nil {
		return err
	}

	// check for duplicate
	if ref.HasTest(test) {
		err := models.NewError(409, "duplicate")
		utils.SpectreLog.Errorf("Error: " + err.Error())
		return err
	}

	// update the app
	test.UpdatedOn = time.Now()
	err := ref.testResults.UpdateId(id, test)
	if err != nil {
		utils.SpectreLog.Errorf("Update app errored: %#v", err)
		return fmt.Errorf("Update app errored: %#v", err)
	}
	return nil
}

func (ref *TestResultDB) RemoveTest(id string) error {
	// verify app - if it's already deleted then return
	if err := verifyTestResults(ref.testResults, id); err != nil {
		if ref.testResults.IsDeleted(err) {
			return nil
		}
		utils.SpectreLog.Errorf("RemoveTest errored: %#v", err)
		return fmt.Errorf("RemoveTest errored: %#v", err)
	}

	// otherwise, mark app as deleted
	update := map[string]interface{}{
		"$set": map[string]interface{}{
			"deletedOn": time.Now(),
		},
	}
	err := ref.testResults.UpdateId(id, update)
	if err != nil {
		utils.SpectreLog.Errorf("RemoveTest errored: %#v", err)
		return fmt.Errorf("RemoveTest errored: %#v", err)
	}

	return nil
}

func verifyTestResults(db mdb.IDb, id string) error {
	results := &models.ResultsStorage{}
	err := db.FindId(id, results)
	if err != nil {
		if db.IsNotFound(err) {
			err := models.NewError(404, "not found")
			utils.SpectreLog.Errorf("App Verification Error: " + err.Error())
			return err
		}
		utils.SpectreLog.Errorf("App Verification Error: " + err.Error())
		return err
	}

	if results.DeletedOn != db.ZeroTime() {
		err := models.NewError(500, "deleted")
		utils.SpectreLog.Errorf("App Verification Error: " + err.Error())
		return err
	}

	return nil
}

type TestResultDBMock struct {
	MockFindTests    func(query interface{}) ([]*models.ResultsStorage, error)
	MockGetUrl      func(id string) (*models.ResultsStorage, error)
	MockGetTest      func(id string) (*models.ResultsStorage, error)
	MockGetAllTests  func(tenantId string) ([]*models.ResultsStorage, error)
	MockGetAllTestResults func(tenantId string) ([]*models.ResultsStorage, error)
	MockHasTest      func(app *models.ResultsStorage) bool
	MockHasTestId    func(id string) bool
	MockRemoveTest func(id string) error
	MockUpdateTest func(id string, app *models.ResultsStorage) error
	MockInsertTest   func(app *models.ResultsStorage) (string, error)
}

func NewTestResultDBMock() ITestResultDB {
	return &TestResultDBMock{}
}

func (ref *TestResultDBMock) FindTests(query interface{}) ([]*models.ResultsStorage, error) {
	return ref.MockFindTests(query)
}
func (ref *TestResultDBMock) GetTest(id string) (*models.ResultsStorage, error) {
	return ref.MockGetTest(id)
}
func (ref *TestResultDBMock) GetUrl(id string) (*models.ResultsStorage, error) {
	return ref.MockGetUrl(id)
}
func (ref *TestResultDBMock) GetAllTests(tenantId string) ([]*models.ResultsStorage, error) {
	return ref.MockGetAllTests(tenantId)
}
func (ref *TestResultDBMock) GetAllTestResults(tenantId string) ([]*models.ResultsStorage, error) {
	return ref.MockGetAllTestResults(tenantId)
}
func (ref *TestResultDBMock) HasTest(test *models.ResultsStorage) bool {
	return ref.MockHasTest(test)
}
func (ref *TestResultDBMock) HasTestId(id string) bool {
	return ref.MockHasTestId(id)
}
func (ref *TestResultDBMock) InsertTest(test *models.ResultsStorage) (string, error) {
	return ref.MockInsertTest(test)
}
func (ref *TestResultDBMock) UpdateTest(id string, test *models.ResultsStorage) error {
	return ref.MockUpdateTest(id, test)
}
func (ref *TestResultDBMock) RemoveTest(id string) error {
	return ref.MockRemoveTest(id)
}
