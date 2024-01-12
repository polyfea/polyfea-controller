package polyfea

import (
	"bytes"
	"crypto/rand"
	_ "embed"
	"encoding/base64"
	"html/template"
	"log"
	"net/http"
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

//go:embed .resources/index.html
var html string

//go:embed .resources/boot.mjs
var bootJs []byte

func NewSinglePageApplication(microFrontendClassRepository repository.PolyfeaRepository[*v1alpha1.MicroFrontendClass]) *SingePageApplication {
	return &SingePageApplication{
		microFrontendClassRepository: microFrontendClassRepository,
	}
}

func (s *SingePageApplication) HandleSinglePageApplication(w http.ResponseWriter, r *http.Request) {

	basePath := getBasePath(r)
	microFrontendClass, err := s.getMicrofrontendClass(basePath)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	nonce, err := generateNonce()

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	extraMeta := ""

	for _, metaTag := range microFrontendClass.Spec.ExtraMetaTags {
		extraMeta += "<meta name=\"" + metaTag.Name + "\" content=\"" + metaTag.Content + "\" >"
	}

	templateVars := TemplateData{
		BaseUri:   basePath,
		Title:     *microFrontendClass.Spec.Title,
		Nonce:     nonce,
		ExtraMeta: template.HTML(extraMeta),
	}

	templatedHtml, err := templateHtml(html, &templateVars)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	cspHeader := strings.ReplaceAll(microFrontendClass.Spec.CspHeader, "{NONCE_VALUE}", nonce)

	for _, header := range microFrontendClass.Spec.ExtraHeaders {
		w.Header().Set(header.Name, header.Value)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Content-Security-Policy", cspHeader)

	w.Write([]byte(templatedHtml))
}

func (s *SingePageApplication) HandleBootJs(w http.ResponseWriter, r *http.Request) {
	basePath := getBasePath(r)
	microFrontendClass, err := s.getMicrofrontendClass(basePath)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	for _, header := range microFrontendClass.Spec.ExtraHeaders {
		w.Header().Set(header.Name, header.Value)
	}

	w.Header().Set("Content-Type", "application/javascript;")

	w.Write(bootJs)
}

func templateHtml(content string, templateVars *TemplateData) (string, error) {

	tmpl, err := template.New("index.html").Parse(content)
	if err != nil {
		log.Println(err)
		return "", err
	}

	var tmplBytes bytes.Buffer
	if err := tmpl.Execute(&tmplBytes, templateVars); err != nil {
		log.Println(err)
		return "", err
	}

	return tmplBytes.String(), nil
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

func (s *SingePageApplication) getMicrofrontendClass(basePath string) (*v1alpha1.MicroFrontendClass, error) {
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
		return nil, err
	}

	if len(microFrontendClasses) == 0 {
		return nil, nil
	}

	return microFrontendClasses[0], nil
}

func getBasePath(r *http.Request) string {
	ctx := r.Context()
	basePathValue := ctx.Value(PolyfeaContextKeyBasePath)
	var basePath string
	if basePathValue == nil {
		basePath = "/"
	} else {
		basePath = basePathValue.(string)
	}

	return basePath
}
