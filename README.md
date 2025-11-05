# AI Anxiety to Action Web App

A Go + Gin web application that helps students reflect on their confidence, anxiety, and attitudes toward AI.  
Based on survey data, the app provides a short personalized message encouraging positive engagement with AI and future learning.

## Features
- Simple three-question survey on AI confidence, job worry, and human skill values  
- Generates a student profile based on responses  
- Fetches short, personalized advice from the Groq API  
- Lightweight and fast (Go backend, minimal HTML/CSS frontend)

## Project Structure
ai-anxiety-webapp/
├── main.go
├── go.mod
├── templates/
│ ├── survey.html
│ └── results.html
├── static/
│ └── style.css
└── .air.toml (optional, for live reload)


## Setup

### Prerequisites
- Go 1.22+
- Groq API key  
  Sign up at [https://console.groq.com/keys](https://console.groq.com/keys)

### Installation
```bash
git clone https://github.com/<your-username>/ai-anxiety-to-action.git
cd ai-anxiety-to-action
go mod tidy

Running

Set your API key:

export GROQ_API_KEY=your_key_here


Start the server:

go run .


Or with Air (live reload):

air


Open the app:

http://localhost:8080/survey

Notes

If the Groq API is unavailable, the app falls back to local guidance.

No database is required; responses are processed in memory.

CSS is minimal for clarity and speed.
