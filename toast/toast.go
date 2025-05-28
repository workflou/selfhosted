package toast

import (
	"context"
	"net/http"
	"selfhosted/html"
)

type Toast struct {
	Type  string `json:"type"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

func New(toastType, title, body string) Toast {
	return Toast{
		Type:  toastType,
		Title: title,
		Body:  body,
	}
}

func Success(title, body string) Toast {
	return New("success", title, body)
}

func Error(title, body string) Toast {
	return New("error", title, body)
}

func Info(title, body string) Toast {
	return New("info", title, body)
}

func Warning(title, body string) Toast {
	return New("warning", title, body)
}

func (t Toast) Send(ctx context.Context, w http.ResponseWriter) error {
	return html.Notification(html.NotificationProps{
		Title: t.Title,
		Body:  t.Body,
	}).Render(ctx, w)
}
