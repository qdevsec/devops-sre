package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/yaml.v3"
)

// view go's internal metrics (memory usage, garbage collection times, etc.)
// http://localhost:8080/metrics

const slackWebhookURL = "https://hooks.slack.com/services/YOUR/WEBHOOK/URL"

type SlackMessage struct {
	Text string `json:"text"`
}

// config structures for parsing config.yaml
type Config struct {
	Rules []RemediationRule `yaml:"rules"`
}

type RemediationRule struct {
	AlertName			string `yaml:"alert_name"`
	Script				string `yaml:"script"`
	VerifyCmd           string `yaml:"verify_cmd"`
	TimeoutSeconds		int	   `yaml:"timeout_seconds"`
	CooldownMinutes		int	   `yaml:"cooldown_minutes"`
}

// capture the finalized state of a closed-loop remediation attempt
type PipelineResult struct {
	Verified bool
	Message  string
}

// Global variable to store our configuration rules
var (
	appConfig Config
	// executionHistory keeps track when an alert was last executed: map[alert_name]last_run_time
	executionHistory = make(map[string]time.Time)
	// historyMutex protects executionHistory from concurrent read/write map panics
	historyMutex sync.Mutex

	// Prometheus Metrics Definitions
	remediationsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "sre_remediations_total",
			Help: "Total number of remediation executions handled by this engine.",
		},
		[]string{"alert_name", "status"}, // Labels to segment our metric
	)

	guardrailSkipsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "sre_guardrail_skips_total",
			Help: "Total number of remediations blocked by the cooldown guardrail.",
		},
		[]string{"alert_name"},
	)

	remediationDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "sre_remediation_duration_seconds",
			Help: "How much time remediation scripts take to complete.",
			Buckets: prometheus.DefBuckets, // default buckets like 0.1s
		},
		[]string{"alert_name"},
	)
)

// AlertmanagerPayload JSON structure mirros the Prometheus Alertmanager
type AlertmanagerPayload struct {
	Status string  `json:"status"` // status is either firing or resolved
	Alerts []Alert `json:"alerts"`
}

type Alert struct {
	Status          string				`json:"status"`
	Labels			map[string]string	`json:"labels"`
//	Annotations		map[string]string   `json:"annotations"`
//	StartsAt        time.Time			`json:"startsAt"`
//	GeneratorURL	string				`json:"generatorURL"`
}

func init() {
	// register custom metrics with Prometheus's default registry
	prometheus.MustRegister(remediationsTotal)
	prometheus.MustRegister(guardrailSkipsTotal)
	prometheus.MustRegister(remediationDuration)
}

func main() {

	// 1. Load the configuration file at startup
	err := loadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	log.Printf("Loaded %d remediation rules from config.yaml", len(appConfig.Rules))

	// Webhook endpoint for Alertmanager
	http.HandleFunc("/webhook", handleWebhook)

	// Expose the standard Prometheus metrics endpoint
	http.Handle("/metrics", promhttp.Handler())

	fmt.Println("SRE Auto-Remediation Engine running on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// retrieve the secret injected by Kubernetes
func GetSlackWebhook() string {
	// Look up the env var injected by your secretKeyRef
	webhookURL, exists := os.LookupEnv("SLACK_WEBHOOK_URL")

	// If missing, return an empty string gracefully instead of killing the engine execution context
	if !exists || webhookURL == "" {
		return ""
	}
	return webhookURL
}

func sendSlackNotification(message string) {
	webhookURL := GetSlackWebhook()

	// OPTIONAL TOGGLE CHECK: If the secret variable wasn't injected or is left blank,
	// smoothly bypass the entire notification engine block without emitting errors.
	if webhookURL == "" {
		// Log it quietly locally so you still see what happened in standard stdout
		log.Printf("[Notification Bypass] %s", message)
		return
	}

	payload := SlackMessage{Text: message}
	bytesData, _ := json.Marshal(payload)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", webhookURL, strings.NewReader(string(bytesData)))
	if err != nil {
		log.Printf("Failed to create Slack request struct: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Failed to send Slack notification over HTTP: %v", err)
		return
	}
	resp.Body.Close()
}

func loadConfig(filename string) error {
	file, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(file, &appConfig)
}

func handleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload AlertmanagerPayload
	// err := json.NewDecoder(r.Body).Decode(&payload)
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Printf("Error decoding JSON payload: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// process each alert inside the payload
	for _, alert := range payload.Alerts {
		alertName := alert.Labels["alertname"]
		// status := alert.Status

		// log.Printf("Received Alert: [%s] Status: [%s]", alertName, status)

		// only remediate if the alert is firing
		if alert.Status == "firing" {
			// find the matching rule dynamically from config
			go findAndExecuteRemediation(alertName)
		}
	}

	w.WriteHeader(http.StatusOK)
	// w.Write([]byte("Alert processed successfully"))
}

// func triggerRemediation(alertName string, labels map[string]string) {
// 	log.Printf("Analyzing remediation path for alert: %s", alertName)

// 	switch alertName {
// 	case "DiskRunningFull":
// 		log.Println("Target Action: Running disk cleanup script...")
// 		executeScript("./remediate_disk.sh")

// 	case "ServiceDown":
// 		service := labels["service"]
// 		log.Printf("Target Action: Attempting to restart service: %s...", service)
// 		executeScript("./remediate_service.sh")
	
// 	default:
// 		log.Printf("No auto-remediation rule defined for alert: %s. Forwarding to on-call engineer.", alertName)
// 	}
// }

func findAndExecuteRemediation(alertName string) {
	var targetRule *RemediationRule

	// loop through loaded rules to find a match
	for _, rule := range appConfig.Rules {
		if rule.AlertName == alertName {
			targetRule = &rule
			break
		}
	}

	if targetRule == nil {
		log.Printf("No auto-remediation rule defined in config for alert: %s.", alertName)
		return
	}

	// This guardrail will check cooldown period
	historyMutex.Lock()
	lastExecuted, exists := executionHistory[alertName]
	historyMutex.Unlock()

	if exists {
		cooldownDuration := time.Duration(targetRule.CooldownMinutes) * time.Minute
		if time.Since(lastExecuted) < cooldownDuration {
			remaining := cooldownDuration - time.Since(lastExecuted)
			log.Printf("[GUARDRAIL] Remediation skipped for %s. Cooldown active for another %v.", alertName, remaining.Round(time.Second))
			
			// Increment Guardrail Skip Metric
			guardrailSkipsTotal.WithLabelValues(alertName).Inc()
			return
		}
	}

	// update the execution history timestamp before running
	historyMutex.Lock()
	executionHistory[alertName] = time.Now()
	historyMutex.Unlock()

	log.Printf("Match found! Script: %s | Timeout: %ds", targetRule.Script, targetRule.TimeoutSeconds)
	// executeScript(targetRule.Script, targetRule.TimeoutSeconds, alertName)
	executeClosedLoopPipeline(*targetRule)
}

func executeClosedLoopPipeline(rule RemediationRule) {
	startTime := time.Now()
	var stdout, stderr bytes.Buffer

	// 1. Core Remediation Step with Enforced Config Timeout Context
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(rule.TimeoutSeconds)*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "/bin/bash", rule.Script)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	duration := time.Since(startTime).Seconds()
	remediationDuration.WithLabelValues(rule.AlertName).Observe(duration)

	// Check if the script execution timed out
	if ctx.Err() == context.DeadlineExceeded {
		msg := fmt.Sprintf("🚨 *Remediation TIMEOUT*: Script `%s` for alert `%s` was killed.", rule.Script, rule.AlertName)
		log.Println(msg)
		sendSlackNotification(msg)
		remediationsTotal.WithLabelValues(rule.AlertName, "timeout").Inc()
		return
	}

	// Handle standard script runtime execution errors
	if err != nil {
		msg := fmt.Sprintf("❌ *Remediation Failed*: Alert `%s` failed to execute.\nError: `%v` | Stderr: %s", rule.AlertName, err, strings.TrimSpace(stderr.String()))
		log.Println(msg)
		sendSlackNotification(msg)
		remediationsTotal.WithLabelValues(rule.AlertName, "failed").Inc()
		return
	}

	// 2. Closed-Loop Post-Check Verification Phase
	if rule.VerifyCmd != "" {
		log.Printf("[VERIFICATION] Running validation query for alert: %s", rule.AlertName)

		// Sub-timeout context dedicated to the automated verification step
		verifyCtx, verifyCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer verifyCancel()

		verifyCmd := exec.CommandContext(verifyCtx, "/bin/bash", "-c", rule.VerifyCmd)
		var vStdout, vStderr bytes.Buffer
		verifyCmd.Stdout = &vStdout
		verifyCmd.Stderr = &vStderr

		if verifyErr := verifyCmd.Run(); verifyErr != nil {
			msg := fmt.Sprintf("⚠️ *Remediation Warning*: Script executed for alert `%s`, but the system *failed* verification verification check.\nError output: `%s`", rule.AlertName, strings.TrimSpace(vStderr.String()))
			log.Println(msg)
			sendSlackNotification(msg)
			remediationsTotal.WithLabelValues(rule.AlertName, "verification_failed").Inc()
			return
		}

		msg := fmt.Sprintf("✅ *Remediation Success*: Alert `%s` resolved dynamically using `%s` and verified healthy in %.2fs. Check output: `%s`", rule.AlertName, rule.Script, duration, strings.TrimSpace(vStdout.String()))
		log.Println(msg)
		sendSlackNotification(msg)
		remediationsTotal.WithLabelValues(rule.AlertName, "success").Inc()
		return
	}

	// Fallback success path if no explicit verification step was designated
	msg := fmt.Sprintf("✅ *Remediation Complete*: Alert `%s` executed `%s` cleanly in %.2fs (no validation configured).", rule.AlertName, rule.Script, duration)
	log.Printf(msg)
	sendSlackNotification(msg)
	remediationsTotal.WithLabelValues(rule.AlertName, "success").Inc()
}

// func executeScript(scriptPath string, timeoutSeconds int, alertName string) {
// 	startTime := time.Now()
// 	// 1. Create a context with a strict 5s timeout
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSeconds)*time.Second)
// 	defer cancel() // make sure resources are cleaned up when the function returns

// 	// 2. Pass the context to the command
// 	// using CommandContext instead of Command so that if the script hands Go will 
// 	// stop process after a set duration, freeing up resources
// 	cmd := exec.CommandContext(ctx, "/bin/bash", scriptPath)
// 	output, err := cmd.CombinedOutput()

// 	// Calculate and record script duration
// 	duration := time.Since(startTime).Seconds()
// 	remediationDuration.WithLabelValues(alertName).Observe(duration)

// 	// 3. Check if the error is due to a timeout
// 	if ctx.Err() == context.DeadlineExceeded {
// 		log.Printf("Remediation TIMEOUT: Script %s killed.", scriptPath)
// 		remediationsTotal.WithLabelValues(alertName, "timeout").Inc()
// 		return
// 	}
// 	// execute bash script on Ubuntu
// 	// cmd := exec.Command("/bin/bash", scriptPath)
// 	// output, err := cmd.CombinedOutput()

// 	// handle other execution errors
// 	if err != nil {
// 		log.Printf("Remediation failed: %v | Output %s", err, string(output))
// 		remediationsTotal.WithLabelValues(alertName, "failed").Inc()
// 		return
// 	}
// 	log.Printf("Remediation script executed successfully! Duration: %.2fs", duration)
// 	remediationsTotal.WithLabelValues(alertName, "success").Inc()

// 	/*********************************************
// 	UNCOMMENT code below FOR SLACK webhook to run
// 	update slackWebhookURL
// 	**********************************************/
	
// 	// if ctx.Err() == context.DeadlineExceeded {
// 	// 	msg := fmt.Sprintf(" *Remediation TIMEOUT*: Script `%s` for alert `%s` was killed.", scriptPath, alertName)
// 	// 	log.Println(msg)
// 	// 	sendSlackNotification(msg)
// 	// 	remediationsTotal.WithLabelValues(alertName, "timeout").Inc()
// 	// 	return
// 	// }

// 	// if err != nil {
// 	// 	msg := fmt.Sprintf("*Remediation Failed*: Alert `%s` failed.\nError: `%v`", alertName, err)
// 	// 	log.Println(msg)
// 	// 	sendSlackNotification(msg)
// 	// 	remediationsTotal.WithLabelValues(alertName, "failed").Inc()
// 	// 	return
// 	// }

// 	// msg := fmt.Sprintf("*Remediation Success*: Alert `%s` resolved dynamically using `%s` in %.2fs.", alertName, scriptPath, duration)
// 	// log.Println(msg)
// 	// sendSlackNotification(msg)
// 	// remediationsTotal.WithLabelValues(alertName, "success").Inc()
// }

