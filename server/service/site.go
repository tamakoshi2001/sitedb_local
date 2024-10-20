package service

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"sort"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/gocolly/colly"
	"github.com/tamakoshi2001/gextension/model"
	"gonum.org/v1/gonum/mat"
)

var OPENAI_API_KEY = os.Getenv("OPENAI_API_KEY")

// A SiteService implements handling REST endpoints.
type SiteService struct {
	sites   *[]model.Site
	vectors *[]model.Vector
}

// NewSiteService returns a new SiteService.
func NewSiteService(sites *[]model.Site, vectors *[]model.Vector) *SiteService {
	return &SiteService{
		sites:   sites,
		vectors: vectors,
	}
}

// Create handles the endpoint that creates the Site.
func (s *SiteService) Create(r *model.CreateSiteRequest) (*model.CreateSiteResponse, error) {
	// get url from request
	url := r.URL
	// chech if url is already registered
	for _, site := range *s.sites {
		if site.URL == url {
			log.Println("url is already registered")
			return &model.CreateSiteResponse{Site: site}, nil
		}
	}
	var title, body, post, question string

	c := colly.NewCollector()
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:92.0) Gecko/20100101 Firefox/92.0")
		r.Headers.Set("Referer", "https://www.google.com/")
	})
	c.OnHTML("title", func(e *colly.HTMLElement) {
		title = e.Text
	})
	c.OnHTML("body", func(e *colly.HTMLElement) {
		body = e.Text
	})
	err := c.Visit(url)
	if err != nil {
		return nil, err
	}

	// Combine title and body
	post = title + CutString(body, 60000)
	question = "200字で要約して。 \n" + post
	summary, err := getCompletionResponse(question)
	if err != nil {
		return nil, err
	}
	embedding, err := getEmbeddingResponse(post)
	if err != nil {
		return nil, err
	}

	site := model.Site{
		URL:       url,
		Title:     title,
		Summary:   summary.Choices[0].Message.Content,
		CreatedAt: time.Now(),
	}
	vector := model.Vector{
		URL:    url,
		Vector: &embedding.Data[0].Embedding,
	}
	*s.sites = append(*s.sites, site)
	*s.vectors = append(*s.vectors, vector)

	response := model.CreateSiteResponse{
		Site: site}

	return &response, nil
}

// Read handles the endpoint that reads the Site.
func (s *SiteService) Read(r *model.ReadSiteRequest) (*model.ReadSiteResponse, error) {
	var wg sync.WaitGroup

	query := r.Query
	e, err := getEmbeddingResponse(CutString(query, 60000))
	if err != nil {
		return nil, err
	}
	embedding := mat.NewVecDense(len(e.Data[0].Embedding), e.Data[0].Embedding)

	similarity := make([]float64, len(*s.sites))
	// calculate cosine similarity for each s.vector by parallel
	for i, v := range *s.vectors {
		wg.Add(1)
		go func(i int, v model.Vector) {
			similarity[i] = CosineSimilarity(embedding, mat.NewVecDense(len(*v.Vector), *v.Vector))
			wg.Done()
		}(i, v)
	}
	wg.Wait()

	// Pairを使用してkeysとvaluesを組み合わせる
	type Pair struct {
		Key   float64
		Value model.Site
	}
	pairs := make([]Pair, len(*s.sites))
	for i := range *s.sites {
		pairs[i] = Pair{similarity[i], (*s.sites)[i]}
	}

	// keysを基にソートする
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Key > pairs[j].Key
	})

	// ソート後の結果をmodel.SiteResponseに格納する
	var sites []model.Site
	for _, pair := range pairs {
		sites = append(sites, pair.Value)
	}
	return &model.ReadSiteResponse{SiteS: sites}, nil
}

// Delete handles the endpoint that deletes the Site.
func (s *SiteService) Delete(r *model.DeleteSiteRequest) (*model.DeleteSiteResponse, error) {
	url := r.URL
	for i, site := range *s.sites {
		if site.URL == url {
			*s.sites = append((*s.sites)[:i], (*s.sites)[i+1:]...)
			*s.vectors = append((*s.vectors)[:i], (*s.vectors)[i+1:]...)
			return &model.DeleteSiteResponse{Response: "success"}, nil
		}
	}
	return &model.DeleteSiteResponse{Response: "this url is not registered"}, nil
}

func getCompletionResponse(question string) (*model.CompletionResponse, error) {
	const completionURL = "https://api.openai.com/v1/chat/completions"

	messages := []model.Message{}
	messages = append(messages, model.Message{
		Role:    "user",
		Content: question,
	})

	requestBody := model.CompletionRequest{
		Model:    "gpt-4o-mini",
		Messages: messages,
	}

	requestJSON, _ := json.Marshal(requestBody)

	req, err := http.NewRequest("POST", completionURL, bytes.NewBuffer(requestJSON))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+OPENAI_API_KEY)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response model.CompletionResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func getEmbeddingResponse(question string) (*model.EmbeddingResponse, error) {
	const embeddingURL = "https://api.openai.com/v1/embeddings"

	requestBody := model.EmbeddingRequest{
		Model: "text-embedding-3-small",
		Input: question,
	}

	requestJSON, _ := json.Marshal(requestBody)

	req, err := http.NewRequest("POST", embeddingURL, bytes.NewBuffer(requestJSON))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+OPENAI_API_KEY)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response model.EmbeddingResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		println("Error: ", err.Error())
		return &model.EmbeddingResponse{}, nil
	}

	return &response, nil
}

func CutString(s string, cnt int64) string {
	if int64(utf8.RuneCountInString(s)) <= cnt {
		return s
	}
	i := 0
	count := int64(0)
	for i < len(s) {
		if count == cnt {
			return s[:i]
		}
		_, size := utf8.DecodeRuneInString(s[i:])
		i += size
		count++
	}
	return s
}

func CosineSimilarity(vecA, vecB *mat.VecDense) float64 {
	dotProduct := mat.Dot(vecA, vecB)

	normA := math.Sqrt(mat.Dot(vecA, vecA))
	normB := math.Sqrt(mat.Dot(vecB, vecB))

	if normA == 0 || normB == 0 {
		return 0
	}

	return dotProduct / (normA * normB)
}
