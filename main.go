package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/apiweb")
	if err != nil {
		fmt.Print(err.Error())
	}
	defer db.Close()
	//membuat koneksi ke database
	err = db.Ping()
	if err != nil {
		fmt.Println(err.Error())
	}
	type Orang struct {
		ID           int
		NamaPertama  string
		NamaTerakhir string
	}
	router := gin.Default()
	//penghandlean API dari sini
	//mengambil detail orang
	router.GET("/orang/:id", func(c *gin.Context) {
		var (
			orang  Orang
			result gin.H
		)
		id := c.Param("id")
		row := db.QueryRow("select id, nama_pertama, nama_terakhir from orang where id =?;", id)
		err = row.Scan(&orang.ID, &orang.NamaPertama, &orang.NamaTerakhir)
		if err != nil {
			// If no results send null
			result = gin.H{
				"Hasil":       nil,
				"Jumlah Data": 0,
			}
		} else {
			result = gin.H{
				"Hasil":       orang,
				"Jumlah Data": 1,
			}
		}
		c.JSON(http.StatusOK, result)
	})
	//mengambil data semua orang
	router.GET("/orang-orang", func(c *gin.Context) {
		var (
			orang      Orang
			orangorang []Orang
		)
		rows, err := db.Query("select id, nama_pertama, nama_terakhir from orang;")

		if err != nil {
			fmt.Print(err.Error())
		}
		for rows.Next() {
			err = rows.Scan(&orang.ID, &orang.NamaPertama, &orang.NamaTerakhir)
			orangorang = append(orangorang, orang)
			if err != nil {
				fmt.Print(err.Error())
			}
		}
		defer rows.Close()
		c.JSON(http.StatusOK, gin.H{
			"hasil":        orangorang,
			"jumalah data": len(orangorang),
		})
	})

	//menambah data orang-orang

	// POST new person details
	router.POST("/orang", func(c *gin.Context) {
		var buffer bytes.Buffer
		NamaPertama := c.PostForm("nama_pertama")
		NamaTerakhir := c.PostForm("nama_terakhir")
		stmt, err := db.Prepare("insert into orang (nama_pertama, nama_terakhir) values(?,?);")
		if err != nil {
			fmt.Print(err.Error())
		}
		_, err = stmt.Exec(NamaPertama, NamaTerakhir)

		if err != nil {
			fmt.Print(err.Error())
		}

		// Fastest way to append strings
		buffer.WriteString(NamaPertama)
		buffer.WriteString(" ")
		buffer.WriteString(NamaTerakhir)
		defer stmt.Close()
		name := buffer.String()
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf(" %s berhasil di tambahkan ke database", name),
		})
	})
	// mengubah data yang sudah ada
	router.PUT("/orang", func(c *gin.Context) {
		var buffer bytes.Buffer
		id := c.Query("id")
		NamaPertama := c.PostForm("nama_pertama")
		NamaTerakhir := c.PostForm("nama_terakhir")
		stmt, err := db.Prepare("update orang set nama_pertama= ?, nama_terakhir= ? where id= ?;")
		if err != nil {
			fmt.Print(err.Error())
		}
		_, err = stmt.Exec(NamaPertama, NamaTerakhir, id)
		if err != nil {
			fmt.Print(err.Error())
		}

		// Fastest way to append strings
		buffer.WriteString(NamaPertama)
		buffer.WriteString(" ")
		buffer.WriteString(NamaTerakhir)
		defer stmt.Close()
		name := buffer.String()
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Data berhasil di rubah %s", name),
		})
	})
	//menghapus Data berdasarkan ID

	router.DELETE("/orang", func(c *gin.Context) {
		id := c.Query("id")
		stmt, err := db.Prepare("delete from orang where id= ?;")
		if err != nil {
			fmt.Print(err.Error())
		}
		_, err = stmt.Exec(id)
		if err != nil {
			fmt.Print(err.Error())
		}
		c.JSON(http.StatusOK, gin.H{
			"Pesan": fmt.Sprintf("Anda Telah berhasil menghapus id: %s", id)})

	})

	router.Run(":8000")

} //end
