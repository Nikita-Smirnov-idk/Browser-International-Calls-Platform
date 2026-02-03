package handlers

import (
	"log/slog"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

var e164Re = regexp.MustCompile(`^\+[1-9]\d{6,14}$`)

type VoiceHandler struct {
	tokenGenerator     TokenGenerator
	voicePublicBaseURL string
	dialCallerID       string
}

type TokenGenerator interface {
	GetToken(identity string, ttlSec int) (string, error)
}

func NewVoiceHandler(tokenGenerator TokenGenerator, voicePublicBaseURL, dialCallerID string) *VoiceHandler {
	return &VoiceHandler{
		tokenGenerator:     tokenGenerator,
		voicePublicBaseURL: voicePublicBaseURL,
		dialCallerID:       dialCallerID,
	}
}

func (h *VoiceHandler) Token(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "Authentication required",
		})
		return
	}
	if h.tokenGenerator == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error":   "voice_not_configured",
			"message": "Voice token is not configured",
		})
		return
	}
	token, err := h.tokenGenerator.GetToken(userID, 3600)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "token_failed",
			"message": "Failed to generate voice token",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *VoiceHandler) TwiML(c *gin.Context) {
	to := c.PostForm("To")
	if to == "" {
		to = c.Query("To")
	}
	slog.Info("twiml request from Twilio", "To", to, "method", c.Request.Method)
	if to == "" || !e164Re.MatchString(to) {
		slog.Warn("twiml invalid or missing To", "To", to)
		c.Data(http.StatusOK, "application/xml", []byte(`<?xml version="1.0" encoding="UTF-8"?><Response><Say language="en-US">Invalid or missing phone number.</Say><Hangup/></Response>`))
		return
	}
	escaped := escapeXML(to)
	var dialAttrs []string
	if h.dialCallerID != "" && e164Re.MatchString(h.dialCallerID) {
		dialAttrs = append(dialAttrs, `callerId="`+escapeXML(h.dialCallerID)+`"`)
	}
	if h.voicePublicBaseURL != "" {
		statusURL := strings.TrimSuffix(h.voicePublicBaseURL, "/") + "/api/voice/status"
		dialAttrs = append(dialAttrs, `statusCallback="`+escapeXML(statusURL)+`"`, `statusCallbackEvent="initiated ringing answered completed"`)
	}
	dialAttrStr := ""
	if len(dialAttrs) > 0 {
		dialAttrStr = " " + strings.Join(dialAttrs, " ")
	}
	twiml := `<?xml version="1.0" encoding="UTF-8"?><Response><Dial` + dialAttrStr + `><Number>` + escaped + `</Number></Dial></Response>`
	c.Header("Content-Type", "application/xml; charset=utf-8")
	c.String(http.StatusOK, twiml)
	slog.Info("twiml returned Dial", "To", to)
}

func (h *VoiceHandler) VoiceStatusCallback(c *gin.Context) {
	callSid := c.PostForm("CallSid")
	dialCallStatus := c.PostForm("DialCallStatus")
	callStatus := c.PostForm("CallStatus")
	to := c.PostForm("To")
	slog.Info("voice status callback from Twilio",
		"CallSid", callSid,
		"DialCallStatus", dialCallStatus,
		"CallStatus", callStatus,
		"To", to)
	c.Status(http.StatusNoContent)
}

func escapeXML(s string) string {
	const (
		amp  = "&amp;"
		lt   = "&lt;"
		gt   = "&gt;"
		quot = "&quot;"
		apos = "&#39;"
	)
	var b []byte
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '&':
			b = append(b, amp...)
		case '<':
			b = append(b, lt...)
		case '>':
			b = append(b, gt...)
		case '"':
			b = append(b, quot...)
		case '\'':
			b = append(b, apos...)
		default:
			b = append(b, s[i])
		}
	}
	return string(b)
}
