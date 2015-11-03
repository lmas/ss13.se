package ss13

import (
	"fmt"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

func OpenSqliteDB(args ...interface{}) *gorm.DB {
	var e error
	db, e := gorm.Open("sqlite3", args...)
	check_error(e)
	return &db
}

func InitSchema(db *gorm.DB) {
	db.AutoMigrate(&Server{})
	db.AutoMigrate(&ServerPopulation{})
}

func AllServers(db *gorm.DB) []*Server {
	var tmp []*Server
	db.Order("players_current desc, last_updated desc, title").Find(&tmp)
	return tmp
}

func GetServer(db *gorm.DB, id int) (*Server, error) {
	var tmp Server
	if db.First(&tmp, id).RecordNotFound() {
		return nil, fmt.Errorf("Server not found")
	}
	return &tmp, nil
}

func GetOldServers(db *gorm.DB, ts time.Time) []*Server {
	var tmp []*Server
	db.Where("last_updated < ?", ts).Find(&tmp)
	return tmp
}

func RemoveOldServers(db *gorm.DB, ts time.Time) {
	db.Where("last_updated < datetime(?, '-7 days')").Delete(Server{})
}

func GetServerPopulation(db *gorm.DB, id int, d time.Duration) []*ServerPopulation {
	var tmp []*ServerPopulation
	t := time.Now().Add(-d)
	db.Order("timestamp desc, server_id").Where("server_id = ? and timestamp > ?", id, t).Find(&tmp)
	return tmp
}

func InsertOrSelect(db *gorm.DB, s *RawServerData) int {
	var tmp Server
	newserver := Server{
		LastUpdated:    s.Timestamp,
		Title:          s.Title,
		GameUrl:        s.Game_url,
		SiteUrl:        s.Site_url,
		PlayersCurrent: s.Players,
		PlayersAvg:     s.Players,
		PlayersMin:     s.Players,
		PlayersMax:     s.Players,
	}
	db.Where("title = ?", s.Title).Attrs(newserver).FirstOrCreate(&tmp)
	return tmp.ID
}

func AddServerPopulation(db *gorm.DB, id int, s *RawServerData) {
	var tmp Server
	db.Where("id = ?", id).First(&tmp)
	pop := ServerPopulation{
		Timestamp: Now(),
		Players:   s.Players,
		Server:    tmp,
	}
	db.Create(&pop)
}

func UpdateServerStats(db *gorm.DB, id int, s *RawServerData) {
	var tmp Server

	period := Now().Add(-time.Duration(30*24) * time.Hour)
	db.Where("id = ?", id).First(&tmp)
	rows, err := db.Table("server_populations").Where("server_id = ? AND timestamp > ?", tmp.ID, period).Select("timestamp, players").Order("timestamp desc").Rows()
	if err != nil {
		log.Panic(err)
		return
	}
	defer rows.Close()

	var timestamp time.Time
	var players, sum, count int
	var day_sums [7]int
	day_counts := [7]int{1, 1, 1, 1, 1, 1, 1} // Using 1's to prevent ZeroDiv error
	min := s.Players
	max := s.Players
	for rows.Next() {
		rows.Scan(&timestamp, &players)
		count++
		sum += players
		if players < min {
			min = players
		}
		if players > max {
			max = players
		}

		day := timestamp.Weekday()
		day_sums[day] += players
		day_counts[day]++
	}

	tmp.LastUpdated = s.Timestamp
	tmp.Title = s.Title
	tmp.GameUrl = s.Game_url
	tmp.SiteUrl = s.Site_url
	tmp.PlayersCurrent = s.Players
	tmp.PlayersAvg = sum / count
	tmp.PlayersMin = min
	tmp.PlayersMax = max
	tmp.PlayersMon = day_sums[1] / day_counts[1]
	tmp.PlayersTue = day_sums[2] / day_counts[2]
	tmp.PlayersWed = day_sums[3] / day_counts[3]
	tmp.PlayersThu = day_sums[4] / day_counts[4]
	tmp.PlayersFri = day_sums[5] / day_counts[5]
	tmp.PlayersSat = day_sums[6] / day_counts[6]
	tmp.PlayersSun = day_sums[0] / day_counts[0]
	db.Save(&tmp)
}
