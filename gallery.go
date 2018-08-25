package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"io/ioutil"
	"gopkg.in/h2non/filetype.v1"
	"os"
	"log"
	"path"
	"github.com/rakyll/statik/fs"
	_ "github.com/blan4/galery/statik"
	"flag"
	"html/template"
)

var flagSrc = flag.String("src", ".", "The path of the source directory.")

func main() {
	flag.Parse()
	runServer(*flagSrc)
}

func runServer(dirname string) {
	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}
	r := gin.Default()
	r.SetHTMLTemplate(loadIndexTemplate())
	//r.Static("/assets", "./assets")
	r.StaticFS("/assets", statikFS)
	r.StaticFS("/img", http.Dir(dirname))
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title":  "Main website",
			"images": ls(dirname),
		})
	})
	r.Run()
}

func loadIndexTemplate() *template.Template {
	t, err := template.New("index.tmpl").Parse(indexTemplate)
	if err != nil {
		log.Fatal(err)
	}
	return t
}

func ls(dirname string) (fnames []string) {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		log.Fatal(err)
		return
	}
	for _, file := range files {
		fileName := path.Join(dirname, file.Name())
		if file.IsDir() {
			continue // TODO: maybe recursive?
		} else {
			buf, err := os.Open(fileName)
			head := make([]byte, 261)
			buf.Read(head)
			if err != nil {
				log.Fatal(err)
			} else {
				if filetype.IsImage(head) {
					fnames = append(fnames, file.Name())
				}
			}
		}
	}
	return
}

var indexTemplate string = `
<html>
    <head>
        <link type="text/css" rel="stylesheet" href="/assets/lightgallery/css/lightgallery.css" />
        <link type="text/css" rel="stylesheet" href="/assets/normalize/normalize.css" />
        <link type="text/css" rel="stylesheet" href="/assets/main.css" />
    </head>
    <body>
    	<div class='gallery'>
            <ul id="lightgallery">
                {{ range .images}}
                <li data-src="img/{{ . }}" >
                    <img src="img/{{ . }}" />
                </li>
                {{ end }}
            </ul>
        </div>
    	<script src="/assets/lightgallery/js/lightgallery.js"></script>
    	<script src="/assets/lightgallery/js/lg-thumbnail.js"></script>
        <script src="/assets/lightgallery/js/lg-fullscreen.js"></script>
    	<script type="text/javascript">
            lightGallery(document.getElementById('lightgallery'));
        </script>
    </body>
</html>
`
