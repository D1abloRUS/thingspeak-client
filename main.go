package main

import (
	"encoding/json"
	"net/http"
	"time"
	"log"
	"io/ioutil"
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"strconv"
)

type Stats struct {
	T1 []int
	T2 []int
	T1_mon int
	T2_mon int
}

type Thingspeak_field2 struct {
	Channel struct {
		ID          int       `json:"id"`
		Name        string    `json:"name"`
		Latitude    string    `json:"latitude"`
		Longitude   string    `json:"longitude"`
		Field1      string    `json:"field1"`
		Field2      string    `json:"field2"`
		Field3      string    `json:"field3"`
		Field4      string    `json:"field4"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
		LastEntryID int       `json:"last_entry_id"`
	} `json:"channel"`
	Feeds []struct {
		CreatedAt time.Time `json:"created_at"`
		EntryID   int       `json:"entry_id"`
		Field2    string    `json:"field2"`
	} `json:"feeds"`
}

type Thingspeak_field3 struct {
	Channel struct {
		ID          int       `json:"id"`
		Name        string    `json:"name"`
		Latitude    string    `json:"latitude"`
		Longitude   string    `json:"longitude"`
		Field1      string    `json:"field1"`
		Field2      string    `json:"field2"`
		Field3      string    `json:"field3"`
		Field4      string    `json:"field4"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
		LastEntryID int       `json:"last_entry_id"`
	} `json:"channel"`
	Feeds []struct {
		CreatedAt time.Time `json:"created_at"`
		EntryID   int       `json:"entry_id"`
		Field3    string    `json:"field3"`
	} `json:"feeds"`
}

type Enviroment struct {
	Url         string
	Channel     int
	T1_Field    int
	T2_Field	int
	ApiKey      string
}

func main() {
	log.Print("Start app")
	http.HandleFunc("/", get_kwt)
	http.ListenAndServe(":3000", nil)
}

func get_kwt(w http.ResponseWriter, r *http.Request) {
	var e Enviroment
	err := envconfig.Process("thingspeak", &e)
	if err != nil {
		log.Fatal(err.Error())
	}

    date_start, date_end := get_date("old")
    t1_start := makerequest(e.Url,e.Channel,e.T1_Field,e.ApiKey,date_start,date_end)
    t2_start := makerequest(e.Url,e.Channel,e.T2_Field,e.ApiKey,date_start,date_end)
	date_start, date_end = get_date("now")
    t1_end := makerequest(e.Url,e.Channel,e.T1_Field,e.ApiKey,date_start,date_end)
	t2_end := makerequest(e.Url,e.Channel,e.T2_Field,e.ApiKey,date_start,date_end)

	t1_mon := parse("T1", t1_end) - parse("T1", t1_start)
	t2_mon := parse("T2", t2_end) - parse("T2", t2_start)

	stats := Stats{ T1:[]int{parse("T1", t1_start), parse("T1", t1_end)},
	T2:[]int{parse("T2", t2_start), parse("T2", t2_end)}, T1_mon: t1_mon, T2_mon: t2_mon}

	js, err := json.Marshal(stats)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("%s", js)

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func makerequest(thingspeak_url string, channel int, field int, api_key string, date_start string, date_end string) []byte {
	url:= fmt.Sprintf("%s/channels/%d/fields/%d.json?api_key=%s&start=%s&end=%s",thingspeak_url,channel,field,api_key,date_start,date_end)

	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	resp, err := client.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return body
}

func parse(t_zone string, resp []byte) int {
	if t_zone == "T1" {
		var t = Thingspeak_field2{}
		err := json.Unmarshal(resp, &t)
		if err != nil {
			log.Fatalln(err)
		}
		t1,_ := strconv.ParseFloat(t.Feeds[0].Field2, 64)
		return int(t1)
	} else {
		var t = Thingspeak_field3{}
		err := json.Unmarshal(resp, &t)
		if err != nil {
			log.Fatalln(err)
		}
		t2,_ := strconv.ParseFloat(t.Feeds[0].Field3, 64)
		return int(t2)
	}
}

func get_date(d string) (string, string) {
	t := time.Now().UTC()
	if d == "old" {
		start := time.Date(t.Year(), t.Month(), 25, 11, 0, 0, 0, time.UTC)
		oneMonthAgoStart := start.AddDate(0, -1, 0)
		oneMonthAgoEnd := oneMonthAgoStart.Add(5 * time.Minute)
		return oneMonthAgoStart.Format("2006-01-02 15:04:05"), oneMonthAgoEnd.Format("2006-01-02 15:04:05")
	} else {
		nowStart := t.Add(time.Duration(-5) * time.Minute)
		nowEnd := t
		return nowStart.Format("2006-01-02 15:04:05"), nowEnd.Format("2006-01-02 15:04:05")
	}
}