# Bamboo HR Go Library 

[![GoDoc](https://godoc.org/github.com/darrenparkinson/bamboohr?status.svg)](https://godoc.org/github.com/darrenparkinson/bamboohr)
[![PkgGoDev](https://pkg.go.dev/badge/darrenparkinson/bamboohr)](https://pkg.go.dev/github.com/darrenparkinson/bamboohr)
[![Go Report Card](https://goreportcard.com/badge/github.com/darrenparkinson/bamboohr)](https://goreportcard.com/report/github.com/darrenparkinson/bamboohr)

This is a simple Go library to utilise the [Bamboo HR API](https://documentation.bamboohr.com/docs).

## Installation

To get the library, simply:

> `go get github.com/darrenparkinson/bamboohr`


## Using the Library

To use the library you can use the `New` helper function and provide your API Key and your Company Domain, or provide these directly to the `Client` struct.

To get an API Key, follow the instructions on the [Bamboo HR site](https://documentation.bamboohr.com/docs#section-authentication), but essentially: 

> *To generate an API key, users should log in and click their name in the upper right-hand corner of any page to get to the user context menu. If they have sufficient permissions, there will be an "API Keys" option in that menu to go to the page.*

Your company domain is the part before bamboohr.com, so for `https://acmecorp.bamboohr.com` it would be `acmecorp`.

```go
package main

import (
    "log"
    "os"

    "github.com/darrenparkinson/bamboohr"
)

func main() {
    apikey := os.Getenv("BAMBOO_API_KEY")
    bamboo, _ := bamboohr.New(apikey, "acmecorp", nil)
    ctx := context.Background()
    people, _ := bamboo.GetEmployeeDirectory(ctx)
    for _, person := range people {
        log.Println(person.ID, person.Displayname)
    }

    // Note: 0 is a special ID meaning the user that created the API Key
    // If no field names are specified after the ID, then all fields will be fetched.
    me, _ := bamboo.GetEmployee(ctx, "0", bamboohr.DisplayName, bamboohr.FirstName, bamboohr.LastName)
    log.Println(me.ID, me.FirstName, me.LastName, me.DisplayName)
}
```

A context is required for each API request, for which you can just provide `context.Background()`, but you can also provide a timeout for your requests if necessary:

```go
ctx := context.Background()
ctx, _ = context.WithTimeout(ctx, 1*time.Second)

people, err := bamboo.GetEmployeeDirectory(ctx)
if err != nil {
	log.Fatal(err)
}
```

This will likely return a `context deadline exceeded` error since the request will take longer than 1 second.

## Documentation

There is an online reference for the package at
[godoc.org/github.com/darrenparkinson/bamboohr][godoc-bamboohr].

Bamboo HR Documentation is available [on their site](https://documentation.bamboohr.com/docs).

So far, the following endpoints have been implemented:

**Employees**

* [Get Employee](https://documentation.bamboohr.com/reference#get-employee)
* [Get Employee Directory](https://documentation.bamboohr.com/reference#get-employees-directory-1)

**Employee Files**

* [List Employee Files and Categories](https://documentation.bamboohr.com/reference#list-employee-files-1)
* [Upload Employee File](https://documentation.bamboohr.com/reference#upload-employee-file-1)

**Account Information**

These have been removed temporarily due to some inconsistencies with the ID field returned from Bamboo.

* [Get A List of Fields](https://documentation.bamboohr.com/reference#metadata-get-a-list-of-fields)
* [Get Details For List Fields](https://documentation.bamboohr.com/reference#metadata-get-details-for-list-fields-1)