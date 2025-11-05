package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatRequest struct {
	Model       string        `json:"model"`
	Messages    []chatMessage `json:"messages"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
	Temperature float32       `json:"temperature,omitempty"`
}

type chatChoice struct {
	Index   int `json:"index"`
	Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
}

type chatResponse struct {
	ID      string       `json:"id"`
	Choices []chatChoice `json:"choices"`
}

func classifyProfile(conf, worry, human int) (string, string) {
	if conf >= 4 && worry >= 4 {
		return "Aware but Anxious", "You understand AI well, but you’re concerned about job impact. Focus on adapting and building human-AI collaboration skills."
	}
	if conf <= 2 {
		return "Low Confidence", "Start small. Try simple, hands-on AI tasks — your confidence will grow quickly through experience."
	}
	if conf >= 4 && worry <= 2 {
		return "Confident and Adaptive", "You’re ready to lead. Share what you know and help others understand AI’s potential."
	}
	if human >= 4 {
		return "Human-Centered Learner", "You value creativity and empathy — keep combining those with AI skills for the best of both worlds."
	}
	return "Curious Learner", "Stay curious. Keep exploring AI and how it fits with your strengths."
}

func groqAdvice(ctx context.Context, profile string) (string, error) {
	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("GROQ_API_KEY is not set")
	}

	endpoint := "https://api.groq.com/openai/v1/chat/completions"
	prompt := fmt.Sprintf("Give concise, positive advice (2–3 sentences) for a high school student whose AI mindset profile is '%s'.", profile)

	reqBody := chatRequest{
		Model: "openai/gpt-oss-20b",
		Messages: []chatMessage{
			{Role: "system", Content: "You are a concise, encouraging career coach for students."},
			{Role: "user", Content: prompt},
		},
		MaxTokens:   200,
		Temperature: 0.6,
	}

	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(reqBody); err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, buf)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("groq api status %d", resp.StatusCode)
	}

	var payload chatResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", err
	}

	if len(payload.Choices) == 0 {
		return "", fmt.Errorf("no choices returned")
	}
	return payload.Choices[0].Message.Content, nil
}

func mustAtoi(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/survey")
	})

	r.GET("/survey", func(c *gin.Context) {
		c.HTML(http.StatusOK, "survey.html", nil)
	})

	r.POST("/results", func(c *gin.Context) {
		conf := mustAtoi(c.PostForm("confidence"))
		worry := mustAtoi(c.PostForm("worry"))
		human := mustAtoi(c.PostForm("human_skills"))

		profile, localTip := classifyProfile(conf, worry, human)

		ctx, cancel := context.WithTimeout(c.Request.Context(), 15*time.Second)
		defer cancel()

		advice, err := groqAdvice(ctx, profile)
		if err != nil {
			log.Printf("groq error: %v", err)
			advice = localTip + "\n\n(Note: AI advice unavailable, showing local guidance.)"
		}

		c.HTML(http.StatusOK, "results.html", gin.H{
			"Profile": profile,
			"Advice":  advice,
			"Conf":    conf,
			"Worry":   worry,
			"Human":   human,
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
