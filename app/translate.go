package app

import "github.com/karintomania/kaigai-go-scraper/db"

func translate(
	dateString string,
	pageRepository *db.PageRepository,
	commentRepository *db.CommentRepository,
) error {
	pages := pageRepository.FindByDate(dateString)

	for _, page := range pages {
		comments := commentRepository.FindByPageId(page.Id)

		translatePage(&page, comments)

		// update page, comments
		pageRepository.Update(&page)

		for _, comment := range comments {
			commentRepository.Update(&comment)
		}
	}

	return nil
}

func translatePage(page *db.Page, comments []db.Comment){
	// translate page

	// translate comments
}
