package entity

type Bounce_mail_detail struct {
	Send_no         int    `json:"username"`
	Msg_id_ses      string `json:"msg_id_ses"`
	Addresser       string `json:"addresser"`
	Timestamp       int    `json:"timestamp"`
	Smtp_command    string `json:"smtp_command"`
	Subject         string `json:"subject"`
	Recipient       string `json:"recipient"`
	Destination     string `json:"destination"`
	Alias           string `json:"alias"`
	Sender_domain   string `json:"sender_domain"`
	Token           string `json:"token"`
	Smtp_agent      string `json:"smtp_agent"`
	Reply_code      string `json:"reply_code"`
	Soft_bounce     int    `json:"soft_bounce"`
	Rhost           string `json:"rhost"`
	Timezone_offset string `json:"timezone_offset"`
	Diagnostic_type string `json:"diagnostic_type"`
	Action          string `json:"action"`
	Listid          string `json:"listid"`
	Feedback_type   string `json:"feedback_type"`
	Diagnostic_code string `json:"diagnostic_code"`
	Reason          string `json:"reason"`
	Delivery_status string `json:"delivery_status"`
	Lhost           string `json:"lhost"`
}
