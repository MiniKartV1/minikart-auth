package email

import "fmt"

type Adapter struct{}

func NewAdapter() *Adapter {
	return &Adapter{}
}

func (email Adapter) SendWelcomeEmail() {
	fmt.Println("Sent Welcome Email")
}
func (email Adapter) SendResetPasswordEmail() {
	fmt.Println("Send Reset Password Email")
}
