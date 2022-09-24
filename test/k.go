package main

import "gopkg.in/gomail.v2"

func main() {
	m := gomail.NewMessage()
	m.SetHeader("From", "eininst@aliyun.com")
	m.SetHeader("To", "eininst@qq.com")

	//抄送
	//m.SetAddressHeader("Cc", "eininst@aliyun.com", "Dan")

	m.SetHeader("Subject", "Hello!")
	m.SetBody("text/html", "Hello <b>Bob</b> and <i>Cora</i>!")
	//m.Attach("/home/Alex/lolcat.jpg")

	d := gomail.NewDialer("smtp.qq.com", 465, "eininst@qq.com", "fbkhlrysrbkbbhjd")

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}
