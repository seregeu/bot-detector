package scriptsdb

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/AlyonaAg/bot-detector/internal/config"
	"github.com/AlyonaAg/bot-detector/internal/model"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"
)

type Repository interface {
	CreateUser(user model.User) (int64, error)
	CreateStatic(static model.Static) (int64, error)
	CreateDynamic(static model.Dynamic) (int64, error)
	GetUser(username string) (*model.User, error)
	GetLastCountDynamic(userId int64, count int) ([]*model.Dynamic, error)
}

type repo struct {
	db *sql.DB
}

func (r *repo) CreateUser(user model.User) (int64, error) {
	if err := r.db.QueryRow(
		`INSERT INTO "user" (username,password,first_name,last_name,phone_number,mail) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id`,
		user.Username,
		user.Password,
		user.FirstName,
		user.LastName,
		user.Phone,
		user.Email,
	).Scan(&user.ID); err != nil {
		return 0, err
	}
	return user.ID, nil
}

func (r *repo) CreateStatic(static model.Static) (int64, error) {
	if err := r.db.QueryRow(
		`INSERT INTO "static" (user_id,battery_charge,battery_status,data_trans_stand,sim_presence) VALUES ($1,$2,$3,$4,$5) RETURNING id`,
		static.UserID,
		static.BattaryCharge,
		static.BattaryStatus,
		static.DataTransStand,
		static.SimPresence,
	).Scan(&static.ID); err != nil {
		return 0, err
	}
	return static.ID, nil
}

func (r *repo) CreateDynamic(dynamic model.Dynamic) (int64, error) {
	if err := r.db.QueryRow(
		`INSERT INTO "dynamic" (user_id,max_device_offs,min_device_offs,max_dev_acceleration,min_dev_acceleration,min_light,max_light,hit_y,hit_x) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING id`,
		dynamic.UserID,
		dynamic.MaxDeviceOffs,
		dynamic.MinDeviceOffs,
		dynamic.MaxDevAcceleration,
		dynamic.MinDevAcceleration,
		dynamic.MinLight,
		dynamic.MaxLight,
		dynamic.HitY,
		dynamic.HitX,
	).Scan(&dynamic.ID); err != nil {
		return 0, err
	}
	return dynamic.ID, nil
}

func (r *repo) GetUser(username string) (*model.User, error) {
	var user = &model.User{}
	if err := r.db.QueryRow(`SELECT id,username,password,first_name,last_name,phone_number,mail FROM "user" WHERE username = $1`, username).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.Phone,
		&user.Email,
	); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *repo) GetLastCountDynamic(userId int64, count int) ([]*model.Dynamic, error) {
	var dynamics = []*model.Dynamic{}
	rows, err := r.db.Query(`SELECT id,user_id,max_device_offs,min_device_offs, max_dev_acceleration, min_dev_acceleration, min_light, max_light, hit_y, hit_x  
	FROM "dynamic" WHERE user_id = $1 ORDER BY id DESC LIMIT $2`, userId, count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		s := model.Dynamic{}
		if err := rows.Scan(
			&s.ID,
			&s.UserID,
			&s.MaxDeviceOffs,
			&s.MinDeviceOffs,
			&s.MaxDevAcceleration,
			&s.MinDevAcceleration,
			&s.MinLight,
			&s.MaxLight,
			&s.HitY,
			&s.HitX,
		); err != nil {
			return nil, err
		}
		fmt.Println(s)

		dynamics = append(dynamics, &s)
	}

	return dynamics, nil
}

/*
func (r *repo) ListScripts(filter model.ListScriptsFilter) (model.Scripts, error) {
	fmt.Println(filter)

	rows, err := r.db.Query(`SELECT id, url, original_script, result, danger_percent, virus_total FROM "scripts" LIMIT $1 OFFSET $2`,
		filter.Limit, filter.Page)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scripts model.Scripts
	for rows.Next() {
		s := model.Script{}
		if err := rows.Scan(
			&s.ID,
			&s.URL,
			&s.Script,
			&s.Result,
			&s.DangerPercent,
			&s.VirusTotal,
		); err != nil {
			continue
		}
		fmt.Println(s)

		scripts = append(scripts, &s)
	}
	return scripts, nil
}*/

func NewRepository() (Repository, error) {
	databaseURL, err := config.GetValue(config.DatabaseURL)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("postgres", databaseURL.(string))
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, err
	}

	pathMigration, err := config.GetValue(config.PathMigration)
	if err != nil {
		return nil, err
	}

	log.Print(pathMigration.(string))
	m, err := migrate.NewWithDatabaseInstance(pathMigration.(string), "postgres", driver)
	if err != nil {
		log.Print("c")
		return nil, err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, err
	}

	log.Print("Store OK.")

	return &repo{db: db}, nil
}
