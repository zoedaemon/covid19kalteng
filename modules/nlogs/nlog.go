package nlogs

import (
	"covid19kalteng/components/logs"
	"covid19kalteng/covid19"
	"covid19kalteng/models"

	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strconv"

	"github.com/dgrijalva/jwt-go"
)

// NLog send log to northstar service
func NLog(level string, tag string, message interface{}, jwttoken *jwt.Token, note string, nouser bool) {
	var (
		uid      string
		username string
		err      error
	)

	if !nouser {
		jti, _ := strconv.ParseUint(jwttoken.Claims.(jwt.MapClaims)["jti"].(string), 10, 64)
		user := models.User{}
		err = user.FindbyID(jti)
		if err == nil {
			uid = fmt.Sprint(user.ID)
			username = user.Username
		}
	}

	jMarshal, _ := json.Marshal(message)

	if flag.Lookup("test.v") == nil {
		err = covid19.App.Northstar.SubmitKafkaLog(logs.Log{
			Level:    level,
			Tag:      tag,
			Messages: string(jMarshal),
			UID:      uid,
			Username: username,
			Note:     note,
		}, "log")
	}

	if err != nil {
		log.Printf("error northstar log : %v", err)
	}
}

// NAudittrail send audit trail log to northstar service
func NAudittrail(ori interface{}, new interface{}, jwttoken *jwt.Token, entity string, entityID string, action string) {
	var (
		uid      string
		username string
		err      error
	)

	jti, _ := strconv.ParseUint(jwttoken.Claims.(jwt.MapClaims)["jti"].(string), 10, 64)
	user := models.User{}
	err = user.FindbyID(jti)
	if err == nil {
		uid = fmt.Sprint(user.ID)
		username = user.Username
	} else {
		uid = "0"
		username = "not found"
	}

	oriMarshal, _ := json.Marshal(ori)
	newMarshal, _ := json.Marshal(new)

	if flag.Lookup("test.v") == nil {
		err = covid19.App.Northstar.SubmitKafkaLog(logs.Audittrail{
			Client:   covid19.App.Northstar.Secret,
			UserID:   uid,
			Username: username,
			Roles:    fmt.Sprint(user.Roles),
			Entity:   entity,
			EntityID: entityID,
			Action:   action,
			Original: fmt.Sprintf(`%s`, string(oriMarshal)),
			New:      fmt.Sprintf(`%s`, string(newMarshal)),
		}, "audittrail")
	}

	if err != nil {
		log.Printf("error northstar log : %v", err)
	}
}
