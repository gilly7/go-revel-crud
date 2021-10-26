package controllers

import (
	"fmt"
	"go-revel-crud/app/db"
	"go-revel-crud/app/entities"
	"go-revel-crud/app/forms"
	"go-revel-crud/app/models"
	"go-revel-crud/app/webutils"
	"net/http"
	"strings"

	"github.com/revel/revel"
)

type Post struct {
	App
}

func (c Post) List() revel.Result {
	var result entities.Response

	ctx := c.Request.Context()

	paginationFilter, err := webutils.FilterFromQuery(c.Params)
	if err != nil {
		c.Log.Errorf("could not filter from params: %v", err)
		result = entities.Response{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: "Failed to parse page filters",
		}
		return c.RenderJSON(result)
	}

	newPost := &models.Post{}
	data, err := newPost.All(ctx, db.DB(), paginationFilter)
	if err != nil {
		c.Log.Errorf("could not get posts: %v", err)
		result = entities.Response{
			Success: false,
			Status:  http.StatusInternalServerError,
			Message: "Could not get posts",
		}
		return c.RenderJSON(result)
	}

	recordsCount, err := newPost.Count(c.Request.Context(), db.DB(), paginationFilter)
	if err != nil {
		c.Log.Errorf("could not get posts count: %v", err)
		result = entities.Response{
			Success: false,
			Status:  http.StatusInternalServerError,
			Message: "Could not get posts count",
		}
		return c.Render(result)
	}

	result = entities.Response{
		Success: true,
		Status:  http.StatusOK,
		Data: map[string]interface{}{
			"posts":      data,
			"pagination": models.NewPagination(recordsCount, paginationFilter.Page, paginationFilter.Per),
		},
	}
	return c.RenderJSON(result)
}

func (c Post) Get(id int64) revel.Result {
	result := &entities.Response{}
	ctx := c.Request.Context()

	newPost := &models.Post{}
	post, err := newPost.ByID(ctx, db.DB(), id)
	if err != nil {
		c.Log.Errorf("could not get post: %v", err)
		result.Success = false
		result.Message = "Could not get the post"
		result.Data = post
		return c.RenderJSON(result)
	}

	result.Success = true
	result.Data = post
	return c.RenderJSON(result)
}

func (c Post) Add() revel.Result {

	var status int
	postForm := forms.Post{}
	c.Params.BindJSON(&postForm)

	ctx := c.Request.Context()

	v := c.Validation
	postForm.Validate(v)
	if v.HasErrors() {
		retErrors := make([]string, 0)
		for _, theErr := range v.Errors {
			retErrors = append(retErrors, theErr.Message)
		}
		status = http.StatusBadRequest
		c.Response.SetStatus(status)
		return c.RenderJSON(entities.Response{
			Message: strings.Join(retErrors, ","),
			Status:  status,
			Success: false,
		})
	}

	newPost := &models.Post{
		Username: postForm.Username,
		Title:    postForm.Title,
		Content:  postForm.Content,
	}
	err := newPost.Save(ctx, db.DB())
	if err != nil {
		c.Log.Errorf("could not save post: %v", err)
		status = http.StatusInternalServerError
		c.Response.SetStatus(status)
		return c.RenderJSON(entities.Response{
			Message: "Encountered an error saving request.",
			Status:  status,
			Success: false,
		})
	}

	status = http.StatusCreated
	c.Response.SetStatus(status)
	return c.RenderJSON(entities.Response{
		Data:    newPost,
		Status:  status,
		Success: true,
	})
}

func (c *Post) Update(id int64) revel.Result {
	var status int
	data := models.Post{}
	c.Params.BindJSON(&data)
	data.ID = id

	err := data.Save(c.Request.Context(), db.DB())
	if err != nil {
		c.Log.Errorf("could not save post: %v", err)
		status = http.StatusInternalServerError
		c.Response.SetStatus(status)
		return c.RenderJSON(entities.Response{
			Message: "Encountered an error saving request.",
			Status:  status,
			Success: false,
		})

	}

	status = http.StatusOK
	c.Response.SetStatus(status)
	return c.RenderJSON(entities.Response{
		Data:    data,
		Status:  status,
		Success: true,
	})
}

func (c *Post) Delete(id int64) revel.Result {
	var status int
	ctx := c.Request.Context()

	data := models.Post{}
	data.ID = id
	_, err := data.Delete(ctx, db.DB())
	if err != nil {
		c.Log.Errorf("could not delete post: %v", err)
		status = http.StatusInternalServerError
		c.Response.SetStatus(status)
		return c.RenderJSON(entities.Response{
			Message: "Encountered an error deleting the post",
			Status:  status,
			Success: false,
		})
	}

	status = http.StatusOK
	c.Response.SetStatus(status)
	return c.RenderJSON(entities.Response{
		Status:  status,
		Success: true,
		Message: fmt.Sprintf("post id =[%v] deleted", id),
	})
}
