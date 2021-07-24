package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/k0kubun/pp"
	_ "github.com/lib/pq"
)

type RawUser struct {
	Id        int64
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func main() {
	db, err := sql.Open("postgres", "host=db port=5432 user=hoge-user password=hoge-pass dbname=hoge-db sslmode=disable")
	if err != nil {
		fmt.Println(err)
		panic("sql.Openに失敗しましたので、終了します")
	}
	defer db.Close()

	// db.QueryRow
	rawUser := RawUser{}
	err2 := db.QueryRow("SELECT users.id, users.name, users.email, users.created_at, users.updated_at FROM users WHERE id = $1", 1).Scan(&rawUser.Id, &rawUser.Name, &rawUser.Email, &rawUser.CreatedAt, &rawUser.UpdatedAt)
	if err2 != nil {
		pp.Println(err2)
		panic("QueryRowに失敗しましたので、終了します")
	}
	fmt.Println("-----------------------------------------")
	pp.Println(rawUser)
	fmt.Println("-----------------------------------------")

	// db.Query
	rows, err3 := db.Query("SELECT users.id, users.name, users.email, users.created_at, users.updated_at FROM users")
	if err3 != nil {
		pp.Println(err3)
		panic("Queryに失敗しましたので、終了します")
	}
	fmt.Println("-----------------------------------------")
	for rows.Next() {
		user := RawUser{}
		err4 := rows.Scan(&user.Id, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
		if err4 != nil {
			pp.Println(err4)
			panic("rows.Scanに失敗しましたので、終了します")
		}
		pp.Println(user)
	}
	fmt.Println("-----------------------------------------")

	// db.Execute
	result, err5 := db.Exec("DELETE FROM users WHERE email = $1;", "hoge@example.com")
	if err5 == nil {
		pp.Println(result.RowsAffected())
	} else {
		fmt.Println("失敗")
		pp.Println(err5)
	}

	// db.QueryRow
	err7 := db.QueryRow("INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id,name,email,created_at,updated_at;", "taro", "hoge@example.com").Scan(&rawUser.Id, &rawUser.Name, &rawUser.Email, &rawUser.CreatedAt, &rawUser.UpdatedAt)
	if err7 != nil {
		pp.Println(err7)
		panic("QueryRow(INSERT)に失敗しましたので、終了します")
	}
	fmt.Println("-----------------------------------------")
	pp.Println(rawUser)
	fmt.Println("-----------------------------------------")

	// db.QueryRow
	err8 := db.QueryRow("UPDATE users SET name = $1, email = $2 WHERE users.id = 2 RETURNING id,name,email,created_at,updated_at;", "樋口一葉", "higuchi@example.com").Scan(&rawUser.Id, &rawUser.Name, &rawUser.Email, &rawUser.CreatedAt, &rawUser.UpdatedAt)
	if err8 != nil {
		pp.Println(err8)
		panic("QueryRow(UPDATE)に失敗しましたので、終了します")
	}
	fmt.Println("-----------------------------------------")
	pp.Println(rawUser)
	fmt.Println("-----------------------------------------")

	// Create a new context, and begin a transaction
	fmt.Println("---------------------------トランザクション(成功例)")
	ctx := context.Background()
	tx, err9 := db.BeginTx(ctx, nil)
	if err9 != nil {
		log.Fatal(err)
		panic("トランザクション生成に失敗しましたので終了します")
	}
	_, err10 := tx.ExecContext(ctx, "INSERT INTO users (name, email) VALUES ('tao01', 'taro01@example.com') RETURNING id,name,email,created_at,updated_at;")
	if err10 != nil {
		fmt.Println("ロールバックします")
		tx.Rollback()
	}
	_, err11 := tx.ExecContext(ctx, "INSERT INTO users (name, email) VALUES ('tao02', 'taro02@example.com') RETURNING id,name,email,created_at,updated_at;")
	if err11 != nil {
		fmt.Println("ロールバックします")
		tx.Rollback()
	}
	err12 := tx.Commit()
	if err12 != nil {
		fmt.Println("トランザクションのコミットに失敗しました")
		log.Fatal(err12)
	}
	fmt.Println("---------------------------トランザクション")
	fmt.Println("---------------------------トランザクション(ロールバック)")
	tx2, err13 := db.BeginTx(ctx, nil)
	if err13 != nil {
		log.Fatal(err)
		panic("トランザクション生成に失敗しましたので終了します")
	}
	_, err14 := tx2.ExecContext(ctx, "INSERT INTO users (name, email) VALUES ('tao03', 'taro03@example.com') RETURNING id,name,email,created_at,updated_at;")
	if err14 != nil {
		fmt.Println("ロールバックします")
		tx.Rollback()
	}
	_, err15 := tx2.ExecContext(ctx, "INSERT INTO users (name, email) VALUES ('tao02', 'taro02@example.com') RETURNING id,name,email,created_at,updated_at;")
	if err15 != nil {
		fmt.Println("ロールバックします")
		tx.Rollback()
	}
	err16 := tx2.Commit()
	if err16 != nil {
		fmt.Println("トランザクションのコミットに失敗しました")
		log.Fatal(err16)
	}
	fmt.Println("---------------------------トランザクション")
}
