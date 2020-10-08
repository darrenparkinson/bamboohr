package bamboohr

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// EmployeeCategoryResponse is the top level response from the API
type EmployeeCategoryResponse struct {
	EmployeeID struct {
		ID int
	} `json:"employee"`
	Categories []EmployeeCategory
}

// EmployeeCategory represents a files category (or folder!)
type EmployeeCategory struct {
	ID                int
	Name              string
	CanRenameCategory string
	CanDeleteCategory string
	CanUploadFiles    string
	DisplayIfEmpty    string
	Files             []File
}

// File represents an individual file
type File struct {
	ID                int
	Name              string
	OriginalFileName  string
	Size              int
	DateCreated       string
	CreatedBy         string
	ShareWithEmployee string
}

// GetEmployeeFilesAndCategories returns a list of employee files and categories
func (c *Client) GetEmployeeFilesAndCategories(ctx context.Context, id string) ([]EmployeeCategory, error) {
	url := fmt.Sprintf("%s/employees/%s/files/view/", c.BaseURL, id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	ec := EmployeeCategoryResponse{}
	if err := c.makeRequest(req, &ec); err != nil {
		return nil, err
	}
	return ec.Categories, nil
}

// UploadEmployeeFile uploads a file to a specific employees files under the given category ID.
// Beware the inconsistent ID types Bamboo uses.  We require all strings here.
func (c *Client) UploadEmployeeFile(ctx context.Context, employeeID, categoryID, fileName, filePath, share string) error {

	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		return err
	}

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	err = writer.WriteField("category", categoryID)
	if err != nil {
		return err
	}
	err = writer.WriteField("fileName", fileName)
	if err != nil {
		return err
	}
	err = writer.WriteField("share", share)
	if err != nil {
		return err
	}

	part4, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return err
	}
	_, err = io.Copy(part4, file)
	if err != nil {
		return err
	}
	err = writer.Close()
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/employees/%s/files/", c.BaseURL, employeeID)

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req = req.WithContext(ctx)
	if err := c.makeRequest(req, nil); err != nil {
		return err
	}
	return nil
}
