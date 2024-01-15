package polyfea

import (
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/polyfea/polyfea-controller/api/v1alpha1"
	"github.com/polyfea/polyfea-controller/repository"
)

const (
	NamespacePathParamName     = "namespace"
	MicrofrontendPathParamName = "microfrontend"
	PathPathParamName          = "path"
)

type PolyfeaProxy struct {
	microfrontendClassRepository repository.PolyfeaRepository[*v1alpha1.MicroFrontendClass]
	microfrontendRepository      repository.PolyfeaRepository[*v1alpha1.MicroFrontend]
	client                       *http.Client
}

func NewPolyfeaProxy(
	microfrontendClassRepository repository.PolyfeaRepository[*v1alpha1.MicroFrontendClass],
	microfrontendRepository repository.PolyfeaRepository[*v1alpha1.MicroFrontend],
	httpClient *http.Client) *PolyfeaProxy {

	return &PolyfeaProxy{
		microfrontendClassRepository: microfrontendClassRepository,
		microfrontendRepository:      microfrontendRepository,
		client:                       httpClient,
	}
}

func (p *PolyfeaProxy) HandleProxy(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	nameSpace := params[NamespacePathParamName]
	nameMicroFrontend := params[MicrofrontendPathParamName]
	path := params[PathPathParamName]

	microfrontends, err := p.microfrontendRepository.GetItems(func(mf *v1alpha1.MicroFrontend) bool {
		return mf.Namespace == nameSpace && mf.Name == nameMicroFrontend
	})

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if len(microfrontends) == 0 {
		log.Println("No microfrontend found for the given namespace and name.")
		http.Error(w, "No microfrontend found for the given namespace and name.", http.StatusNotFound)
		return
	}

	microfrontend := microfrontends[0]

	microfrontendClasses, err := p.microfrontendClassRepository.GetItems(func(mfc *v1alpha1.MicroFrontendClass) bool {
		return mfc.Name == *microfrontend.Spec.FrontendClass
	})

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if len(microfrontendClasses) == 0 {
		log.Println("No microfrontend class found for the given namespace and name.")
		http.Error(w, "No microfrontend class found for the given namespace and name.", http.StatusNotFound)
		return
	}

	microfrontendClass := microfrontendClasses[0]

	proxyUrl := *microfrontend.Spec.Service + path

	if (*microfrontend.Spec.Service)[len(*microfrontend.Spec.Service)-1] != '/' && path[0] != '/' {
		proxyUrl = *microfrontend.Spec.Service + "/" + path
	}

	req, err := http.NewRequest("GET", proxyUrl, r.Body)

	copyHeaders(req.Header, r.Header)

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Proxying request to the module.", "Resolved URL:", proxyUrl)
	resp, err := p.client.Do(req)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	io.Copy(w, resp.Body)

	copyHeaders(w.Header(), resp.Header)

	copyExtraHeaders(w.Header(), microfrontendClass.Spec.ExtraHeaders)

	w.WriteHeader(resp.StatusCode)
}

func copyHeaders(dst, src http.Header) {
	for key, values := range src {
		for _, value := range values {
			dst.Add(key, value)
		}
	}
}

func copyExtraHeaders(dst http.Header, extraHeaders []v1alpha1.Header) {
	for _, extraHeader := range extraHeaders {
		dst.Add(extraHeader.Name, extraHeader.Value)
	}
}
