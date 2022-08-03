package services

import "github.com/marius004/phoenix-algo/models"

func makeProblemFilter(filter *models.ProblemFilter) (query []string, args []interface{}) {
	if filter.AuthorId > 0 {
		query = append(query, "author_id = ?")
		args = append(args, filter.AuthorId)
	}

	if filter.ProblemId > 0 {
		query = append(query, "id = ?")
		args = append(args, filter.ProblemId)
	}

	if filter.Status != "" {
		query = append(query, "status = ?")
		args = append(args, filter.Status)
	}

	return
}

func makeSubmissionFilter(filter models.SubmissionFilter) (query []string, args []interface{}) {
	if filter.UserId > 0 {
		query = append(query, "user_id = ?")
		args = append(args, filter.UserId)
	}

	if filter.ProblemId >= 0 {
		query = append(query, "problem_id = ?")
		args = append(args, filter.ProblemId)
	}

	if filter.Score >= 0 {
		query = append(query, "score = ?")
		args = append(args, filter.Score)
	}

	if filter.Status != "" {
		query = append(query, "status = ?")
		args = append(args, filter.Status)
	}

	if filter.CompiledSuccesfully != nil {
		query = append(query, "compiled_succesfully = ?")
		args = append(args, filter.CompiledSuccesfully)
	}

	return
}

func makeUserFilter(filter *models.UserFilter) (query []string, args []interface{}) {
	if filter.Email != "" {
		query = append(query, "email = ?")
		args = append(args, filter.Email)
	}

	if filter.GithubURL != "" {
		query = append(query, "github_url = ?")
		args = append(args, filter.GithubURL)
	}

	if filter.LinkedInURL != "" {
		query = append(query, "linked_in_url = ?")
		args = append(args, filter.LinkedInURL)
	}

	if filter.UserIconURL != "" {
		query = append(query, "user_icon_url = ?")
		args = append(args, filter.UserIconURL)
	}

	if filter.Username != "" {
		query = append(query, "username = ?")
		args = append(args, filter.Username)
	}

	if filter.WebsiteURL != "" {
		query = append(query, "website_url = ?")
		args = append(args, filter.WebsiteURL)
	}

	if filter.IsAdmin != nil {
		query = append(query, "is_admin = ?")
		args = append(args, filter.IsAdmin)
	}

	if filter.IsProposer != nil {
		query = append(query, "is_proposer = ?")
		args = append(args, filter.IsProposer)
	}

	if filter.UserId > 0 {
		query = append(query, "id = ?")
		args = append(args, filter.UserId)
	}

	return
}
