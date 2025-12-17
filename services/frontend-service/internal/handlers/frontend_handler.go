package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/sirupsen/logrus"
)

type FrontendHandler struct {
	gameServiceURL string
	logger         *logrus.Logger
	client         *http.Client
}

func NewFrontendHandler(gameServiceURL string, logger *logrus.Logger) *FrontendHandler {
	return &FrontendHandler{
		gameServiceURL: gameServiceURL,
		logger:         logger,
		client:         &http.Client{},
	}
}

func (h *FrontendHandler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

func (h *FrontendHandler) Index(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html>
<head>
    <title>Portfolio Space Invaders</title>
    <style>
        body { font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; margin: 0; padding: 0; background: linear-gradient(135deg, #0c0c0c 0%, #1a1a2e 50%, #16213e 100%); color: #fff; min-height: 100vh; }
        .container { max-width: 1200px; margin: 0 auto; padding: 20px; }
        .header { text-align: center; margin-bottom: 30px; padding: 30px 0; background: rgba(255,255,255,0.05); border-radius: 15px; backdrop-filter: blur(10px); }
        .header h1 { font-size: 3em; margin: 0; background: linear-gradient(45deg, #00d4ff, #00ff88); -webkit-background-clip: text; -webkit-text-fill-color: transparent; text-shadow: 0 0 30px rgba(0,212,255,0.5); }
        .header p { font-size: 1.2em; margin: 10px 0; opacity: 0.9; }
        .cv-section { background: rgba(255,255,255,0.08); padding: 30px; margin: 25px 0; border-radius: 15px; backdrop-filter: blur(15px); border: 1px solid rgba(255,255,255,0.1); box-shadow: 0 8px 32px rgba(0,0,0,0.3); }
        .cv-section h2 { color: #00d4ff; font-size: 2em; margin-bottom: 20px; text-shadow: 0 0 10px rgba(0,212,255,0.3); }
        .cv-section h3 { color: #00ff88; margin: 25px 0 15px 0; }
        .game-container { text-align: center; }
        #gameCanvas { border: 3px solid #00d4ff; background: #000; border-radius: 10px; box-shadow: 0 0 30px rgba(0,212,255,0.4), inset 0 0 20px rgba(0,0,0,0.5); }
        .game-stats { display: flex; justify-content: center; gap: 30px; margin: 20px 0; }
        .stat-box { background: rgba(0,212,255,0.1); padding: 15px 25px; border-radius: 10px; border: 1px solid rgba(0,212,255,0.3); }
        .stat-label { font-size: 0.9em; opacity: 0.8; margin-bottom: 5px; }
        .stat-value { font-size: 1.8em; font-weight: bold; color: #00d4ff; }
        .controls { margin: 25px 0; }
        button { padding: 12px 25px; margin: 8px; background: linear-gradient(45deg, #00d4ff, #0099cc); color: #fff; border: none; border-radius: 8px; cursor: pointer; font-weight: bold; transition: all 0.3s ease; box-shadow: 0 4px 15px rgba(0,212,255,0.3); }
        button:hover { transform: translateY(-2px); box-shadow: 0 6px 20px rgba(0,212,255,0.5); background: linear-gradient(45deg, #00ff88, #00cc66); }
        button:active { transform: translateY(0); }
        .skills { display: grid; grid-template-columns: repeat(auto-fit, minmax(280px, 1fr)); gap: 20px; }
        .skill { background: rgba(0,255,136,0.1); padding: 20px; border-radius: 10px; border: 1px solid rgba(0,255,136,0.2); transition: transform 0.3s ease; }
        .skill:hover { transform: translateY(-5px); box-shadow: 0 10px 25px rgba(0,255,136,0.2); }
        .experience { margin: 25px 0; }
        .experience h4 { color: #00ff88; margin: 20px 0 8px 0; font-size: 1.3em; }
        .experience ul { margin: 15px 0; padding-left: 25px; }
        .experience li { margin: 8px 0; line-height: 1.6; }
        .game-instructions { background: rgba(255,255,255,0.05); padding: 20px; border-radius: 10px; margin: 20px 0; border-left: 4px solid #00d4ff; }
    </style>
</head>
<body>
    <div class="container">
    <div class="header">
        <h1>SIHLE NDLOVU</h1>
        <p>AWS-Certified Cloud Engineer | DevOps & Infrastructure Specialist</p>
        <p>📍 Sandton, Gauteng | 📧 ndlovu.code@outlook.com | 📱 0839578644</p>
    </div>
    
    <div class="cv-section">
        <h2>Professional Summary</h2>
        <p>AWS-certified Cloud Engineer with hands-on experience designing, deploying, and operating secure, scalable cloud infrastructure on AWS. Strong expertise in Infrastructure as Code (Terraform), containerized workloads (ECS/Fargate), CI/CD automation, cloud monitoring, security best practices, and cost optimization.</p>
        
        <h3>Work Experience</h3>
        <div class="experience">
            <h4>Cloud Engineer (DevOps & Infrastructure) | Synthesis Software Technologies</h4>
            <p><em>February 2025 – Present | Sandton, South Africa</em></p>
            <ul>
                <li>Led AWS infrastructure provisioning using Terraform across DEV, QA, and Production environments</li>
                <li>Deployed production microservices using Amazon ECS with Fargate, implementing container orchestration and auto-scaling</li>
                <li>Built CI/CD pipelines using GitHub Actions and AWS services, automating deployments</li>
                <li>Implemented centralized monitoring using Amazon CloudWatch, improving incident response</li>
            </ul>
            
            <h4>Solutions Architect (DevOps & Cloud) | CloudZA</h4>
            <p><em>January 2024 – January 2025 | Bellville, Western Cape</em></p>
            <ul>
                <li>Designed cloud-native and serverless architectures using AWS Lambda, API Gateway, DynamoDB</li>
                <li>Built CI/CD pipelines for container builds, security scanning, and multi-environment deployments</li>
                <li>Established production logging and monitoring using CloudWatch and Container Insights</li>
            </ul>
        </div>
        
        <h3>AWS Certifications</h3>
        <div class="skills">
            <div class="skill">🏆 AWS Cloud Practitioner</div>
            <div class="skill">🏆 AWS Solutions Architect Associate</div>
            <div class="skill">🏆 AWS Data Engineer Associate</div>
            <div class="skill">🏆 AWS Developer Associate</div>
            <div class="skill">🏆 AWS Security Specialty</div>
            <div class="skill">🏆 HashiCorp Terraform Associate</div>
        </div>
        
        <h3>Technical Skills</h3>
        <div class="skills">
            <div class="skill"><strong>Cloud & Infrastructure:</strong> AWS (EC2, S3, VPC, RDS, ECS, Lambda), Terraform, CloudFormation</div>
            <div class="skill"><strong>DevOps & Automation:</strong> GitHub Actions, Docker, CI/CD, Monitoring & Alerting</div>
            <div class="skill"><strong>Security & Compliance:</strong> IAM, KMS, CloudTrail, GuardDuty, ISO 27001/SOC2</div>
            <div class="skill"><strong>Development:</strong> Python, Java, C#, Go, Node.js</div>
        </div>
        
        <h3>Education</h3>
        <p><strong>Bachelor of Commerce Honours Information Systems and Technology</strong><br>
        University of KwaZulu-Natal, 2022</p>
    </div>
    
    <div class="cv-section game-container">
        <h2>🎮 Interactive Cloud Architecture Demo</h2>
        <p>This Space Invaders game demonstrates cloud-native microservices architecture:</p>
        <ul style="text-align: left; max-width: 600px; margin: 0 auto;">
            <li><strong>Go Microservices:</strong> Game logic and frontend services</li>
            <li><strong>Docker Containers:</strong> Multi-stage builds with distroless images</li>
            <li><strong>ECS Fargate:</strong> Serverless container orchestration</li>
            <li><strong>Prometheus & Grafana:</strong> Real-time monitoring and metrics</li>
            <li><strong>Terraform IaC:</strong> Infrastructure automation</li>
            <li><strong>Cost Optimization:</strong> ~$80/month vs $150+ traditional setup</li>
        </ul>
        <div class="game-stats">
            <div class="stat-box">
                <div class="stat-label">Score</div>
                <div class="stat-value" id="score">0</div>
            </div>
            <div class="stat-box">
                <div class="stat-label">Level</div>
                <div class="stat-value" id="level">1</div>
            </div>
            <div class="stat-box">
                <div class="stat-label">Status</div>
                <div class="stat-value" id="gameStatus">Ready</div>
            </div>
        </div>
        <canvas id="gameCanvas" width="800" height="600"></canvas>
        <div class="controls">
            <button onclick="startGame()">🎮 New Game</button>
            <button onmousedown="moveLeft()" onmouseup="stopMove()">← Move Left</button>
            <button onmousedown="moveRight()" onmouseup="stopMove()">Move Right →</button>
            <button onclick="shoot()">🚀 Fire Bullet</button>
        </div>
        <div class="game-instructions">
            <p><strong>🎮 Controls:</strong> Arrow keys to move, Spacebar to shoot</p>
            <p><strong>🎯 Objective:</strong> Destroy all enemies to advance to the next level!</p>
            <p><strong>⚡ Features:</strong> Progressive difficulty, level bonuses, real-time scoring</p>
        </div>
    </div>
    
    <script>
        let currentGame = null;
        let canvas = document.getElementById('gameCanvas');
        let ctx = canvas.getContext('2d');
        let gameLoop = null;
        let moveInterval = null;
        
        document.addEventListener('keydown', handleKeyPress);
        
        async function startGame() {
            const response = await fetch('/api/game/start', { method: 'POST' });
            const data = await response.json();
            currentGame = data.game_id;
            
            if (gameLoop) clearInterval(gameLoop);
            gameLoop = setInterval(updateGame, 100);
        }
        
        async function makeMove(action, direction = '') {
            if (!currentGame) return;
            
            const response = await fetch('/api/game/move', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ game_id: currentGame, action: action, direction: direction })
            });
            
            if (response.ok) {
                const game = await response.json();
                renderGame(game);
            }
        }
        
        function moveLeft() {
            if (moveInterval) clearInterval(moveInterval);
            moveInterval = setInterval(() => makeMove('move', 'left'), 100);
        }
        
        function moveRight() {
            if (moveInterval) clearInterval(moveInterval);
            moveInterval = setInterval(() => makeMove('move', 'right'), 100);
        }
        
        function stopMove() {
            if (moveInterval) {
                clearInterval(moveInterval);
                moveInterval = null;
            }
        }
        
        function shoot() {
            makeMove('shoot');
        }
        
        function handleKeyPress(e) {
            switch(e.code) {
                case 'ArrowLeft':
                    makeMove('move', 'left');
                    break;
                case 'ArrowRight':
                    makeMove('move', 'right');
                    break;
                case 'Space':
                    e.preventDefault();
                    shoot();
                    break;
            }
        }
        
        async function updateGame() {
            if (!currentGame) return;
            
            // Trigger game update on server
            await fetch('/api/game/move', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ game_id: currentGame, action: 'update' })
            });
            
            // Get updated game state
            const response = await fetch('/api/game/status?game_id=' + currentGame);
            if (response.ok) {
                const game = await response.json();
                renderGame(game);
            }
        }
        
        function renderGame(game) {
            // Clear canvas
            ctx.fillStyle = '#000';
            ctx.fillRect(0, 0, canvas.width, canvas.height);
            
            // Update score and level
            document.getElementById('score').textContent = game.score;
            document.getElementById('level').textContent = game.level;
            document.getElementById('gameStatus').textContent = game.status === 'active' ? 'Playing' : game.status.toUpperCase();
            
            // Draw player (green rectangle)
            if (game.player && game.player.active) {
                ctx.fillStyle = '#00ff00';
                ctx.fillRect(game.player.position.x, game.player.position.y, 30, 20);
                
                // Player ship details
                ctx.fillStyle = '#ffffff';
                ctx.fillRect(game.player.position.x + 12, game.player.position.y - 5, 6, 8);
            }
            
            // Draw enemies (red rectangles)
            ctx.fillStyle = '#ff0000';
            if (game.enemies) {
                game.enemies.forEach(enemy => {
                    if (enemy.active) {
                        ctx.fillRect(enemy.position.x, enemy.position.y, 25, 20);
                        // Enemy details
                        ctx.fillStyle = '#ffffff';
                        ctx.fillRect(enemy.position.x + 10, enemy.position.y + 20, 5, 8);
                        ctx.fillStyle = '#ff0000';
                    }
                });
            }
            
            // Draw bullets (yellow lines)
            ctx.fillStyle = '#ffff00';
            if (game.bullets) {
                game.bullets.forEach(bullet => {
                    if (bullet.active) {
                        ctx.fillRect(bullet.position.x + 2, bullet.position.y, 4, 12);
                    }
                });
            }
            
            // Draw game borders
            ctx.strokeStyle = '#ffffff';
            ctx.lineWidth = 2;
            ctx.strokeRect(0, 0, canvas.width, canvas.height);
            
            // Check game status
            if (game.status === 'lost') {
                ctx.fillStyle = 'rgba(255, 0, 0, 0.8)';
                ctx.fillRect(0, 0, canvas.width, canvas.height);
                ctx.fillStyle = '#ffffff';
                ctx.font = 'bold 48px Arial';
                ctx.textAlign = 'center';
                ctx.fillText('GAME OVER', canvas.width/2, canvas.height/2 - 20);
                ctx.font = '24px Arial';
                ctx.fillText('Final Score: ' + game.score + ' | Level: ' + game.level, canvas.width/2, canvas.height/2 + 30);
                ctx.fillText('Click New Game to restart', canvas.width/2, canvas.height/2 + 60);
                if (gameLoop) clearInterval(gameLoop);
            }
            
            ctx.textAlign = 'left'; // Reset text alignment
        }
    </script>
</body>
</html>`
	
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

func (h *FrontendHandler) ProxyStartGame(w http.ResponseWriter, r *http.Request) {
	h.proxyRequest(w, r, "/game/start")
}

func (h *FrontendHandler) ProxyMove(w http.ResponseWriter, r *http.Request) {
	h.proxyRequest(w, r, "/game/move")
}

func (h *FrontendHandler) ProxyStatus(w http.ResponseWriter, r *http.Request) {
	gameID := r.URL.Query().Get("game_id")
	targetURL := h.gameServiceURL + "/game/status?game_id=" + url.QueryEscape(gameID)
	
	resp, err := h.client.Get(targetURL)
	if err != nil {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func (h *FrontendHandler) proxyRequest(w http.ResponseWriter, r *http.Request, path string) {
	targetURL := h.gameServiceURL + path
	
	req, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	
	req.Header = r.Header
	
	resp, err := h.client.Do(req)
	if err != nil {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}