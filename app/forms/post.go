package forms

import "github.com/revel/revel"

type (
	Post struct {
		Username string
		Title    string
		Content  string
	}
)

func (post *Post) Validate(v *revel.Validation) {
	ValidatePostTitle(v, post.Title)
	ValidatePostContent(v, post.Content)
}

func ValidatePostTitle(v *revel.Validation, title string) *revel.ValidationResult {
	result := v.Required(title).Message("Post title is required")
	if !result.Ok {
		return result
	}

	return result
}

func ValidatePostContent(v *revel.Validation, content string) *revel.ValidationResult {
	result := v.Required(content).Message("Post content is required")
	if !result.Ok {
		return result
	}

	return result
}
