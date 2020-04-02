package migration

import (
	"covid19kalteng/covid19"
	"covid19kalteng/models"

	"fmt"
	"strings"

	"github.com/lib/pq"
)

//Seed func
func Seed() {
	if covid19.App.ENV == "development" {
		seedClients()
		seedRoles()
		seedUsers()
		seedEdu()
	}
}

//TestSeed func
func TestSeed() {
	if covid19.App.ENV == "development" {
		seedClients()
		seedRoles()
		seedUsers()
		seedEdu()
	}
}

// Truncate defined tables. []string{"all"} to truncate all tables.
func Truncate(tableList []string) (err error) {
	if len(tableList) > 0 {
		if tableList[0] == "all" {
			tableList = []string{
				"clients",
				"roles",
				"users",
				"edus",
			}
		}

		tables := strings.Join(tableList, ", ")
		sqlQuery := fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", tables)
		err = covid19.App.DB.Exec(sqlQuery).Error
		return err
	}

	return fmt.Errorf("define tables that you want to truncate")
}

// seed clients
func seedClients() {
	clients := []models.Client{
		models.Client{
			Name:   "admin",
			Key:    "adminkey",
			Secret: "adminsecret",
		},
		models.Client{
			Name:   "reporter",
			Key:    "reactkey",
			Secret: "reactsecret",
		},
		models.Client{
			Name:   "reporter",
			Key:    "clientkey",
			Secret: "clientsecret",
		},
	}
	for _, client := range clients {
		client.Create()
	}
}

func seedRoles() {
	roles := []models.Roles{
		models.Roles{
			Name:        "Administrator",
			Status:      "active",
			Description: "Super Admin",
			System:      "Core",
			Permissions: pq.StringArray{"all"},
		},
		models.Roles{
			Name:        "Ops",
			Status:      "active",
			Description: "Ops",
			System:      "Core",
			Permissions: pq.StringArray{
				"core_create_client", "core_view_image", "core_role_list", "core_role_details", "core_role_new", "core_role_patch", "core_role_range", "core_permission_list", "core_user_list", "core_user_details", "core_user_new", "core_user_patch"},
		},
		models.Roles{
			Name:        "Reporter",
			Status:      "active",
			Description: "ini untuk Reporter",
			System:      "Dashboard",
			Permissions: pq.StringArray{
				"reporter_profile", "reporter_profile_edit"},
		},
	}
	for _, role := range roles {
		role.Create()
	}
}

func seedUsers() {
	users := []models.User{
		models.User{
			Roles:      pq.Int64Array{1},
			Username:   "adminkey",
			Password:   "adminsecret",
			Email:      "admin@covid19kalteng.com",
			Phone:      "081234567890",
			Status:     "active",
			FirstLogin: false,
		},
		models.User{
			Roles:      pq.Int64Array{2},
			Username:   "viewer",
			Password:   "password",
			Email:      "covid19kalteng@covid19kalteng.com",
			Phone:      "081234567891",
			Status:     "active",
			FirstLogin: false,
		},
		models.User{
			Roles:      pq.Int64Array{3},
			Username:   "responder1",
			Email:      "responder1@covid19kalteng.com",
			Phone:      "081234567892",
			Password:   "password",
			Status:     "active",
			FirstLogin: false,
		},
		models.User{
			Roles:      pq.Int64Array{3},
			Username:   "responder2",
			Email:      "responder2@covid19kalteng.com",
			Phone:      "081234567893",
			Password:   "password",
			Status:     "active",
			FirstLogin: false,
		},
	}
	for _, user := range users {
		user.Create()
	}
}

func seedEdu() {

	edus := []models.Edu{
		models.Edu{
			Title: "Pencegahan agar Terhindar dari Virus SARS-CoV-2",
			Description: `
			<p>menggosok tangan dengan air mengalir dan <b>sabun</b> atau jika tidak ada gunakan hand sanitizer</p>`,
		},
		models.Edu{
			Title: "Penularan wabah Covid19",
			Description: `
			<p>penyebaran yang masif....</p>`,
		},
	}
	for _, edu := range edus {
		edu.Create()
	}
}
