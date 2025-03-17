package mapepire

import (
	"fmt"
	"log"
	"testing"
)

func initSQLTable(command string, queryops2 QueryOptions) (*SQLJob, *Query) {
	job := NewSQLJob("test")
	err := job.Connect(server)
	if err != nil {
		log.Fatal(err)
	}

	query, _ := job.Query("CREATE TABLE qtemp.TEMPTEST (ID CHAR(8) NOT NULL, DESCRIPTION VARCHAR(60) NOT NULL, SERIALNO CHAR(12) NOT NULL)")
	resp := query.Execute()
	if resp.Error != nil {
		log.Fatal(resp.Error)
	}

	insertquery, _ := job.Query(`INSERT INTO TEMPTEST VALUES (1, 'Lorem ipsum', 121212),
	(2, 'dolor sit amet', 232323),
	(3, 'consetetur sadipscing elitr', 343434),
	(4, 'sed diam nonumy', 454545),
	(5, 'eirmod tempor', 565656)`)
	resp2 := insertquery.Execute()
	log.Println(resp2)

	query2, err := job.QueryWithOptions(command, queryops2)
	if err != nil {
		log.Fatal(err)
	}
	return job, query2
}

// Execute SQL
func TestNewExecute(t *testing.T) {
	queryops := QueryOptions{
		Rows: 5,
	}
	_, query := initSQLTable("SELECT * FROM TEMPTEST", queryops)
	have := query.Execute()

	want := &ServerResponse{
		ID:         "3",
		Success:    true,
		IsDone:     false,
		Error:      nil,
		HasResults: true,
	}

	if have.ID != want.ID {
		t.Errorf("got %s, want %s", have.ID, want.ID)
	}
	if have.Success != want.Success {
		t.Errorf("got %t, want %t", have.Success, want.Success)
	}
	if have.Error != want.Error {
		t.Errorf("got %s, want %s", have.Error, want.Error)
	}
	if have.IsDone != want.IsDone {
		t.Errorf("got %t, want %t", have.IsDone, want.IsDone)
	}
	if have.HasResults != want.HasResults {
		t.Errorf("got %v, want %v", have.HasResults, want.HasResults)
	}
	if have.Data == nil {
		t.Errorf("should have data")
	}
	if have.Metadata == nil {
		t.Errorf("should have metadata")
	}
}

// Execute SQL, response in terse
func TestNewExecuteTerse(t *testing.T) {
	queryops := QueryOptions{
		Rows:        5,
		TerseResult: true,
	}
	_, query := initSQLTable("SELECT * FROM TEMPTEST", queryops)
	have := query.Execute()

	want := &ServerResponse{
		ID:         "3",
		Success:    true,
		IsDone:     false,
		Error:      nil,
		HasResults: true,
	}

	if have.ID != want.ID {
		t.Errorf("got %s, want %s", have.ID, want.ID)
	}
	if have.Success != want.Success {
		t.Errorf("got %t, want %t", have.Success, want.Success)
	}
	if have.Error != want.Error {
		t.Errorf("got %s, want %s", have.Error, want.Error)
	}
	if have.IsDone != want.IsDone {
		t.Errorf("got %t, want %t", have.IsDone, want.IsDone)
	}
	if have.HasResults != want.HasResults {
		t.Errorf("got %v, want %v", have.HasResults, want.HasResults)
	}
	if have.Data != nil {
		t.Errorf("should be in Terse_data")
	}
	if have.TerseData == nil {
		t.Errorf("should have Terse_data")
	}
	if have.Metadata == nil {
		t.Errorf("should have metadata")
	}
}

// Execute SQL invalid
func TestNewExecuteInvalid(t *testing.T) {
	queryops := QueryOptions{
		Rows: 5,
	}
	_, query := initSQLTable("SELECT * FROM INVALIDTABLE", queryops)
	have := query.Execute()

	want := &ServerResponse{
		ID:         "3",
		Success:    false,
		IsDone:     false,
		HasResults: false,
	}

	if have.ID != want.ID {
		t.Errorf("got %s, want %s", have.ID, want.ID)
	}
	if have.Success != want.Success {
		t.Errorf("got %t, want %t", have.Success, want.Success)
	}
	if have.Error == nil {
		t.Errorf("should throw Error")
	}
	if have.IsDone != want.IsDone {
		t.Errorf("got %t, want %t", have.IsDone, want.IsDone)
	}
	if have.HasResults != want.HasResults {
		t.Errorf("got %v, want %v", have.HasResults, want.HasResults)
	}
	if have.Data != nil {
		t.Errorf("should not have data")
	}
	if have.Metadata != nil {
		t.Errorf("should not have metadata")
	}
}

// Execute SQL insert
func TestNewExecuteInsert(t *testing.T) {
	job := NewSQLJob("test")
	job.Connect(server)

	query, _ := job.Query("DECLARE GLOBAL TEMPORARY TABLE TEMPTEST (ID CHAR(8) NOT NULL, DESCRIPTION VARCHAR(60) NOT NULL, SERIALNO CHAR(12) NOT NULL)")
	resp := query.Execute()
	if resp.Error != nil {
		log.Fatal(resp.Error)
	}

	insertquery, _ := job.Query("INSERT INTO TEMPTEST VALUES (1, 'testtest', 121212)")
	have := insertquery.Execute()

	want := ServerResponse{
		ID:          "2",
		Success:     true,
		HasResults:  false,
		UpdateCount: 1,
	}

	if have.ID != want.ID {
		t.Errorf("got %s, want %s", have.ID, want.ID)
	}
	if have.Success != want.Success {
		t.Errorf("got %t, want %t", have.Success, want.Success)
	}
	if have.Error != nil {
		t.Errorf("should not throw Error")
	}
	if have.UpdateCount != want.UpdateCount {
		t.Errorf("have %v, want %v", have.UpdateCount, want.UpdateCount)
	}
	if have.HasResults != want.HasResults {
		t.Errorf("have %v, want %v", have.HasResults, want.HasResults)
	}
}

// Execute SQL update
func TestNewExecuteUpdate(t *testing.T) {
	queryops := QueryOptions{}

	job, query := initSQLTable("UPDATE TEMPTEST SET ID = 5, DESCRIPTION = 'test', SERIALNO = 545454 WHERE ID = 1", queryops)
	update := query.Execute()
	if update.Error != nil {
		t.Errorf("should not throw error")
	}

	query2, _ := job.Query("SELECT * FROM TEMPTEST WHERE ID = 5")
	have := query2.Execute()
	want := ServerResponse{
		ID:          "4",
		Success:     true,
		HasResults:  true,
		IsDone:      true,
		UpdateCount: 1,
	}

	if have.ID != want.ID {
		t.Errorf("got %s, want %s", have.ID, want.ID)
	}
	if have.Success != want.Success {
		t.Errorf("got %t, want %t", have.Success, want.Success)
	}
	if have.Error != nil {
		t.Errorf("should not throw Error")
	}
	if have.IsDone != want.IsDone {
		t.Errorf("got %t, want %t", have.IsDone, want.IsDone)
	}
	if have.HasResults != want.HasResults {
		t.Errorf("got %v, want %v", have.HasResults, want.HasResults)
	}
	if have.Data == nil {
		t.Errorf("should have data")
	}
	if have.Metadata == nil {
		t.Errorf("should have metadata")
	}
	if update.UpdateCount != want.UpdateCount {
		t.Errorf("have %v, want %v", have.UpdateCount, want.UpdateCount)
	}
}

// Execute SQL delete
func TestNewExecuteDelete(t *testing.T) {
	queryops := QueryOptions{}

	_, query := initSQLTable("DELETE FROM TEMPTEST WHERE ID = 1", queryops)
	have := query.Execute()

	want := ServerResponse{
		ID:          "3",
		Success:     true,
		Error:       nil,
		UpdateCount: 1,
	}

	if have.ID != want.ID {
		t.Errorf("have %v, want %v", have.ID, want.ID)
	}
	if have.Error != want.Error {
		t.Errorf("should not throw error")
	}
	if have.Success != want.Success {
		t.Errorf("have %v, want %v", have.Success, want.Success)
	}
	if have.UpdateCount != want.UpdateCount {
		t.Errorf("have %v, want %v", have.UpdateCount, want.UpdateCount)
	}
}

func TestNewExecuteSelectParam(t *testing.T) {
	queryops := QueryOptions{
		Parameters: [][]string{{"1"}},
	}
	_, query := initSQLTable("SELECT * FROM TEMPTEST WHERE ID = ?", queryops)

	have := query.Execute()

	want := &ServerResponse{
		ID:         "3",
		Success:    true,
		IsDone:     true,
		Error:      nil,
		HasResults: true,
	}

	if have.ID != want.ID {
		t.Errorf("got %s, want %s", have.ID, want.ID)
	}
	if have.Success != want.Success {
		t.Errorf("got %t, want %t", have.Success, want.Success)
	}
	if have.Error != want.Error {
		t.Errorf("should not throw error")
	}
	if have.IsDone != want.IsDone {
		t.Errorf("got %t, want %t", have.IsDone, want.IsDone)
	}
	if have.HasResults != want.HasResults {
		t.Errorf("got %v, want %v", have.HasResults, want.HasResults)
	}
	if have.Data == nil {
		t.Errorf("should have data")
	}
	if have.Metadata == nil {
		t.Errorf("should have metadata")
	}
}

func TestNewExecuteSelectParamInvalid(t *testing.T) {
	queryops := QueryOptions{
		Parameters: [][]string{{"1"}, {"2"}},
	}
	_, query := initSQLTable("SELECT * FROM TEMPTEST WHERE ID = ?", queryops)

	have := query.Execute()

	want := &ServerResponse{
		ID:       "3",
		Success:  false,
		SqlRC:    -99999,
		SqlState: "HY010",
	}

	if have.ID != want.ID {
		t.Errorf("got %s, want %s", have.ID, want.ID)
	}
	if have.Success != want.Success {
		t.Errorf("got %t, want %t", have.Success, want.Success)
	}
	if have.SqlRC != want.SqlRC {
		t.Errorf("got %v, want %v", have.SqlRC, want.SqlRC)
	}
	if have.SqlState != want.SqlState {
		t.Errorf("got %v, want %v", have.SqlState, want.SqlState)
	}
}

// Prepare and Execute SQL with Params
func TestNewExecuteUpdateParam(t *testing.T) {
	queryops := QueryOptions{
		Parameters: [][]string{{"123", "test", "345678", "1"}, {"321", "testtest", "876543", "2"}},
		Rows:       5,
	}
	_, query := initSQLTable("UPDATE TEMPTEST SET ID = ?, DESCRIPTION = ?, SERIALNO = ? WHERE ID = ?", queryops)

	have := query.Execute()

	want := &ServerResponse{
		ID:             "3",
		Success:        true,
		IsDone:         true,
		Error:          nil,
		ParameterCount: 4,
	}

	if have.ID != want.ID {
		t.Errorf("got %s, want %s", have.ID, want.ID)
	}
	if have.Success != want.Success {
		t.Errorf("got %t, want %t", have.Success, want.Success)
	}
	if have.Error != want.Error {
		t.Errorf("got %s, want %s", have.Error, want.Error)
	}
	if have.IsDone != want.IsDone {
		t.Errorf("got %t, want %t", have.IsDone, want.IsDone)
	}
	if have.ParameterCount != want.ParameterCount {
		t.Errorf("got %v, want %v", have.ParameterCount, want.ParameterCount)
	}
	if have.Data == nil {
		t.Errorf("should have data")
	}
	if have.Metadata == nil {
		t.Errorf("should have metadata")
	}
}

func TestNewExecuteInsertParam(t *testing.T) {
	queryops := QueryOptions{
		Parameters: [][]string{{"123", "test", "345678"}, {"321", "testtest", "876543"}},
		Rows:       5,
	}
	_, query := initSQLTable("INSERT INTO TEMPTEST (ID, DESCRIPTION, SERIALNO) VALUES(?,?,?)", queryops)

	have := query.Execute()

	want := &ServerResponse{
		ID:      "3",
		Success: true,
		Error:   nil,
	}

	if have.ID != want.ID {
		t.Errorf("got %s, want %s", have.ID, want.ID)
	}
	if have.Success != want.Success {
		t.Errorf("got %t, want %t", have.Success, want.Success)
	}
	if have.Error != want.Error {
		t.Errorf("got %s, want %s", have.Error, want.Error)
	}
}

func TestNewExecuteDeleteParam(t *testing.T) {
	queryops := QueryOptions{
		Parameters: [][]string{{"1"}, {"2"}},
	}
	_, query := initSQLTable("DELETE FROM TEMPTEST WHERE ID = ?", queryops)

	have := query.Execute()

	want := &ServerResponse{
		ID:          "3",
		Success:     true,
		Error:       nil,
		UpdateCount: 2,
	}

	if have.ID != want.ID {
		t.Errorf("got %s, want %s", have.ID, want.ID)
	}
	if have.Success != want.Success {
		t.Errorf("got %t, want %t", have.Success, want.Success)
	}
	if have.Error != want.Error {
		t.Errorf("got %s, want %s", have.Error, want.Error)
	}
	if have.UpdateCount != want.UpdateCount {
		t.Errorf("got %v, want %v", have.UpdateCount, want.UpdateCount)
	}
}

// Execute CL command
func TestNewExecuteCL(t *testing.T) {
	queryops := QueryOptions{
		IsCLcommand: true,
	}

	job := NewSQLJob("test")
	job.Connect(server)

	command := fmt.Sprintf("DSPUSRPRF %s", server.User)
	query, err := job.QueryWithOptions(command, queryops)
	if err != nil {
		t.Errorf("should not throw error")
	}

	have := query.Execute()

	want := &ServerResponse{
		ID:      "1",
		Success: true,
		Error:   nil,
		IsDone:  true,
	}

	if have.ID != want.ID {
		t.Errorf("got %s, want %s", have.ID, want.ID)
	}
	if have.Success != want.Success {
		t.Errorf("got %t, want %t", have.Success, want.Success)
	}
	if have.Error != want.Error {
		t.Errorf("got %s, want %s", have.Error, want.Error)
	}
	if have.IsDone != want.IsDone {
		t.Errorf("got %t, want %t", have.IsDone, want.IsDone)
	}
}

// Execute CL command invalid
func TestNewExecuteCLinvalid(t *testing.T) {
	queryops := QueryOptions{
		IsCLcommand: true,
	}
	job := NewSQLJob("test")
	job.Connect(server)

	query, err := job.QueryWithOptions("INVALIDCOMMAND", queryops)
	if err != nil {
		t.Errorf("should not throw error")
	}

	have := query.Execute()

	want := &ServerResponse{
		ID:       "1",
		Success:  false,
		SqlState: "42806",
		SqlRC:    -303,
	}

	if have.ID != want.ID {
		t.Errorf("got %s, want %s", have.ID, want.ID)
	}
	if have.Success != want.Success {
		t.Errorf("got %t, want %t", have.Success, want.Success)
	}
	if have.SqlState != want.SqlState {
		t.Errorf("got %v, want %v", have.SqlState, want.SqlState)
	}
	if have.SqlRC != want.SqlRC {
		t.Errorf("got %v, want %v", have.SqlRC, want.SqlRC)
	}
}

// FetchMore valid
func TestFetchMore(t *testing.T) {
	queryops := QueryOptions{
		Rows: 1,
	}
	_, query := initSQLTable("SELECT * FROM TEMPTEST", queryops)

	query.Execute()
	have := query.FetchMore(query.ID, "5")

	want := &ServerResponse{
		ID:         "3",
		Success:    true,
		IsDone:     true,
		Error:      nil,
		HasResults: true,
	}

	if have.ID != want.ID {
		t.Errorf("got %s, want %s", have.ID, want.ID)
	}
	if have.Success != want.Success {
		t.Errorf("got %t, want %t", have.Success, want.Success)
	}
	if have.Error != want.Error {
		t.Errorf("got %v, want %v", have.Error, want.Error)
	}
	if have.IsDone != want.IsDone {
		t.Errorf("got %t, want %t", have.IsDone, want.IsDone)
	}
	if have.HasResults != want.HasResults {
		t.Errorf("got %v, want %v", have.HasResults, want.HasResults)
	}
	if have.Data == nil {
		t.Errorf("should have data")
	}
}

// Fetch more in terse format
func TestFetchMoreTerse(t *testing.T) {
	queryops := QueryOptions{
		Rows:        1,
		TerseResult: true,
	}
	_, query := initSQLTable("SELECT * FROM TEMPTEST", queryops)

	query.Execute()
	have := query.FetchMore(query.ID, "1")

	want := &ServerResponse{
		ID:         "3",
		Success:    true,
		IsDone:     false,
		Error:      nil,
		HasResults: true,
	}

	if have.ID != want.ID {
		t.Errorf("got %s, want %s", have.ID, want.ID)
	}
	if have.Success != want.Success {
		t.Errorf("got %t, want %t", have.Success, want.Success)
	}
	if have.Error != want.Error {
		t.Errorf("got %s, want %s", have.Error, want.Error)
	}
	if have.IsDone != want.IsDone {
		t.Errorf("got %t, want %t", have.IsDone, want.IsDone)
	}
	if have.HasResults != want.HasResults {
		t.Errorf("got %v, want %v", have.HasResults, want.HasResults)
	}
	if have.Data != nil {
		t.Errorf("should be in Terse_data")
	}
	if have.TerseData == nil {
		t.Errorf("should have Terse_data")
	}
}

func TestFetchMoreinvalid(t *testing.T) {
	queryops := QueryOptions{
		Rows: 1,
	}

	_, query := initSQLTable("SELECT * FROM TEMPTEST", queryops)
	query.Execute()
	have := query.FetchMore("invalid", "1")

	want := &ServerResponse{
		Success:    false,
		HasResults: false,
	}

	if have.Success != want.Success {
		t.Errorf("got %t, want %t", have.Success, want.Success)
	}
	if have.Error == nil {
		t.Errorf("should throw error")
	}
	if have.HasResults != want.HasResults {
		t.Errorf("got %v, want %v", have.HasResults, want.HasResults)
	}
}

// Fetch more invalid 2
func TestFetchMoreInvalid2(t *testing.T) {
	queryops := QueryOptions{
		Rows: 1,
	}

	_, query := initSQLTable("SELECT * FROM TEMPTEST", queryops)
	query.Execute()
	have := query.FetchMore("", "1")

	want := &ServerResponse{
		Success:    false,
		HasResults: false,
	}

	if have.Success != want.Success {
		t.Errorf("got %t, want %t", have.Success, want.Success)
	}
	if have.Error == nil {
		t.Errorf("should throw error")
	}
	if have.HasResults != want.HasResults {
		t.Errorf("got %v, want %v", have.HasResults, want.HasResults)
	}
}

func TestFetchMoreinvalid3(t *testing.T) {
	queryops := QueryOptions{
		Rows: 1,
	}

	_, query := initSQLTable("SELECT * FROM TEMPTEST", queryops)
	have := query.FetchMore("5", "1")

	want := &ServerResponse{
		Success:    false,
		HasResults: false,
	}

	if have.Success != want.Success {
		t.Errorf("got %t, want %t", have.Success, want.Success)
	}
	if have.Error == nil {
		t.Errorf("should throw error")
	}
	if have.HasResults != want.HasResults {
		t.Errorf("got %v, want %v", have.HasResults, want.HasResults)
	}
}

// SQL Close
func TestSQLClose(t *testing.T) {
	queryops := QueryOptions{
		Rows: 5,
	}
	_, query := initSQLTable("SELECT * FROM TEMPTEST", queryops)
	query.Execute()
	err := query.SQLClose(query.ID)
	if err != nil {
		t.Errorf("should not throw error")
	}
}

// SQLClose null
func TestSQLCloseNULL(t *testing.T) {
	queryops := QueryOptions{
		Rows: 5,
	}
	_, query := initSQLTable("SELECT * FROM TEMPTEST", queryops)
	query.Execute()
	err := query.SQLClose("")
	if err == nil {
		t.Errorf("should throw error")
	}
}

// SQLClose invalid
func TestSQLCloseInvalid(t *testing.T) {
	queryops := QueryOptions{
		Rows: 50,
	}
	_, query := initSQLTable("SELECT * FROM TEMPTEST", queryops)
	query.Execute()
	err := query.SQLClose("invalid")
	if err == nil {
		t.Errorf("should throw error")
	}
}
