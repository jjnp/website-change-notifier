package main

type Notifier interface {
	NotifyOfStart(message string, url string)
	NotifyOfChange(title string, message string, url string)
	SendSummary(title string, message string, url string)
}
