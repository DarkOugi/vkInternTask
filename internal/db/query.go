package db

import (
	"context"
	"errors"
	"fmt"
	"vk/internal/entity"

	"github.com/jackc/pgx/v5"
)

func (db *PostgresDB) GetUserInfo(ctx context.Context, login string) (*entity.User, bool, error) {
	var u entity.User

	selectUserInfo := "SELECT login,password FROM Users WHERE login = $1"
	err := db.conn.QueryRow(ctx, selectUserInfo, login).Scan(&u.Password, &u.Login)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &u, false, nil
		}
		return &u, false, fmt.Errorf("select user error: %w", err)
	}

	return &u, true, nil
}

func (db *PostgresDB) InitAdv(ctx context.Context, login, header, about, picture string,price float32) error {
	insertUser := "INSERT INTO Advertisement (login, header, about, picture, price) VALUES ($1,$2,$3,$4,$5);"
	_, err := db.conn.Exec(ctx, insertUser, login, header, about, picture, price)

	if err != nil {
		return fmt.Errorf("error create adv: %w", err)
	}
	return nil
}

func (db *PostgresDB) InitUser(ctx context.Context, login, password string) error {
	insertUser := "INSERT INTO Users (login,password) VALUES ($1,$2);"
	_, err := db.conn.Exec(ctx, insertUser, login, password)

	if err != nil {
		return fmt.Errorf("error create user: %w", err)
	}
	return nil
}

func GetAdvs(ctx context.Context, page, pageSize int, colSort, typeSort string, filter *entity.Filter) []*entity.Advertisement,error{
	advs := []*entity.Advertisement{}
	tx, err := db.conn.Begin(ctx)
	if err != nil {
		return 0, nil, nil, nil, fmt.Errorf("can't start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	selectSql := " SELECT a.header,a.about,a.pucture,a.price,a.login FROM Advertisement AS a "
	if filter != nil {
		pagerSql := "ORDER BY $3 $4
			LIMIT $5 OFFSET $6
		"
		selectSql += "WHERE a.price BETWEEN ($1,$2)" + pagerSql
		advRows, errAdvSql := tx.Query(ctx, filter.min, filter.max,colSort,typeSort,pageSize,page-1)
		if errAdvSql != nil {
			return nil,fmt.Errorf("can't get adv: %w", errAdvSql)
		}
		for advRows.Next() {
			var advHeader string
			var advAbout string
			var advPicture string
			var advPrice float32
			var advLogin string
	
			errScan := advRows.Scan(&advHeader, &advAbout,&advPicture,&advPrice,&advLogin)
			if errScan != nil {
				return nil, fmt.Errorf("error scan adv: %w", errScan)
			}
	
			advs = append(advs, &entity.Advertisement{
				header: advHeader,
				about:  advAbout,
				picture:advPicture,
				price:advPrice,
				login:advLogin,
			})
		}
		err = tx.Commit(ctx)
		if err != nil {
			return nil, fmt.Errorf("can't complete transaction: %w", err)
		}
		return advs,nil
	} 
	pagerSql := "ORDER BY $1 $2
		LIMIT $3 OFFSET $4
	"
	selectSql += pagerSql
	advRows, errAdvSql := tx.Query(ctx, colSort,typeSort,pageSize,page-1)
	if errAdvSql != nil {
		return nil,fmt.Errorf("can't get adv: %w", errAdvSql)
	}
	for advRows.Next() {
		var advHeader string
		var advAbout string
		var advPicture string
		var advPrice float32
		var advLogin string

		errScan := advRows.Scan(&advHeader, &advAbout,&advPicture,&advPrice,&advLogin)
		if errScan != nil {
			return nil, fmt.Errorf("error scan adv: %w", errScan)
		}

		advs = append(advs, &entity.Advertisement{
			header: advHeader,
			about:  advAbout,
			picture:advPicture,
			price:advPrice,
			login:advLogin,
		})
	}
	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't complete transaction: %w", err)
	}
	return advs,nil
}
