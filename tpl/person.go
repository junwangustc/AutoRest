package main

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Person struct {
	Id         int
	First_Name string
	Last_Name  string
}

func GetPerson(c *gin.Context) {
	var (
		person Person
		result gin.H
	)
	id := c.Param("id")
	row := db.QueryRow("select id, first_name, last_name from mytesttable where id = ?;", id)
	err := row.Scan(&person.Id, &person.First_Name, &person.Last_Name)
	if err != nil {
		// If no results send null
		result = gin.H{
			"result": nil,
			"count":  0,
		}
	} else {
		result = gin.H{
			"result": person,
			"count":  1,
		}
	}
	c.JSON(http.StatusOK, result)
}
func GetPersons(c *gin.Context) {
	var (
		person  Person
		persons []Person
	)
	rows, err := db.Query("select id, first_name, last_name from mytesttable;")
	if err != nil {
		fmt.Print(err.Error())
	}
	for rows.Next() {
		err = rows.Scan(&person.Id, &person.First_Name, &person.Last_Name)
		persons = append(persons, person)
		if err != nil {
			fmt.Print(err.Error())
		}
	}
	defer rows.Close()
	c.JSON(http.StatusOK, gin.H{
		"result": persons,
		"count":  len(persons),
	})

}

func PostPerson(c *gin.Context) {
	var buffer bytes.Buffer
	first_name := c.PostForm("first_name")
	last_name := c.PostForm("last_name")
	fmt.Println("==data==", first_name, last_name)
	stmt, err := db.Prepare("insert into mytesttable (first_name, last_name) values(?,?);")
	if err != nil {
		fmt.Print(err.Error())
	}
	_, err = stmt.Exec(first_name, last_name)

	if err != nil {
		fmt.Print(err.Error())
	}

	// Fastest way to append strings
	buffer.WriteString(first_name)
	buffer.WriteString(" ")
	buffer.WriteString(last_name)
	defer stmt.Close()
	name := buffer.String()
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf(" %s successfully created", name),
	})
}
func PutPerson(c *gin.Context) {
	var buffer bytes.Buffer
	id := c.Query("id")
	first_name := c.PostForm("first_name")
	last_name := c.PostForm("last_name")
	stmt, err := db.Prepare("update mytesttable set first_name= ?, last_name= ? where id= ?;")
	if err != nil {
		fmt.Print(err.Error())
	}
	_, err = stmt.Exec(first_name, last_name, id)
	if err != nil {
		fmt.Print(err.Error())
	}

	// Fastest way to append strings
	buffer.WriteString(first_name)
	buffer.WriteString(" ")
	buffer.WriteString(last_name)
	defer stmt.Close()
	name := buffer.String()
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Successfully updated to %s", name),
	})

}
func DeletePerson(c *gin.Context) {
	id := c.Query("id")
	stmt, err := db.Prepare("delete from mytesttable where id= ?;")
	if err != nil {
		fmt.Print(err.Error())
	}
	_, err = stmt.Exec(id)
	if err != nil {
		fmt.Print(err.Error())
	}
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Successfully deleted user: %s", id),
	})

}
