package helper

import "bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/utils"

const (
	KEEPALIVE string = "Keep Alive"
	GZIP string = "GZip"
	COMPRESS string = "Compress"
	CACHE string = "Cache"
	CDN string = "CDN"
	COMBINE string = "Combine"
	COOKIES string = "Cookies"
	ETAGS string = "ETags"
	MINIFY string = "Minify"
	GRADE_TEXT              string = " optimization grade is "
	//TTFB_SUGGESTION         string = ""
	//TTFB_SUGG_DETAILS string = ""
	//KEEP_ALIVE_SUGGESTION   string = "Performance review score for using persistent connections. "
	//KEEP_ALIVE_SUGG_DETAILS string = "Each request for a piece of content on the page (image, javascript, css, flash, etc) needs to be made over a connection to the web server.  Setting up new connections can take a lot of time so it is best to re-use connections when you can and keep-alive is the way that is done.  They are enabled by default on most configurations and are part of the HTTP 1.1 standard but there are times when they will be broken (sometimes unintentionally).  Enabling keep-alives is usually just a configuration change on the server and does not require any changes to the page itself and can usually reduce the time to load a page by 40-50%. "
	//COMPRESS_TXFER_SUGGESTION string = "performance review score for using gzip compression for transferring compressable responses. "
	//COMPRESS_TXFER_SUGG_DETAILS string = "Enable server side configuration, to compress resources like HTML, Javascript and CSS resulting in a smaller file size and quicker download. Just about everything on a page that isn't an image or video is text of some kind (html, javascript, css).  Text compresses really well and HTTP provides a way to transfer the files in compressed form.  Enabling compression for text resources is usually just a server configuration change without requiring any changes to the page itself and can both improve the performance and reduce the costs of serving the content (by reducing the amount of data transmitted).  Since text resources are usually downloaded at the beginning of the page (javascript and css), delivering them faster has a much larger impact on the user experience than excessive bytes on images or other content. "
	//COMPRESS_IMG_SUGGESTION string = "Performance review score for compressing images."
	//COMPRESS_IMG_SUGG_DETAILS string = "The image compression check just looks at photo images (JPEG files) and makes sure the quality isn't set too high.  JPEG images can usually be compressed pretty substantially without any noticeable reduction in visual quality.  We use a standard of compressing the images at a quality level of '50' in Photoshop's 'Save for Web' mode but generally you should compress them as much as you can before they start to look bad.  It's also not uncommon for other data to be included in photos, particularly if they came from a digital camera (information about the camera, lens, location, even thumbnail images) and some of that should be removed from images before being published to a web page (be careful to retain any copyright information). "
	//CACHE_SUGGESTION string = "Performance review score for leveraging browser caching of requested resource "
	//CACHE_SUGG_DETAILS = "The grade can be improved by caching static resources like images. Static Content are the pieces of content on your page that don't change frequently (images, javascript, css).  You can configure them so that the user's browser will store them in a cache so if the user comes back to the page (or visits another page that uses the same file) they can just use the copy they already have instead of requesting the file from the web server.  Successfully caching static content in the user's browser can significantly improve the performance of a repeat view visit (up to 80+% depending on the page) and reduces the load on the web servers.  It can sometimes be tricky to implement caching without breaking a page though so don't just enable it blindly (you need to be able to change the file name for any files that you expect to change). If the lifetime of a static resource in Cache is greater than a week the score will be better."
	//CDN_SUGGESTION	string = "Performance review score for using CDN for all static assets.
	//CDN_SUGG_DETAILS  string = "A content delivery network(CDN) is a system for distributing resourcesto servers geographically closer to users. Each request for a piece of content to the web server has to travel from the user's browser all the way to the server and back.  As you get further and further from the server this can become a significant amount of time (which adds up quickly as there are more requests on the page).  Ultimately the time it takes is limited by the speed of light so there's not much you can do except to move your server closer to the users.  That is exactly what a Content Distribution Network (CDN) does.  They have servers all over the world that are close to users and they can serve a web site's static content from servers close to users.  The only case where it doesn't make sense to use a CDN is if all of the users for a web site are close to the web server already (like for a community web site). Use CDN to serve static resources. The more resources served from CDN, the better the optimization grade."
	//COMBINE_SUGGESTION string = "Performance review score for bundling JavaScript and/or CSS assets."
	//COMBINE_SUGG_DETAILS string = "Improving performance usually means reducing the number of requests for content and one of the easiest (and most significant) ways to do that is to reduce the number of individual css and javascript files that load at the beginning of the page (in the <head> which blocks the page from displaying to the user).  An easy way to achieve this is to just merge all of the javascript code into a single file and the css into a single file so you have one of each (with the css preferably being loaded before the javascript).  This can have a very substantial impact on the user experience by reducing the amount of time before they see something appear on the screen. "
	//ETAG_SUGGESTION = "Performance review score for disabling *ETag*s."
	//ETAG_SUGG_DETAILS = "By removing the ETag header, you disable caches and browsers from being able to validate files, so they are forced to rely on your Cache-Control and Expires header. Basically you can remove If-Modified-Since and If-None-Match requests and their 304 Not Modified Responses.Entity tags (ETags) are a mechanism to check for a newer version of a cached file."
	//COOKIES_SUGGESTION = "Performance review score for not using cookies on static assets"
	//COOKIES_SUGG_DETAILS = "Do not serve static content like images and stylesheets from a domain that sets cookies. The easiest way to set up a cookieless domain for your static content is to create a CNAME record aliasing your static domain to your main domain.  Your static files remain in the same place but you can now reference them at a different domain."
	//MINIFY_SUGGESTION = "Performance review score for minifying text static assets"
	//MINIFY_SUGG_DETAILS = "Your website has to load a lot of files, including your HTML, CSS, and JavaScript. It includes extra white space, comments, and formatting that computers don’t need in order to run the code. Minifying removes unnecessary characters that are not required for the code to execute."
	//PROGJPEG_SUGG = "Performance review score for image resource using progressive JPEG."
	//PROGJPEG_SUGG_DETAILS = ""
)

type CategorySugg struct  {
	Suggestion    string `json:"suggestion"`
	SuggestionDetails string `json:"suggdetails"`
}
var (
	GradeMap = map[string]string{}
	CategoryMap = map[string]CategorySugg{}
)

func InitGradeMap()  {
	GradeMap["A"] = "Excellent"
	GradeMap["B"] = "Above Average"
	GradeMap["C"] = "Average"
	GradeMap["D"] = "Below Average"
	GradeMap["F"] = "Failed"
	GradeMap["N"] = "Not Applicable"
}

func InitCategoryMap() {

	var category CategorySugg

	category.Suggestion = "Performance review score for using persistent connections. "
	category.SuggestionDetails = "Each request for a piece of content on the page (image, javascript, css, flash, etc) needs to be made over a connection to the web server.  Setting up new connections can take a lot of time so it is best to re-use connections when you can and keep-alive is the way that is done.  They are enabled by default on most configurations and are part of the HTTP 1.1 standard but there are times when they will be broken (sometimes unintentionally).  Enabling keep-alives is usually just a configuration change on the server and does not require any changes to the page itself and can usually reduce the time to load a page by 40-50%. "
	CategoryMap[KEEPALIVE] = category

	category.Suggestion = "Performance review score for using gzip compression for transferring compressable responses. "
	category.SuggestionDetails = "Enable server side configuration, to compress resources like HTML, Javascript and CSS resulting in a smaller file size and quicker download. Just about everything on a page that isn't an image or video is text of some kind (html, javascript, css).  Text compresses really well and HTTP provides a way to transfer the files in compressed form.  Enabling compression for text resources is usually just a server configuration change without requiring any changes to the page itself and can both improve the performance and reduce the costs of serving the content (by reducing the amount of data transmitted).  Since text resources are usually downloaded at the beginning of the page (javascript and css), delivering them faster has a much larger impact on the user experience than excessive bytes on images or other content. "
	CategoryMap[GZIP] = category

	category.Suggestion = "Performance review score for compressing images. "
	category.SuggestionDetails = "The image compression check just looks at photo images (JPEG files) and makes sure the quality isn't set too high.  JPEG images can usually be compressed pretty substantially without any noticeable reduction in visual quality.  We use a standard of compressing the images at a quality level of '50' in Photoshop's 'Save for Web' mode but generally you should compress them as much as you can before they start to look bad.  It's also not uncommon for other data to be included in photos, particularly if they came from a digital camera (information about the camera, lens, location, even thumbnail images) and some of that should be removed from images before being published to a web page (be careful to retain any copyright information)."
	CategoryMap[COMPRESS] = category

	category.Suggestion = "Performance review score for using CDN for all static assets. "
	category.SuggestionDetails = "A content delivery network (CDN) is a system for distributing resources to servers geographically closer to users. Each static resource is checked to see if its host server was one such provider. The more resources served from a CDN, the better the value."
	CategoryMap[CDN] = category

	category.Suggestion = "Performance review score for leveraging browser caching of requested resource. "
	category.SuggestionDetails = "The grade can be improved by caching static resources like images. Static Content are the pieces of content on your page that don't change frequently (images, javascript, css).  You can configure them so that the user's browser will store them in a cache so if the user comes back to the page (or visits another page that uses the same file) they can just use the copy they already have instead of requesting the file from the web server.  Successfully caching static content in the user's browser can significantly improve the performance of a repeat view visit (up to 80+% depending on the page) and reduces the load on the web servers.  It can sometimes be tricky to implement caching without breaking a page though so don't just enable it blindly (you need to be able to change the file name for any files that you expect to change). If the lifetime of a static resource in Cache is greater than a week the score will be better."
	CategoryMap[CACHE] = category

	category.Suggestion = "Performance review score for bundling JavaScript and/or CSS assets. "
	category.SuggestionDetails = "Improving performance usually means reducing the number of requests for content and one of the easiest (and most significant) ways to do that is to reduce the number of individual css and javascript files that load at the beginning of the page (in the <head> which blocks the page from displaying to the user).  An easy way to achieve this is to just merge all of the javascript code into a single file and the css into a single file so you have one of each (with the css preferably being loaded before the javascript).  This can have a very substantial impact on the user experience by reducing the amount of time before they see something appear on the screen. "
	CategoryMap[COMBINE] = category

	category.Suggestion = "Performance review score for disabling *ETag*s. "
	category.SuggestionDetails = "By removing the ETag header, you disable caches and browsers from being able to validate files, so they are forced to rely on your Cache-Control and Expires header. Basically you can remove If-Modified-Since and If-None-Match requests and their 304 Not Modified Responses.Entity tags (ETags) are a mechanism to check for a newer version of a cached file."
	CategoryMap[ETAGS] = category

	category.Suggestion = "Performance review score for not using cookies on static assets. "
	category.SuggestionDetails = "Do not serve static content like images and stylesheets from a domain that sets cookies. The easiest way to set up a cookieless domain for your static content is to create a CNAME record aliasing your static domain to your main domain.  Your static files remain in the same place but you can now reference them at a different domain."
	CategoryMap[COOKIES] = category

	category.Suggestion = "Performance review score for minifying text static assets. "
	category.SuggestionDetails = "Your website has to load a lot of files, including your HTML, CSS, and JavaScript. It includes extra white space, comments, and formatting that computers don’t need in order to run the code. Minifying removes unnecessary characters that are not required for the code to execute."
	CategoryMap[MINIFY] = category
}

func BuildSuggestion(category string, percent int) (string) {
	utils.SpectreLog.Debug("Getting grade for category %s ", category)
	utils.SpectreLog.Debugln(" with percent %d.", percent)
	if grade, err := GetGrade(percent); err == nil {
		if grade == "N" {
			return GradeMap["N"]
		} else {
			return CategoryMap[category].Suggestion + GradeMap[grade] + ": " + category + GRADE_TEXT + grade
		}
	} else {
		utils.SpectreLog.Warnf("Error %v getting grades for category %s", err, category)
	}
	return ""
}

func BuildSuggestionDetails(category string, percent int) (string) {
	if grade, err := GetGrade(percent); err == nil {
		if grade == "A" {
			return "Achieved optimum results for " + category + ". For improving performance further, optimize performance grades for other metrics."
		}
	}

	if percent == -1 {
		return "Not Applicable"
	}

	return CategoryMap[category].SuggestionDetails
}

