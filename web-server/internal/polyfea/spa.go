package polyfea

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/polyfea/polyfea-controller/api/v1alpha1"
	"github.com/polyfea/polyfea-controller/repository"
)

type SingePageApplication struct {
	microFrontendClassRepository repository.PolyfeaRepository[*v1alpha1.MicroFrontendClass]
}

type TemplateData struct {
	BaseUri   string
	Title     string
	Nonce     string
	ExtraMeta template.HTML
}

func NewSinglePageApplication(microFrontendClassRepository repository.PolyfeaRepository[*v1alpha1.MicroFrontendClass]) *SingePageApplication {
	return &SingePageApplication{
		microFrontendClassRepository: microFrontendClassRepository,
	}
}

func (s *SingePageApplication) HandleSinglePageApplication(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	basePathValue := ctx.Value(PolyfeaContextKeyBasePath)
	var basePath string
	if basePathValue == nil {
		basePath = "/"
	} else {
		basePath = basePathValue.(string)
	}

	microFrontendClasses, err := s.microFrontendClassRepository.GetItems(func(mfc *v1alpha1.MicroFrontendClass) bool {
		frontendClassBasePath := *mfc.Spec.BaseUri
		if frontendClassBasePath[0] != '/' {
			frontendClassBasePath = "/" + frontendClassBasePath
		}
		if frontendClassBasePath[len(frontendClassBasePath)-1] != '/' {
			frontendClassBasePath += "/"
		}

		return basePath == frontendClassBasePath
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(microFrontendClasses) == 0 {
		http.NotFound(w, r)
		return
	}

	microFrontendClass := microFrontendClasses[0]

	data, err := os.ReadFile("../../.template/index.html")

	if err != nil {
		if os.IsNotExist(err) {
			log.Println("index.html does not exist!")
			http.NotFound(w, r)
			return
		}
		log.Panic(err)
	}

	nonce, err := generateNonce()

	if err != nil {
		log.Panic(err)
	}

	fileContent := string(data)

	extraMeta := ""

	for _, metaTag := range microFrontendClass.Spec.ExtraMetaTags {
		extraMeta += "<meta name=\"" + metaTag.Name + "\" content=\"" + metaTag.Content + "\" >"
	}

	templateVars := TemplateData{
		BaseUri:   basePath,
		Title:     *microFrontendClasses[0].Spec.Title,
		Nonce:     nonce,
		ExtraMeta: template.HTML(extraMeta),
	}

	templatedHtml := templateHtml(fileContent, &templateVars)

	cspHeader := strings.ReplaceAll(microFrontendClass.Spec.CspHeader, "{NONCE_VALUE}", nonce)

	for _, header := range microFrontendClass.Spec.ExtraHeaders {
		w.Header().Set(header.Name, header.Value)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Content-Security-Policy", cspHeader)

	w.Write([]byte(templatedHtml))
}

func templateHtml(content string, templateVars *TemplateData) string {

	tmpl, err := template.New("index.html").Parse(content)
	if err != nil {
		log.Panic(err)
	}

	var tmplBytes bytes.Buffer
	if err := tmpl.Execute(&tmplBytes, templateVars); err != nil {
		log.Panic(err)
	}

	return tmplBytes.String()
}

func generateNonce() (string, error) {
	codes := make([]byte, 128)
	_, err := rand.Read(codes)
	if err != nil {
		return "", err
	}

	text := string(codes)
	nonce := base64.StdEncoding.EncodeToString([]byte(text))

	return nonce, nil
}
