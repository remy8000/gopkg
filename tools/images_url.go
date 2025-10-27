package tools

import (
	"net/url"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

// clean any url from its added query And unescape characters
// https://www.lesoir.be/sites/default/files/dpistyles_v2/2022/11/03/B9732577116Z.1_20221103201825_000+GS7LJJH34.1-0.jpg?itok=UruRYvZc1667503111
// => https://www.lesoir.be/sites/default/files/dpistyles_v2//2022/11/03/B9732577116Z.1_20221103201825_000+GS7LJJH34.1-0.jpg
func RemoveImgQueryParams(imgUrl string) (string, error) {
	parsedUrl, err := url.Parse(imgUrl)
	if err != nil {
		return "", err
	}
	parsedUrl.RawQuery = ""
	return parsedUrl.String(), nil
}

// extract form any url string an img url matching
func ExtractImgURL(input string) string {
	re := regexp.MustCompile(`(?i)https?://\S+?\.(?:jpg|jpeg|png|gif)`)
	match := re.FindString(input)
	return match
}

// extract the first img url from a html content
func ExtractFirstImgURLFromHtmlContent(htmlContent string) (string, error) {

	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return "", err
	}

	var extract func(*html.Node) string
	extract = func(n *html.Node) string {
		if n.Type == html.ElementNode && n.Data == "img" {
			for _, attr := range n.Attr {
				if attr.Key == "src" {
					return attr.Val
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			result := extract(c)
			if result != "" {
				return result
			}
		}

		return ""
	}

	imgURL := extract(doc)

	return imgURL, nil
}

// extract all img urls from html content
func ExtractAllImgURLsFromHtmlContent(htmlContent string) ([]string, error) {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return nil, err
	}

	var imgURLs []string
	var extract func(*html.Node)
	extract = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "img" {
			for _, attr := range n.Attr {
				if attr.Key == "src" && attr.Val != "" {
					imgURLs = append(imgURLs, attr.Val)
					break
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extract(c)
		}
	}

	extract(doc)

	return imgURLs, nil
}
